package main

import (
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

func TestVersionConst(t *testing.T) {
	Test(t, "main: version is 1.0.0", func(t *T) {
		t.Equal(version, "1.0.0")
		t.End()
	})
}
