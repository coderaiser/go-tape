package formatter_time_test

import (
	"regexp"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_time"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestTimeEventTestEndReturnsEmpty(t *testing.T) {
	Test(t, "formatter-time: test-end event returns empty string", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		result := f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1, Total: 10})
		t.Equal(result, "")
		t.End()
	})
}

func TestTimeEventTestEndWritesCR(t *testing.T) {
	Test(t, "formatter-time: test-end writes CR-prefixed progress line", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1, Total: 10})
		t.Match(buf.String(), regexp.MustCompile(`^\r`))
		t.End()
	})
}

func TestTimeEventTestEndWritesPct(t *testing.T) {
	Test(t, "formatter-time: test-end writes percentage", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1, Total: 10})
		t.Match(buf.String(), "10%")
		t.End()
	})
}

func TestTimeEventTestEndWritesCount(t *testing.T) {
	Test(t, "formatter-time: test-end writes count/total", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1, Total: 10})
		t.Match(buf.String(), "1/10")
		t.End()
	})
}

func TestTimeEventTestEndWritesTestName(t *testing.T) {
	Test(t, "formatter-time: test-end writes test name", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: foo", Count: 1, Total: 10})
		t.Match(buf.String(), "scope: foo")
		t.End()
	})
}

func TestTimeEventTestEndWithFail(t *testing.T) {
	Test(t, "formatter-time: test-end with failures writes red count", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: bar", Count: 1, Total: 10, Failed: 1})
		t.Match(buf.String(), regexp.MustCompile(`\033\[31m1\033\[0m`))
		t.End()
	})
}

func TestTimeEventTestEndWithFailScope(t *testing.T) {
	Test(t, "formatter-time: test-end with failures writes test name", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: bar", Count: 1, Total: 10, Failed: 1})
		t.Match(buf.String(), "scope: bar")
		t.End()
	})
}

func TestTimeEventNonTestEndDelegates(t *testing.T) {
	Test(t, "formatter-time: non-test-end events delegate to progress-bar", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		result := f.Event(stream.Event{Type: stream.TypeComment, Message: "note"})
		t.Equal(result, "# note\n")
		t.End()
	})
}

func TestTimeClockEnv(t *testing.T) {
	Test(t, "formatter-time: custom clock emoji appears in writer output", func(t *T) {
		t.TB().Setenv("TAPE_TIME_CLOCK", "X")
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Event(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1, Total: 10})
		t.Match(buf.String(), regexp.MustCompile(`X \d\d:\d\d`))
		t.End()
	})
}

func TestTimeClockEnvConstructs(t *testing.T) {
	Test(t, "formatter-time: New uses TAPE_TIME_CLOCK env var", func(t *T) {
		t.TB().Setenv("TAPE_TIME_CLOCK", "\U0001f550")
		f := formatter_time.New(10, &strings.Builder{})
		t.Ok(f != nil)
		t.End()
	})
}

func TestTimeEndDelegates(t *testing.T) {
	Test(t, "formatter-time: End delegates to progress-bar", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		result := f.End(5, 0, 0)
		t.Match(result, "# tests 5")
		t.End()
	})
}
