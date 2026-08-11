package formatter_fail_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_fail"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestFailFormatterTestEndSuppressed(t *testing.T) {
	Test(t, "formatter-fail: test-end event returns empty (passing suppressed)", func(t *T) {
		f := formatter_fail.New()
		result := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1})
		t.Equal(result, "")
		t.End()
	})
}

func TestFailFormatterTestEndStoresName(t *testing.T) {
	Test(t, "formatter-fail: fail after test-end prepends stored test name", func(t *T) {
		f := formatter_fail.New()
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1})
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 2,
			Operator: "Equal", Result: "got", Expected: "want",
		})
		t.Match(result, "# scope: foo")
		t.End()
	})
}

func TestFailFormatterFailPrependsEmptyWhenNoTestEnd(t *testing.T) {
	Test(t, "formatter-fail: fail prepends header even when no prior test-end", func(t *T) {
		f := formatter_fail.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Operator: "Equal", Result: "got", Expected: "want",
		})
		t.Match(result, "# \n")
		t.End()
	})
}

func TestFailFormatterKeepsErrorStack(t *testing.T) {
	Test(t, "formatter-fail: fail keeps error stack (not stripped like short)", func(t *T) {
		f := formatter_fail.New()
		result := f.Event(stream.Event{
			Type:       stream.TypeFail,
			Test:       "scope: foo",
			Count:      1,
			Operator:   "Equal",
			Result:     "got",
			Expected:   "want",
			ErrorStack: "my stack trace",
		})
		t.Match(result, "my stack trace")
		t.End()
	})
}

func TestFailFormatterUnknownFail(t *testing.T) {
	Test(t, "formatter-fail: unknown-fail event includes raw output", func(t *T) {
		f := formatter_fail.New()
		result := f.Event(stream.Event{
			Type:   stream.TypeUnknownFail,
			Test:   "scope: foo",
			Count:  1,
			Output: "panic: boom\n",
		})
		t.Match(result, "panic: boom")
		t.End()
	})
}

func TestFailFormatterEndDelegates(t *testing.T) {
	Test(t, "formatter-fail: End delegates to tap", func(t *T) {
		f := formatter_fail.New()
		result := f.End(3, 1, 0)
		t.Match(result, "# fail 1")
		t.End()
	})
}
