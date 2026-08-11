package formatter_tap_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestTAPTestEnd(t *testing.T) {
	Test(t, "formatter-tap: test-end event returns ok line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1})
		t.Equal(result, "ok 1 scope: foo\n")
		t.End()
	})
}

func TestTAPFailEvent(t *testing.T) {
	Test(t, "formatter-tap: fail event returns not ok line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Operator: "should equal", Result: "got", Expected: "want",
		})
		t.Match(result, "not ok 1 scope: foo")
		t.End()
	})
}

func TestTAPFailEventOperator(t *testing.T) {
	Test(t, "formatter-tap: fail event includes operator", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Operator: "should equal", Result: "got", Expected: "want",
		})
		t.Match(result, "operator: should equal")
		t.End()
	})
}

func TestTAPFailEventDiff(t *testing.T) {
	Test(t, "formatter-tap: fail event includes diff block", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Operator: "equal", Result: "hello", Expected: "world",
		})
		t.Match(result, "diff: |-")
		t.End()
	})
}

func TestTAPFailEventDiffMinusLine(t *testing.T) {
	Test(t, "formatter-tap: fail event diff has minus line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Result: "hello", Expected: "world",
		})
		t.Match(result, `- "world"`)
		t.End()
	})
}

func TestTAPFailEventDiffPlusLine(t *testing.T) {
	Test(t, "formatter-tap: fail event diff has plus line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Result: "hello", Expected: "world",
		})
		t.Match(result, `+ "hello"`)
		t.End()
	})
}

func TestTAPFailWithOutput(t *testing.T) {
	Test(t, "formatter-tap: fail event with output uses it directly", func(t *T) {
		f := formatter_tap.New()
		cut := "operator: should equal\n        expected: want\n        result: got\n"
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Operator: "should equal", Result: "got", Expected: "want",
			Output: cut,
		})
		t.Equal(result, "not ok 1 scope: foo\n"+cut+"\n")
		t.End()
	})
}

func TestTAPFailWithAt(t *testing.T) {
	Test(t, "formatter-tap: fail event includes at when non-empty", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			At: "at file:///proj/file.go:10",
		})
		t.Match(result, "file:///proj/file.go:10")
		t.End()
	})
}

func TestTAPFailWithErrorStack(t *testing.T) {
	Test(t, "formatter-tap: fail event includes stack when non-empty", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			ErrorStack: "stack trace",
		})
		t.Match(result, "stack trace")
		t.End()
	})
}

func TestTAPFailNilExpected(t *testing.T) {
	Test(t, "formatter-tap: fail with nil expected omits expected line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{Type: stream.TypeFail, Test: "scope: foo", Count: 1})
		t.NotMatch(result, "expected:")
		t.End()
	})
}

func TestTAPFailNilResult(t *testing.T) {
	Test(t, "formatter-tap: fail with nil result omits result line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{Type: stream.TypeFail, Test: "scope: foo", Count: 1})
		t.NotMatch(result, "result:")
		t.End()
	})
}

func TestTAPEqualValuesShowsExpected(t *testing.T) {
	Test(t, "formatter-tap: fail with equal values shows expected", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Result: "same", Expected: "same",
		})
		t.Match(result, "expected:")
		t.End()
	})
}

func TestTAPEqualValuesShowsResult(t *testing.T) {
	Test(t, "formatter-tap: fail with equal values shows result", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Result: "same", Expected: "same",
		})
		t.Match(result, "result:")
		t.End()
	})
}

func TestTAPUnknownFail(t *testing.T) {
	Test(t, "formatter-tap: unknown-fail event returns not ok line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeUnknownFail, Test: "scope: foo", Count: 1,
			Output: "panic: something went wrong\n",
		})
		t.Match(result, "not ok 1 scope: foo")
		t.End()
	})
}

func TestTAPUnknownFailOutput(t *testing.T) {
	Test(t, "formatter-tap: unknown-fail event includes raw output", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type: stream.TypeUnknownFail, Test: "scope: foo", Count: 1,
			Output: "panic: something went wrong\n",
		})
		t.Match(result, "panic: something went wrong")
		t.End()
	})
}

func TestTAPBuildError(t *testing.T) {
	Test(t, "formatter-tap: build-error event includes build-error prefix", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{
			Type:    stream.TypeBuildError,
			Package: "example.com/foo",
			Output:  "foo.go:1: undefined\n",
		})
		t.Match(result, "build-error")
		t.End()
	})
}

func TestTAPComment(t *testing.T) {
	Test(t, "formatter-tap: comment event returns comment line", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{Type: stream.TypeComment, Message: "hello"})
		t.Equal(result, "# hello\n")
		t.End()
	})
}

func TestTAPUnknownEventType(t *testing.T) {
	Test(t, "formatter-tap: unknown event type returns empty string", func(t *T) {
		f := formatter_tap.New()
		result := f.Event(stream.Event{Type: "bogus"})
		t.Equal(result, "")
		t.End()
	})
}

func TestTAPEndAllPass(t *testing.T) {
	Test(t, "formatter-tap: End with no failures shows ok", func(t *T) {
		f := formatter_tap.New()
		result := f.End(5, 0, 0)
		t.Match(result, "# ok")
		t.End()
	})
}

func TestTAPEndWithSkipped(t *testing.T) {
	Test(t, "formatter-tap: End includes skip count when skipped > 0", func(t *T) {
		f := formatter_tap.New()
		result := f.End(4, 0, 1)
		t.Match(result, "# skip 1")
		t.End()
	})
}

func TestTAPEndWithFailed(t *testing.T) {
	Test(t, "formatter-tap: End includes fail count when failed > 0", func(t *T) {
		f := formatter_tap.New()
		result := f.End(4, 1, 0)
		t.Match(result, "# fail 1")
		t.End()
	})
}

func TestTAPEndTotal(t *testing.T) {
	Test(t, "formatter-tap: End 1..N uses passed+failed+skipped", func(t *T) {
		f := formatter_tap.New()
		result := f.End(3, 1, 1)
		t.Match(result, "1..5")
		t.End()
	})
}
