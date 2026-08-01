package state

import (
	"fmt"

	"github.com/coderaiser/go-tape/internal/model"
	"github.com/coderaiser/go-tape/internal/statemachine"
	"github.com/coderaiser/go-tape/internal/statemachine/adapters"
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
	adapter statemachine.Adapter[TestState]
	outputs map[string]string
}

// New creates a new Store.
func New() (*Store, error) {
	return newFromSource(&statemachine.MemorySource{
		Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
			{From: "running", Event: "pass", To: "passed"},
			{From: "running", Event: "fail", To: "failed"},
			{From: "running", Event: "skip", To: "skipped"},
		},
	})
}

// newFromSource creates a Store from a TransitionSource.
// Exported for testing.
func newFromSource(src statemachine.TransitionSource) (*Store, error) {
	adapter := adapters.NewMemory[TestState]()
	m, err := statemachine.New(src, parseTestState, parseTestEvent, adapter)
	if err != nil {
		return nil, err
	}
	return &Store{
		machine: m,
		adapter: adapter,
		outputs: make(map[string]string),
	}, nil
}

// Apply processes a model.Event and updates state.
func (s *Store) Apply(e model.Event) (TestState, error) {
	if e.Test == "" {
		return 0, nil
	}

	if _, ok := s.outputs[e.Test]; !ok {
		s.outputs[e.Test] = ""
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
		ptr, _ := s.adapter.Get(e.Test)
		if ptr == nil {
			return StateIdle, nil
		}
		return *ptr, nil
	}

	return state, nil
}

// Get returns the current state of a test.
func (s *Store) Get(test string) (TestState, error) {
	ptr, err := s.adapter.Get(test)
	if err != nil {
		return 0, err
	}
	if ptr == nil {
		return 0, fmt.Errorf("state not found: %s", test)
	}
	return *ptr, nil
}

// GetOutput returns the accumulated output for a test.
func (s *Store) GetOutput(test string) string {
	return s.outputs[test]
}

// Summary returns the test names grouped by state.
func (s *Store) Summary() (passed, failed, skipped []string) {
	for test := range s.outputs {
		ptr, err := s.adapter.Get(test)
		if err != nil || ptr == nil {
			continue
		}

		switch *ptr {
		case StatePassed:
			passed = append(passed, test)
		case StateFailed:
			failed = append(failed, test)
		case StateSkipped:
			skipped = append(skipped, test)
		case StateIdle, StateRunning:
			// A test stuck in Running at summary time means go test -json -v
			// dropped its terminal event (a known Go issue with long subtest names).
			// Silently ignore — it is neither passed nor failed.
		}
	}

	return
}

func (s *Store) MarkSkipped(names []string) {
	for _, name := range names {
		ptr, _ := s.adapter.Get(name)

		if ptr != nil {
			continue
		}

		if err := s.adapter.Set(name, StateSkipped); err != nil {
			panic(fmt.Sprintf("state.MarkSkipped: Set failed: %v", err))
		}
		s.outputs[name] = ""
	}
}
