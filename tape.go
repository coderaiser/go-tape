package tape

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/coderaiser/go-tape/internal/config"
	"github.com/coderaiser/go-tape/internal/scope"
)

var (
	mu    sync.Mutex
	count = make(map[*testing.T]int)
)

// assertOne resets the assertion counter for a test and removes it on cleanup.
func assertOne(t *testing.T) {
	t.Helper()
	mu.Lock()
	count[t] = 0
	mu.Unlock()
	// cleanup prevents memory leak — removes entry after test finishes
	t.Cleanup(func() {
		mu.Lock()
		delete(count, t)
		mu.Unlock()
	})
}

// hit increments the assertion counter and fails if it exceeds one.
func hit(t *testing.T) {
	t.Helper()
	mu.Lock()
	count[t]++
	c := count[t]
	mu.Unlock()
	if c <= 1 {
		return
	}
	if !config.CheckAssertionsCount() {
		return
	}
	_, file, line, _ := runtime.Caller(2)
	t.Fatalf("too many assertions: got %d, expected 1\nat %s:%d", c, file, line)
}

// Test runs a subtest with guards: scope check, assertion count, timeout, End check.
func Test(t *testing.T, name string, fn func(t *T)) {
	t.Helper()

	// guard 1: scope check
	if config.CheckScopes() && !scope.Valid(name) {
		t.Fatalf("tape: invalid scope name: %q — expected 'scope: message'", name)
	}

	t.Run(name, func(t *testing.T) {
		t.Helper()

		// guard 2: assertion count
		assertOne(t)

		tt := newT(t)

		done := make(chan struct{})
		go func() {
			defer close(done)
			fn(tt)
		}()

		// guard 3: timeout
		select {
		case <-done:
		case <-time.After(config.Timeout()):
			t.Fatalf("tape: test timed out after %s", config.Timeout())
		}

		// guard 4: t.End() check
		if config.CheckEnd() && !tt.ended {
			t.Fatal("tape: t.End() not called")
		}
	})
}

// Only runs a single test.
func Only(t *testing.T, name string, fn func(t *T)) {
	Test(t, name, fn)
}

// Skip marks a test as skipped.
func Skip(_ *testing.T, _ string, _ func(t *T)) {}
