package skipped_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestSkippedFixture(t *testing.T) {
	Test.Skip(t, "skipped: this test is skipped", func(t *T) {
		t.Ok(true)
		t.End()
	})
}
