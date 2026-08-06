// Package report formats go test -json output using tape's formatters.
// It provides a single entry point for any tool that runs go test and wants
// tape-style output without depending on tape's internal packages directly.
package report

import (
	"bufio"
	"io"

	"github.com/coderaiser/go-tape/internal/formatter"
	"github.com/coderaiser/go-tape/internal/parser"
	"github.com/coderaiser/go-tape/internal/state"
)

// Run reads go test -json output from r, formats events using the named
// formatter, and writes to w.
//
// format selects the reporter: "fail", "tap", "short", "progress-bar",
// "json-lines", "time". Empty string defaults to "progress-bar" locally
// and "fail" on CI (when CI=1).
//
// total is the expected test count for the progress bar; pass 0 if unknown —
// formatters degrade gracefully.
//
// Returns an error only on scanner failure; test failures are reflected in
// the formatted output, not as a returned error.
func Run(r io.Reader, w io.Writer, format string, total int) error {
	f := formatter.New(format, w, total)

	store, err := state.New()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		e, err := parser.Parse(scanner.Text())
		if err != nil {
			continue
		}
		if _, err := store.Apply(e); err != nil {
			continue
		}
		f.FromEvent(e)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	passed, failed, skipped := store.Summary()
	f.End(len(passed), len(failed), len(skipped))

	return nil
}
