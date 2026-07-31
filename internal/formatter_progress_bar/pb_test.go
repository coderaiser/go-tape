package formatter_progress_bar

import (
	"strings"
	"testing"
)

func TestNewDefaultColor(t *testing.T) {
	pb := New(100)
	if pb.Color != DEFAULT_COLOR {
		t.Fatalf("want %s, got %s", DEFAULT_COLOR, pb.Color)
	}
}

func TestNewCustomColorEnv(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR_COLOR", "#ff0000")
	pb := New(100)
	if pb.Color != "#ff0000" {
		t.Fatalf("want #ff0000, got %s", pb.Color)
	}
}

func TestNewProgressBarMin(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR_MIN", "5")
	pb := New(5)
	if !pb.show {
		t.Fatal("expected show=true when total >= min")
	}
}

func TestNewProgressBarMinInvalid(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR_MIN", "bad")
	pb := New(5)
	if pb.show {
		t.Fatal("expected show=false when min invalid and total < 100")
	}
}

func TestNewForceShow(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR", "1")
	pb := New(1)
	if !pb.show {
		t.Fatal("expected show=true when TAPE_PROGRESS_BAR=1")
	}
}

func TestNewForceHide(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR", "0")
	pb := New(1000)
	if pb.show {
		t.Fatal("expected show=false when TAPE_PROGRESS_BAR=0")
	}
}

func TestHexToANSIValid(t *testing.T) {
	result := hexToANSI("#f9d472")
	if !strings.HasPrefix(result, "\033[38;2;") {
		t.Fatalf("expected ANSI prefix, got %q", result)
	}
}

func TestHexToANSINonHex(t *testing.T) {
	result := hexToANSI("\033[33m")
	if result != "\033[33m" {
		t.Fatalf("expected passthrough, got %q", result)
	}
}

func TestHexToANSWrongLength(t *testing.T) {
	result := hexToANSI("#fff")
	if result != "#fff" {
		t.Fatalf("expected passthrough, got %q", result)
	}
}

func TestDecodeRuneAtASCII(t *testing.T) {
	r, size := decodeRuneAt([]byte("A"), 0)
	if r != 'A' || size != 1 {
		t.Fatalf("want 'A'/1, got %c/%d", r, size)
	}
}
func TestDecodeRuneAt2Byte(t *testing.T) {
	r, size := decodeRuneAt([]byte("\u00e9"), 0)
	if r != '\u00e9' || size != 2 {
		t.Fatalf("want '\\u00e9'/2, got %c/%d", r, size)
	}
}

func TestDecodeRuneAt3Byte(t *testing.T) {
	r, size := decodeRuneAt([]byte("\u2588"), 0)
	if r != '\u2588' || size != 3 {
		t.Fatalf("want '\\u2588'/3, got %c/%d", r, size)
	}
}

func TestDecodeRuneAt4Byte(t *testing.T) {
	s := "\U0001d11e"
	r, size := decodeRuneAt([]byte(s), 0)
	if size != 4 {
		t.Fatalf("want size 4, got %d", size)
	}
	_ = r
}

func TestDecodeRuneAtInvalidLead(t *testing.T) {
	r, size := decodeRuneAt([]byte{0x80}, 0)
	if r != '\uFFFD' || size != 1 {
		t.Fatalf("want '\\uFFFD'/1, got %c/%d", r, size)
	}
}

func TestDecodeRuneAtTruncated(t *testing.T) {
	r, size := decodeRuneAt([]byte{0xE2, 0x96}, 0)
	if r != '\uFFFD' || size != 1 {
		t.Fatalf("want '\\uFFFD'/1, got %c/%d", r, size)
	}
}

func TestDecodeRuneAtInvalidContinuation(t *testing.T) {
	r, size := decodeRuneAt([]byte{0xC3, 0x00}, 0)
	if r != '\uFFFD' || size != 1 {
		t.Fatalf("want '\\uFFFD'/1, got %c/%d", r, size)
	}
}

func TestTruncateANSIASCII(t *testing.T) {
	result := truncateANSI("hello", 3)
	if result != "hel" {
		t.Fatalf("want hel, got %s", result)
	}
}

func TestTruncateANSIWithCodes(t *testing.T) {
	result := truncateANSI("\033[33mhello\033[0m", 3)
	if result != "\033[33mhel" {
		t.Fatalf("want \\033[33mhel, got %s", result)
	}
}

func TestTruncateANSIMultiByte(t *testing.T) {
	result := truncateANSI("\u2588\u2588\u2588", 2)
	if result != "\u2588\u2588" {
		t.Fatalf("want \\u2588\\u2588, got %s", result)
	}
}

func TestRenderBarZeroTotal(t *testing.T) {
	result := renderBar(0, 0, DEFAULT_COLOR)
	if !strings.Contains(result, string(barEmpty)) {
		t.Fatal("expected empty bar when total is 0")
	}
}

func TestRenderBarPartial(t *testing.T) {
	result := renderBar(5, 10, DEFAULT_COLOR)
	if !strings.Contains(result, string(barComplete)) {
		t.Fatal("expected complete chars in bar")
	}
}

func TestRenderBarOverfill(t *testing.T) {
	result := renderBar(1000, 10, DEFAULT_COLOR)
	if !strings.Contains(result, string(barComplete)) {
		t.Fatal("expected all complete chars")
	}
}

func TestTruncateUnderLimit(t *testing.T) {
	result := truncate("hello", 10)
	if result != "hello" {
		t.Fatalf("want hello, got %s", result)
	}
}

func TestTruncateOverLimit(t *testing.T) {
	result := truncate("hello world!", 8)
	if result != "hello..." {
		t.Fatalf("want hello..., got %s", result)
	}
}

func TestStartReturnsEmpty(t *testing.T) {
	pb := New(10)
	result := pb.Start(10)
	if result != "" {
		t.Fatalf("want empty, got %s", result)
	}
}

func TestTestReturnsEmpty(t *testing.T) {
	pb := New(10)
	result := pb.Test("scope: x")
	if result != "" {
		t.Fatalf("want empty, got %s", result)
	}
}

func TestSuccessReturnsEmpty(t *testing.T) {
	pb := New(10)
	result := pb.Success(1, "scope: x")
	if result != "" {
		t.Fatalf("want empty, got %s", result)
	}
}

func TestCommentReturnsLine(t *testing.T) {
	pb := New(10)
	result := pb.Comment("msg")
	if result != "# msg\n" {
		t.Fatalf("want '# msg\\n', got %s", result)
	}
}

func TestTestEndNoShow(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR", "0")
	pb := New(1000)
	result := pb.TestEnd(1, 10, 0, "scope: x")
	if result != "" {
		t.Fatalf("want empty, got %s", result)
	}
}

func failOutput(pb *ProgressBarFormatter, count int, message, operator string, result, expected any, output, at, errorStack string) string {
	pb.Fail(count, message, operator, result, expected, output, at, errorStack)
	return pb.End(count, 0, 1, 0)
}

func TestFailWithOperator(t *testing.T) {
	result := failOutput(New(10), 1, "scope: x", "Equal", "got", "want", "", "", "")
	if !strings.Contains(result, "not ok 1") {
		t.Fatalf("expected not ok 1, got %s", result)
	}
}

func TestFailWithoutOperator(t *testing.T) {
	result := failOutput(New(10), 1, "scope: x", "", nil, nil, "raw output", "", "")
	if !strings.Contains(result, "raw output") {
		t.Fatalf("expected raw output, got %s", result)
	}
}

func TestFailWithAt(t *testing.T) {
	result := failOutput(New(10), 1, "scope: x", "Equal", "got", "want", "", "file.go:10:", "")
	if !strings.Contains(result, "file.go:10:") {
		t.Fatalf("expected file.go:10:, got %s", result)
	}
}

func TestFailWithErrorStack(t *testing.T) {
	result := failOutput(New(10), 1, "scope: x", "Equal", "got", "want", "", "", "stack trace")
	if !strings.Contains(result, "stack trace") {
		t.Fatalf("expected stack trace, got %s", result)
	}
}

func TestFailStackEnvDisabled(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR_STACK", "0")
	pb := New(10)
	result := pb.Fail(1, "scope: x", "Equal", "got", "want", "", "", "stack trace")
	if strings.Contains(result, "stack trace") {
		t.Fatal("expected no stack trace when env is 0")
	}
}

func TestEndShowFalse(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR", "0")
	pb := New(1000)
	result := pb.End(5, 5, 0, 0)
	if !strings.Contains(result, "1..5") {
		t.Fatalf("expected 1..5, got %s", result)
	}
}

func TestEndWithFailed(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR", "0")
	pb := New(10)
	result := pb.End(5, 4, 1, 0)
	if !strings.Contains(result, "\u274c") {
		t.Fatal("expected fail emoji")
	}
}

func TestEndWithSkipped(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR", "0")
	pb := New(10)
	result := pb.End(5, 5, 0, 1)
	if !strings.Contains(result, "\u26a0") {
		t.Fatal("expected skip emoji")
	}
}

func TestEndShowTrue(t *testing.T) {
	t.Setenv("TAPE_PROGRESS_BAR", "1")
	pb := New(10)
	pb.TestEnd(1, 10, 0, "scope: x")
	result := pb.End(5, 5, 0, 0)
	if !strings.Contains(result, "1..5") {
		t.Fatalf("expected 1..5, got %s", result)
	}
}

func TestVisibleLen(t *testing.T) {
	length := visibleLen("hello")
	if length != 5 {
		t.Fatalf("want 5, got %d", length)
	}
}

func TestVisibleLenWithANSI(t *testing.T) {
	length := visibleLen("\033[33mhello\033[0m")
	if length != 5 {
		t.Fatalf("want 5, got %d", length)
	}
}

func TestTermWidth(t *testing.T) {
	t.Setenv("TAPE_TERM_WIDTH", "120")
	result := termWidth()
	if result != 120 {
		t.Fatalf("want 120, got %d", result)
	}
}

func TestTermWidthInvalid(t *testing.T) {
	t.Setenv("TAPE_TERM_WIDTH", "0")
	result := termWidth()
	if result <= 0 {
		t.Fatalf("want > 0, got %d", result)
	}
}
