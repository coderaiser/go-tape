package formatter_debug_test

import (
	"bytes"
	"testing"

	tape "github.com/coderaiser/go-tape"
	debug "github.com/coderaiser/go-tape/internal/formatter_debug"
	"github.com/coderaiser/go-tape/internal/formatter_tap"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestDebugEventTestEndWritesToWriter(t *testing.T) {
	tape.Test(t, "formatter-debug: test-end writes to writer", func(t *tape.T) {
		var buf bytes.Buffer
		f := debug.New(&buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1})
		t.Match(buf.String(), "test-end")
		t.End()
	})
}

func TestDebugEventTestEndContainsCount(t *testing.T) {
	tape.Test(t, "formatter-debug: test-end line contains count", func(t *tape.T) {
		var buf bytes.Buffer
		f := debug.New(&buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 3})
		t.Match(buf.String(), "count=3")
		t.End()
	})
}

func TestDebugEventFailContainsOperator(t *testing.T) {
	tape.Test(t, "formatter-debug: fail line contains operator", func(t *tape.T) {
		var buf bytes.Buffer
		f := debug.New(&buf)
		f.Event(stream.Event{Type: stream.TypeFail, Test: "scope: x", Operator: "should equal"})
		t.Match(buf.String(), "should equal")
		t.End()
	})
}

func TestDebugEventBuildError(t *testing.T) {
	tape.Test(t, "formatter-debug: build-error line contains package", func(t *tape.T) {
		var buf bytes.Buffer
		f := debug.New(&buf)
		f.Event(stream.Event{Type: stream.TypeBuildError, Package: "example.com/foo"})
		t.Match(buf.String(), "example.com/foo")
		t.End()
	})
}

func TestDebugEventReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-debug: Event returns empty string", func(t *tape.T) {
		var buf bytes.Buffer
		f := debug.New(&buf)
		out := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1})
		t.Equal(out, "")
		t.End()
	})
}

func TestDebugEndWritesSummary(t *testing.T) {
	tape.Test(t, "formatter-debug: End writes passed/failed/skipped", func(t *tape.T) {
		var buf bytes.Buffer
		f := debug.New(&buf)
		f.End(5, 1, 2)
		t.Match(buf.String(), "passed=5")
		t.End()
	})
}

func TestDebugEndReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-debug: End returns empty string", func(t *tape.T) {
		var buf bytes.Buffer
		f := debug.New(&buf)
		out := f.End(1, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestDebugWrappingDelegatesEventOutput(t *testing.T) {
	tape.Test(t, "formatter-debug: wrapping returns inner formatter stdout output", func(t *tape.T) {
		var dbgBuf bytes.Buffer
		inner := formatter_tap.New()
		f := debug.NewWrapping(inner, &dbgBuf)
		out := f.Event(stream.Event{
			Type: stream.TypeTestEnd, Test: "scope: x", Count: 1, Total: 1,
		})
		t.Match(out, "ok 1")
		t.End()
	})
}

func TestDebugWrappingStillWritesDebugLines(t *testing.T) {
	tape.Test(t, "formatter-debug: wrapping writes debug lines to stderr writer", func(t *tape.T) {
		var dbgBuf bytes.Buffer
		inner := formatter_tap.New()
		f := debug.NewWrapping(inner, &dbgBuf)
		f.Event(stream.Event{
			Type: stream.TypeTestEnd, Test: "scope: x", Count: 1, Total: 1,
		})
		t.Match(dbgBuf.String(), "test-end")
		t.End()
	})
}

func TestDebugWrappingEndDelegates(t *testing.T) {
	tape.Test(t, "formatter-debug: wrapping End returns inner formatter output", func(t *tape.T) {
		var dbgBuf bytes.Buffer
		inner := formatter_tap.New()
		f := debug.NewWrapping(inner, &dbgBuf)
		out := f.End(1, 0, 0)
		t.Match(out, "# ok")
		t.End()
	})
}
