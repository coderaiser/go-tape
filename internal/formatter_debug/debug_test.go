package formatter_debug_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	debug "github.com/coderaiser/go-tape/internal/formatter_debug"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestDebugEventTestEndContainsPass(t *testing.T) {
	tape.Test(t, "formatter-debug: test-end event contains pass label", func(t *tape.T) {
		f := debug.New(nil, 10)
		out := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1, Total: 10})
		t.Match(out, "pass")
		t.End()
	})
}

func TestDebugEventTestEndContainsCount(t *testing.T) {
	tape.Test(t, "formatter-debug: test-end line contains count/total", func(t *tape.T) {
		f := debug.New(nil, 10)
		out := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 3, Total: 10})
		t.Match(out, "count=3/10")
		t.End()
	})
}

func TestDebugEventTestEndContainsTestName(t *testing.T) {
	tape.Test(t, "formatter-debug: test-end line contains test name", func(t *tape.T) {
		f := debug.New(nil, 1)
		out := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: my test", Count: 1})
		t.Match(out, "scope: my test")
		t.End()
	})
}

func TestDebugEventFailContainsOperator(t *testing.T) {
	tape.Test(t, "formatter-debug: fail line contains operator", func(t *tape.T) {
		f := debug.New(nil, 1)
		out := f.Event(stream.Event{Type: stream.TypeFail, Test: "scope: x", Count: 1, Operator: "should equal"})
		t.Match(out, "should equal")
		t.End()
	})
}

func TestDebugEventFailContainsCountTotal(t *testing.T) {
	tape.Test(t, "formatter-debug: fail line contains count/total", func(t *tape.T) {
		f := debug.New(nil, 5)
		out := f.Event(stream.Event{Type: stream.TypeFail, Test: "scope: x", Count: 2, Operator: "ok"})
		t.Match(out, "count=2/5")
		t.End()
	})
}

func TestDebugEventBuildErrorContainsPackage(t *testing.T) {
	tape.Test(t, "formatter-debug: build-error line contains package", func(t *tape.T) {
		f := debug.New(nil, 0)
		out := f.Event(stream.Event{Type: stream.TypeBuildError, Package: "example.com/foo", Output: "foo.go:1:1: err\n"})
		t.Match(out, "example.com/foo")
		t.End()
	})
}

func TestDebugEventBuildErrorContainsOutput(t *testing.T) {
	tape.Test(t, "formatter-debug: build-error line contains compiler output", func(t *tape.T) {
		f := debug.New(nil, 0)
		out := f.Event(stream.Event{Type: stream.TypeBuildError, Package: "example.com/foo", Output: "foo.go:1:1: declared and not used\n"})
		t.Match(out, "declared and not used")
		t.End()
	})
}

func TestDebugEventPackageError(t *testing.T) {
	tape.Test(t, "formatter-debug: package-error line contains package", func(t *tape.T) {
		f := debug.New(nil, 0)
		out := f.Event(stream.Event{Type: stream.TypePackageError, Package: "example.com/bar", Output: "403 Forbidden\n"})
		t.Match(out, "example.com/bar")
		t.End()
	})
}

func TestDebugEventUnknownFail(t *testing.T) {
	tape.Test(t, "formatter-debug: unknown-fail line contains test name", func(t *tape.T) {
		f := debug.New(nil, 1)
		out := f.Event(stream.Event{Type: stream.TypeUnknownFail, Test: "scope: x"})
		t.Match(out, "scope: x")
		t.End()
	})
}

func TestDebugEventComment(t *testing.T) {
	tape.Test(t, "formatter-debug: comment line contains message", func(t *tape.T) {
		f := debug.New(nil, 1)
		out := f.Event(stream.Event{Type: stream.TypeComment, Message: "hello"})
		t.Match(out, "hello")
		t.End()
	})
}

func TestDebugEventUnknownTypeReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-debug: unknown event type returns empty string", func(t *tape.T) {
		f := debug.New(nil, 0)
		out := f.Event(stream.Event{Type: "unknown-type"})
		t.Equal(out, "")
		t.End()
	})
}

func TestDebugEndContainsPassed(t *testing.T) {
	tape.Test(t, "formatter-debug: End output contains passed count", func(t *tape.T) {
		f := debug.New(nil, 0)
		out := f.End(5, 1, 2)
		t.Match(out, "passed=5")
		t.End()
	})
}

func TestDebugEndContainsFailed(t *testing.T) {
	tape.Test(t, "formatter-debug: End output contains failed count", func(t *tape.T) {
		f := debug.New(nil, 0)
		out := f.End(5, 1, 2)
		t.Match(out, "failed=1")
		t.End()
	})
}

func TestDebugEndContainsSkipped(t *testing.T) {
	tape.Test(t, "formatter-debug: End output contains skipped count", func(t *tape.T) {
		f := debug.New(nil, 0)
		out := f.End(5, 1, 2)
		t.Match(out, "skipped=2")
		t.End()
	})
}
