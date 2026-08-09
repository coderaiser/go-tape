package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/coderaiser/go-tape"

	"github.com/BurntSushi/toml"
	"github.com/coderaiser/go-coverage"
	tapeast "github.com/coderaiser/go-tape/internal/ast"
	"github.com/coderaiser/go-tape/internal/config"
	"github.com/coderaiser/go-tape/internal/formatter"
	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/runner"
	"github.com/coderaiser/go-tape/internal/state"
	"github.com/coderaiser/go-tape/internal/tapeconfig"
)

//go:embed help.toml
var helpToml []byte

type helpConfig struct {
	Usage   struct{ Text string }
	Options []struct {
		Flag string
		Desc string
		Note string
	}
	Env []struct {
		Name string
		Desc string
	}
}

func loadUsage() string {
	var cfg helpConfig
	if err := toml.Unmarshal(helpToml, &cfg); err != nil {
		return `Usage: tape [options] [path]`
	}
	var sb strings.Builder
	sb.WriteString(cfg.Usage.Text + "\n\nOptions:\n")
	for _, o := range cfg.Options {
		fmt.Fprintf(&sb, "  %-28s %s\n", o.Flag, o.Desc)
		if o.Note != "" {
			fmt.Fprintf(&sb, "                               %s\n", o.Note)
		}
	}
	sb.WriteString("\nEnvironment variables:\n")
	for _, e := range cfg.Env {
		fmt.Fprintf(&sb, "  %-36s %s\n", e.Name, e.Desc)
	}
	sb.WriteString("\n")
	return sb.String()
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	covOpts, args := preprocess(args)

	flags := flag.NewFlagSet("tape", flag.ContinueOnError)
	format := flags.String("f", "", "output format: tap|progress-bar|short|fail|time|json-lines")
	flags.StringVar(format, "format", "", "output format (alias for -f)")

	var help bool
	flags.BoolVar(&help, "h", false, "display this help and exit")
	flags.BoolVar(&help, "help", false, "display this help and exit")

	var ver bool
	flags.BoolVar(&ver, "v", false, "output version information and exit")
	flags.BoolVar(&ver, "version", false, "output version information and exit")

	noCheckScopes := flags.Bool("no-check-scopes", false, "do not check scope format")
	noCheckAssertions := flags.Bool("no-check-assertions-count", false, "do not check assertion count")
	noCheckDuplicates := flags.Bool("no-check-duplicates", false, "do not check for duplicates")
	if err := flags.Parse(args); err != nil {
		if _, werr := fmt.Fprintf(stderr, "tape: %v\n", err); werr != nil {
			log.Fatal(werr)
		}
		return 2
	}

	if help {
		if _, err := fmt.Fprint(stdout, loadUsage()); err != nil {
			log.Fatal(err)
		}
		return 0
	}
	if ver {
		version := tape.TapeVersionLine()
		if _, err := fmt.Fprintln(stdout, version); err != nil {
			log.Fatal(err)
		}
		return 0
	}
	if *noCheckScopes {
		if err := os.Setenv("TAPE_CHECK_SCOPES", "0"); err != nil {
			_, _ = fmt.Fprintf(stderr, "tape: setenv: %v\n", err)
			return 1
		}
	}
	if *noCheckAssertions {
		if err := os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0"); err != nil {
			_, _ = fmt.Fprintf(stderr, "tape: setenv: %v\n", err)
			return 1
		}
	}

	path := "./..."
	if flags.NArg() > 0 {
		path = flags.Arg(0)
	}
	dir := "."

	tcfg := tapeconfig.Load(".tape.toml")
	exclude := tcfg.Test.Exclude

	// check duplicates before running
	if !*noCheckDuplicates && config.CheckDuplicates() {
		dups, err := tapeast.FindDuplicates(dir, exclude)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "tape: scan duplicates: %v\n", err)
			return 1
		}
		if len(dups) > 0 {
			f := formatter_tap.New()
			count := 0
			for _, d := range dups {
				first := d.Locations[0]
				second := d.Locations[1]
				firstURI := fmt.Sprintf("file://%s:%d:1", first.File, first.Line)
				secondURI := fmt.Sprintf("file://%s:%d:1", second.File, second.Line)
				message := fmt.Sprintf("Duplicate at %s", firstURI)
				at := fmt.Sprintf("at %s", secondURI)
				stack := fmt.Sprintf("Error: Duplicate at %s\n    at findDuplicates (tape)", firstURI)
				count++
				fmt.Fprint(stdout, f.Fail(count, message, "fail", nil, nil, "", at, stack))
			}
			fmt.Fprintf(stdout, "\n1..%d\n# tests %d\n# pass 0\n# fail %d\n\n",
				count, count, count)
			return 1
		}
	}

	// count tests for progress bar total
	total, err := tapeast.CountTestsInTestFiles(dir, exclude)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "tape: count tests: %v\n", err)
		return 1
	}

	// find Only calls — restrict run if any found
	onlyCalls, err := tapeast.FindOnlyCalls(dir, exclude)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "tape: scan Only: %v\n", err)
		return 1
	}

	goArgs := []string{"test", "-json", "-v", path}
	if pattern := tapeast.BuildRunPattern(onlyCalls); pattern != "" {
		goArgs = append(goArgs, "-run", pattern)
	}

	var coverTmp *os.File
	if covOpts.enabled {
		var terr error
		coverTmp, terr = os.CreateTemp("", "tape-coverage-*.out")
		if terr != nil {
			_, _ = fmt.Fprintf(stderr, "tape: create coverprofile: %v\n", terr)
			return 1
		}
		defer os.Remove(coverTmp.Name())
		defer coverTmp.Close()
		goArgs = append(goArgs, "-coverprofile="+coverTmp.Name(), "-covermode=atomic")
	}

	f := formatter.New(*format, stdout, total)
	store, err := state.New()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "tape: init state: %v\n", err)
		return 1
	}
	r := runner.New(runner.NewOSExecutor())

	ch, err := r.Run(goArgs...)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "tape: %v\n", err)
		return 1
	}

	for event := range ch {
		// go test -json emits events for both the outer TestXxx wrapper and the
		// inner t.Run subtest. Only the subtest (containing "/") is a tape.Test
		// call — skip the wrapper to avoid double-counting.
		if event.Test != "" && !strings.Contains(event.Test, "/") {
			continue
		}
		if _, err := store.Apply(event); err != nil {
			_, _ = fmt.Fprintf(stderr, "tape: apply event: %v\n", err)
			continue
		}
		f.FromEvent(event)
	}

	passed, failed, _ := store.Summary()
	buildFailedCount := store.BuildFailedCount()

	// skipped = declared tests that never ran (Skip calls + Only filtering).
	// If any package failed to build, we can't attribute tests accurately —
	// suppress skipped entirely so they don't show as skipped.
	skipped := 0
	if buildFailedCount == 0 {
		skipped = total - len(passed) - len(failed)
		if skipped < 0 {
			skipped = 0
		}
	}

	// When Only calls are present, recompute with full name list
	if len(onlyCalls) > 0 {
		allNames, err := tapeast.FindAllTestNames(dir, exclude)
		if err == nil {
			if err := store.MarkSkipped(allNames); err != nil {
				_, _ = fmt.Fprintf(stderr, "tape: mark skipped: %v\n", err)
				return 1
			}
			passed, failed, _ = store.Summary()
			skipped = total - len(passed) - len(failed)
			if skipped < 0 {
				skipped = 0
			}
		}
	}

	f.End(len(passed), len(failed), skipped)

	if covOpts.enabled && len(failed) == 0 && buildFailedCount == 0 {
		reportPath := ""
		if covOpts.report {
			reportPath = covOpts.path
		}
		if _, err := coverTmp.Seek(0, 0); err != nil {
			_, _ = fmt.Fprintf(stderr, "tape: seek coverprofile: %v\n", err)
			return 1
		}
		if err := coverage.ProcessProfileWithConfig(coverTmp, covOpts.format, reportPath, tcfg.Coverage.Exclude, stdout); err != nil {
			if !errors.Is(err, coverage.ErrUncovered) {
				_, _ = fmt.Fprintf(stderr, "tape: coverage: %v\n", err)
				return 1
			}
			return 1
		}
		if _, err := fmt.Fprintln(stdout, "💪 coverage 100%, good job!"); err != nil {
			return 1
		}
	}

	if len(failed) > 0 || buildFailedCount > 0 {
		return 1
	}
	if config.CheckSkipped() && skipped > 0 {
		return 5
	}
	return 0
}
