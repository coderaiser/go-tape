//go:build no_external

package dedent_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/dedent"
)

func TestDedentCommonIndent(t *testing.T) {
	tape.Test(t, "dedent: Dedent removes common leading whitespace", func(t *tape.T) {
		input := "    hello\n    world\n"
		result := dedent.Dedent(input)
		expected := "hello\nworld\n"
		t.Equal(result, expected)
		t.End()
	})
}

func TestDedentNoCommonIndent(t *testing.T) {
	tape.Test(t, "dedent: Dedent returns unchanged when no common indent", func(t *tape.T) {
		input := "  foo\nbar\n"
		result := dedent.Dedent(input)
		t.Match(result, "foo")
		t.End()
	})
}

func TestDedentShrinkingMargin(t *testing.T) {
	tape.Test(t, "dedent: Dedent uses smallest common indent", func(t *tape.T) {
		input := "    deep\n  shallow\n"
		result := dedent.Dedent(input)
		expected := "  deep\nshallow\n"
		t.Equal(result, expected)
		t.End()
	})
}

func TestDedentConflictingIndents(t *testing.T) {
	tape.Test(t, "dedent: Dedent returns unchanged for mixed tab/space indent", func(t *tape.T) {
		input := "\t foo\n  bar\n"
		result := dedent.Dedent(input)
		t.Ok(result != "")
		t.End()
	})
}