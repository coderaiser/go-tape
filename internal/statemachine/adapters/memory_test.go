package adapters_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/statemachine/adapters"
)

func TestNewMemory(t *testing.T) {
	tape.Test(t, "adapters: NewMemory returns non-nil adapter", func(t *tape.T) {
		m := adapters.NewMemory[string]()
		t.Ok(m != nil)
		t.End()
	})
}

func TestMemorySetAndGet(t *testing.T) {
	tape.Test(t, "adapters: Set and Get round-trips a value", func(t *tape.T) {
		m := adapters.NewMemory[string]()
		m.Set("k", "v")
		ptr, err := m.Get("k")
		t.Ok(err == nil && ptr != nil && *ptr == "v")
		t.End()
	})
}

func TestMemoryGetNotFound(t *testing.T) {
	tape.Test(t, "adapters: Get returns nil for unknown key", func(t *tape.T) {
		m := adapters.NewMemory[string]()
		ptr, err := m.Get("nonexistent")
		t.Ok(err == nil && ptr == nil)
		t.End()
	})
}
