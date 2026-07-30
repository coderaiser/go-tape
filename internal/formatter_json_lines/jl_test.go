package formatter_json_lines_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_json_lines"
)

func TestJLTestEnd(t *testing.T) {
	tape.Test(t, "formatter-json-lines: TestEnd returns JSON with count and test", func(t *tape.T) {
		f := formatter_json_lines.New(10)
		result := f.TestEnd(1, 10, 0, "scope: foo")
		t.Match(result, `"count":1`)
		t.End()
	})
}

func TestJLFail(t *testing.T) {
	tape.Test(t, "formatter-json-lines: Fail returns JSON with test name", func(t *tape.T) {
		f := formatter_json_lines.New(10)
		result := f.Fail(1, "scope: foo", "Equal", "got", "want", "", "", "")
		t.Match(result, `"test":"scope: foo"`)
		t.End()
	})
}

func TestJLEnd(t *testing.T) {
	tape.Test(t, "formatter-json-lines: End returns JSON with passed and failed", func(t *tape.T) {
		f := formatter_json_lines.New(10)
		result := f.End(5, 5, 0, 0)
		t.Match(result, `"passed":5`)
		t.End()
	})
}

func TestJLStartReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-json-lines: Start returns empty string", func(t *tape.T) {
		f := formatter_json_lines.New(10)
		result := f.Start(10)
		t.Equal(result, "")
		t.End()
	})
}

func TestJLTestReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-json-lines: Test returns empty string", func(t *tape.T) {
		f := formatter_json_lines.New(10)
		result := f.Test("scope: foo")
		t.Equal(result, "")
		t.End()
	})
}

func TestJLCommentReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-json-lines: Comment returns empty string", func(t *tape.T) {
		f := formatter_json_lines.New(10)
		result := f.Comment("hello")
		t.Equal(result, "")
		t.End()
	})
}

func TestJLSuccessReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-json-lines: Success returns empty string", func(t *tape.T) {
		f := formatter_json_lines.New(10)
		result := f.Success(1, "scope: foo")
		t.Equal(result, "")
		t.End()
	})
}