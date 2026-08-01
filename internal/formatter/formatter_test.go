package formatter

import (
	"os"
	"strings"
	"testing"

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

	s, cf := newCaptureState(total)

	// Simulate only 3 out of 10 tests streaming (the rest were cached).
	for i := 0; i < 3; i++ {
		s.FromEvent(model.Event{Action: "run", Test: "scope: x"})
		s.FromEvent(model.Event{Action: "pass", Test: "scope: x"})
	}

	s.End(3, 0, 0)

	if len(cf.testEndCalls) == 0 {
		t.Fatal("expected at least one TestEnd call from End()")
	}

	last := cf.testEndCalls[len(cf.testEndCalls)-1]
	if last.count != total || last.total != total {
		t.Fatalf("expected final TestEnd(%d,%d,...), got TestEnd(%d,%d,...)",
			total, total, last.count, last.total)
	}
}

// TestFullRunNoExtraTestEnd verifies that End does NOT emit an extra synthetic
// TestEnd when s.count already equals s.total (non-cached run).
func TestFullRunNoExtraTestEnd(t *testing.T) {
	const total = 3

	s, cf := newCaptureState(total)

	for i := 0; i < total; i++ {
		s.FromEvent(model.Event{Action: "run", Test: "scope: x"})
		s.FromEvent(model.Event{Action: "pass", Test: "scope: x"})
	}

	callsBefore := len(cf.testEndCalls)
	s.End(total, 0, 0)
	callsAfter := len(cf.testEndCalls)

	if callsAfter != callsBefore {
		t.Fatalf("expected no extra TestEnd when count==total, got %d extra call(s)",
			callsAfter-callsBefore)
	}
}

// TestFromEventRoutesCorrectly is a smoke test that FromEvent increments count
// for pass/fail/skip and ignores run and output.
func TestFromEventRoutesCorrectly(t *testing.T) {
	var buf strings.Builder

	// Use a real progress-bar formatter with show forced off so we get no stderr noise.
	t.Setenv("TAPE_PROGRESS_BAR", "0")
	s := New("progress-bar", &buf, 5)

	actions := []string{"run", "output", "pass", "fail", "skip"}
	for _, a := range actions {
		s.FromEvent(model.Event{Action: a, Test: "scope: x", Output: "# ok\n"})
	}

	// pass + fail + skip = 3
	if s.count != 3 {
		t.Fatalf("expected count=3, got %d", s.count)
	}
	if s.failed != 1 {
		t.Fatalf("expected failed=1, got %d", s.failed)
	}
}

// TestFromEventIgnoresEmptyTest confirms that events with no Test name are skipped.
func TestFromEventIgnoresEmptyTest(t *testing.T) {
	s, cf := newCaptureState(5)
	s.FromEvent(model.Event{Action: "pass", Test: ""})
	if s.count != 0 {
		t.Fatalf("expected count=0, got %d", s.count)
	}
	if len(cf.testEndCalls) != 0 {
		t.Fatalf("expected no TestEnd calls, got %d", len(cf.testEndCalls))
	}
}
