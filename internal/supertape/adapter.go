package supertape

import "github.com/coderaiser/go-tape/internal/stream"

// ToStreamEvent converts a supertape Event into the stream.Event type that the
// tape formatters consume. supertape already mirrors stream.Event's fields and
// type constants (TypeTestEnd/TypeFail), so this is a direct field copy — it
// lets the interp mode feed the same downstream formatters as the go-test mode
// with zero changes to the formatter layer.
func ToStreamEvent(ev Event) stream.Event {
	return stream.Event{
		Type:     ev.Type,
		Test:     ev.Test,
		Message:  ev.Message,
		Operator: ev.Operator,
		Result:   ev.Result,
		Expected: ev.Expected,
		Output:   ev.Output,
		At:       ev.At,
	}
}
