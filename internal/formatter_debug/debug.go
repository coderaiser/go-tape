package formatter_debug

import (
	"fmt"
	"io"

	"github.com/coderaiser/go-tape/internal/stream"
)

// DebugFormatter writes structured debug lines to w (stdout when -f debug).
// Event returns the debug line as a string so the dispatcher can write it
// to stdout via the normal formatter pipeline.
type DebugFormatter struct {
	w     io.Writer
	total int
}

// New returns a DebugFormatter that writes all output to w.
func New(w io.Writer, total int) *DebugFormatter {
	return &DebugFormatter{w: w, total: total}
}

func (f *DebugFormatter) Event(e stream.Event) string {
	return f.format(e)
}

func (f *DebugFormatter) End(passed, failed, skipped int) string {
	return fmt.Sprintf("[tape:debug] result: passed=%d failed=%d skipped=%d\n",
		passed, failed, skipped)
}

func (f *DebugFormatter) format(e stream.Event) string {
	switch e.Type {
	case stream.TypeTestEnd:
		return fmt.Sprintf("[tape:debug] pass  %s  count=%d/%d failed=%d\n",
			e.Test, e.Count, f.total, e.Failed)
	case stream.TypeFail:
		return fmt.Sprintf("[tape:debug] fail  %s  count=%d/%d operator=%q\n",
			e.Test, e.Count, f.total, e.Operator)
	case stream.TypeBuildError:
		return fmt.Sprintf("[tape:debug] build-error   %s\n%s\n", e.Package, e.Output)
	case stream.TypePackageError:
		return fmt.Sprintf("[tape:debug] package-error %s\n%s\n", e.Package, e.Output)
	case stream.TypeUnknownFail:
		return fmt.Sprintf("[tape:debug] unknown-fail  %s\n", e.Test)
	case stream.TypeComment:
		return fmt.Sprintf("[tape:debug] comment       %s\n", e.Message)
	}
	return ""
}
