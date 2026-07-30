package diff_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/diff"
)

func TestDiffIsNonEmpty(t *testing.T) {
	tape.Test(t, "diff: Diff returns non-empty string for different values", func(t *tape.T) {
		result := diff.Diff("want", "got")
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffContainsExpected(t *testing.T) {
	tape.Test(t, "diff: Diff output contains expected value", func(t *tape.T) {
		result := diff.Diff("want", "got")
		t.Match(result, "want")
		t.End()
	})
}