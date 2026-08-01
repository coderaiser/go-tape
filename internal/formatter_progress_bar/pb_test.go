package formatter_progress_bar

import (
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestNewDefaultColor(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: new default color", func(t *tape.T) {
		pb := New(100)
		t.Ok(pb.Color == DEFAULT_COLOR)
		t.End()
	})
}

func TestNewCustomColorEnv(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: new custom color env", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR_COLOR", "#ff0000")
		pb := New(100)
		t.Ok(pb.Color == "#ff0000")
		t.End()
	})
}

func TestNewProgressBarMin(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: new progress bar min", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR_MIN", "5")
		pb := New(5)
		t.Ok(pb.show)
		t.End()
	})
}

func TestNewProgressBarMinInvalid(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: new progress bar min invalid", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR_MIN", "bad")
		pb := New(5)
		t.NotOk(pb.show)
		t.End()
	})
}

func TestNewForceShow(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: new force show", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		pb := New(1)
		t.Ok(pb.show)
		t.End()
	})
}

func TestNewForceHide(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: new force hide", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(1000)
		t.NotOk(pb.show)
		t.End()
	})
}

func TestHexToANSIValid(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: hex to ansi valid", func(t *tape.T) {
		result := hexToANSI("#f9d472")
		t.Ok(strings.HasPrefix(result, "\033[38;2;"))
		t.End()
	})
}

func TestHexToANSINonHex(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: hex to ansi non hex", func(t *tape.T) {
		result := hexToANSI("\033[33m")
		t.Equal(result, "\033[33m")
		t.End()
	})
}

func TestHexToANSWrongLength(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: hex to ansi wrong length", func(t *tape.T) {
		result := hexToANSI("#fff")
		t.Equal(result, "#fff")
		t.End()
	})
}

func TestDecodeRuneAtASCII(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: decode rune at ascii", func(t *tape.T) {
		r, size := decodeRuneAt([]byte("A"), 0)
		t.Ok(r == 'A' && size == 1)
		t.End()
	})
}

func TestDecodeRuneAt2Byte(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: decode rune at 2 byte", func(t *tape.T) {
		r, size := decodeRuneAt([]byte("\u00e9"), 0)
		t.Ok(r == '\u00e9' && size == 2)
		t.End()
	})
}

func TestDecodeRuneAt3Byte(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: decode rune at 3 byte", func(t *tape.T) {
		r, size := decodeRuneAt([]byte("\u2588"), 0)
		t.Ok(r == '\u2588' && size == 3)
		t.End()
	})
}

func TestDecodeRuneAt4Byte(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: decode rune at 4 byte", func(t *tape.T) {
		s := "\U0001d11e"
		r, size := decodeRuneAt([]byte(s), 0)
		_ = r
		t.Equal(size, 4)
		t.End()
	})
}

func TestDecodeRuneAtInvalidLead(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: decode rune at invalid lead", func(t *tape.T) {
		r, size := decodeRuneAt([]byte{0x80}, 0)
		t.Ok(r == '\uFFFD' && size == 1)
		t.End()
	})
}

func TestDecodeRuneAtTruncated(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: decode rune at truncated", func(t *tape.T) {
		r, size := decodeRuneAt([]byte{0xE2, 0x96}, 0)
		t.Ok(r == '\uFFFD' && size == 1)
		t.End()
	})
}

func TestDecodeRuneAtInvalidContinuation(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: decode rune at invalid continuation", func(t *tape.T) {
		r, size := decodeRuneAt([]byte{0xC3, 0x00}, 0)
		t.Ok(r == '\uFFFD' && size == 1)
		t.End()
	})
}

func TestTruncateANSIASCII(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: truncate ansi ascii", func(t *tape.T) {
		result := truncateANSI("hello", 3)
		t.Equal(result, "hel")
		t.End()
	})
}

func TestTruncateANSIWithCodes(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: truncate ansi with codes", func(t *tape.T) {
		result := truncateANSI("\033[33mhello\033[0m", 3)
		t.Equal(result, "\033[33mhel")
		t.End()
	})
}

func TestTruncateANSIMultiByte(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: truncate ansi multi byte", func(t *tape.T) {
		result := truncateANSI("\u2588\u2588\u2588", 2)
		t.Equal(result, "\u2588\u2588")
		t.End()
	})
}

func TestRenderBarZeroTotal(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: render bar zero total", func(t *tape.T) {
		result := renderBar(0, 0, DEFAULT_COLOR)
		t.Ok(strings.Contains(result, string(barEmpty)))
		t.End()
	})
}

func TestRenderBarPartial(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: render bar partial", func(t *tape.T) {
		result := renderBar(5, 10, DEFAULT_COLOR)
		t.Ok(strings.Contains(result, string(barComplete)))
		t.End()
	})
}

func TestRenderBarOverfill(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: render bar overfill", func(t *tape.T) {
		result := renderBar(1000, 10, DEFAULT_COLOR)
		t.Ok(strings.Contains(result, string(barComplete)))
		t.End()
	})
}

func TestTruncateUnderLimit(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: truncate under limit", func(t *tape.T) {
		result := truncate("hello", 10)
		t.Equal(result, "hello")
		t.End()
	})
}

func TestTruncateOverLimit(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: truncate over limit", func(t *tape.T) {
		result := truncate("hello world!", 8)
		t.Equal(result, "hello...")
		t.End()
	})
}

func TestStartReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: start returns empty", func(t *tape.T) {
		pb := New(10)
		result := pb.Start(10)
		t.Equal(result, "")
		t.End()
	})
}

func TestTestReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: test returns empty", func(t *tape.T) {
		pb := New(10)
		result := pb.Test("scope: x")
		t.Equal(result, "")
		t.End()
	})
}

func TestSuccessReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: success returns empty", func(t *tape.T) {
		pb := New(10)
		result := pb.Success(1, "scope: x")
		t.Equal(result, "")
		t.End()
	})
}

func TestCommentReturnsLine(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: comment returns line", func(t *tape.T) {
		pb := New(10)
		result := pb.Comment("msg")
		t.Equal(result, "# msg\n")
		t.End()
	})
}

func TestTestEndNoShow(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: test end no show", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(1000)
		result := pb.TestEnd(1, 10, 0, "scope: x")
		t.Equal(result, "")
		t.End()
	})
}

func failOutput(pb *ProgressBarFormatter, count int, message, operator string, result, expected any, output, at, errorStack string) string {
	pb.Fail(count, message, operator, result, expected, output, at, errorStack)
	return pb.End(count, 0, 1, 0)
}

func TestFailWithOperator(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail with operator", func(t *tape.T) {
		result := failOutput(New(10), 1, "scope: x", "Equal", "got", "want", "", "", "")
		t.Ok(strings.Contains(result, "not ok 1"))
		t.End()
	})
}

func TestFailWithoutOperator(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail without operator", func(t *tape.T) {
		result := failOutput(New(10), 1, "scope: x", "", nil, nil, "raw output", "", "")
		t.Ok(strings.Contains(result, "raw output"))
		t.End()
	})
}

func TestFailWithAt(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail with at", func(t *tape.T) {
		result := failOutput(New(10), 1, "scope: x", "Equal", "got", "want", "", "file.go:10:", "")
		t.Ok(strings.Contains(result, "file.go:10:"))
		t.End()
	})
}

func TestFailWithErrorStack(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail with error stack", func(t *tape.T) {
		result := failOutput(New(10), 1, "scope: x", "Equal", "got", "want", "", "", "stack trace")
		t.Ok(strings.Contains(result, "stack trace"))
		t.End()
	})
}

func TestFailStackEnvDisabled(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail stack env disabled", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR_STACK", "0")
		pb := New(10)
		result := pb.Fail(1, "scope: x", "Equal", "got", "want", "", "", "stack trace")
		t.NotOk(strings.Contains(result, "stack trace"))
		t.End()
	})
}

func TestEndShowFalse(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end show false", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(1000)
		result := pb.End(5, 5, 0, 0)
		t.Ok(strings.Contains(result, "1..5"))
		t.End()
	})
}

func TestEndWithFailed(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end with failed", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		result := pb.End(5, 4, 1, 0)
		t.Ok(strings.Contains(result, "\u274c"))
		t.End()
	})
}

func TestEndWithSkipped(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end with skipped", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		result := pb.End(5, 5, 0, 1)
		t.Ok(strings.Contains(result, "\u26a0"))
		t.End()
	})
}

func TestEndShowTrue(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end show true", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		pb := New(10)
		pb.TestEnd(1, 10, 0, "scope: x")
		result := pb.End(5, 5, 0, 0)
		t.Ok(strings.Contains(result, "1..5"))
		t.End()
	})
}

func TestVisibleLen(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: visible len", func(t *tape.T) {
		length := visibleLen("hello")
		t.Equal(length, 5)
		t.End()
	})
}

func TestVisibleLenWithANSI(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: visible len with ansi", func(t *tape.T) {
		length := visibleLen("\033[33mhello\033[0m")
		t.Equal(length, 5)
		t.End()
	})
}

func TestTermWidth(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: term width", func(t *tape.T) {
		t.TB().Setenv("TAPE_TERM_WIDTH", "120")
		result := termWidth()
		t.Equal(result, 120)
		t.End()
	})
}

func TestTermWidthInvalid(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: term width invalid", func(t *tape.T) {
		t.TB().Setenv("TAPE_TERM_WIDTH", "0")
		result := termWidth()
		t.Ok(result > 0)
		t.End()
	})
}

func TestHexToANSIUpper(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: hex to ansi upper", func(t *tape.T) {
		result := hexToANSI("#ABCDEF")
		t.Ok(strings.HasPrefix(result, "\033[38;2;171;205;239m"))
		t.End()
	})
}

func TestTestEndWithFailOverfill(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: test end with fail overfill", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		t.TB().Setenv("TAPE_TERM_WIDTH", "5")
		pb := New(10)
		result := pb.TestEnd(100, 10, 3, "a very long scope name that will be truncated")
		t.Equal(result, "")
		t.End()
	})
}
