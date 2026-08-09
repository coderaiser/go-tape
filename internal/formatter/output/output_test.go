package output_test

import (
	"os"
	"strings"
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
		fields := output.ParseOutput([]string{"    operator: should equal\n"})
		t.Equal(fields.Operator, "should equal")
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
			"should equal", "should not equal", "should be truthy", "should be falsy",
			"should deep equal", "should not deep equal", "should match", "should not match",
			"pass", "fail",
		}
		reached := true
		for _, op := range operators {
			fields := output.ParseOutput([]string{"    operator: " + op + "\n"})
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
			"    operator: should equal\n",
			"    result: 1\n",
			"    expected: 2\n",
			"    tape_test.go:10:\n",
		}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Operator, "should equal")
		t.Equal(fields.Result, "1")
		t.Equal(fields.Expected, "2")
		t.Equal(fields.At, "tape_test.go:10:")
		t.End()
	})
}

func TestParseOutputCutStripsAtPrefix(t *testing.T) {
	tape.Test(t, "output: Cut strips file:line prefix from first content line", func(t *tape.T) {
		lines := []string{
			"=== RUN   TestFoo\n",
			"    foo_test.go:42: operator: should be truthy\n",
			"        expected: truthy\n",
			"        result: false\n",
			"--- FAIL: TestFoo (0.00s)\n",
		}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Cut, "operator: should be truthy\n        expected: truthy\n        result: false\n")
		t.End()
	})
}

func TestParseOutputCutSkipsNoise(t *testing.T) {
	tape.Test(t, "output: Cut omits Go test runner noise lines", func(t *tape.T) {
		lines := []string{
			"=== RUN   TestFoo\n",
			"=== RUN   TestFoo/scope:_bar\n",
			"    foo_test.go:5: operator: should equal\n",
			"--- FAIL: TestFoo/scope:_bar (0.00s)\n",
			"--- FAIL: TestFoo (0.00s)\n",
		}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Cut, "operator: should equal\n")
		t.End()
	})
}

func TestParseOutputCutEmptyWhenOnlyNoise(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	tape.Test(t, "output: Cut is empty when all lines are noise", func(t *tape.T) {
		lines := []string{
			"=== RUN   TestFoo\n",
			"--- FAIL: TestFoo (0.00s)\n",
		}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Cut, "")
		t.End()
	})
}

func TestParseOutputDiffBlock(t *testing.T) {
	tape.Test(t, "output: parses diff block into Diff field", func(t *tape.T) {
		lines := []string{
			"- expected\n",
			"+ received\n",
			"\n",
			"- 1\n",
			"+ 2\n",
			"    operator: should equal\n",
		}
		fields := output.ParseOutput(lines)
		t.Ok(fields.Diff == "- 1\n+ 2" && fields.Operator == "should equal")
		t.End()
	})
}

func TestParseOutputTapeEndNoiseSuppressed(t *testing.T) {
	tape.Test(t, "output: tape: t.End() not called is filtered from Cut", func(t *tape.T) {
		lines := []string{
			"    extend.go:40: parse error: 2:1: expected 'package', found t\n",
			"    tape.go:99: tape: t.End() not called\n",
		}
		fields := output.ParseOutput(lines)
		t.Ok(!strings.Contains(fields.Cut, "t.End() not called"))
		t.End()
	})
}

func TestParseOutputTapeEndNoisePrimaryAtKept(t *testing.T) {
	tape.Test(t, "output: real error At kept when tape noise follows", func(t *tape.T) {
		lines := []string{
			"    extend.go:40: parse error: 2:1: expected 'package', found t\n",
			"    tape.go:99: tape: t.End() not called\n",
		}
		fields := output.ParseOutput(lines)
		t.Equal(fields.At, "extend.go:40:")
		t.End()
	})
}

func TestParseOutputOperatorWithFilePrefix(t *testing.T) {
	tape.Test(t, "output: parses operator from line with file:line: prefix", func(t *tape.T) {
		lines := []string{"    tape_test.go:37: operator: transform\n"}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Operator, "transform")
		t.End()
	})
}

func TestParseOutputAtAndOperatorSameLine(t *testing.T) {
	tape.Test(t, "output: At and Operator both set when on same line", func(t *tape.T) {
		lines := []string{"    tape_test.go:37: operator: noTransform\n"}
		fields := output.ParseOutput(lines)
		t.Equal(fields.At, "tape_test.go:37:")
		t.End()
	})
}

func TestParseOutputResultWithFilePrefix(t *testing.T) {
	tape.Test(t, "output: parses result from line with file:line: prefix", func(t *tape.T) {
		lines := []string{"    tape_test.go:10: result: 42\n"}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Result, "42")
		t.End()
	})
}

func TestParseOutputExpectedWithFilePrefix(t *testing.T) {
	tape.Test(t, "output: parses expected from line with file:line: prefix", func(t *tape.T) {
		lines := []string{"    tape_test.go:10: expected: 43\n"}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Expected, "43")
		t.End()
	})
}

func TestParseOutputFullErrorfBlock(t *testing.T) {
	tape.Test(t, "output: full Errorf block produces correct Operator and Diff", func(t *tape.T) {
		lines := []string{
			"    t_test.go:37: operator: transform\n",
			"        expected: hello\n",
			"        result: world\n",
			"        - expected\n",
			"        + received\n",
			"        \n",
			"        - hello\n",
			"        + world\n",
		}
		fields := output.ParseOutput(lines)
		t.Equal(fields.Operator, "transform")
		t.End()
	})
}

func TestParseOutputFullErrorfBlockDiff(t *testing.T) {
	tape.Test(t, "output: full Errorf block has non-empty Diff", func(t *tape.T) {
		lines := []string{
			"    t_test.go:37: operator: transform\n",
			"        expected: hello\n",
			"        result: world\n",
			"        - expected\n",
			"        + received\n",
			"        \n",
			"        - hello\n",
			"        + world\n",
		}
		fields := output.ParseOutput(lines)
		t.Ok(fields.Diff != "")
		t.End()
	})
}
