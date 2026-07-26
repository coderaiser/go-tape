package assert

import (
	"errors"
	"fmt"
	"testing"
)

// mockT implements TB without calling runtime.Goexit().
// Used to test failure paths without failing the outer test.
type mockT struct {
	failed  bool
	message string
}

func (m *mockT) Helper()                      {}
func (m *mockT) Errorf(f string, args ...any) { m.failed = true; m.message = fmt.Sprintf(f, args...) }
func (m *mockT) Fatalf(f string, args ...any) { m.failed = true; m.message = fmt.Sprintf(f, args...) }
func (m *mockT) Fatal(args ...any)            { m.failed = true }

// -- One / hit --

func TestOneResetsCount(t *testing.T) {
	One(t)
	hit(t)
	One(t) // reset
	hit(t) // should not fail
}

// hit failure path is tested indirectly via tape_test.go which
// exercises the assertion count guard with TAPE_CHECK_ASSERTIONS_COUNT.

// -- Equal --

func TestEqualMatch(t *testing.T) {
	m := &mockT{}
	Equal(m, 42, 42)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestEqualMismatch(t *testing.T) {
	m := &mockT{}
	Equal(m, 1, 2)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- NoError --

func TestNoErrorNil(t *testing.T) {
	m := &mockT{}
	NoError(m, nil)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNoErrorNonNil(t *testing.T) {
	m := &mockT{}
	NoError(m, errors.New("oops"))
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- Error --

func TestErrorNonNil(t *testing.T) {
	m := &mockT{}
	Error(m, errors.New("oops"))
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestErrorNil(t *testing.T) {
	m := &mockT{}
	Error(m, nil)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- Ok --

func TestOkTrue(t *testing.T) {
	m := &mockT{}
	Ok(m, true)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestOkFalse(t *testing.T) {
	m := &mockT{}
	Ok(m, false)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- NotOk --

func TestNotOkFalse(t *testing.T) {
	m := &mockT{}
	NotOk(m, false)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNotOkTrue(t *testing.T) {
	m := &mockT{}
	NotOk(m, true)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- Contains --

func TestContainsMatch(t *testing.T) {
	m := &mockT{}
	Contains(m, "hello world", "world")
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestContainsNoMatch(t *testing.T) {
	m := &mockT{}
	Contains(m, "hello", "xyz")
	if !m.failed {
		t.Fatal("expected fail")
	}
}
