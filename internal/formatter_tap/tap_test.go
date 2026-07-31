package formatter_tap_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_tap"
)

func TestTAPStart(t *testing.T) {
	tape.Test(t, "formatter-tap: Start returns TAP version header", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Start(10)
		t.Equal(result, "TAP version 13\n")
		t.End()
	})
}

func TestTAPTest(t *testing.T) {
	tape.Test(t, "formatter-tap: Test returns comment line", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Test("scope: foo")
		t.Equal(result, "# scope: foo\n")
		t.End()
	})
}

func TestTAPTestEnd(t *testing.T) {
	tape.Test(t, "formatter-tap: TestEnd returns empty string", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.TestEnd(1, 10, 0, "scope: foo")
		t.Equal(result, "")
		t.End()
	})
}

func TestTAPSuccess(t *testing.T) {
	tape.Test(t, "formatter-tap: Success returns ok line", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Success(1, "scope: foo")
		t.Equal(result, "ok 1 scope: foo\n")
		t.End()
	})
}

func TestTAPFailWithOperator(t *testing.T) {
	tape.Test(t, "formatter-tap: Fail with operator shows expected and result", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "", "")
		t.Match(result, "not ok 1")
		t.End()
	})
}

func TestTAPFailWithoutOperator(t *testing.T) {
	tape.Test(t, "formatter-tap: Fail without operator shows raw output", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Fail(1, "scope: foo", "", nil, nil, "raw output", "", "")
		t.Match(result, "raw output")
		t.End()
	})
}

func TestTAPFailWithAt(t *testing.T) {
	tape.Test(t, "formatter-tap: Fail includes at when non-empty", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "file.go:10:", "")
		t.Match(result, "file.go:10:")
		t.End()
	})
}

func TestTAPFailWithErrorStack(t *testing.T) {
	tape.Test(t, "formatter-tap: Fail includes stack when non-empty", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "", "stack trace")
		t.Match(result, "stack trace")
		t.End()
	})
}

func TestTAPComment(t *testing.T) {
	tape.Test(t, "formatter-tap: Comment returns comment line", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.Comment("hello")
		t.Equal(result, "# hello\n")
		t.End()
	})
}

func TestTAPEndAllPass(t *testing.T) {
	tape.Test(t, "formatter-tap: End with no failures shows ok", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.End(5, 5, 0, 0)
		t.Match(result, "# ok")
		t.End()
	})
}

func TestTAPEndWithSkipped(t *testing.T) {
	tape.Test(t, "formatter-tap: End includes skip count when skipped > 0", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.End(5, 4, 0, 1)
		t.Match(result, "# skip 1")
		t.End()
	})
}

func TestTAPEndWithFailed(t *testing.T) {
	tape.Test(t, "formatter-tap: End includes fail count when failed > 0", func(t *tape.T) {
		f := formatter_tap.New()
		result := f.End(5, 4, 1, 0)
		t.Match(result, "# fail 1")
		t.End()
	})
}
