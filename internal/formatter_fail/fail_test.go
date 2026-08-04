package formatter_fail_test

import (
	"regexp"
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

func TestFailFormatterAlwaysPrependsTestName(t *testing.T) {
	tape.Test(t, "formatter-fail: Fail always prepends test name (even when empty)", func(t *tape.T) {
		f := formatter_fail.New()
		// no prior Test() call — currentTest is ""
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "", "")
		t.Match(result, regexp.MustCompile(`^# `))
		t.End()
	})
}

func TestFailFormatterKeepsErrorStack(t *testing.T) {
	tape.Test(t, "formatter-fail: Fail includes error stack (not stripped like formatter-short)", func(t *tape.T) {
		f := formatter_fail.New()
		f.Test("scope: foo")
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "", "my stack trace")
		t.Match(result, "my stack trace")
		t.End()
	})
}

func TestFailFormatterSuccessReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-fail: Success suppresses passing test output", func(t *tape.T) {
		f := formatter_fail.New()
		result := f.Success(1, "scope: foo")
		t.Equal(result, "")
		t.End()
	})
}
