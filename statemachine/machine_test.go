package statemachine

import (
	"errors"
	"testing"

	"github.com/coderaiser/go-tape/statemachine/adapters"
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

func TestNewValidMachine(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
			{From: "running", Event: "finish", To: "done"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, err := New(src, parseState, parseEvent, adapter, false)
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("expected machine")
	}
}

func TestApplyValidTransition(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	next, err := m.Apply("test-1", eventRun, nil)
	if err != nil {
		t.Fatal(err)
	}
	if next != stateRunning {
		t.Errorf("want %v, got %v", stateRunning, next)
	}
}

func TestApplyInvalidTransitionNonStrict(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	_, err := m.Apply("test-1", eventFinish, nil)
	if err == nil {
		t.Fatal("expected error for invalid transition")
	}
}

func TestApplyInvalidTransitionStrictPanics(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, true)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for strict invalid transition")
		}
	}()

	m.Apply("test-1", eventFinish, nil)
}

func TestHookCalledOnTransition(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	called := false
	m.Hook(stateIdle, eventRun, func(ctx Context[testState, testEvent]) error {
		called = true
		return nil
	})

	m.Apply("test-1", eventRun, nil)
	if !called {
		t.Error("expected hook to be called")
	}
}

func TestHookErrorReturned(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	m.Hook(stateIdle, eventRun, func(ctx Context[testState, testEvent]) error {
		return errors.New("hook failed")
	})

	_, err := m.Apply("test-1", eventRun, nil)
	if err == nil || err.Error() != "hook failed" {
		t.Fatalf("expected hook error, got %v", err)
	}
}

func TestApplyStoresState(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	m.Apply("test-1", eventRun, nil)
	stored, _ := adapter.Get("test-1")
	if stored != stateRunning {
		t.Errorf("want %v, got %v", stateRunning, stored)
	}
}

func TestUnknownIdGetsZeroState(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	New(src, parseState, parseEvent, adapter, false)

	_, err := adapter.Get("unknown")
	if err == nil {
		t.Error("expected error for unknown id")
	}
}

func TestNewWithParseStateError(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "invalid", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	_, err := New(src, parseState, parseEvent, adapter, false)
	if err == nil {
		t.Fatal("expected error for invalid state")
	}
}

func TestNewWithParseEventError(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "invalid", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	_, err := New(src, parseState, parseEvent, adapter, false)
	if err == nil {
		t.Fatal("expected error for invalid event")
	}
}

func TestNewWithParseToStateError(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "invalid"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	_, err := New(src, parseState, parseEvent, adapter, false)
	if err == nil {
		t.Fatal("expected error for invalid to state")
	}
}

func TestMemoryAdapterGetSet(t *testing.T) {
	adapter := adapters.NewMemory[testState]()

	err := adapter.Set("key1", stateRunning)
	if err != nil {
		t.Fatal(err)
	}

	got, err := adapter.Get("key1")
	if err != nil {
		t.Fatal(err)
	}
	if got != stateRunning {
		t.Errorf("want %v, got %v", stateRunning, got)
	}
}

func TestMemoryAdapterGetNotFound(t *testing.T) {
	adapter := adapters.NewMemory[testState]()
	_, err := adapter.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestValidatePasses(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestHookNotCalledForUnknownTransition(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	called := false
	m.Hook(stateRunning, eventFinish, func(ctx Context[testState, testEvent]) error {
		called = true
		return nil
	})

	m.Apply("test-1", eventFinish, nil)
	if called {
		t.Error("hook should not be called for invalid transition")
	}
}
