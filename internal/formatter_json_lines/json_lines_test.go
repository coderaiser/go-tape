package formatter_json_lines_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_json_lines"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestJLTestEnd(t *testing.T) {
	Test(t, "formatter-json-lines: test-end event has type field", func(t *T) {
		f := formatter_json_lines.New(10)
		result := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1, Total: 10})
		t.Match(result, `"type":"test-end"`)
		t.End()
	})
}

func TestJLTestEndCount(t *testing.T) {
	Test(t, "formatter-json-lines: test-end event has count", func(t *T) {
		f := formatter_json_lines.New(10)
		result := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1, Total: 10})
		t.Match(result, `"count":1`)
		t.End()
	})
}

func TestJLFail(t *testing.T) {
	Test(t, "formatter-json-lines: fail event has type field", func(t *T) {
		f := formatter_json_lines.New(10)
		result := f.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: foo", Count: 1,
			Operator: "Equal", Result: "got", Expected: "want",
		})
		t.Match(result, `"type":"fail"`)
		t.End()
	})
}

func TestJLFailTest(t *testing.T) {
	Test(t, "formatter-json-lines: fail event has test name", func(t *T) {
		f := formatter_json_lines.New(10)
		result := f.Event(stream.Event{Type: stream.TypeFail, Test: "scope: foo", Count: 1})
		t.Match(result, `"test":"scope: foo"`)
		t.End()
	})
}

func TestJLBuildError(t *testing.T) {
	Test(t, "formatter-json-lines: build-error event has type field", func(t *T) {
		f := formatter_json_lines.New(0)
		result := f.Event(stream.Event{Type: stream.TypeBuildError, Package: "example.com/foo", Output: "err\n"})
		t.Match(result, `"type":"build-error"`)
		t.End()
	})
}

func TestJLEnd(t *testing.T) {
	Test(t, "formatter-json-lines: End has passed field", func(t *T) {
		f := formatter_json_lines.New(10)
		result := f.End(5, 0, 0)
		t.Match(result, `"passed":5`)
		t.End()
	})
}

func TestJLEndFailed(t *testing.T) {
	Test(t, "formatter-json-lines: End has failed field", func(t *T) {
		f := formatter_json_lines.New(10)
		result := f.End(4, 1, 0)
		t.Match(result, `"failed":1`)
		t.End()
	})
}

func TestJLEndType(t *testing.T) {
	Test(t, "formatter-json-lines: End has type end", func(t *T) {
		f := formatter_json_lines.New(10)
		result := f.End(5, 0, 0)
		t.Match(result, `"type":"end"`)
		t.End()
	})
}
