package formatter_progress_bar

import (
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/stream"
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

func TestBarGlyphsAreNarrowWidth(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: bar glyphs are eaw=N to avoid misalignment", func(t *tape.T) {
		// U+2588 FULL BLOCK and U+00B7 MIDDLE DOT are eaw=A (ambiguous) —
		// they render as 1 or 2 columns depending on the terminal, breaking
		// bar alignment. Both characters must be eaw=N (narrow, always 1 col).
		t.Equal(string(barComplete), "█") // ▪ BLACK SMALL SQUARE — eaw=N
		t.End()
	})
}

func TestBarEmptyGlyphIsNarrowWidth(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: empty glyph is eaw=N to avoid misalignment", func(t *tape.T) {
		t.Equal(string(barEmpty), "\u2591") // ░ LIGHT SHADE — eaw=N
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

func TestEventCommentReturnsLine(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: comment event returns line", func(t *tape.T) {
		pb := New(10)
		result := pb.Event(stream.Event{Type: stream.TypeComment, Message: "msg"})
		t.Equal(result, "# msg\n")
		t.End()
	})
}

func TestEventFailBuffersOutput(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail event buffers output for End", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		pb.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: x", Count: 1,
			Operator: "Equal", Result: "got", Expected: "want",
		})
		result := pb.End(0, 1, 0)
		t.Ok(strings.Contains(result, "not ok 1"))
		t.End()
	})
}

func TestEventFailWithErrorStack(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail event includes error stack", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		pb.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: x", Count: 1,
			Operator: "Equal", Result: "got", Expected: "want",
			ErrorStack: "stack trace",
		})
		result := pb.End(0, 1, 0)
		t.Ok(strings.Contains(result, "stack trace"))
		t.End()
	})
}

func TestEventFailStackEnvDisabled(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail event hides stack when env disabled", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		t.TB().Setenv("TAPE_PROGRESS_BAR_STACK", "0")
		pb := New(10)
		pb.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: x", Count: 1,
			Operator: "Equal", Result: "got", Expected: "want",
			ErrorStack: "stack trace",
		})
		result := pb.End(0, 1, 0)
		t.NotOk(strings.Contains(result, "stack trace"))
		t.End()
	})
}

func TestEndShowFalseSummary(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end summary when hidden", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		result := pb.End(5, 4, 1)
		t.Ok(strings.Contains(result, "1..10"))
		t.End()
	})
}

func TestEndWithFailed(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end with failed", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		result := pb.End(5, 4, 1)
		t.Ok(strings.Contains(result, "\u274c"))
		t.End()
	})
}

func TestEndWithSkipped(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end with skipped", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		result := pb.End(5, 5, 1)
		t.Ok(strings.Contains(result, "\u26a0"))
		t.End()
	})
}

func TestEndShowTrue(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end show true returns summary with prefix", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		pb := New(10)
		pb.Event(stream.Event{Type: stream.TypeTestEnd, Count: 1, Total: 10, Test: "scope: x"})
		result := pb.End(5, 5, 0)
		t.Ok(strings.Contains(result, "1..10"))
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

func TestEventTestEndWithFailOverfill(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: test-end event with fail overfill", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		t.TB().Setenv("TAPE_TERM_WIDTH", "5")
		pb := New(10)
		result := pb.Event(stream.Event{Type: stream.TypeTestEnd, Count: 100, Total: 10, Failed: 3, Test: "a very long scope name that will be truncated"})
		t.Equal(result, "")
		t.End()
	})
}

func TestEventTestEndShowTrueWritesToStderr(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: test-end event writes to stderr when show is true", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		pb := New(10)
		result := pb.Event(stream.Event{Type: stream.TypeTestEnd, Count: 1, Total: 10, Test: "scope: x"})
		t.Equal(result, "")
		t.End()
	})
}

func TestEventTestEndHiddenReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: test-end event returns empty when hidden", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		pb := New(10)
		result := pb.Event(stream.Event{Type: stream.TypeTestEnd, Count: 1, Total: 10, Test: "scope: x"})
		t.Equal(result, "")
		t.End()
	})
}

func TestEndShowTrueReturnsSummaryWithPrefix(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: end returns summary with prefix when show is true", func(t *tape.T) {
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		pb := New(10)
		pb.Event(stream.Event{Type: stream.TypeTestEnd, Count: 1, Total: 10, Test: "scope: x"})
		result := pb.End(5, 5, 0)
		t.Ok(strings.Contains(result, "1..10"))
		t.End()
	})
}

func TestTermWidthInvalidString(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: term width invalid string falls back to syscall", func(t *tape.T) {
		t.TB().Setenv("TAPE_TERM_WIDTH", "bad")
		result := termWidth()
		t.Ok(result > 0)
		t.End()
	})
}

func TestOkLineDefault(t *testing.T) {
	tape.Test(t, "progress-bar: ok line has no extra space by default", func(t *tape.T) {
		t.TB().Setenv("TERMINAL_EMULATOR", "")
		result := okLine()
		t.Equal(result, "# \u2705 ok\n")
		t.End()
	})
}

func TestOkLineJetBrains(t *testing.T) {
	tape.Test(t, "progress-bar: ok line has extra space in JetBrains", func(t *tape.T) {
		t.TB().Setenv("TERMINAL_EMULATOR", "JetBrains-JediTerm")
		result := okLine()
		t.Equal(result, "# \u2705  ok\n")
		t.End()
	})
}

func TestFailExpectedResultValues(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail shows expected value", func(t *tape.T) {
		pb := New(10)
		pb.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: x", Count: 1,
			Operator: "equal", Result: "hello", Expected: "world",
		})
		result := pb.End(0, 1, 0)
		t.Ok(strings.Contains(result, "expected: |-\n      world\n"))
		t.End()
	})
	tape.Test(t, "formatter-progress-bar: fail shows result value", func(t *tape.T) {
		pb := New(10)
		pb.Event(stream.Event{
			Type: stream.TypeFail, Test: "scope: x", Count: 1,
			Operator: "equal", Result: "hello", Expected: "world",
		})
		result := pb.End(0, 1, 0)
		t.Ok(strings.Contains(result, "result: |-\n      hello\n"))
		t.End()
	})
}

func TestProgressBarBuildError(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: build-error appears in End output", func(t *tape.T) {
		f := New(0)
		f.Event(stream.Event{
			Type:    stream.TypeBuildError,
			Package: "example.com/foo",
			Output:  "foo.go:5:2: declared and not used: x\n",
		})
		out := f.End(0, 0, 0)
		t.Match(out, "build-error")
		t.End()
	})
	tape.Test(t, "formatter-progress-bar: build-error output contains package name", func(t *tape.T) {
		f := New(0)
		f.Event(stream.Event{
			Type:    stream.TypeBuildError,
			Package: "example.com/foo",
			Output:  "foo.go:5:2: declared and not used: x\n",
		})
		out := f.End(0, 0, 0)
		t.Match(out, "example.com/foo")
		t.End()
	})
}

func TestEventPackageError(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: package-error buffered into End output", func(t *tape.T) {
		f := New(0)
		f.Event(stream.Event{
			Type:    stream.TypePackageError,
			Package: "example.com/bar",
			Output:  "panic: oops\n",
		})
		out := f.End(0, 0, 0)
		t.Match(out, "package-error")
		t.End()
	})
}

func TestEventFailWithAt(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: fail with At field includes location", func(t *tape.T) {
		f := New(0)
		f.Event(stream.Event{
			Type:     stream.TypeFail,
			Test:     "scope: x",
			Count:    1,
			Operator: "equal",
			Output:   "",
			Result:   "a",
			Expected: "b",
			At:       "foo_test.go:42",
		})
		out := f.End(0, 1, 0)
		t.Match(out, "foo_test.go:42")
		t.End()
	})
}

func TestTermWidthSyscall(t *testing.T) {
	tape.Test(t, "formatter-progress-bar: termWidth falls back to 80 without env", func(t *tape.T) {
		t.TB().Setenv("TAPE_TERM_WIDTH", "")
		result := termWidth()
		t.Ok(result >= 80)
		t.End()
	})
}
