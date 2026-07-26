package assert

import (
	"fmt"
	"reflect"
	"strings"
)

// TB is the subset of testing.TB used by assert functions.
// Allows mockT injection in tests.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Fatal(args ...any)
}

func Equal(t TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %#v\nwant: %#v", got, want)
	}
}

func NoError(t TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Error(t TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal(fmt.Errorf("expected an error, got nil"))
	}
}

func Ok(t TB, v bool) {
	t.Helper()
	if !v {
		t.Errorf("expected true, got false")
	}
}

func NotOk(t TB, v bool) {
	t.Helper()
	if v {
		t.Errorf("expected false, got true")
	}
}

func Contains(t TB, s, sub string) {
	t.Helper()
	if !strings.Contains(s, sub) {
		t.Errorf("%q does not contain %q", s, sub)
	}
}
