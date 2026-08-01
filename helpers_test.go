package tape

import "testing"

// TestHitFatalfWhenCountExceeded verifies hit() Fatalfs on second assertion when
// the assertion count check is enabled (the default). The Fatalf fires on the
// inner subtest, so the outer test observes the failure via inner.Failed().
func TestHitFatalfWhenCountExceeded(t *testing.T) {
	failed := false
	t.Run("inner", func(inner *testing.T) {
		defer func() { failed = inner.Failed() }()
		assertOne(inner)
		hit(inner) // count = 1 — ok
		hit(inner) // count = 2 — Fatalf fires on inner, not t
	})
	if !failed {
		t.Fatal("expected second hit to fail the test")
	}
}