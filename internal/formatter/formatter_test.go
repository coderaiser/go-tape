package formatter

import (
	"os"
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"
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

	tape.Test(t, "formatter: cached run pushes bar to 100%", func(t *tape.T) {
		s, cf := newCaptureState(total)

		// Simulate only 3 out of 10 tests streaming (the rest were cached).
		for i := 0; i < 3; i++ {
			s.FromEvent(model.Event{Action: "run", Test: "scope: x"})
			s.FromEvent(model.Event{Action: "pass", Test: "scope: x"})
		}

		s.End(3, 0, 0)

		t.Ok(len(cf.testEndCalls) > 0 &&
			cf.testEndCalls[len(cf.testEndCalls)-1].count == total &&
			cf.testEndCalls[len(cf.testEndCalls)-1].total == total)
		t.End()
	})
}

// TestFullRunNoExtraTestEnd verifies that End does NOT emit an extra synthetic
// TestEnd when s.count already equals s.total (non-cached run).
func TestFullRunNoExtraTestEnd(t *testing.T) {
	const total = 3

	tape.Test(t, "formatter: full run adds no extra TestEnd", func(t *tape.T) {
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
	tape.Test(t, "formatter: FromEvent routes count correctly", func(t *tape.T) {
		var buf strings.Builder

		// Use a real progress-bar formatter with show forced off so we get no stderr noise.
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		s := New("progress-bar", &buf, 5)

		actions := []string{"run", "output", "pass", "fail", "skip"}
		for _, a := range actions {
			s.FromEvent(model.Event{Action: a, Test: "scope: x", Output: "# ok\n"})
		}

		// pass + fail + skip = 3
		t.Ok(s.count == 3 && s.failed == 1)
		t.End()
	})
}

// TestFromEventIgnoresEmptyTest confirms that events with no Test name are skipped.
func TestFromEventIgnoresEmptyTest(t *testing.T) {
	tape.Test(t, "formatter: empty test events are ignored", func(t *tape.T) {
		s, cf := newCaptureState(5)
		s.FromEvent(model.Event{Action: "pass", Test: ""})
		t.Ok(s.count == 0 && len(cf.testEndCalls) == 0)
		t.End()
	})
}

// TestNewCIForcesTap verifies that New selects tap when CI env is set.
func TestNewCIForcesTap(t *testing.T) {
	tape.Test(t, "formatter: CI env forces tap format", func(t *tape.T) {
		var buf strings.Builder
		t.TB().Setenv("CI", "true")
		s := New("whatever", &buf, 1)
		s.FromEvent(model.Event{Action: "run", Test: "scope: x"})
		t.Ok(strings.Contains(buf.String(), "TAP version"))
		t.End()
	})
}

// TestNewEmptyFormatDefaultsToProgressBar verifies default format selection.
func TestNewEmptyFormatDefaultsToProgressBar(t *testing.T) {
	tape.Test(t, "formatter: empty format defaults to progress-bar", func(t *tape.T) {
		var buf strings.Builder
		t.TB().Setenv("CI", "0")
		t.TB().Setenv("TAPE_PROGRESS_BAR", "1")
		s := New("", &buf, 1)
		t.Ok(s.formatter != nil)
		t.End()
	})
}

// TestNewAllFormats verifies each named format constructs without panic.
func TestNewAllFormats(t *testing.T) {
	tape.Test(t, "formatter: all named formats construct", func(t *tape.T) {
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

// TestWriteNonEmptyString exercises the discard-write path in write().
func TestWriteNonEmptyString(t *testing.T) {
	tape.Test(t, "formatter: write emits non-empty output", func(t *tape.T) {
		var buf strings.Builder
		s := New("tap", &buf, 3)
		s.FromEvent(model.Event{Action: "pass", Test: "scope: x"})
		t.Match(buf.String(), "ok 1")
		t.End()
	})
}
