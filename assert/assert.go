package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// TB is the subset of testing.TB used by assert functions.
// Allows mockT injection in tests.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Fatal(args ...any)
}

var (
	mu    sync.Mutex
	count = make(map[*testing.T]int)
)

func One(t *testing.T) {
	t.Helper()
	mu.Lock()
	count[t] = 0
	mu.Unlock()
}

func hit(t *testing.T) {
	t.Helper()
	mu.Lock()
	defer mu.Unlock()
	count[t]++
	if count[t] <= 1 {
		return
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		t.Fatalf(
			"too many assertions: got %d, expected 1\nat %s:%d",
			count[t], file, line,
		)
	}
	t.Fatalf("too many assertions: got %d, expected 1", count[t])
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
