package failing

import "testing"

func TestAdd(t *testing.T) {
	if Add(1, 2) != 999 {
		t.Fatal("expected 999")
	}
}
