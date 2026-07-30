package formatter

import (
	"io"
	"os"

	"github.com/coderaiser/go-tape/internal/formatter/output"
	"github.com/coderaiser/go-tape/internal/formatter_fail"
	"github.com/coderaiser/go-tape/internal/formatter_json_lines"
	"github.com/coderaiser/go-tape/internal/formatter_progress_bar"
	"github.com/coderaiser/go-tape/internal/formatter_short"
	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/formatter_time"
	"github.com/coderaiser/go-tape/internal/model"
)

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
}

// New returns the formatter for the given format string.
func New(format string, w io.Writer, total int) *State {
	if os.Getenv("CI") == "true" {
		format = "tap"
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
		f = formatter_time.New(total)
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
	write(w, f.Start(total))
	return s
}

// FromEvent routes a model.Event to the appropriate formatter method.
func (s *State) FromEvent(e model.Event) {
	if e.Test == "" {
		return
	}
	switch e.Action {
	case "run":
		write(s.w, s.formatter.Test(e.Test))
	case "output":
		s.outputs[e.Test] = append(s.outputs[e.Test], e.Output)
	case "pass":
		s.count++
		write(s.w, s.formatter.Success(s.count, e.Test))
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, e.Test))
	case "fail":
		s.count++
		s.failed++
		lines := s.outputs[e.Test]
		fields := output.ParseOutput(lines)
		write(s.w, s.formatter.Fail(
			s.count, e.Test,
			fields.Operator, fields.Result, fields.Expected,
			fields.Raw, fields.At, fields.ErrorStack,
		))
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, e.Test))
	case "skip":
		s.count++
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, e.Test))
	}
}

// End writes the final summary.
func (s *State) End(passed, failed, skipped int) {
	write(s.w, s.formatter.End(s.count, passed, failed, skipped))
}

func write(w io.Writer, s string) {
	if s != "" {
		w.Write([]byte(s))
	}
}
