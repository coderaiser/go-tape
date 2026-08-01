package state

import (
	"errors"
	"testing"

	"github.com/coderaiser/go-tape/internal/model"
	"github.com/coderaiser/go-tape/internal/statemachine"
	"github.com/coderaiser/go-tape/internal/statemachine/adapters"
)

func TestRunEventCreatesRunningState(t *testing.T) {
	s, _ := New()
	_, err := s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	if err != nil {
		t.Fatal(err)
	}
	st, _ := s.Get("TestFoo")
	if st != StateRunning {
		t.Errorf("want %v, got %v", StateRunning, st)
	}
}

func TestPassEventMarksTestPassed(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "pass", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StatePassed {
		t.Errorf("want %v, got %v", StatePassed, st)
	}
}

func TestFailEventMarksTestFailed(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "fail", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StateFailed {
		t.Errorf("want %v, got %v", StateFailed, st)
	}
}

func TestSkipEventMarksTestSkipped(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "skip", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StateSkipped {
		t.Errorf("want %v, got %v", StateSkipped, st)
	}
}

func TestOutputAppendedToLogs(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "output", Test: "TestFoo", Output: "line1\n"})
	s.Apply(model.Event{Action: "output", Test: "TestFoo", Output: "line2\n"})
	got := s.GetOutput("TestFoo")
	if got != "line1\nline2\n" {
		t.Errorf("want line1\\nline2\\n, got %s", got)
	}
}

func TestPackageEventNoTestID(t *testing.T) {
	s, _ := New()
	_, err := s.Apply(model.Event{Action: "pass", Package: "mypkg"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidActionError(t *testing.T) {
	s, _ := New()
	_, err := s.Apply(model.Event{Test: "TestFoo", Action: "invalid"})
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestRunTwice(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StateRunning {
		t.Errorf("want %v, got %v", StateRunning, st)
	}
}

func TestGetNonExistent(t *testing.T) {
	s, _ := New()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent test")
	}
}

func TestApplyInvalidTransitionReturnsCurrentState(t *testing.T) {
	s, _ := New()
	// pass without run — invalid transition
	state, err := s.Apply(model.Event{Action: "pass", Test: "TestA"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != StateIdle {
		t.Errorf("want StateIdle, got %v", state)
	}
}

func TestSummaryKeepsSubtestName(t *testing.T) {
	s, _ := New()

	s.Apply(model.Event{
		Action: "run",
		Test:   "TestOnlyRuns/tape:_Only_delegates_to_Test",
	})

	s.Apply(model.Event{
		Action: "pass",
		Test:   "TestOnlyRuns/tape:_Only_delegates_to_Test",
	})

	passed, _, _ := s.Summary()

	if len(passed) != 1 {
		t.Fatalf("want 1 passed, got %d", len(passed))
	}

	if passed[0] != "TestOnlyRuns/tape:_Only_delegates_to_Test" {
		t.Fatalf("unexpected name: %s", passed[0])
	}
}

func TestSummaryPassedFailedSkipped(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "run", Test: "TestA"})
	s.Apply(model.Event{Action: "pass", Test: "TestA"})
	s.Apply(model.Event{Action: "output", Test: "TestA", Output: "ok"})
	s.Apply(model.Event{Action: "run", Test: "TestB"})
	s.Apply(model.Event{Action: "fail", Test: "TestB"})
	s.Apply(model.Event{Action: "output", Test: "TestB", Output: "fail"})
	s.Apply(model.Event{Action: "run", Test: "TestC"})
	s.Apply(model.Event{Action: "skip", Test: "TestC"})
	s.Apply(model.Event{Action: "output", Test: "TestC", Output: "skip"})
	passed, failed, skipped := s.Summary()
	if len(passed) != 1 {
		t.Errorf("want 1 passed, got %d", len(passed))
	}
	if len(failed) != 1 {
		t.Errorf("want 1 failed, got %d", len(failed))
	}
	if len(skipped) != 1 {
		t.Errorf("want 1 skipped, got %d", len(skipped))
	}
}

func TestParseTestStateUnknown(t *testing.T) {
	_, err := parseTestState("unknown")
	if err == nil {
		t.Fatal("expected error for unknown state")
	}
}

func TestParseTestEventUnknown(t *testing.T) {
	_, err := parseTestEvent("unknown")
	if err == nil {
		t.Fatal("expected error for unknown event")
	}
}

func TestSummaryWithOutputOnly(t *testing.T) {
	s, _ := New()
	// output-only test — never applied with a state transition
	s.Apply(model.Event{Action: "output", Test: "orphan", Output: "some output"})
	passed, failed, skipped := s.Summary()
	if len(passed)+len(failed)+len(skipped) != 0 {
		t.Errorf("expected 0 total, got pass=%d fail=%d skip=%d", len(passed), len(failed), len(skipped))
	}
}

func TestNew(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

// errSource always fails Load
type errSource struct{}

func (s errSource) Load() ([]statemachine.TransitionDef, error) {
	return nil, errors.New("source failed")
}

func TestNewFromSourceError(t *testing.T) {
	_, err := newFromSource(errSource{})
	if err == nil {
		t.Fatal("expected error from bad source")
	}
}

// errAdapter always fails Get
type errAdapter struct{ *adapters.Memory[TestState] }

func (a errAdapter) Get(id string) (*TestState, error) {
	return nil, errors.New("adapter error")
}

func TestGetAdapterError(t *testing.T) {
	s, _ := New()
	mem, ok := s.adapter.(*adapters.Memory[TestState])
	if !ok {
		t.Fatal("expected memory adapter")
	}
	s.adapter = errAdapter{mem}
	_, err := s.Get("TestFoo")
	if err == nil {
		t.Fatal("expected error from adapter")
	}
}

func TestSummaryAdapterError(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "output", Test: "TestFoo", Output: "ok"})
	mem, ok := s.adapter.(*adapters.Memory[TestState])
	if !ok {
		t.Fatal("expected memory adapter")
	}
	s.adapter = errAdapter{mem}
	// Summary silently skips on adapter error — should not panic
	passed, failed, skipped := s.Summary()
	_ = passed
	_ = failed
	_ = skipped
}

func TestMarkSkippedAddsUnseen(t *testing.T) {
	s, _ := New()
	if err := s.MarkSkipped([]string{"scope: foo", "scope: bar"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, _, skipped := s.Summary()
	if len(skipped) != 2 {
		t.Fatalf("want 2 skipped, got %d", len(skipped))
	}
}

func TestMarkSkippedDoesNotOverwritePassed(t *testing.T) {
	s, _ := New()
	s.Apply(model.Event{Action: "run", Test: "scope: foo"})
	s.Apply(model.Event{Action: "pass", Test: "scope: foo"})
	if err := s.MarkSkipped([]string{"scope: foo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	passed, _, skipped := s.Summary()
	if len(passed) != 1 {
		t.Fatalf("want 1 passed, got %d", len(passed))
	}
	if len(skipped) != 0 {
		t.Fatalf("want 0 skipped, got %d", len(skipped))
	}
}

func TestSummaryIgnoresRunningState(t *testing.T) {
	s, _ := New()
	// "run" leaves the test in StateRunning — this happens when go test -json -v
	// drops the terminal event for subtests with long names (a known Go toolchain issue).
	// Summary must silently ignore such tests rather than panicking.
	if _, err := s.Apply(model.Event{Action: "run", Test: "TestFoo"}); err != nil {
		t.Fatalf("Apply run: %v", err)
	}
	passed, failed, skipped := s.Summary()
	if len(passed) != 0 || len(failed) != 0 || len(skipped) != 0 {
		t.Fatalf("expected empty summary for running test, got passed=%v failed=%v skipped=%v", passed, failed, skipped)
	}
}

// setErrAdapter fails Set like errAdapter fails Get.
type setErrAdapter struct{ *adapters.Memory[TestState] }

func (a setErrAdapter) Set(id string, state TestState) error {
	return errors.New("adapter set error")
}

func TestMarkSkippedSetError(t *testing.T) {
	s, _ := New()
	// remove existing entry so Set is attempted
	mem, ok := s.adapter.(*adapters.Memory[TestState])
	if !ok {
		t.Fatal("expected memory adapter")
	}
	s.adapter = setErrAdapter{mem}
	if err := s.MarkSkipped([]string{"scope: foo"}); err == nil {
		t.Fatal("expected error from Set failure in MarkSkipped")
	}
}
