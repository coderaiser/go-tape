package formatter_time_test

import (
	"strings"
	"testing"
	"regexp"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_time"
)

func TestTimeFormatterStart(t *testing.T) {
	Test(t, "formatter-time: Start returns empty string", func(t *T) {
		f := formatter_time.New(10, &strings.Builder{})
		result := f.Start(10)
		t.Equal(result, "")
		t.End()
	})
}

func TestTimeFormatterTestEndReturnsEmpty(t *testing.T) {
	Test(t, "formatter-time: TestEnd returns empty string", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		result := f.TestEnd(1, 10, 0, "scope: foo")
		t.Equal(result, "")
		t.End()
	})
}

func TestTimeFormatterTestEndWritesToWriter(t *testing.T) {
	Test(t, "formatter-time: TestEnd writes progress line to writer: ^\r", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		f.TestEnd(1, 10, 0, "scope: foo")
		out := buf.String()
		t.Match(out, regexp.MustCompile(`^\r`))
		t.End()
	})
	Test(t, "formatter-time: TestEnd writes progress line to writer: 10%", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		f.TestEnd(1, 10, 0, "scope: foo")
		out := buf.String()
		t.Match(out, `10%`)
		t.End()
	})
	Test(t, "formatter-time: TestEnd writes progress line to writer: 1/10", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		f.TestEnd(1, 10, 0, "scope: foo")
		out := buf.String()
		t.Match(out, `1/10`)
		t.End()
	})
	Test(t, "formatter-time: TestEnd writes progress line to writer: scope", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		f.TestEnd(1, 10, 0, "scope: foo")
		out := buf.String()
		t.Match(out, `scope: foo`)
		t.End()
	})
}

func TestTimeFormatterTestEndWithFail(t *testing.T) {
	Test(t, "formatter-time: TestEnd formats failure count in red", func(t *T) {
		f := formatter_time.New(10, &strings.Builder{})
		f.Start(10)
		result := f.TestEnd(1, 10, 1, "scope: foo")
		t.Equal(result, "")
		t.End()
	})
}

func TestTimeFormatterTestEndFailWritesToWriter(t *testing.T) {
	Test(t, "formatter-time: TestEnd with failures writes red count to writer", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		f.TestEnd(1, 10, 1, "scope: bar")
		out := buf.String()
		t.Match(out, regexp.MustCompile(`\033\[31m1\033\[0m`))
		t.End()
	})
	Test(t, "formatter-time: TestEnd with failures writes red count to writer: scope", func(t *T) {
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		f.TestEnd(1, 10, 1, "scope: bar")
		out := buf.String()
		t.Match(out, `scope: bar`)
		t.End()
	})
}

func TestTimeFormatterClockEnv(t *testing.T) {
	Test(t, "formatter-time: New uses TAPE_TIME_CLOCK env var", func(t *T) {
		t.TB().Setenv("TAPE_TIME_CLOCK", "\U0001f550")
		f := formatter_time.New(10, &strings.Builder{})
		t.Ok(f != nil)
		t.End()
	})
}

func TestTimeFormatterClockEnvAppearsInOutput(t *testing.T) {
	Test(t, "formatter-time: custom clock emoji appears in writer output", func(t *T) {
		t.TB().Setenv("TAPE_TIME_CLOCK", "X")
		var buf strings.Builder
		f := formatter_time.New(10, &buf)
		f.Start(10)
		f.TestEnd(1, 10, 0, "scope: x")
		t.Match(buf.String(), regexp.MustCompile(`X \d\d:\d\d`))
		t.End()
	})
}
