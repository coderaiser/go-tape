package supertape_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
	st "github.com/coderaiser/go-tape/internal/supertape"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestToStreamEventFail(t *testing.T) {
	Test(t, "supertape: ToStreamEvent maps a fail event", func(t *T) {
		ev := st.Event{
			Type:     st.TypeFail,
			Test:     "scope: x",
			Message:  "should equal",
			Operator: "equal",
			Result:   1,
			Expected: 2,
			Output:   "diff",
			At:       "file.go:1",
		}
		s := st.ToStreamEvent(ev)
		t.DeepEqual(s, stream.Event{
			Type:     stream.TypeFail,
			Test:     "scope: x",
			Message:  "should equal",
			Operator: "equal",
			Result:   1,
			Expected: 2,
			Output:   "diff",
			At:       "file.go:1",
		})
		t.End()
	})
}

func TestToStreamEventEnd(t *testing.T) {
	Test(t, "supertape: ToStreamEvent maps a test-end event", func(t *T) {
		ev := st.Event{Type: st.TypeTestEnd, Test: "scope: done"}
		s := st.ToStreamEvent(ev)
		t.Equal(s.Type == stream.TypeTestEnd && s.Test == "scope: done", true)
		t.End()
	})
}