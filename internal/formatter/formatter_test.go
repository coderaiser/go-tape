package formatter

import (
	"os"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/model"
)

// captureFormatter is a minimal Formatter that records every TestEnd call.
type captureFormatter struct {
	testEndCalls []testEndArgs
}

type testEndArgs struct {
	count, total, failed int
	name                 string
}

func (c *captureFormatter) Start(total int) string                   { return "" }
func (c *captureFormatter) Test(name string) string                  { return "" }
func (c *captureFormatter) Success(count int, message string) string { return "" }
func (c *captureFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	return ""
}
func (c *captureFormatter) Comment(message string) string { return "" }
func (c *captureFormatter) End(count, passed, failed, skipped int) string {
	return ""
}
func (c *captureFormatter) TestEnd(count, total, failed int, name string) string {
	c.testEndCalls = append(c.testEndCalls, testEndArgs{count, total, failed, name})
	return ""
}

func newCaptureState(total int) (*State, *captureFormatter) {
	cf := &captureFormatter{}
	s := &State{
		total:     total,
		outputs:   make(map[string][]string),
		formatter: cf,
		w:         os.Stderr, // discarded; captureFormatter returns ""
	}
	return s, cf
}

// TestCachedRunBarCompletion verifies that End emits a synthetic TestEnd(total,
// total) when s.count < s.total, so the progress bar reaches 100% on a cached
// run where plain-Go-subtest packages emit no events.
func TestCachedRunBarCompletion(t *testing.T) {
	const total = 10

	Test(t, "formatter: cached run pushes bar to 100%", func(t *T) {
		s, cf := newCaptureState(total)

		// Simulate only 3 out of 10 tests streaming (the rest were cached).
		for i := 0; i < 3; i++ {
			s.FromEvent(model.Event{Action: "run", Test: "scope: x"})
			s.FromEvent(model.Event{Action: "pass", Test: "scope: x"})
		}

		s.End(3, 0, 0)

		t.TB().Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		t.Ok(len(cf.testEndCalls) > 0)
		last := cf.testEndCalls[len(cf.testEndCalls)-1]
		t.Equal(last.count, total)
		t.Equal(last.total, total)
		t.End()
	})
}

// TestFullRunNoExtraTestEnd verifies that End does NOT emit an extra synthetic
// TestEnd when s.count already equals s.total (non-cached run).
func TestFullRunNoExtraTestEnd(t *testing.T) {
	const total = 3

	Test(t, "formatter: full run adds no extra TestEnd", func(t *T) {
		s, cf := newCaptureState(total)

		for i := 0; i < total; i++ {
			s.FromEvent(model.Event{Action: "run", Test: "scope: x"})
			s.FromEvent(model.Event{Action: "pass", Test: "scope: x"})
		}

		callsBefore := len(cf.testEndCalls)
		s.End(total, 0, 0)
		callsAfter := len(cf.testEndCalls)

		t.Equal(callsAfter, callsBefore)
		t.End()
	})
}

// TestFromEventRoutesCorrectly is a smoke test that FromEvent increments count
// for pass/fail/skip and ignores run and output.
func TestFromEventRoutesCorrectly(t *testing.T) {
	Test(t, "formatter: FromEvent routes count correctly", func(t *T) {
		var buf strings.Builder

		// Use a real progress-bar formatter with show forced off so we get no stderr noise.
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		s := New("progress-bar", &buf, 5)

		actions := []string{"run", "output", "pass", "fail", "skip"}
		for _, a := range actions {
			s.FromEvent(model.Event{Action: a, Test: "scope: x", Output: "# ok\n"})
		}

		// pass + fail + skip = 3
		t.TB().Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		t.Equal(s.count, 3)
		t.Equal(s.failed, 1)
		t.End()
	})
}

// TestFromEventIgnoresEmptyTest confirms that events with no Test name are skipped.
func TestFromEventIgnoresEmptyTest(t *testing.T) {
	Test(t, "formatter: empty test events are ignored", func(t *T) {
		s, cf := newCaptureState(5)
		s.FromEvent(model.Event{Action: "pass", Test: ""})
		t.TB().Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		t.Equal(s.count, 0)
		t.Equal(len(cf.testEndCalls), 0)
		t.End()
	})
}

// TestNewCIForcesTap verifies that New selects tap when CI env is set.
func TestNewCIForcesTap(t *testing.T) {
	Test(t, "formatter: CI env forces tap format", func(t *T) {
		var buf strings.Builder
		t.TB().Setenv("CI", "1")
		s := New("whatever", &buf, 1)
		s.FromEvent(model.Event{Action: "run", Test: "scope: x"})
		t.Ok(strings.Contains(buf.String(), "TAP version"))
		t.End()
	})
}

// TestNewEmptyFormatDefaultsToProgressBar verifies default format selection.
func TestNewEmptyFormatDefaultsToProgressBar(t *testing.T) {
	Test(t, "formatter: empty format defaults to progress-bar", func(t *T) {
		var buf strings.Builder
		t.TB().Setenv("CI", "0")
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		s := New("", &buf, 1)
		t.Ok(s.formatter != nil)
		t.End()
	})
}

func TestNewTapFormat(t *testing.T) {
	Test(t, "formatter: tap format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")

		var buf strings.Builder
		t.Ok(New("tap", &buf, 3) != nil)

		t.End()
	})

	Test(t, "formatter: short format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")

		var buf strings.Builder
		t.Ok(New("short", &buf, 3) != nil)

		t.End()
	})
	Test(t, "formatter: fail format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")

		var buf strings.Builder
		t.Ok(New("fail", &buf, 3) != nil)
		t.End()
	})

	Test(t, "formatter: time format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")

		var buf strings.Builder
		t.Ok(New("time", &buf, 3) != nil)

		t.End()
	})

	Test(t, "formatter: json-lines format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")

		var buf strings.Builder
		t.Ok(New("json-lines", &buf, 3) != nil)

		t.End()
	})

	Test(t, "formatter: progress-bar format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")

		var buf strings.Builder
		t.Ok(New("progress-bar", &buf, 3) != nil)

		t.End()
	})
	Test(t, "formatter: unknown format uses default", func(t *T) {
		t.TB().Setenv("CI", "false")

		var buf strings.Builder
		t.Ok(New("unknown-format", &buf, 3) != nil)

		t.End()
	})
}

// TestNewAllFormats verifies each named format constructs without panic.
func TestNewAllFormats(t *testing.T) {
	Test(t, "formatter: all named formats construct", func(t *T) {
		var buf strings.Builder
		reached := true
		for _, f := range []string{"tap", "short", "fail", "time", "json-lines"} {
			s := New(f, &buf, 3)
			if s == nil {
				reached = false
			}
		}
		t.Ok(reached)
		t.End()
	})
}

// TestTestLabelWithSlash covers the branch where the test name contains a slash.
func TestTestLabelWithSlash(t *testing.T) {
	Test(t, "formatter: testLabel strips prefix before slash", func(t *T) {
		t.Equal(testLabel("TestFoo/scope:_bar_baz"), "scope: bar baz")
		t.End()
	})
}

// TestTestLabelWithoutSlash covers the branch where the test name has no slash.
func TestTestLabelWithoutSlash(t *testing.T) {
	Test(t, "formatter: testLabel returns name unchanged when no slash", func(t *T) {
		t.Equal(testLabel("scope:_bar"), "scope: bar")
		t.End()
	})
}

// TestFileLinkEmptyAt covers the early-return branch in fileLink.
func TestFileLinkEmptyAt(t *testing.T) {
	Test(t, "formatter: fileLink returns empty string for empty at", func(t *T) {
		t.Equal(fileLink("", "/some/dir"), "")
		t.End()
	})
}

// TestFileLinkWithDir covers the branch where dir is set and at is relative.
func TestFileLinkWithDir(t *testing.T) {
	Test(t, "formatter: fileLink prepends dir to relative at", func(t *T) {
		t.Equal(fileLink("file.go:10:", "/proj"), "file:///proj/file.go:10")
		t.End()
	})
}

// TestFileLinkAbsoluteAt covers the branch where at is already absolute.
func TestFileLinkAbsoluteAt(t *testing.T) {
	Test(t, "formatter: fileLink does not prepend dir to absolute at", func(t *T) {
		t.Equal(fileLink("/abs/file.go:10:", "/proj"), "file:///abs/file.go:10")
		t.End()
	})
}

// TestFileLinkNoDir covers the branch where dir is empty.
func TestFileLinkNoDir(t *testing.T) {
	Test(t, "formatter: fileLink with no dir uses at as-is", func(t *T) {
		t.Equal(fileLink("file.go:5:", ""), "file://file.go:5")
		t.End()
	})
}

// TestWriteNonEmptyString exercises the discard-write path in write().
func TestWriteNonEmptyString(t *testing.T) {
	Test(t, "formatter: write emits non-empty output", func(t *T) {
		var buf strings.Builder
		s := New("tap", &buf, 3)
		s.FromEvent(model.Event{Action: "pass", Test: "scope: x"})
		t.Match(buf.String(), "ok 1")
		t.End()
	})
}