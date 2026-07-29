package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	tapeast "github.com/coderaiser/go-tape/cmd/tape/ast"
	"github.com/coderaiser/go-tape/formatter"
	"github.com/coderaiser/go-tape/runner"
	"github.com/coderaiser/go-tape/state"
)

const version = "1.0.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("go-tape", flag.ExitOnError)
	format := flags.String("f", "", "output format: tap|progress-bar|short|fail|time|json-lines")
	help := flags.Bool("h", false, "display this help and exit")
	ver := flags.Bool("v", false, "output version information and exit")
	noCheckScopes := flags.Bool("no-check-scopes", false, "do not check scope format")
	noCheckAssertions := flags.Bool("no-check-assertions-count", false, "do not check assertion count")
	noCheckDuplicates := flags.Bool("no-check-duplicates", false, "do not check for duplicates")
	flags.Parse(args)

	if *help {
		fmt.Fprint(stdout, usage)
		return 0
	}
	if *ver {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if *noCheckScopes {
		os.Setenv("TAPE_CHECK_SCOPES", "0")
	}
	if *noCheckAssertions {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
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
			fmt.Fprintf(stderr, "tape: scan duplicates: %v\n", err)
			return 1
		}
		if len(dups) > 0 {
			fmt.Fprintf(stderr, "tape: duplicate test names found:\n")
			for _, d := range dups {
				fmt.Fprintf(stderr, "  %s\n", d)
			}
			return 1
		}
	}

	// count tests for progress bar total
	total, err := tapeast.CountTestsInTestFiles(dir)
	if err != nil {
		fmt.Fprintf(stderr, "tape: count tests: %v\n", err)
		return 1
	}

	// find Only calls — restrict run if any found
	onlyCalls, err := tapeast.FindOnlyCalls(dir)
	if err != nil {
		fmt.Fprintf(stderr, "tape: scan Only: %v\n", err)
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
		fmt.Fprintf(stderr, "tape: %v\n", err)
		return 1
	}

	for event := range ch {
		// go test -json emits events for both the outer TestXxx wrapper and the
		// inner t.Run subtest. Only the subtest (containing "/") is a tape.Test
		// call — skip the wrapper to avoid double-counting.
		if event.Test != "" && !strings.Contains(event.Test, "/") {
			continue
		}
		store.Apply(event)
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

const usage = `Usage: go-tape [options] [path]

Options:
  -h                           display this help and exit
  -v                           output version information and exit
  -f format                    use a specific output format
                               default: progress-bar (tap on CI)
                               values: tap|progress-bar|short|fail|time|json-lines
  --no-check-scopes            do not check that messages contain scope: 'scope: message'
  --no-check-assertions-count  do not check that assertion count is no more than 1
  --no-check-duplicates        do not check messages for duplicates

Environment variables:
  TAPE_PROGRESS_BAR=1|0        force progress bar on (1) or off (0)
  TAPE_PROGRESS_BAR_MIN=10     minimum test count to show progress bar (default: 10)
  TAPE_PROGRESS_BAR_COLOR      progress bar color (ANSI escape or hex, default: yellow)
  TAPE_PROGRESS_BAR_STACK=0    set to 0 to hide error stack traces
`
