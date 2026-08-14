package skipped_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestSkippedFixture(t *testing.T) {
	Test(t, "skipped: ran", func(t *T) {
		t.Ok(true)
		t.End()
	})
	t.Run("skipped: skipped", func(sub *testing.T) {
		sub.Skip()
	})
}
