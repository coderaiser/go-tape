package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet("go-tape", flag.ExitOnError)
}

func TestUsage(t *testing.T) {
	if !strings.Contains(usage, "go-tape") {
		t.Errorf("usage should contain 'go-tape', got %q", usage)
	}
	if !strings.Contains(usage, "-f") {
		t.Errorf("usage should contain '-f', got %q", usage)
	}
}

func TestVersionConst(t *testing.T) {
	if version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", version)
	}
}

func TestCheckScopesEnv(t *testing.T) {
	os.Setenv("TAPE_CHECK_SCOPES", "0")
	defer os.Unsetenv("TAPE_CHECK_SCOPES")
	if os.Getenv("TAPE_CHECK_SCOPES") != "0" {
		t.Error("expected TAPE_CHECK_SCOPES=0")
	}
}

func TestCheckAssertionsEnv(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	if os.Getenv("TAPE_CHECK_ASSERTIONS_COUNT") != "0" {
		t.Error("expected TAPE_CHECK_ASSERTIONS_COUNT=0")
	}
}

func TestVersionFlagParsing(t *testing.T) {
	resetFlags()
	os.Args = []string{"go-tape", "-v"}
	ver := flag.Bool("v", false, "")
	flag.Parse()
	if !*ver {
		t.Error("expected -v flag")
	}
}

func TestHelpFlagParsing(t *testing.T) {
	resetFlags()
	os.Args = []string{"go-tape", "-h"}
	help := flag.Bool("h", false, "")
	flag.Parse()
	if !*help {
		t.Error("expected -h flag")
	}
}

func TestFormatFlagParsing(t *testing.T) {
	resetFlags()
	os.Args = []string{"go-tape", "-f", "tap"}
	format := flag.String("f", "", "")
	flag.Parse()
	if *format != "tap" {
		t.Errorf("expected 'tap', got %q", *format)
	}
}

func TestNoCheckDuplicatesFlagParsing(t *testing.T) {
	resetFlags()
	os.Args = []string{"go-tape", "--no-check-duplicates"}
	noCheckDuplicates := flag.Bool("no-check-duplicates", false, "")
	flag.Parse()
	if !*noCheckDuplicates {
		t.Error("expected --no-check-duplicates flag")
	}
}

func TestNoCheckScopesFlagParsing(t *testing.T) {
	resetFlags()
	os.Args = []string{"go-tape", "--no-check-scopes"}
	noCheckScopes := flag.Bool("no-check-scopes", false, "")
	flag.Parse()
	if !*noCheckScopes {
		t.Error("expected --no-check-scopes flag")
	}
}

func TestNoCheckAssertionsCountFlagParsing(t *testing.T) {
	resetFlags()
	os.Args = []string{"go-tape", "--no-check-assertions-count"}
	noCheckAssertions := flag.Bool("no-check-assertions-count", false, "")
	flag.Parse()
	if !*noCheckAssertions {
		t.Error("expected --no-check-assertions-count flag")
	}
}

func TestPathArgParsing(t *testing.T) {
	resetFlags()
	os.Args = []string{"go-tape", "./mypkg"}
	flag.Parse()
	if flag.NArg() != 1 {
		t.Errorf("expected 1 arg, got %d", flag.NArg())
	}
	if flag.Arg(0) != "./mypkg" {
		t.Errorf("expected './mypkg', got %q", flag.Arg(0))
	}
}
