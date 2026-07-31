package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	tapeast "github.com/coderaiser/go-tape/internal/ast"
	"github.com/coderaiser/go-tape/internal/formatter"
	"github.com/coderaiser/go-tape/internal/runner"
	"github.com/coderaiser/go-tape/internal/state"
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

const version = "1.0.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
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

	// check duplicates before running
	if !*noCheckDuplicates {
		dups, err := tapeast.FindDuplicates(dir)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "tape: scan duplicates: %v\n", err)
			return 1
		}
		if len(dups) > 0 {
			_, _ = fmt.Fprintf(stderr, "tape: duplicate test names found:\n")
			for _, d := range dups {
				_, _ = fmt.Fprintf(stderr, "  %s\n", d)
			}
			return 1
		}
	}

	// count tests for progress bar total
	total, err := tapeast.CountTestsInTestFiles(dir)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "tape: count tests: %v\n", err)
		return 1
	}

	// find Only calls — restrict run if any found
	onlyCalls, err := tapeast.FindOnlyCalls(dir)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "tape: scan Only: %v\n", err)
		return 1
	}

	goArgs := []string{"test", "-json", "-v", path}
	if pattern := tapeast.BuildRunPattern(onlyCalls); pattern != "" {
		goArgs = append(goArgs, "-run", pattern)
	}

	f := formatter.New(*format, stdout, total)
	store := state.New()
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

	// skipped = declared tests that never ran (Skip calls + Only filtering)
	skipped := total - len(passed) - len(failed)
	if skipped < 0 {
		skipped = 0
	}

	// When Only calls are present, recompute with full name list
	if len(onlyCalls) > 0 {
		allNames, err := tapeast.FindAllTestNames(dir)
		if err == nil {
			store.MarkSkipped(allNames)
			passed, failed, _ = store.Summary()
			skipped = total - len(passed) - len(failed)
			if skipped < 0 {
				skipped = 0
			}
		}
	}

	f.End(len(passed), len(failed), skipped)

	if len(failed) > 0 {
		return 1
	}
	return 0
}
