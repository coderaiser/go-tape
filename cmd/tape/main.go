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
	"github.com/coderaiser/go-tape/internal/stream"
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
	dir := dirFromPattern(path)

	tcfg := tapeconfig.Load(dir)
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
			for i, d := range dups {
				first := d.Locations[0]
				second := d.Locations[1]
				firstURI := fmt.Sprintf("file://%s:%d:1", first.File, first.Line)
				secondURI := fmt.Sprintf("file://%s:%d:1", second.File, second.Line)
				message := fmt.Sprintf("duplicate test name: %q", d.Name)
				at := fmt.Sprintf("at %s (first)\n    at %s (second)", firstURI, secondURI)
				stack := fmt.Sprintf("Error: duplicate test name: %q\n    at findDuplicates (tape)", d.Name)
				fmt.Fprint(stdout, f.Event(stream.Event{
					Type:       stream.TypeFail,
					Count:      i + 1,
					Test:       d.Name,
					Message:    message,
					Operator:   "fail",
					At:         at,
					ErrorStack: stack,
				}))
			}
			fmt.Fprintf(stdout, "\n1..%d\n# tests %d\n# pass 0\n# fail %d\n\n",
				len(dups), len(dups), len(dups))
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

	d := formatter.New(*format, stdout, total)

	ch, err := stream.Run(total, goArgs...)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "tape: %v\n", err)
		return 1
	}

	var lastCount, lastFailed int
	buildFailed := 0
	var buildErrors []stream.Event

	for e := range ch {
		if e.Type == stream.TypeTestEnd {
			lastCount = e.Count
			lastFailed = e.Failed
		}
		if e.Type == stream.TypeBuildError {
			buildFailed++
			buildErrors = append(buildErrors, e)
		}
		d.Emit(e)
	}

	passedCount := lastCount - lastFailed

	// skipped = declared tests that never ran.
	// If any package failed to build, suppress skipped count.
	skipped := 0
	if buildFailed == 0 {
		skipped = total - lastCount
		if skipped < 0 {
			skipped = 0
		}
	}

	// When Only calls are present, recompute skipped with full name list.
	if len(onlyCalls) > 0 {
		allNames, err := tapeast.FindAllTestNames(dir, exclude)
		if err == nil {
			skipped = total - lastCount
			if skipped < 0 {
				skipped = 0
			}
			_ = allNames
		}
	}

	d.End(passedCount, lastFailed, skipped)

	if covOpts.enabled && lastFailed == 0 && buildFailed == 0 {
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

	if buildFailed > 0 {
		_, _ = fmt.Fprintf(stderr, "tape: build failed — tests did not run\n\n")
		for _, e := range buildErrors {
			_, _ = fmt.Fprintf(stderr, "  %s:\n%s\n", e.Package, e.Output)
		}
		return 1
	}
	if lastFailed > 0 {
		return 1
	}
	if config.CheckSkipped() && skipped > 0 {
		return 5
	}
	return 0
}
