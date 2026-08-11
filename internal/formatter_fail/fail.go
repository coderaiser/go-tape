package formatter_fail

import (
	"fmt"

	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/stream"
)

// FailFormatter suppresses passing test output and prepends the test name
// to every fail block (used in CI mode).
type FailFormatter struct {
	*formatter_tap.TAPFormatter
	currentTest string
}

func New() *FailFormatter {
	return &FailFormatter{TAPFormatter: formatter_tap.New()}
}

// Event stores the current test name and suppresses passing output.
func (f *FailFormatter) Event(e stream.Event) string {
	switch e.Type {
	case stream.TypeTestEnd:
		f.currentTest = e.Test
		return "" // suppress passing output
	case stream.TypeFail, stream.TypeUnknownFail:
		header := fmt.Sprintf("# %s\n", f.currentTest)
		return header + f.TAPFormatter.Event(e)
	default:
		return f.TAPFormatter.Event(e)
	}
}
