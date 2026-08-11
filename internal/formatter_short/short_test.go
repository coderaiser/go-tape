package formatter_short_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_short"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestShortFailOmitsErrorStack(t *testing.T) {
	Test(t, "formatter-short: fail event omits error stack", func(t *T) {
		f := formatter_short.New()
		result := f.Event(stream.Event{
			Type:       stream.TypeFail,
			Test:       "scope: x",
			Count:      1,
			Operator:   "Equal",
			Result:     "got",
			Expected:   "want",
			ErrorStack: "stack here",
		})
		t.NotMatch(result, "stack here")
		t.End()
	})
}

func TestShortPassDelegates(t *testing.T) {
	Test(t, "formatter-short: test-end event delegates to tap", func(t *T) {
		f := formatter_short.New()
		result := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1})
		t.Equal(result, "ok 1 scope: x\n")
		t.End()
	})
}

func TestShortEndDelegates(t *testing.T) {
	Test(t, "formatter-short: End delegates to tap", func(t *T) {
		f := formatter_short.New()
		result := f.End(3, 0, 0)
		t.Match(result, "# tests 3")
		t.End()
	})
}
