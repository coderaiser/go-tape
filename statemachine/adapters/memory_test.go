package adapters

import (
	"testing"
)

func TestNewMemory(t *testing.T) {
	m := NewMemory[string]()
	if m == nil {
		t.Fatal("expected non-nil adapter")
	}
}

func TestMemorySetAndGet(t *testing.T) {
	m := NewMemory[string]()
	err := m.Set("key1", "value1")
	if err != nil {
		t.Fatal(err)
	}
	got, err := m.Get("key1")
	if err != nil {
		t.Fatal(err)
	}
	if got != "value1" {
		t.Errorf("want value1, got %s", got)
	}
}

func TestMemoryGetNotFound(t *testing.T) {
	m := NewMemory[string]()
	_, err := m.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}
