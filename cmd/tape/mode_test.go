package main

import (
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
)

// TestRunInterpModeDefault verifies the default mode is interp and a passing
// supertape source exits 0.
func TestRunInterpModeDefault(t *testing.T) {
	Test(t, "mode: default interp passes a passing tape file", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"testdata/interp/pass_tape.go"}, &out, &errOut)
		t.Ok(code == 0 && strings.Contains(out.String(), "ok"))
		t.End()
	})
}

// TestRunInterpFlag verifies --mode interp works explicitly.
func TestRunInterpFlag(t *testing.T) {
	Test(t, "mode: --mode interp runs a failing tape file", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"--mode", "interp", "testdata/interp/fail_tape.go"}, &out, &errOut)
		t.Ok(code == 1 && strings.Contains(out.String(), "not ok"))
		t.End()
	})
}

// TestRunGotestMode verifies --mode gotest still uses the gotest subprocess path.
func TestRunGotestMode(t *testing.T) {
	Test(t, "mode: --mode gotest runs go test subprocess", func(t *T) {
		var out, errOut strings.Builder
		code := run([]string{"--mode", "gotest", "./testdata/covered/..."}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

// TestHelpContainsMode verifies the --mode flag is documented.
func TestHelpContainsMode(t *testing.T) {
	Test(t, "mode: -h documents --mode", func(t *T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), "--mode")
		t.End()
	})
}