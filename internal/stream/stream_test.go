package stream_test

import (
	"encoding/json"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
	. "github.com/coderaiser/go-tape/internal/stream"
)

// collect drains the Parse channel into a slice.
func collect(ch <-chan Event) []Event {
	var out []Event
	for e := range ch {
		out = append(out, e)
	}
	return out
}

// lines joins strings with newlines — keeps fixtures readable.
func lines(ll ...string) string {
	return strings.Join(ll, "\n") + "\n"
}

const (
	pkg      = "github.com/coderaiser/go-tape/internal/stream"
	outerFn  = "TestStream"
	subtest  = "TestStream/scope:_my_test"
	testName = "scope: my test"
)

func tapeOutput(tapeJSON string) string {
	// go test -json JSON-escapes the Output field; replicate that here so the
	// TAPE: payload's quotes survive the outer json.Unmarshal in Parse.
	payload := "    TAPE:" + tapeJSON + "\n"
	enc, _ := json.Marshal(payload)
	return `{"Action":"output","Package":"` + pkg + `","Test":"` + subtest + `","Output":` + string(enc) + `}`
}

func runLine(test string) string {
	return `{"Action":"run","Package":"` + pkg + `","Test":"` + test + `"}`
}

func passLine(test string) string {
	return `{"Action":"pass","Package":"` + pkg + `","Test":"` + test + `","Elapsed":0.01}`
}

func failLine(test string) string {
	return `{"Action":"fail","Package":"` + pkg + `","Test":"` + test + `","Elapsed":0.01}`
}

func skipLine(test string) string {
	return `{"Action":"skip","Package":"` + pkg + `","Test":"` + test + `","Elapsed":0.0}`
}

func outputLine(test, output string) string {
	return `{"Action":"output","Package":"` + pkg + `","Test":"` + test + `","Output":"` + output + `"}`
}

func pkgOutputLine(output string) string {
	return `{"Action":"output","Package":"` + pkg + `","Output":"` + output + `"}`
}

func pkgFailLine() string {
	return `{"Action":"fail","Package":"` + pkg + `"}`
}

// ---------------------------------------------------------------------------
// pass
// ---------------------------------------------------------------------------

func TestStreamPassEmitsOneEvent(t *testing.T) {
	Test(t, "stream: pass emits one event", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":1,"expected":1,"output":""}`),
			passLine(subtest),
			passLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(len(events), 1)
		t.End()
	})
}

func TestStreamPassEventType(t *testing.T) {
	Test(t, "stream: pass event has type test-end", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			passLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Type, TypeTestEnd)
		t.End()
	})
}

func TestStreamPassEventName(t *testing.T) {
	Test(t, "stream: pass event carries decoded test name", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			passLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Test, testName)
		t.End()
	})
}

func TestStreamPassEventCounts(t *testing.T) {
	Test(t, "stream: pass event has correct count, total, failed", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			passLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.DeepEqual(events[0], Event{
			Type:   TypeTestEnd,
			Test:   testName,
			Count:  1,
			Total:  1,
			Failed: 0,
		})
		t.End()
	})
}

// ---------------------------------------------------------------------------
// fail
// ---------------------------------------------------------------------------

func TestStreamFailEmitsTwoEvents(t *testing.T) {
	Test(t, "stream: fail emits two events (fail + test-end)", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":2,"expected":1,"output":""}`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(len(events), 2)
		t.End()
	})
}

func TestStreamFailFirstEventType(t *testing.T) {
	Test(t, "stream: fail first event has type fail", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":2,"expected":1,"output":""}`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Type, TypeFail)
		t.End()
	})
}

func TestStreamFailSecondEventType(t *testing.T) {
	Test(t, "stream: fail second event has type test-end", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":2,"expected":1,"output":""}`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[1].Type, TypeTestEnd)
		t.End()
	})
}

func TestStreamFailEventFields(t *testing.T) {
	Test(t, "stream: fail event carries operator, result, expected from TAPE: line", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":2,"expected":1,"output":""}`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.DeepEqual(events[0], Event{
			Type:     TypeFail,
			Test:     testName,
			Count:    1,
			Message:  "should equal",
			Operator: "should equal",
			Result:   float64(2), // JSON numbers decode as float64
			Expected: float64(1),
			Output:   "",
		})
		t.End()
	})
}

func TestStreamFailTestEndHasFailedOne(t *testing.T) {
	Test(t, "stream: fail test-end has failed count 1", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":2,"expected":1,"output":""}`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[1].Failed, 1)
		t.End()
	})
}

func TestStreamFailWithOutput(t *testing.T) {
	Test(t, "stream: fail event carries non-empty output field", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":[],"expected":[],"output":"values not equal, but deepEqual"}`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Output, "values not equal, but deepEqual")
		t.End()
	})
}

func TestStreamFailPrefixedTAPELine(t *testing.T) {
	Test(t, "stream: TAPE: line with file.go:N: prefix is parsed and At is captured", func(t *T) {
		// go test -json wraps t.Log output as "    foo_test.go:12: TAPE:{...}".
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			outputLine(subtest, `    foo_test.go:12: TAPE:{\"message\":\"should equal\",\"operator\":\"should equal\",\"result\":2,\"expected\":1,\"output\":\"\"}\n`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.DeepEqual(events[0], Event{
			Type:     TypeFail,
			Test:     testName,
			Count:    1,
			Message:  "should equal",
			Operator: "should equal",
			Result:   float64(2),
			Expected: float64(1),
			At:       "foo_test.go:12",
		})
		t.End()
	})
}

func TestStreamFailAtOnSeparateLine(t *testing.T) {
	Test(t, "stream: At is captured from a preceding t.Log line", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			outputLine(subtest, `    foo_test.go:12: \n`),
			tapeOutput(`{"message":"should equal","operator":"should equal","result":2,"expected":1,"output":""}`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].At, "foo_test.go:12")
		t.End()
	})
}

func TestStreamTAPEMalformedFallsThrough(t *testing.T) {
	Test(t, "stream: malformed TAPE payload falls through to unknown-fail", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			outputLine(subtest, `    some plain text without a file link\n`),
			tapeOutput(`not-json-at-all`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Ok(len(events) == 2 && events[0].Type == TypeUnknownFail)
		t.End()
	})
}

// ---------------------------------------------------------------------------
// skip
// ---------------------------------------------------------------------------

func TestStreamSkipEmitsOneEvent(t *testing.T) {
	Test(t, "stream: skip emits one event", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			skipLine(subtest),
			skipLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(len(events), 1)
		t.End()
	})
}

func TestStreamSkipEventType(t *testing.T) {
	Test(t, "stream: skip event has type test-end", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			skipLine(subtest),
			skipLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Type, TypeTestEnd)
		t.End()
	})
}

func TestStreamSkipDoesNotIncrementFailed(t *testing.T) {
	Test(t, "stream: skip does not increment failed counter", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			skipLine(subtest),
			skipLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Failed, 0)
		t.End()
	})
}

// ---------------------------------------------------------------------------
// multiple tests — counter accumulation
// ---------------------------------------------------------------------------

func TestStreamTwoPassesCounterAccumulates(t *testing.T) {
	Test(t, "stream: two passes produce count 1 then 2", func(t *T) {
		sub2 := "TestStream/scope:_second_test"
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			runLine(sub2),
			passLine(sub2),
			passLine(outerFn),
		))
		events := collect(Parse(r, 2))
		t.Ok(events[0].Count == 1 && events[1].Count == 2)
		t.End()
	})
}

func TestStreamPassThenFailFailedCounter(t *testing.T) {
	Test(t, "stream: pass then fail: failed counter is 0 then 1", func(t *T) {
		sub2 := "TestStream/scope:_second_test"
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			runLine(sub2),
			tapeOutput(`{"message":"should be truthy","operator":"ok","result":false,"expected":true,"output":""}`),
			failLine(sub2),
			failLine(outerFn),
		))
		events := collect(Parse(r, 2))
		// events: test-end(pass), fail, test-end(fail)
		t.Ok(events[0].Failed == 0 && events[2].Failed == 1)
		t.End()
	})
}

func TestStreamTotalFlowsThrough(t *testing.T) {
	Test(t, "stream: total from Parse arg appears on every test-end", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			passLine(outerFn),
		))
		events := collect(Parse(r, 42))
		t.Equal(events[0].Total, 42)
		t.End()
	})
}

// ---------------------------------------------------------------------------
// outer TestXxx wrapper — ignored
// ---------------------------------------------------------------------------

func TestStreamOuterWrapperIgnored(t *testing.T) {
	Test(t, "stream: outer TestXxx events do not produce stream events", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			passLine(outerFn),
		))
		events := collect(Parse(r, 0))
		t.Equal(len(events), 0)
		t.End()
	})
}

// ---------------------------------------------------------------------------
// pause / cont — ignored
// ---------------------------------------------------------------------------

func TestStreamPauseIgnored(t *testing.T) {
	Test(t, "stream: pause events produce no output", func(t *T) {
		r := strings.NewReader(lines(
			`{"Action":"pause","Package":"` + pkg + `","Test":"` + subtest + `"}`,
		))
		events := collect(Parse(r, 1))
		t.Equal(len(events), 0)
		t.End()
	})
}

func TestStreamContIgnored(t *testing.T) {
	Test(t, "stream: cont events produce no output", func(t *T) {
		r := strings.NewReader(lines(
			`{"Action":"cont","Package":"` + pkg + `","Test":"` + subtest + `"}`,
		))
		events := collect(Parse(r, 1))
		t.Equal(len(events), 0)
		t.End()
	})
}

// ---------------------------------------------------------------------------
// build error
// ---------------------------------------------------------------------------

func TestStreamBuildErrorEmitsEvent(t *testing.T) {
	Test(t, "stream: build failure emits build-error event", func(t *T) {
		r := strings.NewReader(lines(
			pkgOutputLine(`foo.go:5:2: declared and not used: x\n`),
			pkgOutputLine(`FAIL\t`+pkg+` [build failed]\n`),
			pkgFailLine(),
		))
		events := collect(Parse(r, 0))
		t.Equal(len(events), 1)
		t.End()
	})
}

func TestStreamBuildErrorEventType(t *testing.T) {
	Test(t, "stream: build-error event has correct type", func(t *T) {
		r := strings.NewReader(lines(
			pkgOutputLine(`foo.go:5:2: declared and not used: x\n`),
			pkgOutputLine(`FAIL\t`+pkg+` [build failed]\n`),
			pkgFailLine(),
		))
		events := collect(Parse(r, 0))
		t.Equal(events[0].Type, TypeBuildError)
		t.End()
	})
}

func TestStreamBuildErrorEventPackage(t *testing.T) {
	Test(t, "stream: build-error event carries package name", func(t *T) {
		r := strings.NewReader(lines(
			pkgOutputLine(`foo.go:5:2: declared and not used: x\n`),
			pkgOutputLine(`FAIL\t`+pkg+` [build failed]\n`),
			pkgFailLine(),
		))
		events := collect(Parse(r, 0))
		t.Equal(events[0].Package, pkg)
		t.End()
	})
}

func TestStreamBuildErrorEventOutput(t *testing.T) {
	Test(t, "stream: build-error event carries compiler output", func(t *T) {
		r := strings.NewReader(lines(
			pkgOutputLine(`foo.go:5:2: declared and not used: x\n`),
			pkgOutputLine(`FAIL\t`+pkg+` [build failed]\n`),
			pkgFailLine(),
		))
		events := collect(Parse(r, 0))
		t.Match(events[0].Output, "declared and not used")
		t.End()
	})
}

func TestStreamTwoBuildErrorsTwoEvents(t *testing.T) {
	Test(t, "stream: two packages failing to build emit two build-error events", func(t *T) {
		pkg2 := "github.com/coderaiser/go-tape/internal/other"
		r := strings.NewReader(lines(
			pkgOutputLine(`foo.go:1:1: declared and not used: x\n`),
			pkgOutputLine(`FAIL\t`+pkg+` [build failed]\n`),
			pkgFailLine(),
			`{"Action":"output","Package":"`+pkg2+`","Output":"bar.go:1:1: declared and not used: y\n"}`,
			`{"Action":"output","Package":"`+pkg2+`","Output":"FAIL\t`+pkg2+` [build failed]\n"}`,
			`{"Action":"fail","Package":"`+pkg2+`"}`,
		))
		events := collect(Parse(r, 0))
		t.Equal(len(events), 2)
		t.End()
	})
}

// ---------------------------------------------------------------------------
// unknown-fail — non-tape test with no TAPE: sentinel
// ---------------------------------------------------------------------------

func TestStreamUnknownFailEmitsTwoEvents(t *testing.T) {
	Test(t, "stream: non-tape fail emits unknown-fail and test-end", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			outputLine(subtest, `    t.Fatal called from somewhere\n`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(len(events), 2)
		t.End()
	})
}

func TestStreamUnknownFailFirstEventType(t *testing.T) {
	Test(t, "stream: non-tape fail first event has type unknown-fail", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			outputLine(subtest, `    t.Fatal called from somewhere\n`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Type, TypeUnknownFail)
		t.End()
	})
}

func TestStreamUnknownFailCarriesRawOutput(t *testing.T) {
	Test(t, "stream: unknown-fail event carries raw output", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			outputLine(subtest, `    t.Fatal called from somewhere\n`),
			failLine(subtest),
			failLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Match(events[0].Output, "t.Fatal")
		t.End()
	})
}

// ---------------------------------------------------------------------------
// invalid / malformed input
// ---------------------------------------------------------------------------

func TestStreamInvalidJSONSkipped(t *testing.T) {
	Test(t, "stream: invalid JSON lines are silently skipped", func(t *T) {
		r := strings.NewReader(lines(
			"not json at all",
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			passLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(len(events), 1)
		t.End()
	})
}

func TestStreamEmptyInputNoEvents(t *testing.T) {
	Test(t, "stream: empty input produces no events", func(t *T) {
		events := collect(Parse(strings.NewReader(""), 0))
		t.Equal(len(events), 0)
		t.End()
	})
}

func TestStreamChannelClosedAfterInput(t *testing.T) {
	Test(t, "stream: channel is closed after input is exhausted", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			passLine(outerFn),
		))
		ch := Parse(r, 1)
		for range ch {
		}
		// if channel is not closed this test times out
		t.Ok(true)
		t.End()
	})
}

// ---------------------------------------------------------------------------
// test name decoding
// ---------------------------------------------------------------------------

func TestStreamTestNameDecodesUnderscores(t *testing.T) {
	Test(t, "stream: underscores in subtest name decoded to spaces", func(t *T) {
		sub := "TestStream/scope:_bar_baz"
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(sub),
			passLine(sub),
			passLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.Equal(events[0].Test, "scope: bar baz")
		t.End()
	})
}

func TestStreamTestNameStripsOuterFunction(t *testing.T) {
	Test(t, "stream: outer TestXxx prefix is stripped from test name", func(t *T) {
		r := strings.NewReader(lines(
			runLine(outerFn),
			runLine(subtest),
			passLine(subtest),
			passLine(outerFn),
		))
		events := collect(Parse(r, 1))
		t.NotMatch(events[0].Test, "TestStream")
		t.End()
	})
}

// ---------------------------------------------------------------------------
// run — exec path (covers the goroutine wiring in Run, used by cmd/tape).
// ---------------------------------------------------------------------------

func TestStreamRunVersion(t *testing.T) {
	Test(t, "stream: Run executes go and closes channel", func(t *T) {
		ch, err := Run(0, "version")
		if err != nil {
			t.Ok(false)
			t.End()
			return
		}
		for range ch {
		}
		// draining without locking proves the channel is closed at EOF
		t.Ok(true)
		t.End()
	})
}

func TestStreamRunStartError(t *testing.T) {
	Test(t, "stream: Run returns error when go is not executable", func(t *T) {
		t.TB().Setenv("PATH", t.TB().TempDir())
		_, err := Run(0, "version")
		t.Ok(err != nil)
		t.End()
	})
}
