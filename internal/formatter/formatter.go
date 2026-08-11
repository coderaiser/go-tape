package formatter

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coderaiser/go-tape/internal/formatter_fail"
	"github.com/coderaiser/go-tape/internal/formatter_json_lines"
	"github.com/coderaiser/go-tape/internal/formatter_progress_bar"
	"github.com/coderaiser/go-tape/internal/formatter_short"
	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/formatter_time"
	"github.com/coderaiser/go-tape/internal/stream"
)

// Formatter is the interface every formatter must implement.
// Event is called for each stream event; End is called once at the end.
type Formatter interface {
	Event(e stream.Event) string
	End(passed, failed, skipped int) string
}

// Dispatcher wraps a Formatter and resolves the At field to a file URI
// before dispatch, so individual formatters never need to do path logic.
type Dispatcher struct {
	f       Formatter
	w       io.Writer
	dir     string
	modName string
}

// New returns a Dispatcher for the given format string, writer, and total
// test count. Mirrors the old formatter.New signature so main.go barely changes.
func New(format string, w io.Writer, total int) *Dispatcher {
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

	d := &Dispatcher{f: f, w: w}
	if dir, err := os.Getwd(); err == nil {
		d.dir = dir
		d.modName = readModuleName(dir)
	}
	return d
}

// newWithDir is like New but injects a specific dir (for testing).
func newWithDir(format string, w io.Writer, total int, dir string) *Dispatcher {
	d := New(format, w, total)
	d.dir = dir
	d.modName = readModuleName(dir)
	return d
}

// Emit routes a stream.Event through the formatter and writes the result.
func (d *Dispatcher) Emit(e stream.Event) {
	// resolve At → URI before dispatch
	if e.At != "" {
		e.At = fileLink(e.At, pkgDir(e.Package, d.modName, d.dir))
	}
	write(d.w, d.f.Event(e))
}

// End writes the final summary.
func (d *Dispatcher) End(passed, failed, skipped int) {
	write(d.w, d.f.End(passed, failed, skipped))
}

// --- unexported helpers ---

func testLabel(test string) string {
	if i := strings.LastIndex(test, "/"); i >= 0 {
		test = test[i+1:]
	}
	return strings.ReplaceAll(test, "_", " ")
}

func fileLink(at, dir string) string {
	if at == "" {
		return ""
	}
	at = strings.TrimRight(at, ":")
	if dir != "" && !strings.HasPrefix(at, "/") {
		at = dir + "/" + at
	}
	return "at file://" + at
}

func pkgDir(pkg, modName, dir string) string {
	if pkg == "" || dir == "" {
		return dir
	}
	rel := strings.TrimPrefix(pkg, modName)
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" {
		return dir
	}
	return filepath.Join(dir, rel)
}

func readModuleName(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.SplitN(string(data), "\n", 10) {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func write(w io.Writer, s string) {
	if s == "" {
		return
	}
	_, _ = w.Write([]byte(s))
}
