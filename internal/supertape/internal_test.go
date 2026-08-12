package supertape

import "testing"

// TestItoa covers itoa for zero, positive and negative values.
func TestItoa(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{
		{0, "0"},
		{1, "1"},
		{12345, "12345"},
		{-7, "-7"},
		{-12345, "-12345"},
	}
	for _, c := range cases {
		if got := itoa(c.in); got != c.want {
			t.Fatalf("itoa(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestCallerUnavailable covers the runtime.Caller failure branch of caller.
func TestCallerUnavailable(t *testing.T) {
	if got := caller(1000); got != "" {
		t.Fatalf("caller(1000) = %q, want empty", got)
	}
}

// TestFailDefault covers both branches of failDefault directly.
func TestFailDefault(t *testing.T) {
	if got := failDefault("boom"); got != "boom" {
		t.Fatalf("failDefault(string) = %q, want boom", got)
	}
	if got := failDefault(42); got != "fail" {
		t.Fatalf("failDefault(other) = %q, want fail", got)
	}
}
