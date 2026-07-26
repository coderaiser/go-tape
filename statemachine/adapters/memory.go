package adapters

import (
	"fmt"
	"sync"
)

// Memory is an in-memory adapter for statemachine.Adapter.
type Memory[S any] struct {
	mu     sync.RWMutex
	states map[string]S
}

func NewMemory[S any]() *Memory[S] {
	return &Memory[S]{
		states: make(map[string]S),
	}
}

func (m *Memory[S]) Get(id string) (S, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.states[id]
	if !ok {
		var zero S
		return zero, fmt.Errorf("state not found: %s", id)
	}

	return s, nil
}

func (m *Memory[S]) Set(id string, state S) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[id] = state
	return nil
}
