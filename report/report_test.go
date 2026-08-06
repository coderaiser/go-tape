package report_test

import (
	"bufio"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/report"
)

const passingJSON = `{"Action":"run","Test":"TestFoo/scope:_ok"}
{"Action":"output","Test":"TestFoo/scope:_ok","Output":"    ok\n"}
{"Action":"pass","Test":"TestFoo/scope:_ok","Elapsed":0.001}
`

const failingJSON = `{"Action":"run","Test":"TestFoo/scope:_fail"}
{"Action":"output","Test":"TestFoo/scope:_fail","Output":"    operator: Equal\n"}
{"Action":"output","Test":"TestFoo/scope:_fail","Output":"    expected: |-\n"}
{"Action":"output","Test":"TestFoo/scope:_fail","Output":"      1\n"}
{"Action":"output","Test":"TestFoo/scope:_fail","Output":"    result: |-\n"}
{"Action":"output","Test":"TestFoo/scope:_fail","Output":"      2\n"}
{"Action":"fail","Test":"TestFoo/scope:_fail","Elapsed":0.001}
`

func TestRunNoError(t *testing.T) {
	Test(t, "report: run returns no error on valid input", func(t *T) {
		var sb strings.Builder
		result := report.Run(strings.NewReader(passingJSON), &sb, "tap", 1)
		t.NotOk(result)
		t.End()
	})
}

func TestRunPassingOutputsOk(t *testing.T) {
	Test(t, "report: passing test produces ok line in tap format", func(t *T) {
		var sb strings.Builder
		report.Run(strings.NewReader(passingJSON), &sb, "tap", 1)
		t.Match(sb.String(), "ok 1")
		t.End()
	})
}

func TestRunFailingOutputsNotOk(t *testing.T) {
	Test(t, "report: failing test produces not ok line in tap format", func(t *T) {
		var sb strings.Builder
		report.Run(strings.NewReader(failingJSON), &sb, "tap", 1)
		t.Match(sb.String(), "not ok")
		t.End()
	})
}

func TestRunFailFormatSuppressesPassing(t *testing.T) {
	Test(t, "report: fail formatter suppresses passing test output", func(t *T) {
		var sb strings.Builder
		report.Run(strings.NewReader(passingJSON), &sb, "fail", 1)
		t.NotMatch(sb.String(), "ok 1")
		t.End()
	})
}

func TestRunFailFormatShowsFailures(t *testing.T) {
	Test(t, "report: fail formatter shows failing tests", func(t *T) {
		var sb strings.Builder
		report.Run(strings.NewReader(failingJSON), &sb, "fail", 1)
		t.Match(sb.String(), "not ok")
		t.End()
	})
}

func TestRunEmptyInput(t *testing.T) {
	Test(t, "report: empty input produces no error", func(t *T) {
		var sb strings.Builder
		result := report.Run(strings.NewReader(""), &sb, "tap", 0)
		t.NotOk(result)
		t.End()
	})
}

func TestRunJSONLinesFormat(t *testing.T) {
	Test(t, "report: json-lines format emits JSON output", func(t *T) {
		var sb strings.Builder
		report.Run(strings.NewReader(failingJSON), &sb, "json-lines", 1)
		t.Match(sb.String(), `"`)
		t.End()
	})
}

func TestRunGarbageLinesIgnored(t *testing.T) {
	Test(t, "report: non-JSON lines are silently skipped", func(t *T) {
		input := "not json at all\n" + passingJSON
		var sb strings.Builder
		report.Run(strings.NewReader(input), &sb, "tap", 1)
		t.Match(sb.String(), "ok 1")
		t.End()
	})
}

func TestRunMixedPassFailNotOk(t *testing.T) {
	Test(t, "report: mixed results show failures in fail format", func(t *T) {
		var sb strings.Builder
		report.Run(strings.NewReader(passingJSON+failingJSON), &sb, "fail", 2)
		t.Match(sb.String(), "not ok")
		t.End()
	})
}

func TestRunMixedPassFailSuppressPass(t *testing.T) {
	Test(t, "report: mixed results suppress passing in fail format", func(t *T) {
		var sb strings.Builder
		report.Run(strings.NewReader(passingJSON+failingJSON), &sb, "fail", 2)
		t.NotMatch(sb.String(), "ok 1")
		t.End()
	})
}

func TestRunUnknownActionSkipped(t *testing.T) {
	Test(t, "report: unknown action is skipped without error", func(t *T) {
		var sb strings.Builder
		result := report.Run(strings.NewReader(`{"Action":"bogus","Test":"TestFoo"}`), &sb, "tap", 1)
		t.NotOk(result)
		t.End()
	})
}

func TestRunScannerError(t *testing.T) {
	Test(t, "report: scanner error returns error", func(t *T) {
		long := strings.Repeat("a", bufio.MaxScanTokenSize+1)
		var sb strings.Builder
		result := report.Run(strings.NewReader(long), &sb, "tap", 1)
		t.Ok(result != nil)
		t.End()
	})
}
