package formatter_json_lines

import (
	"encoding/json"

	"github.com/coderaiser/go-tape/internal/stream"
)

// JSONLinesFormatter emits one JSON object per event.
type JSONLinesFormatter struct {
	total int
}

func New(total int) *JSONLinesFormatter {
	return &JSONLinesFormatter{total: total}
}

// Event marshals the stream.Event directly — the schema is already json-lines.
func (f *JSONLinesFormatter) Event(e stream.Event) string {
	b, _ := json.Marshal(e)
	return string(b) + "\n"
}

// End emits the final summary object.
func (f *JSONLinesFormatter) End(passed, failed, skipped int) string {
	b, _ := json.Marshal(map[string]any{
		"type":    stream.TypeEnd,
		"count":   passed + failed + skipped,
		"passed":  passed,
		"failed":  failed,
		"skipped": skipped,
	})
	return string(b) + "\n"
}
