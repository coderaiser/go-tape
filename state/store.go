package state

import (
	"fmt"
	"strings"

	"github.com/coderaiser/go-tape/model"
	"github.com/coderaiser/go-tape/statemachine"
	"github.com/coderaiser/go-tape/statemachine/adapters"
)

// TestState represents the state of a test.
type TestState int

const (
	StateIdle TestState = iota
	StateRunning
	StatePassed
	StateFailed
	StateSkipped
)

// TestEvent represents events that trigger state transitions.
type TestEvent int

const (
	EventRun TestEvent = iota
	EventPass
	EventFail
	EventSkip
)

func parseTestState(s string) (TestState, error) {
	switch s {
	case "idle":
		return StateIdle, nil
	case "running":
		return StateRunning, nil
	case "passed":
		return StatePassed, nil
	case "failed":
		return StateFailed, nil
	case "skipped":
		return StateSkipped, nil
	default:
		return 0, fmt.Errorf("unknown state: %s", s)
	}
}

func parseTestEvent(e string) (TestEvent, error) {
	switch e {
	case "run":
		return EventRun, nil
	case "pass":
		return EventPass, nil
	case "fail":
		return EventFail, nil
	case "skip":
		return EventSkip, nil
	default:
		return 0, fmt.Errorf("unknown event: %s", e)
	}
}

// Store manages test state using the statemachine.
type Store struct {
	machine *statemachine.Machine[TestState, TestEvent]
	adapter *adapters.Memory[TestState]
	outputs map[string]string
}

// New creates a new Store.
func New() *Store {
	adapter := adapters.NewMemory[TestState]()
	src := &statemachine.MemorySource{
		Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
			{From: "running", Event: "pass", To: "passed"},
			{From: "running", Event: "fail", To: "failed"},
			{From: "running", Event: "skip", To: "skipped"},
		},
	}
	m, err := statemachine.New(src, parseTestState, parseTestEvent, adapter, false)
	if err != nil {
		panic(err)
	}
	return &Store{
		machine: m,
		adapter: adapter,
		outputs: make(map[string]string),
	}
}

// Apply processes a model.Event and updates state.
func (s *Store) Apply(e model.Event) (TestState, error) {
	if e.Test == "" {
		return 0, nil
	}

	if e.Action == "output" {
		s.outputs[e.Test] += e.Output
		return 0, nil
	}

	var ev TestEvent
	switch e.Action {
	case "run":
		ev = EventRun
	case "pass":
		ev = EventPass
	case "fail":
		ev = EventFail
	case "skip":
		ev = EventSkip
	default:
		return 0, fmt.Errorf("unknown action: %s", e.Action)
	}

	state, err := s.machine.Apply(e.Test, ev, nil)
	if err != nil {
		// For invalid transitions (e.g., run on running), just return current state
		current, _ := s.adapter.Get(e.Test)
		return current, nil
	}

	return state, nil
}

// Get returns the current state of a test.
func (s *Store) Get(test string) (TestState, error) {
	return s.adapter.Get(test)
}

// GetOutput returns the accumulated output for a test.
func (s *Store) GetOutput(test string) string {
	return s.outputs[test]
}

// Summary returns the test names grouped by state.
func (s *Store) Summary() (passed, failed, skipped []string) {
	allStates := map[string]struct{}{}
	for test := range s.outputs {
		if idx := strings.LastIndex(test, "/"); idx >= 0 {
			allStates[test[:idx]] = struct{}{}
		} else {
			allStates[test] = struct{}{}
		}
	}
	for test := range allStates {
		st, err := s.adapter.Get(test)
		if err != nil {
			continue
		}
		switch st {
		case StatePassed:
			passed = append(passed, test)
		case StateFailed:
			failed = append(failed, test)
		case StateSkipped:
			skipped = append(skipped, test)
		}
	}
	return
}

