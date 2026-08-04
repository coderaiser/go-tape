package covered

import "testing"

func TestCovered(t *testing.T) {
	if Covered() != "covered" {
		t.Fatal("expected covered")
	}
}
