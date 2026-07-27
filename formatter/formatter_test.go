package formatter

import (
	"strings"
	"testing"

	"github.com/coderaiser/go-tape/model"
)

// -- output.go --

func TestParseOutputAt(t *testing.T) {
	t.Parallel()
	fields := ParseOutput([]string{"    t.go:42: \n"})
	if fields.At != "t.go:42:" {
		t.Errorf("expected .t.go:42., got %q", fields.At)
	}
}

func TestParseOutputResult(t *testing.T) {
	t.Parallel()
	fields := ParseOutput([]string{"    result:   1\n"})
	if fields.Result != "1" {
		t.Errorf("expected '1', got %q", fields.Result)
	}
}

func TestParseOutputExpected(t *testing.T) {
	t.Parallel()
	fields := ParseOutput([]string{"    expected: 2\n"})
	if fields.Expected != "2" {
		t.Errorf("expected '2', got %q", fields.Expected)
	}
}

func TestParseOutputRaw(t *testing.T) {
	t.Parallel()
	fields := ParseOutput([]string{"some output\n"})
	if !strings.Contains(fields.Raw, "some output") {
		t.Errorf("expected raw to contain 'some output', got %q", fields.Raw)
	}
}

func TestParseOutputEmpty(t *testing.T) {
	t.Parallel()
	fields := ParseOutput(nil)
	if fields.At != "" {
		t.Errorf("expected empty At, got %q", fields.At)
	}
}

func TestParseOutputOperator(t *testing.T) {
	t.Parallel()
	fields := ParseOutput([]string{"    Equal\n"})
	if fields.Operator != "Equal" {
		t.Errorf("expected 'Equal', got %q", fields.Operator)
	}
}

// -- tap.go --

func TestTAPStart(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Start(10)
	if out != "TAP version 13\n" {
		t.Errorf("expected 'TAP version 13\\n', got %q", out)
	}
}

func TestTAPTest(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Test("parser: run")
	if out != "# parser: run\n" {
		t.Errorf("expected '# parser: run\\n', got %q", out)
	}
}

func TestTAPSuccess(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Success(1, "parser: run")
	if out != "ok 1 parser: run\n" {
		t.Errorf("expected 'ok 1 parser: run\\n', got %q", out)
	}
}

func TestTAPFail(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Fail(1, "parser: bad", "Ok", false, true, "", "t.go:42:", "")
	if !strings.Contains(out, "not ok 1 parser: bad") {
		t.Errorf("expected 'not ok 1 parser: bad', got %q", out)
	}
}

func TestTAPFailOperator(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Fail(1, "parser: bad", "Ok", false, true, "", "", "")
	if !strings.Contains(out, "operator: Ok") {
		t.Errorf("expected 'operator: Ok', got %q", out)
	}
}

func TestTAPFailAt(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Fail(1, "x", "Ok", false, true, "", "t.go:42:", "")
	if !strings.Contains(out, "t.go:42:") {
		t.Errorf("expected .t.go:42., got %q", out)
	}
}

func TestTAPFailErrorStack(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Fail(1, "x", "Ok", false, true, "", "", "stack trace here")
	if !strings.Contains(out, "stack trace here") {
		t.Errorf("expected 'stack trace here', got %q", out)
	}
}

func TestTAPFailRawOutput(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Fail(1, "x", "", false, true, "diff output\n", "", "")
	if !strings.Contains(out, "diff output") {
		t.Errorf("expected 'diff output', got %q", out)
	}
}

func TestTAPComment(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.Comment("hello")
	if out != "# hello\n" {
		t.Errorf("expected '# hello\\n', got %q", out)
	}
}

func TestTAPEndOk(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.End(2, 2, 0, 0)
	if !strings.Contains(out, "# ok") {
		t.Errorf("expected '# ok', got %q", out)
	}
}

func TestTAPEndPlan(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.End(2, 2, 0, 0)
	if !strings.Contains(out, "1..2") {
		t.Errorf("expected '1..2', got %q", out)
	}
}

func TestTAPEndFail(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.End(2, 1, 1, 0)
	if !strings.Contains(out, "# fail 1") {
		t.Errorf("expected '# fail 1', got %q", out)
	}
}

func TestTAPEndSkip(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.End(1, 1, 0, 1)
	if !strings.Contains(out, "# skip 1") {
		t.Errorf("expected '# skip 1', got %q", out)
	}
}

func TestTAPTestEnd(t *testing.T) {
	t.Parallel()
	f := NewTAP()
	out := f.TestEnd(1, 10, 0, "x")
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

// -- short.go --

func TestShortNoStack(t *testing.T) {
	t.Parallel()
	f := NewShort()
	out := f.Fail(1, "x", "Ok", false, true, "", "", "stack trace here")
	if strings.Contains(out, "stack trace here") {
		t.Errorf("expected no stack trace, got %q", out)
	}
}

func TestShortHasOperator(t *testing.T) {
	t.Parallel()
	f := NewShort()
	out := f.Fail(1, "x", "Ok", false, true, "", "", "")
	if !strings.Contains(out, "operator: Ok") {
		t.Errorf("expected 'operator: Ok', got %q", out)
	}
}

// -- fail.go --

func TestFailFormatterPrefixesTestName(t *testing.T) {
	t.Parallel()
	f := NewFail()
	f.Test("parser: run")
	out := f.Fail(1, "should be truthy", "Ok", false, true, "", "", "")
	if !strings.Contains(out, "# parser: run") {
		t.Errorf("expected '# parser: run', got %q", out)
	}
}

func TestFailFormatterTestReturnsEmpty(t *testing.T) {
	t.Parallel()
	f := NewFail()
	out := f.Test("parser: run")
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

// -- json_lines.go --

func TestJSONLinesTestEnd(t *testing.T) {
	t.Parallel()
	f := NewJSONLines(3)
	out := f.TestEnd(1, 3, 0, "parser: run")
	if !strings.Contains(out, "\"count\":1") {
		t.Errorf("expected 'count:1', got %q", out)
	}
}

func TestJSONLinesTestEndTotal(t *testing.T) {
	t.Parallel()
	f := NewJSONLines(3)
	out := f.TestEnd(1, 3, 0, "parser: run")
	if !strings.Contains(out, "\"total\":3") {
		t.Errorf("expected 'total:3', got %q", out)
	}
}

func TestJSONLinesFail(t *testing.T) {
	t.Parallel()
	f := NewJSONLines(3)
	out := f.Fail(2, "parser: bad", "Ok", false, true, "", "t.go:1:", "")
	if !strings.Contains(out, "\"operator\":\"Ok\"") {
		t.Errorf("expected 'operator:Ok', got %q", out)
	}
}

func TestJSONLinesEnd(t *testing.T) {
	t.Parallel()
	f := NewJSONLines(2)
	out := f.End(2, 1, 1, 0)
	if !strings.Contains(out, "\"passed\":1") {
		t.Errorf("expected 'passed:1', got %q", out)
	}
}

func TestJSONLinesEndSkipped(t *testing.T) {
	t.Parallel()
	f := NewJSONLines(1)
	out := f.End(0, 0, 0, 1)
	if !strings.Contains(out, "\"skipped\":1") {
		t.Errorf("expected 'skipped:1', got %q", out)
	}
}

func TestJSONLinesStartEmpty(t *testing.T) {
	t.Parallel()
	f := NewJSONLines(3)
	out := f.Start(3)
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestJSONLinesCommentEmpty(t *testing.T) {
	t.Parallel()
	f := NewJSONLines(3)
	out := f.Comment("hello")
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

// -- progress_bar.go --

func TestProgressBarStart(t *testing.T) {
	t.Parallel()
	f := NewProgressBar(10)
	out := f.Start(10)
	if out != "TAP version 13\n" {
		t.Errorf("expected 'TAP version 13\\n', got %q", out)
	}
}

func TestProgressBarTestReturnsEmpty(t *testing.T) {
	t.Parallel()
	f := NewProgressBar(10)
	out := f.Test("parser: run")
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestProgressBarSuccessReturnsEmpty(t *testing.T) {
	t.Parallel()
	f := NewProgressBar(10)
	out := f.Success(1, "parser: run")
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestProgressBarEndOk(t *testing.T) {
	t.Parallel()
	f := NewProgressBar(1)
	out := f.End(1, 1, 0, 0)
	if !strings.Contains(out, "✅") {
		t.Errorf("expected ok mark, got %q", out)
	}
}

func TestProgressBarEndFail(t *testing.T) {
	t.Parallel()
	f := NewProgressBar(1)
	out := f.End(1, 0, 1, 0)
	if !strings.Contains(out, "❌") {
		t.Errorf("expected fail emoji, got %q", out)
	}
}

func TestProgressBarEndSkip(t *testing.T) {
	t.Parallel()
	f := NewProgressBar(1)
	out := f.End(1, 1, 0, 1)
	if !strings.Contains(out, "⚠️") {
		t.Errorf("expected skip emoji, got %q", out)
	}
}

func TestRenderBar(t *testing.T) {
	t.Parallel()
	bar := RenderBar(5, 10, "")
	if strings.Count(bar, "█")+strings.Count(bar, "░") != 40 {
		t.Errorf("expected 40 chars, got %d", strings.Count(bar, "█")+strings.Count(bar, "░"))
	}
}

func TestRenderBarZeroTotal(t *testing.T) {
	t.Parallel()
	bar := RenderBar(0, 0, "")
	if strings.Count(bar, "░") != 40 {
		t.Errorf("expected 40 empty chars, got %d", strings.Count(bar, "░"))
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()
	out := Truncate("hello world this is a very long test name indeed", 20)
	if len([]rune(out)) > 20 {
		t.Errorf("expected <=20 runes, got %d", len([]rune(out)))
	}
}

func TestTruncateShort(t *testing.T) {
	t.Parallel()
	out := Truncate("hello", 20)
	if out != "hello" {
		t.Errorf("expected 'hello', got %q", out)
	}
}

// -- New factory --

func TestNewDefaultCI(t *testing.T) {
	t.Setenv("CI", "true")
	var buf strings.Builder
	s := New("", &buf, 0)
	if s == nil {
		t.Fatal("expected non-nil state")
	}
}

func TestNewExplicitFormat(t *testing.T) {
	var buf strings.Builder
	s := New("tap", &buf, 5)
	if s == nil {
		t.Fatal("expected non-nil state")
	}
}

func TestNewWritesHeader(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	New("tap", &buf, 5)
	if !strings.Contains(buf.String(), "TAP version 13") {
		t.Errorf("expected 'TAP version 13', got %q", buf.String())
	}
}

func TestNewShortFormat(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("short", &buf, 5)
	if s == nil {
		t.Fatal("expected non-nil state")
	}
}

func TestNewFailFormat(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("fail", &buf, 5)
	if s == nil {
		t.Fatal("expected non-nil state")
	}
}

func TestNewTimeFormat(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("time", &buf, 5)
	if s == nil {
		t.Fatal("expected non-nil state")
	}
}

func TestNewJSONLinesFormat(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("json-lines", &buf, 5)
	if s == nil {
		t.Fatal("expected non-nil state")
	}
}

func TestNewProgressBarDefault(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("progress-bar", &buf, 5)
	if s == nil {
		t.Fatal("expected non-nil state")
	}
}

// -- FromEvent --

func TestFromEventRun(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("tap", &buf, 1)
	s.FromEvent(model.Event{Action: "run", Test: "parser: run"})
	if !strings.Contains(buf.String(), "# parser: run") {
		t.Errorf("expected '# parser: run', got %q", buf.String())
	}
}

func TestFromEventPass(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("tap", &buf, 1)
	s.FromEvent(model.Event{Action: "pass", Test: "parser: run"})
	if !strings.Contains(buf.String(), "ok 1 parser: run") {
		t.Errorf("expected 'ok 1 parser: run', got %q", buf.String())
	}
}

func TestFromEventSkip(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("tap", &buf, 1)
	s.FromEvent(model.Event{Action: "skip", Test: "parser: skip"})
	if strings.Contains(buf.String(), "ok)") || strings.Contains(buf.String(), "not ok") {
		t.Errorf("expected no ok/not ok on skip, got %q", buf.String())
	}
}

func TestFromEventEmptyTest(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("tap", &buf, 1)
	initialLen := buf.Len()
	s.FromEvent(model.Event{Action: "output", Package: "mypkg"})
	if buf.Len() != initialLen {
		t.Errorf("expected no write when Test is empty, initial=%d, after=%d", initialLen, buf.Len())
	}
}

func TestFromEventEnd(t *testing.T) {
	t.Parallel()
	var buf strings.Builder
	s := New("tap", &buf, 1)
	s.End(0, 0, 0)
	if !strings.Contains(buf.String(), "# ok") {
		t.Errorf("expected '# ok', got %q", buf.String())
	}
}
