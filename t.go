package tape

import (
	"testing"
)

// T wraps *testing.T with a thinner, more expressive API.
type T struct {
	t     *testing.T
	ended bool
}

func newT(t *testing.T) *T {
	return &T{t: t}
}

func (tt *T) Equal(want, got any) {
	tt.t.Helper()
	hit(tt.t)
	if !deepEqual(want, got) {
		tt.t.Errorf("\nwant: %#v\n got: %#v", want, got)
	}
}

func (tt *T) Ok(v bool) {
	tt.t.Helper()
	hit(tt.t)
	if !v {
		tt.t.Errorf("expected true, got false")
	}
}

func (tt *T) NotOk(v bool) {
	tt.t.Helper()
	hit(tt.t)
	if v {
		tt.t.Errorf("expected false, got true")
	}
}

func (tt *T) DeepEqual(want, got any) {
	tt.t.Helper()
	hit(tt.t)
	if !deepEqual(want, got) {
		tt.t.Errorf("\nwant: %#v\n got: %#v", want, got)
	}
}

func (tt *T) Contains(s, sub string) {
	tt.t.Helper()
	hit(tt.t)
	if !contains(s, sub) {
		tt.t.Errorf("%q does not contain %q", s, sub)
	}
}

func (tt *T) Error(err error) {
	tt.t.Helper()
	hit(tt.t)
	if err == nil {
		tt.t.Fatal("expected an error, got nil")
	}
}

func (tt *T) NoError(err error) {
	tt.t.Helper()
	hit(tt.t)
	if err != nil {
		tt.t.Fatalf("unexpected error: %v", err)
	}
}

func (tt *T) Match(s, pattern string) {
	tt.t.Helper()
	hit(tt.t)
	if !match(s, pattern) {
		tt.t.Errorf("%q does not match %q", s, pattern)
	}
}

// End marks the test as having completed its single assertion.
func (tt *T) End() {
	tt.t.Helper()
	tt.ended = true
}
