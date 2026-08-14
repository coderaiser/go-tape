package duplicates_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestDupA(t *testing.T) {
	Test(t, "dup: same", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestDupB(t *testing.T) {
	Test(t, "dup: same", func(t *T) {
		t.Ok(true)
		t.End()
	})
}
