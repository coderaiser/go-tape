package statemachine_test

import (
	"errors"
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/statemachine"
	"github.com/coderaiser/go-tape/internal/statemachine/adapters"
)

type testState int

const (
	stateIdle testState = iota
	stateRunning
	stateDone
)

type testEvent int

const (
	eventRun testEvent = iota
	eventFinish
	eventFail
)

func parseState(s string) (testState, error) {
	switch s {
	case "idle":
		return stateIdle, nil
	case "running":
		return stateRunning, nil
	case "done":
		return stateDone, nil
	default:
		return 0, errors.New("unknown state")
	}
}

func parseEvent(e string) (testEvent, error) {
	switch e {
	case "run":
		return eventRun, nil
	case "finish":
		return eventFinish, nil
	case "fail":
		return eventFail, nil
	default:
		return 0, errors.New("unknown event")
	}
}

type MachineT struct{ *tape.T }

func (t *MachineT) NewMachine(src statemachine.TransitionSource, adapter statemachine.Adapter[testState]) *statemachine.Machine[testState, testEvent] {
	t.TB().Helper()
	m, error := statemachine.New(src, parseState, parseEvent, adapter)
	if error != nil {
		t.TB().Fatalf("New: %v", error)
	}
	return m
}

func (t *MachineT) NewMachineError(src statemachine.TransitionSource, adapter statemachine.Adapter[testState]) error {
	t.TB().Helper()
	_, error := statemachine.New(src, parseState, parseEvent, adapter)
	return error
}

func (t *MachineT) Apply(m *statemachine.Machine[testState, testEvent], id string, e testEvent) (testState, error) {
	t.TB().Helper()
	return m.Apply(id, e, nil)
}

func (t *MachineT) ApplyResultIs(m *statemachine.Machine[testState, testEvent], id string, e testEvent, expected testState) {
	t.TB().Helper()
	next, error := m.Apply(id, e, nil)
	if error != nil {
		t.TB().Fatalf("Apply: %v", error)
	}
	t.Equal(next, expected)
}

func (t *MachineT) StoredIs(adapter statemachine.Adapter[testState], id string, expected testState) {
	t.TB().Helper()
	ptr, error := adapter.Get(id)
	if error != nil {
		t.TB().Fatalf("Get: %v", error)
	}
	if ptr == nil {
		t.TB().Fatal("expected stored state, got nil")
	}
	t.Equal(*ptr, expected)
}

func (t *MachineT) HookCalledIs(m *statemachine.Machine[testState, testEvent], from testState, e testEvent, wantCalled bool) {
	t.TB().Helper()
	called := false
	m.Hook(from, e, func(statemachine.Context[testState, testEvent]) error {
		called = true
		return nil
	})
	_, error := m.Apply("test-1", e, nil)
	if error != nil {
		t.TB().Fatalf("Apply: %v", error)
	}
	t.Equal(called, wantCalled)
}

var MachineTest = tape.Extend(func(base *tape.T) *MachineT {
	return &MachineT{T: base}
})

func TestNewValidMachine(t *testing.T) {
	MachineTest(t, "statemachine: New builds a machine", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
			{From: "running", Event: "finish", To: "done"},
		}}
		adapter := adapters.NewMemory[testState]()
		t.Ok(t.NewMachine(src, adapter) != nil)
		t.End()
	})
}

func TestApplyValidTransition(t *testing.T) {
	MachineTest(t, "statemachine: Apply follows a valid transition", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		}}
		adapter := adapters.NewMemory[testState]()
		m := t.NewMachine(src, adapter)
		t.ApplyResultIs(m, "test-1", eventRun, stateRunning)
		t.End()
	})
}

func TestApplyInvalidTransitionNonStrict(t *testing.T) {
	MachineTest(t, "statemachine: Apply errors on invalid transition", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		}}
		adapter := adapters.NewMemory[testState]()
		m := t.NewMachine(src, adapter)
		_, error := t.Apply(m, "test-1", eventFinish)
		t.Equal(error.Error(), "invalid transition: from 0 event 1")
		t.End()
	})
}

func TestHookCalledOnTransition(t *testing.T) {
	MachineTest(t, "statemachine: hook runs on valid transition", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		}}
		adapter := adapters.NewMemory[testState]()
		m := t.NewMachine(src, adapter)
		t.HookCalledIs(m, stateIdle, eventRun, true)
		t.End()
	})
}

func TestHookErrorReturned(t *testing.T) {
	MachineTest(t, "statemachine: Apply returns hook error", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		}}
		adapter := adapters.NewMemory[testState]()
		m := t.NewMachine(src, adapter)
		m.Hook(stateIdle, eventRun, func(statemachine.Context[testState, testEvent]) error {
			return errors.New("hook failed")
		})
		_, error := t.Apply(m, "test-1", eventRun)
		t.Equal(error.Error(), "hook failed")
		t.End()
	})
}

func TestApplyStoresState(t *testing.T) {
	MachineTest(t, "statemachine: Apply persists next state", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		}}
		adapter := adapters.NewMemory[testState]()
		m := t.NewMachine(src, adapter)
		if _, error := m.Apply("test-1", eventRun, nil); error != nil {
			t.TB().Fatalf("Apply: %v", error)
		}
		t.StoredIs(adapter, "test-1", stateRunning)
		t.End()
	})
}

func TestUnknownIdGetsNil(t *testing.T) {
	MachineTest(t, "statemachine: unknown id Get returns nil", func(t *MachineT) {
		adapter := adapters.NewMemory[testState]()
		ptr, error := adapter.Get("unknown")
		if error != nil {
			t.TB().Fatalf("Get: %v", error)
		}
		t.Ok(ptr == nil)
		t.End()
	})
}

func TestNewWithParseStateError(t *testing.T) {
	MachineTest(t, "statemachine: New errors on invalid From", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "invalid", Event: "run", To: "running"},
		}}
		error := t.NewMachineError(src, adapters.NewMemory[testState]())
		t.Equal(error.Error(), "invalid state \"invalid\": unknown state")
		t.End()
	})
}

func TestNewWithParseEventError(t *testing.T) {
	MachineTest(t, "statemachine: New errors on invalid Event", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "invalid", To: "running"},
		}}
		error := t.NewMachineError(src, adapters.NewMemory[testState]())
		t.Equal(error.Error(), "invalid event \"invalid\": unknown event")
		t.End()
	})
}

func TestNewWithParseToStateError(t *testing.T) {
	MachineTest(t, "statemachine: New errors on invalid To", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "invalid"},
		}}
		error := t.NewMachineError(src, adapters.NewMemory[testState]())
		t.Equal(error.Error(), "invalid state \"invalid\": unknown state")
		t.End()
	})
}

func TestValidatePasses(t *testing.T) {
	MachineTest(t, "statemachine: Validate passes for well-formed machine", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		}}
		m := t.NewMachine(src, adapters.NewMemory[testState]())
		t.NotOk(m.Validate())
		t.End()
	})
}

func TestHookNotCalledForUnknownTransition(t *testing.T) {
	MachineTest(t, "statemachine: hook not called on invalid transition", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		}}
		m := t.NewMachine(src, adapters.NewMemory[testState]())
		called := false
		m.Hook(stateRunning, eventFinish, func(statemachine.Context[testState, testEvent]) error {
			called = true
			return nil
		})
		_, error := t.Apply(m, "test-1", eventFinish)
		if error == nil {
			t.TB().Fatal("expected Apply error")
		}
		t.Equal(called, false)
		t.End()
	})
}

// error adapters

type errGetAdapter struct{}

func (a errGetAdapter) Get(id string) (*testState, error) { return nil, errors.New("get failed") }
func (a errGetAdapter) Set(id string, s testState) error  { return nil }

type errSetAdapter struct{}

func (a errSetAdapter) Get(id string) (*testState, error) { s := stateIdle; return &s, nil }
func (a errSetAdapter) Set(id string, s testState) error  { return errors.New("set failed") }

func TestApplyGetError(t *testing.T) {
	MachineTest(t, "statemachine: Apply propagates Get error", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{{From: "idle", Event: "run", To: "running"}}}
		m := t.NewMachine(src, errGetAdapter{})
		_, error := t.Apply(m, "x", eventRun)
		t.Equal(error.Error(), "adapter.Get: get failed")
		t.End()
	})
}

func TestApplySetError(t *testing.T) {
	MachineTest(t, "statemachine: Apply propagates Set error", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{{From: "idle", Event: "run", To: "running"}}}
		m := t.NewMachine(src, errSetAdapter{})
		_, error := t.Apply(m, "x", eventRun)
		t.Equal(error.Error(), "adapter.Set: set failed")
		t.End()
	})
}

func TestWithInitialUsedForUnknownId(t *testing.T) {
	MachineTest(t, "statemachine: WithInitial applies to unknown id", func(t *MachineT) {
		src := &statemachine.MemorySource{Defs: []statemachine.TransitionDef{{From: "idle", Event: "run", To: "running"}}}
		m := t.NewMachine(src, adapters.NewMemory[testState]())
		m.WithInitial(stateIdle)
		t.ApplyResultIs(m, "brand-new", eventRun, stateRunning)
		t.End()
	})
}

func TestValidateEmptyTransitions(t *testing.T) {
	MachineTest(t, "statemachine: Validate errors on empty transitions", func(t *MachineT) {
		m := statemachine.EmptyTransitionsMachine[testState]()
		t.Equal(m.Validate().Error(), "state 0 has no transitions out")
		t.End()
	})
}

// -- FileSource tests --

func TestFileSourceLoadsTransitions(t *testing.T) {
	MachineTest(t, "statemachine: FileSource loads TOML transitions", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/runner.toml"}
		defs, error := src.Load()
		if error != nil {
			t.TB().Fatalf("Load: %v", error)
		}
		t.Ok(len(defs) > 0)
		t.End()
	})
}

func TestFileSourceJSON(t *testing.T) {
	MachineTest(t, "statemachine: FileSource loads JSON transitions", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/runner.json"}
		defs, error := src.Load()
		if error != nil {
			t.TB().Fatalf("Load: %v", error)
		}
		t.Ok(len(defs) > 0)
		t.End()
	})
}

// errSource returns an error on Load.
type errSource struct{}

func (s errSource) Load() ([]statemachine.TransitionDef, error) {
	return nil, errors.New("source failed")
}

func TestNewWithSourceError(t *testing.T) {
	MachineTest(t, "statemachine: New propagates source error", func(t *MachineT) {
		adapter := adapters.NewMemory[testState]()
		error := t.NewMachineError(errSource{}, adapter)
		t.Equal(error.Error(), "source failed")
		t.End()
	})
}

func TestFileSourceUnsupportedExtension(t *testing.T) {
	MachineTest(t, "statemachine: FileSource errors on unsupported extension", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/runner.yaml"}
		_, error := src.Load()
		t.Equal(error.Error(), "FileSource: unsupported extension \".yaml\"")
		t.End()
	})
}

func TestFileSourceInvalidJSON(t *testing.T) {
	MachineTest(t, "statemachine: FileSource errors on invalid JSON", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/invalid.json"}
		_, error := src.Load()
		t.Ok(error)
		t.End()
	})
}

func TestFileSourceMissingJSON(t *testing.T) {
	MachineTest(t, "statemachine: FileSource errors on missing JSON file", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/nonexistent.json"}
		_, error := src.Load()
		t.Ok(error)
		t.End()
	})
}

func TestFileSourceMissingTOML(t *testing.T) {
	MachineTest(t, "statemachine: FileSource errors on missing TOML file", func(t *MachineT) {
		src := statemachine.FileSource{Path: "nonexistent.toml"}
		_, error := src.Load()
		t.Ok(error)
		t.End()
	})
}
