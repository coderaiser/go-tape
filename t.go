package tape

import (
	"reflect"
	"testing"
)

// T wraps *testing.T with supertape-compatible assertions.
type T struct {
	t     *testing.T
	ended bool
}

func newT(t *testing.T) *T { return &T{t: t} }

// TB returns the underlying *testing.T for use with helpers
// that require it directly (e.g. t.TB().TempDir(), t.TB().Fatal(...)).
func (tt *T) TB() *testing.T {
	return tt.t
}

// Setenv sets an environment variable for the duration of the test.
// It is automatically restored when the test ends.
func (tt *T) Setenv(key, value string) {
	tt.t.Helper()
	tt.t.Setenv(key, value)
}

// Equal asserts result == expected using strict equality.
// For primitives and pointers only.
// Use DeepEqual for structs, slices, and maps.
func (tt *T) Equal(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	if !isPrimitive(result) {
		tt.t.Fatalf("Equal: use DeepEqual for %T — Equal is for primitives and pointers only", result)
		return
	}
	if result != expected {
		tt.t.Errorf("\nresult:   %#v\nexpected: %#v", result, expected)
	}
}

// NotEqual asserts result != expected.
// For primitives and pointers only.
func (tt *T) NotEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	if !isPrimitive(result) {
		tt.t.Fatalf("NotEqual: use NotDeepEqual for %T — NotEqual is for primitives and pointers only", result)
		return
	}
	if result == expected {
		tt.t.Errorf("expected values to differ, both are: %#v", result)
	}
}

// DeepEqual asserts deep equality using reflect.DeepEqual.
// Use for structs, slices, and maps.
func (tt *T) DeepEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	if !reflect.DeepEqual(result, expected) {
		tt.t.Errorf("\nresult:   %#v\nexpected: %#v", result, expected)
	}
}

// NotDeepEqual asserts values are not deeply equal.
func (tt *T) NotDeepEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	if reflect.DeepEqual(result, expected) {
		tt.t.Errorf("expected values to differ, both deep-equal: %#v", result)
	}
}

// Ok asserts result is truthy.
func (tt *T) Ok(result any) {
	tt.t.Helper()
	hit(tt.t)
	if !truthy(result) {
		tt.t.Errorf("expected truthy, got: %#v", result)
	}
}

// NotOk asserts result is falsy.
func (tt *T) NotOk(result any) {
	tt.t.Helper()
	hit(tt.t)
	if truthy(result) {
		tt.t.Errorf("expected falsy, got: %#v", result)
	}
}

// Match asserts result matches pattern.
// pattern may be a string or *regexp.Regexp.
func (tt *T) Match(result string, pattern any) {
	tt.t.Helper()
	hit(tt.t)
	re, err := toRegexp(pattern)
	if err != nil {
		tt.t.Errorf("Match: invalid pattern: %v", err)
		return
	}
	if !re.MatchString(result) {
		tt.t.Errorf("%q does not match %q", result, re)
	}
}

// NotMatch asserts result does not match pattern.
// pattern may be a string or *regexp.Regexp.
func (tt *T) NotMatch(result string, pattern any) {
	tt.t.Helper()
	hit(tt.t)
	re, err := toRegexp(pattern)
	if err != nil {
		tt.t.Errorf("NotMatch: invalid pattern: %v", err)
		return
	}
	if re.MatchString(result) {
		tt.t.Errorf("%q should not match %q", result, re)
	}
}

// Pass generates an unconditional passing assertion.
func (tt *T) Pass(message string) {
	tt.t.Helper()
	hit(tt.t)
	tt.t.Log("pass:", message)
}

// Fail generates an unconditional failing assertion.
func (tt *T) Fail(message string) {
	tt.t.Helper()
	hit(tt.t)
	tt.t.Errorf("fail: %s", message)
}

// Comment prints a TAP comment without counting as an assertion.
func (tt *T) Comment(message string) {
	tt.t.Helper()
	tt.t.Log("#", message)
}

// Error asserts err is non-nil.
func (tt *T) Error(err error) {
	tt.t.Helper()
	hit(tt.t)
	if err == nil {
		tt.t.Fatal("expected an error, got nil")
	}
}

// NoError asserts err is nil.
func (tt *T) NoError(err error) {
	tt.t.Helper()
	hit(tt.t)
	if err != nil {
		tt.t.Fatalf("unexpected error: %v", err)
	}
}

// End marks the test as intentionally complete.
// Required when TAPE_CHECK_END is enabled (default: on).
func (tt *T) End() {
	tt.ended = true
}
