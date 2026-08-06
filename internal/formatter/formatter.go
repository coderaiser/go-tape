package formatter

import (
	"io"
	"os"
	"strings"

	"github.com/coderaiser/go-tape/internal/formatter/output"
	"github.com/coderaiser/go-tape/internal/formatter_fail"
	"github.com/coderaiser/go-tape/internal/formatter_json_lines"
	"github.com/coderaiser/go-tape/internal/formatter_progress_bar"
	"github.com/coderaiser/go-tape/internal/formatter_short"
	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/formatter_time"
	"github.com/coderaiser/go-tape/internal/model"
)

// testLabel converts a go test -json Test field like
// "TestFoo/scope:_bar_baz" into a clean label "scope: bar baz".
func testLabel(test string) string {
	if i := strings.LastIndex(test, "/"); i >= 0 {
		test = test[i+1:]
	}
	return strings.ReplaceAll(test, "_", " ")
}

// fileLink converts a relative "file.go:N:" at string into a
// "file:///abs/path/file.go:N" terminal-clickable URI.
func fileLink(at, dir string) string {
	if at == "" {
		return ""
	}
	// at is "file.go:N:" — strip trailing colon
	at = strings.TrimRight(at, ":")
	if dir != "" && !strings.HasPrefix(at, "/") {
		at = dir + "/" + at
	}
	return "file://" + at
}

// Formatter matches supertape's formatter event API exactly.
type Formatter interface {
	Start(total int) string
	Test(name string) string
	TestEnd(count, total, failed int, name string) string
	Success(count int, message string) string
	Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string
	Comment(message string) string
	End(count, passed, failed, skipped int) string
}

// State tracks streaming state for FromEvent.
type State struct {
	count     int
	failed    int
	total     int
	outputs   map[string][]string
	formatter Formatter
	w         io.Writer
	dir       string
}

// New returns the formatter for the given format string.
func New(format string, w io.Writer, total int) *State {
	if ci := os.Getenv("CI"); ci == "1" || ci == "true" {
		format = "fail"
	} else if format == "" {
		format = "progress-bar"
	}
	var f Formatter
	switch format {
	case "tap":
		f = formatter_tap.New()
	case "short":
		f = formatter_short.New()
	case "fail":
		f = formatter_fail.New()
	case "time":
		f = formatter_time.New(total, os.Stderr)
	case "json-lines":
		f = formatter_json_lines.New(total)
	default:
		f = formatter_progress_bar.New(total)
	}
	s := &State{
		total:     total,
		outputs:   make(map[string][]string),
		formatter: f,
		w:         w,
	}
	if dir, err := os.Getwd(); err == nil {
		s.dir = dir
	}
	write(w, f.Start(total))
	return s
}

// FromEvent routes a model.Event to the appropriate formatter method.
func (s *State) FromEvent(e model.Event) {
	if e.Test == "" {
		return
	}
	label := testLabel(e.Test)
	switch e.Action {
	case "run":
		write(s.w, s.formatter.Test(label))
	case "output":
		s.outputs[e.Test] = append(s.outputs[e.Test], e.Output)
	case "pass":
		s.count++
		write(s.w, s.formatter.Success(s.count, label))
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, label))
	case "fail":
		s.count++
		s.failed++
		lines := s.outputs[e.Test]
		fields := output.ParseOutput(lines)
		write(s.w, s.formatter.Fail(
			s.count, label,
			fields.Operator, fields.Result, fields.Expected,
			fields.Cut, fileLink(fields.At, s.dir), fields.ErrorStack,
		))
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, label))
	case "skip":
		s.count++
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, label))
	}
}

// End writes the final summary.
func (s *State) End(passed, failed, skipped int) {
	// On a cached run, packages that use plain Go subtests (not tape.Test) emit
	// no subtest events, so s.count may be short of s.total. Emit one synthetic
	// TestEnd at 100% so the progress bar reaches completion before it is cleared.
	if s.count < s.total {
		write(s.w, s.formatter.TestEnd(s.total, s.total, s.failed, ""))
	}
	write(s.w, s.formatter.End(s.count, passed, failed, skipped))
}

func write(w io.Writer, s string) {
	if s == "" {
		return
	}
	// Discard the error: the writer is os.Stdout, and if it fails the
	// process environment is already broken — there is nowhere to report
	// the failure to. No recovery is possible, so be explicit about it.
	_, _ = w.Write([]byte(s))
}
