package formatter_short

import (
	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/stream"
)

// ShortFormatter is TAP with the error stack stripped from fail blocks.
type ShortFormatter struct {
	*formatter_tap.TAPFormatter
}

func New() *ShortFormatter {
	return &ShortFormatter{TAPFormatter: formatter_tap.New()}
}

// Event delegates to TAPFormatter but clears ErrorStack before dispatch.
func (f *ShortFormatter) Event(e stream.Event) string {
	if e.Type == stream.TypeFail {
		e.ErrorStack = ""
	}
	return f.TAPFormatter.Event(e)
}
