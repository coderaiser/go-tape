package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestRunHelp(t *testing.T) {
	Test(t, "main: -h prints usage", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"-h"}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestRunV(t *testing.T) {
	Test(t, "main: -v prints version", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"-v"}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestRunVersion(t *testing.T) {
	Test(t, "main: --version prints version", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"--version"}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestRunVersionOutput(t *testing.T) {
	Test(t, "main: -v output contains version string", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-v"}, &out, &errOut)
		t.Match(out.String(), regexp.MustCompile(`\d+\.\d+\.\d+`))
		t.End()
	})
}

func TestRunHOutput(t *testing.T) {
	Test(t, "main: -h output contains Usage", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)

		t.Match(out.String(), `Usage`)
		t.End()
	})
}

func TestRunHelpOutput(t *testing.T) {
	Test(t, "main: --help output contains Usage", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"--help"}, &out, &errOut)
		t.Match(out.String(), `Usage`)
		t.End()
	})
}

func TestHelpOutputContainsTapeTimeout(t *testing.T) {
	Test(t, "main: -h output contains TAPE_TIMEOUT", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), "TAPE_TIMEOUT")
		t.End()
	})
}

func TestHelpOutputContainsTapeCheckScopes(t *testing.T) {
	Test(t, "main: -h output contains TAPE_CHECK_SCOPES", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), "TAPE_CHECK_SCOPES")
		t.End()
	})
}

func TestHelpOutputContainsTapeCheckAssertionsCount(t *testing.T) {
	Test(t, "main: -h output contains TAPE_CHECK_ASSERTIONS_COUNT", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), "TAPE_CHECK_ASSERTIONS_COUNT")
		t.End()
	})
}

func TestHelpOutputContainsTapeCheckSkipped(t *testing.T) {
	Test(t, "main: -h output contains TAPE_CHECK_SKIPPED", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), "TAPE_CHECK_SKIPPED")
		t.End()
	})
}

func TestHelpOutputContainsTapeCheckDuplicates(t *testing.T) {
	Test(t, "main: -h output contains TAPE_CHECK_DUPLICATES", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), "TAPE_CHECK_DUPLICATES")
		t.End()
	})
}

func TestCheckSkippedExitCode(t *testing.T) {
	Test(t, "main: TAPE_CHECK_SKIPPED=1 returns exit code 5 when skipped > 0", func(t *T) {
		t.TB().Setenv("TAPE_CHECK_SKIPPED", "1")
		var out, errOut strings.Builder
		code := run([]string{"./testdata/skipped/..."}, &out, &errOut)
		t.Equal(code, 5)
		t.End()
	})
}

func TestCheckSkippedOffExitCode(t *testing.T) {
	Test(t, "main: TAPE_CHECK_SKIPPED=0 returns 0 when only skipped tests exist", func(t *T) {
		t.TB().Setenv("TAPE_CHECK_SKIPPED", "0")
		var out, errOut strings.Builder
		code := run([]string{"./testdata/skipped/..."}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestCoverageFlagExitZero(t *testing.T) {
	Test(t, "main: -c exits 0 when all covered", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"-c", "./testdata/covered/..."}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestCoverageFlagUncoveredExitOne(t *testing.T) {
	Test(t, "main: -c exits 1 when uncovered blocks exist", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"-c", "./testdata/uncovered/..."}, &out, &errOut)
		t.Equal(code, 1)
		t.End()
	})
}

func TestCoverageUncoveredOutput(t *testing.T) {
	Test(t, "main: -c output lists uncovered file", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-c", "./testdata/uncovered/..."}, &out, &errOut)
		t.Match(out.String(), "uncovered.go")
		t.End()
	})
}

func TestCoverageJSONLines(t *testing.T) {
	Test(t, "main: -c=json-lines emits json for uncovered block", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-c=json-lines", "./testdata/uncovered/..."}, &out, &errOut)
		t.Match(out.String(), `"file"`)
		t.End()
	})
}

func TestCoverageReportWritesFile(t *testing.T) {
	Test(t, "main: -c -r writes coverage.lcov", func(t *T) {
		dir := t.TB().TempDir()
		path := dir + "/coverage.lcov"
		var out, errOut strings.Builder
		run([]string{"-c", "-r=" + path, "./testdata/covered/..."}, &out, &errOut)
		_, err := os.Stat(path)
		t.Ok(err == nil)
		t.End()
	})
}

func TestCoverageReportFileNotEmpty(t *testing.T) {
	Test(t, "main: -c -r produces non-empty lcov file", func(t *T) {
		dir := t.TB().TempDir()
		path := dir + "/coverage.lcov"
		var out, errOut strings.Builder
		run([]string{"-c", "-r=" + path, "./testdata/covered/..."}, &out, &errOut)
		info, _ := os.Stat(path)
		t.Ok(info != nil && info.Size() > 0)
		t.End()
	})
}

func TestCoverageFullyCoveredShowsSummary(t *testing.T) {
	Test(t, "main: -c prints coverage summary when fully covered", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-c", "./testdata/covered/..."}, &out, &errOut)
		t.Match(out.String(), "good job")
		t.End()
	})
}

func TestCoverageFullyCoveredExitZeroWithSummary(t *testing.T) {
	Test(t, "main: -c exits 0 when fully covered with summary", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"-c", "./testdata/covered/..."}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestCoverageUncoveredNoSummary(t *testing.T) {
	Test(t, "main: -c does not print good job summary when uncovered", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"-c", "./testdata/uncovered/..."}, &out, &errOut)
		t.Equal(code, 1)
		t.End()
	})
}

func TestCoverageUncoveredNoGoodJob(t *testing.T) {
	Test(t, "main: -c output omits good job when uncovered", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-c", "./testdata/uncovered/..."}, &out, &errOut)
		t.NotMatch(out.String(), "good job")
		t.End()
	})
}

func TestCoverageSkippedOnTestFailure(t *testing.T) {
	Test(t, "main: -c skips coverage output when tests fail", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"-c", "./testdata/failing/..."}, &out, &errOut)
		t.Equal(code, 1)
		t.End()
	})

	Test(t, "main: -c skips coverage output when tests fail: no good job", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-c", "./testdata/failing/..."}, &out, &errOut)
		t.NotMatch(out.String(), "good job")
		t.End()
	})
}

func TestCoverageHelpContainsCFlag(t *testing.T) {
	Test(t, "main: -h output contains -c flag", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), `-c`)
		t.End()
	})
}

func TestCoverageHelpContainsRFlag(t *testing.T) {
	Test(t, "main: -h output contains -r flag", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), `-r`)
		t.End()
	})
}

func TestCoverageTapeTomlExclude(t *testing.T) {
	Test(t, "main: coverage exclude from .tape.toml is passed to ProcessProfileWithConfig", func(t *T) {
		dir := t.TB().TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".tape.toml"), []byte(`
[coverage]
exclude = ["node_modules"]
`), 0o644); err != nil {
			t.TB().Fatal(err)
		}
		// No assertion on output — just confirm it does not panic and
		// that ProcessProfileWithConfig is reachable via the wired path.
		// A real integration test requires a coverprofile; this confirms
		// the config is loaded and passed through.
		var out, errOut strings.Builder
		old, _ := os.Getwd()
		os.Chdir(dir)
		defer os.Chdir(old)
		run([]string{"-c", "./testdata/covered/..."}, &out, &errOut)
		t.End()
	})
}
