package timeout

import (
	"testing"
	"time"
)

func runWithTimeout(t *testing.T, timeout time.Duration, fn func()) {
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatalf("tape: test timed out after %s", timeout)
	}
}
