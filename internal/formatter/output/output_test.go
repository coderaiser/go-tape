package output_test

import (
	"os"
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter/output"
)

func TestParseOutputEmpty(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	tape.Test(t, "output: empty lines returns empty fields", func(t *tape.T) {
		fields := output.ParseOutput([]string{})
		t.Equal(fields.Raw, "")
		t.Equal(fields.Operator, "")
		t.Equal(fields.Result, "")
		t.Equal(fields.Expected, "")
		t.Equal(fields.At, "")
		t.Equal(fields.ErrorStack, "")
		t.End()
	})
}

func TestParseOutputOperator(t *testing.T) {
	tape.Test(t, "output: parses operator line", func(t *tape.T) {
		fields := output.ParseOutput([]string{"    Equal\n"})
		t.Equal(fields.Operator, "Equal")
		t.End()
	})
}

func TestParseOutputResult(t *testing.T) {
	tape.Test(t, "output: parses result line", func(t *tape.T) {
		fields := output.ParseOutput([]string{"    result: 42\n"})
		t.Equal(fields.Result, "42")
		t.End()
	})
}

func TestParseOutputExpected(t *testing.T) {
	tape.Test(t, "output: parses expected line", func(t *tape.T) {
		fields := output.ParseOutput([]string{"    expected: 43\n"})
		t.Equal(fields.Expected, "43")
		t.End()
	})
}

func TestParseOutputAt(t *testing.T) {
	tape.Test(t, "output: parses at line", func(t *tape.T) {
		fields := output.ParseOutput([]string{"    tape_test.go:42:\n"})
		t.Equal(fields.At, "tape_test.go:42:")
		t.End()
	})
}

func TestParseOutputRaw(t *testing.T) {
	tape.Test(t, "output: raw is all lines joined", func(t *tape.T) {
		fields := output.ParseOutput([]string{"line1\n", "line2\n"})
		t.Equal(fields.Raw, "line1\nline2\n")
		t.End()
	})
}

func TestParseOutputAllOperators(t *testing.T) {
	tape.Test(t, "output: all known operators parse", func(t *tape.T) {
		operators := []string{
			"Equal", "NotEqual", "Ok", "NotOk",
			"DeepEqual", "NotDeepEqual", "Match", "NotMatch",
			"Error", "NoError", "Pass", "Fail",
		}
		reached := true
		for _, op := range operators {
			fields := output.ParseOutput([]string{"    " + op + "\n"})
			if fields.Operator != op {
				reached = false
			}
		}
		t.Ok(reached)
		t.End()
	})
}

func TestParseOutputUnrecognizedLine(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	tape.Test(t, "output: unrecognized lines go into raw only", func(t *tape.T) {
		fields := output.ParseOutput([]string{"some random line\n"})
		t.Equal(fields.Raw, "some random line\n")
		t.Equal(fields.Operator, "")
		t.Equal(fields.At, "")
		t.End()
	})
}

func TestParseOutputFullBlock(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	tape.Test(t, "output: full failure block parses all fields", func(t *tape.T) {
		lines := []string{
			"    Equal\n",
			"    result: 1\n",
			"    expected: 2\n",
			"    tape_test.go:10:\n",
		}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Operator, "Equal")
		t.Equal(fields.Result, "1")
		t.Equal(fields.Expected, "2")
		t.Equal(fields.At, "tape_test.go:10:")
		t.End()
	})
}
