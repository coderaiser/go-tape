package main

import (
	"flag"
	"fmt"
	"os"

	tapeast "github.com/coderaiser/go-tape/cmd/tape/ast"
	"github.com/coderaiser/go-tape/formatter"
	"github.com/coderaiser/go-tape/runner"
	"github.com/coderaiser/go-tape/state"
)

const version = "1.0.0"

func main() {
	format            := flag.String("f", "", "output format: tap|progress-bar|short|fail|time|json-lines")
	help              := flag.Bool("h", false, "display this help and exit")
	ver               := flag.Bool("v", false, "output version information and exit")
	noCheckScopes     := flag.Bool("no-check-scopes", false, "do not check scope format")
	noCheckAssertions := flag.Bool("no-check-assertions-count", false, "do not check assertion count")
	noCheckDuplicates := flag.Bool("no-check-duplicates", false, "do not check for duplicates")
	flag.Parse()

	if *help {
		fmt.Print(usage)
		return
	}
	if *ver {
		fmt.Println(version)
		return
	}
	if *noCheckScopes {
		os.Setenv("TAPE_CHECK_SCOPES", "0")
	}
	if *noCheckAssertions {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	}

	path := "./..."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}
	dir := "."

	// check duplicates before running
	if !*noCheckDuplicates {
		dups, err := tapeast.FindDuplicates(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tape: scan duplicates: %v\n", err)
			os.Exit(1)
		}
		if len(dups) > 0 {
			fmt.Fprintf(os.Stderr, "tape: duplicate test names found:\n")
			for _, d := range dups {
				fmt.Fprintf(os.Stderr, "  %s\n", d)
			}
			os.Exit(1)
		}
	}

	// count tests for progress bar total
	total, err := tapeast.CountTests(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tape: count tests: %v\n", err)
		os.Exit(1)
	}

	// find Only calls — restrict run if any found
	onlyCalls, err := tapeast.FindOnlyCalls(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tape: scan Only: %v\n", err)
		os.Exit(1)
	}

	args := []string{"test", "-json", "-v", path}
	if pattern := tapeast.BuildRunPattern(onlyCalls); pattern != "" {
		fmt.Fprintf(os.Stderr, "tape: Only mode — %s\n", pattern)
		args = append(args, "-run", pattern)
	}

	f := formatter.New(*format, os.Stdout, total)
	store := state.New()
	r := runner.New(runner.NewOSExecutor())

	ch, err := r.Run(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tape: %v\n", err)
		os.Exit(1)
	}

	for event := range ch {
		store.Apply(event)
		f.FromEvent(event)
	}

	passed, failed, skipped := store.Summary()
	f.End(len(passed), len(failed), len(skipped))

	if len(failed) > 0 {
		os.Exit(1)
	}
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
`
