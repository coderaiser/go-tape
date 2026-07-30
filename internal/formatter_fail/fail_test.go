package formatter_fail_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_fail"
)

func TestFailFormatterTestStoresName(t *testing.T) {
	tape.Test(t, "formatter-fail: Test stores current test name", func(t *tape.T) {
		f := formatter_fail.New()
		result := f.Test("scope: foo")
		t.Equal(result, "")
		t.End()
	})
}

func TestFailFormatterFailWithCurrentTest(t *testing.T) {
	tape.Test(t, "formatter-fail: Fail prepends current test name", func(t *tape.T) {
		f := formatter_fail.New()
		f.Test("scope: foo")
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "", "")
		t.Match(result, "# scope: foo")
		t.End()
	})
}

func TestFailFormatterFailWithoutCurrentTest(t *testing.T) {
	tape.Test(t, "formatter-fail: Fail has no prefix when no current test set", func(t *tape.T) {
		f := formatter_fail.New()
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "", "")
		t.NotMatch(result, "^#")
		t.End()
	})
}