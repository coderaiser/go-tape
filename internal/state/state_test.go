package state_test

import (
	"errors"
	"sort"
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/model"
	"github.com/coderaiser/go-tape/internal/state"
	"github.com/coderaiser/go-tape/internal/statemachine"
	"github.com/coderaiser/go-tape/internal/statemachine/adapters"
)

type StateT struct{ *tape.T }

func (t *StateT) NewStore() *state.Store {
	t.TB().Helper()
	s, error := state.New()
	if error != nil {
		t.TB().Fatalf("New: %v", error)
	}
	return s
}

// Apply applies events, failing on infrastructure errors.
func (t *StateT) Apply(s *state.Store, events ...model.Event) {
	t.TB().Helper()
	for _, e := range events {
		if _, error := s.Apply(e); error != nil {
			t.TB().Fatalf("Apply: %v", error)
		}
	}
}

// ApplyState applies events and returns the resulting state, or the error.
func (t *StateT) ApplyState(s *state.Store, e model.Event) (state.TestState, error) {
	t.TB().Helper()
	return s.Apply(e)
}

func (t *StateT) GetError(s *state.Store, name string) error {
	t.TB().Helper()
	_, error := s.Get(name)
	return error
}

func (t *StateT) MarkSkippedErr(s *state.Store, names []string) error {
	t.TB().Helper()
	return s.MarkSkipped(names)
}

func (t *StateT) ParseState(s string) error {
	t.TB().Helper()
	_, error := state.ParseTestState(s)
	return error
}

func (t *StateT) ParseEvent(e string) error {
	t.TB().Helper()
	_, error := state.ParseTestEvent(e)
	return error
}

func (t *StateT) NewFromSource(src statemachine.TransitionSource) (*state.Store, error) {
	t.TB().Helper()
	return state.NewFromSource(src)
}

func (t *StateT) StateIs(s *state.Store, name string, expected state.TestState) {
	t.TB().Helper()
	st, error := s.Get(name)
	if error != nil {
		t.TB().Fatalf("Get: %v", error)
	}
	t.Equal(st, expected)
}

func (t *StateT) OutputIs(s *state.Store, name string, expected string) {
	t.TB().Helper()
	t.Equal(s.GetOutput(name), expected)
}

// SummaryIs asserts the full summary via a single deep comparison.
func (t *StateT) SummaryIs(s *state.Store, passed, failed, skipped []string) {
	t.TB().Helper()
	p, f, sk := s.Summary()
	got := summary{p, f, sk}
	got.sort()
	want := summary{passed, failed, skipped}
	want.sort()
	t.DeepEqual(got, want)
}

func (t *StateT) MarkSkippedIs(s *state.Store, names, passed, failed, skipped []string) {
	t.TB().Helper()
	error := s.MarkSkipped(names)
	if error != nil {
		t.TB().Fatalf("MarkSkipped: %v", error)
	}
	t.SummaryIs(s, passed, failed, skipped)
}

type summary struct {
	passed, failed, skipped []string
}

func (s *summary) sort() {
	sort.Strings(s.passed)
	sort.Strings(s.failed)
	sort.Strings(s.skipped)
}

var StateTest = tape.Extend(func(base *tape.T) *StateT {
	return &StateT{T: base}
})

// errSource always fails Load.
type errSource struct{}

func (s errSource) Load() ([]statemachine.TransitionDef, error) {
	return nil, errors.New("source failed")
}

// errAdapter always fails Get.
type errAdapter struct {
	*adapters.Memory[state.TestState]
}

func (a errAdapter) Get(id string) (*state.TestState, error) {
	return nil, errors.New("adapter error")
}

// setErrAdapter fails Set.
type setErrAdapter struct {
	*adapters.Memory[state.TestState]
}

func (a setErrAdapter) Set(id string, st state.TestState) error {
	return errors.New("adapter set error")
}

func TestRunEventCreatesRunningState(t *testing.T) {
	StateTest(t, "state: run event creates running state", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s, model.Event{Action: "run", Test: "TestFoo"})
		t.StateIs(s, "TestFoo", state.StateRunning)
		t.End()
	})
}

func TestPassEventMarksTestPassed(t *testing.T) {
	StateTest(t, "state: pass event marks passed", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestFoo"},
			model.Event{Action: "pass", Test: "TestFoo"},
		)
		t.StateIs(s, "TestFoo", state.StatePassed)
		t.End()
	})
}

func TestFailEventMarksTestFailed(t *testing.T) {
	StateTest(t, "state: fail event marks failed", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestFoo"},
			model.Event{Action: "fail", Test: "TestFoo"},
		)
		t.StateIs(s, "TestFoo", state.StateFailed)
		t.End()
	})
}

func TestSkipEventMarksTestSkipped(t *testing.T) {
	StateTest(t, "state: skip event marks skipped", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestFoo"},
			model.Event{Action: "skip", Test: "TestFoo"},
		)
		t.StateIs(s, "TestFoo", state.StateSkipped)
		t.End()
	})
}

func TestOutputAppendedToLogs(t *testing.T) {
	StateTest(t, "state: output lines appended", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "output", Test: "TestFoo", Output: "line1\n"},
			model.Event{Action: "output", Test: "TestFoo", Output: "line2\n"},
		)
		t.OutputIs(s, "TestFoo", "line1\nline2\n")
		t.End()
	})
}

func TestPackageEventNoTestID(t *testing.T) {
	StateTest(t, "state: event without test id is a no-op", func(t *StateT) {
		s := t.NewStore()
		_, error := t.ApplyState(s, model.Event{Action: "pass", Package: "mypkg"})
		t.NotOk(error)
		t.End()
	})
}

func TestInvalidActionError(t *testing.T) {
	StateTest(t, "state: unknown action returns error", func(t *StateT) {
		s := t.NewStore()
		_, error := t.ApplyState(s, model.Event{Test: "TestFoo", Action: "invalid"})
		t.Equal(error.Error(), "unknown action: invalid")
		t.End()
	})
}

func TestRunTwice(t *testing.T) {
	StateTest(t, "state: running twice keeps running state", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestFoo"},
			model.Event{Action: "run", Test: "TestFoo"},
		)
		t.StateIs(s, "TestFoo", state.StateRunning)
		t.End()
	})
}

func TestGetNonExistent(t *testing.T) {
	StateTest(t, "state: Get returns error for unknown test", func(t *StateT) {
		s := t.NewStore()
		error := t.GetError(s, "nonexistent")
		t.Equal(error.Error(), "state not found: nonexistent")
		t.End()
	})
}

func TestApplyInvalidTransitionReturnsCurrentState(t *testing.T) {
	StateTest(t, "state: invalid transition returns current state", func(t *StateT) {
		s := t.NewStore()
		st, error := t.ApplyState(s, model.Event{Action: "pass", Test: "TestA"})
		t.NotOk(error)
		if st != state.StateIdle {
			t.TB().Fatalf("want StateIdle, got %v", st)
		}
		t.End()
	})
}

func TestSummaryKeepsSubtestName(t *testing.T) {
	StateTest(t, "state: summary keeps subtest name", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestOnlyRuns/tape:_Only_delegates_to_Test"},
			model.Event{Action: "pass", Test: "TestOnlyRuns/tape:_Only_delegates_to_Test"},
		)
		t.SummaryIs(s, []string{"TestOnlyRuns/tape:_Only_delegates_to_Test"}, nil, nil)
		t.End()
	})
}

func TestSummaryPassedFailedSkipped(t *testing.T) {
	StateTest(t, "state: summary groups passed failed and skipped", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestA"},
			model.Event{Action: "pass", Test: "TestA"},
			model.Event{Action: "output", Test: "TestA", Output: "ok"},
			model.Event{Action: "run", Test: "TestB"},
			model.Event{Action: "fail", Test: "TestB"},
			model.Event{Action: "output", Test: "TestB", Output: "fail"},
			model.Event{Action: "run", Test: "TestC"},
			model.Event{Action: "skip", Test: "TestC"},
			model.Event{Action: "output", Test: "TestC", Output: "skip"},
		)
		t.SummaryIs(s, []string{"TestA"}, []string{"TestB"}, []string{"TestC"})
		t.End()
	})
}

func TestParseTestStateUnknown(t *testing.T) {
	StateTest(t, "state: parse unknown state errors", func(t *StateT) {
		error := t.ParseState("unknown")
		t.Equal(error.Error(), "unknown state: unknown")
		t.End()
	})
}

func TestParseTestEventUnknown(t *testing.T) {
	StateTest(t, "state: parse unknown event errors", func(t *StateT) {
		error := t.ParseEvent("unknown")
		t.Equal(error.Error(), "unknown event: unknown")
		t.End()
	})
}

func TestSummaryWithOutputOnly(t *testing.T) {
	StateTest(t, "state: output-only tests ignored in summary", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s, model.Event{Action: "output", Test: "orphan", Output: "some output"})
		t.SummaryIs(s, nil, nil, nil)
		t.End()
	})
}

func TestNew(t *testing.T) {
	StateTest(t, "state: New returns non-nil store", func(t *StateT) {
		s := t.NewStore()
		t.Ok(s != nil)
		t.End()
	})
}

func TestNewFromSourceError(t *testing.T) {
	StateTest(t, "state: NewFromSource propagates source error", func(t *StateT) {
		_, error := t.NewFromSource(errSource{})
		t.Equal(error.Error(), "source failed")
		t.End()
	})
}

func TestGetAdapterError(t *testing.T) {
	StateTest(t, "state: Get propagates adapter error", func(t *StateT) {
		s := t.NewStore()
		mem := adapters.NewMemory[state.TestState]()
		s.SetAdapter(errAdapter{mem})
		error := t.GetError(s, "TestFoo")
		t.Equal(error.Error(), "adapter error")
		t.End()
	})
}

func TestSummaryAdapterError(t *testing.T) {
	StateTest(t, "state: summary skips adapter errors without panic", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s, model.Event{Action: "output", Test: "TestFoo", Output: "ok"})
		mem := adapters.NewMemory[state.TestState]()
		s.SetAdapter(errAdapter{mem})
		t.SummaryIs(s, nil, nil, nil)
		t.End()
	})
}

func TestMarkSkippedAddsUnseen(t *testing.T) {
	StateTest(t, "state: MarkSkipped adds unseen tests", func(t *StateT) {
		s := t.NewStore()
		t.MarkSkippedIs(s, []string{"scope: foo", "scope: bar"}, nil, nil, []string{"scope: foo", "scope: bar"})
		t.End()
	})
}

func TestMarkSkippedDoesNotOverwritePassed(t *testing.T) {
	StateTest(t, "state: MarkSkipped keeps passed tests", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "scope: foo"},
			model.Event{Action: "pass", Test: "scope: foo"},
		)
		t.MarkSkippedIs(s, []string{"scope: foo"}, []string{"scope: foo"}, nil, nil)
		t.End()
	})
}

func TestSummaryIgnoresRunningState(t *testing.T) {
	StateTest(t, "state: summary ignores running tests", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s, model.Event{Action: "run", Test: "TestFoo"})
		t.SummaryIs(s, nil, nil, nil)
		t.End()
	})
}

func TestMarkSkippedSetError(t *testing.T) {
	StateTest(t, "state: MarkSkipped propagates Set error", func(t *StateT) {
		s := t.NewStore()
		mem := adapters.NewMemory[state.TestState]()
		s.SetAdapter(setErrAdapter{mem})
		error := t.MarkSkippedErr(s, []string{"scope: foo"})
		t.Equal(error.Error(), "state.MarkSkipped: Set failed: adapter set error")
		t.End()
	})
}

func TestBuildFailedCountIsOne(t *testing.T) {
	StateTest(t, "state: build-failed package increments BuildFailedCount to 1", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "start", Package: "mypkg"},
			model.Event{Action: "output", Package: "mypkg", Output: "FAIL\tmypkg [build failed]\n"},
			model.Event{Action: "fail", Package: "mypkg"},
		)
		t.Equal(s.BuildFailedCount(), 1)
		t.End()
	})
}

func TestBuildFailedPackageNotInPassed(t *testing.T) {
	StateTest(t, "state: build-failed package produces no passed tests", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "output", Package: "mypkg", Output: "FAIL\tmypkg [build failed]\n"},
			model.Event{Action: "fail", Package: "mypkg"},
		)
		passed, _, _ := s.Summary()
		t.Equal(len(passed), 0)
		t.End()
	})
}

func TestBuildFailedPackageNotInFailed(t *testing.T) {
	StateTest(t, "state: build-failed package produces no failed tests", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "output", Package: "mypkg", Output: "FAIL\tmypkg [build failed]\n"},
			model.Event{Action: "fail", Package: "mypkg"},
		)
		_, failed, _ := s.Summary()
		t.Equal(len(failed), 0)
		t.End()
	})
}

func TestBuildFailedCountZeroByDefault(t *testing.T) {
	StateTest(t, "state: BuildFailedCount is zero with no build failures", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestFoo/scope: x"},
			model.Event{Action: "pass", Test: "TestFoo/scope: x"},
		)
		t.Equal(s.BuildFailedCount(), 0)
		t.End()
	})
}

func TestBuildFailedCountMultiplePackages(t *testing.T) {
	StateTest(t, "state: BuildFailedCount accumulates across packages", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "output", Package: "pkgA", Output: "FAIL\tpkgA [build failed]\n"},
			model.Event{Action: "fail", Package: "pkgA"},
			model.Event{Action: "output", Package: "pkgB", Output: "FAIL\tpkgB [build failed]\n"},
			model.Event{Action: "fail", Package: "pkgB"},
		)
		t.Equal(s.BuildFailedCount(), 2)
		t.End()
	})
}

func TestNormalPackageFailNotCountedAsBuildFailed(t *testing.T) {
	StateTest(t, "state: normal package fail does not increment BuildFailedCount", func(t *StateT) {
		s := t.NewStore()
		t.Apply(s,
			model.Event{Action: "run", Test: "TestFoo/scope: x"},
			model.Event{Action: "fail", Test: "TestFoo/scope: x"},
			model.Event{Action: "fail", Package: "mypkg"},
		)
		t.Equal(s.BuildFailedCount(), 0)
		t.End()
	})
}
