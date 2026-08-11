// Package report formats go test -json output using tape's formatters.
// It provides a single entry point for any tool that runs go test and wants
// tape-style output without depending on tape's internal packages directly.
package report

import (
	"io"

	"github.com/coderaiser/go-tape/internal/formatter"
	"github.com/coderaiser/go-tape/internal/stream"
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
	d := formatter.New(format, w, total)

	var lastCount, lastFailed int
	for e := range stream.Parse(r, total) {
		d.Emit(e)
		if e.Type == stream.TypeTestEnd {
			lastCount = e.Count
			lastFailed = e.Failed
		}
	}

	passedCount := lastCount - lastFailed
	skipped := total - lastCount
	if skipped < 0 {
		skipped = 0
	}

	d.End(passedCount, lastFailed, skipped)
	return nil
}
