package formatter_debug

import (
	"fmt"
	"io"

	"github.com/coderaiser/go-tape/internal/formatter"
	"github.com/coderaiser/go-tape/internal/stream"
)

// DebugFormatter writes structured debug lines to w (intended: os.Stderr).
// Event always returns "" — no stdout output is produced.
type DebugFormatter struct {
	w io.Writer
}

// New returns a standalone DebugFormatter that writes all output to w.
func New(w io.Writer) *DebugFormatter {
	return &DebugFormatter{w: w}
}

func (f *DebugFormatter) Event(e stream.Event) string {
	f.log(e)
	return ""
}

func (f *DebugFormatter) End(passed, failed, skipped int) string {
	fmt.Fprintf(f.w, "[tape:debug] result: passed=%d failed=%d skipped=%d\n",
		passed, failed, skipped)
	return ""
}

func (f *DebugFormatter) log(e stream.Event) {
	switch e.Type {
	case stream.TypeTestEnd:
		fmt.Fprintf(f.w, "[tape:debug] test-end    %s  count=%d failed=%d\n",
			e.Test, e.Count, e.Failed)
	case stream.TypeFail:
		fmt.Fprintf(f.w, "[tape:debug] fail        %s  operator=%q\n",
			e.Test, e.Operator)
	case stream.TypeBuildError:
		fmt.Fprintf(f.w, "[tape:debug] build-error %s\n", e.Package)
	case stream.TypeUnknownFail:
		fmt.Fprintf(f.w, "[tape:debug] unknown-fail %s\n", e.Test)
	case stream.TypeComment:
		fmt.Fprintf(f.w, "[tape:debug] comment      %s\n", e.Message)
	}
}

// WrappingDebugFormatter delegates Event/End to an inner formatter (for
// stdout output) while also writing debug lines to w (stderr). This lets
// -f debug show both normal progress-bar output and full debug info.
type WrappingDebugFormatter struct {
	inner formatter.Formatter
	dbg   *DebugFormatter
}

// NewWrapping wraps inner, writing debug lines to w alongside inner's output.
func NewWrapping(inner formatter.Formatter, w io.Writer) *WrappingDebugFormatter {
	return &WrappingDebugFormatter{
		inner: inner,
		dbg:   New(w),
	}
}

func (f *WrappingDebugFormatter) Event(e stream.Event) string {
	f.dbg.log(e)
	return f.inner.Event(e)
}

func (f *WrappingDebugFormatter) End(passed, failed, skipped int) string {
	f.dbg.End(passed, failed, skipped)
	return f.inner.End(passed, failed, skipped)
}
