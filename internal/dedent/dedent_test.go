//go:build !no_external

package dedent_test

import (
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/dedent"
)

func TestDedent(t *testing.T) {
	Test(t, "dedent: removes common leading whitespace", func(t *T) {
		got := dedent.Dedent("\n    hello\n    world\n")
		t.Equal(got, "\nhello\nworld\n")
		t.End()
	})
}
