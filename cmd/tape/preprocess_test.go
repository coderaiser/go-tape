package main

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestPreprocessCFlag(t *testing.T) {
	Test(t, "preprocess: -c sets coverage enabled", func(t *T) {
		opts, _ := preprocess([]string{"-c"})
		t.Ok(opts.enabled)
		t.End()
	})
}

func TestPreprocessCoverageFlag(t *testing.T) {
	Test(t, "preprocess: --coverage sets coverage enabled", func(t *T) {
		opts, _ := preprocess([]string{"--coverage"})
		t.Ok(opts.enabled)
		t.End()
	})
}

func TestPreprocessCFlagDefaultFormat(t *testing.T) {
	Test(t, "preprocess: -c defaults format to lines", func(t *T) {
		opts, _ := preprocess([]string{"-c"})
		t.Equal(opts.format, "lines")
		t.End()
	})
}

func TestPreprocessCFlagEqualFormat(t *testing.T) {
	Test(t, "preprocess: -c=json-lines sets format", func(t *T) {
		opts, _ := preprocess([]string{"-c=json-lines"})
		t.Equal(opts.format, "json-lines")
		t.End()
	})
}

func TestPreprocessCFlagSpaceFormat(t *testing.T) {
	Test(t, "preprocess: -c code-frame sets format", func(t *T) {
		opts, _ := preprocess([]string{"-c", "code-frame"})
		t.Equal(opts.format, "code-frame")
		t.End()
	})
}

func TestPreprocessCoverageEqualFormat(t *testing.T) {
	Test(t, "preprocess: --coverage=json-lines sets format", func(t *T) {
		opts, _ := preprocess([]string{"--coverage=json-lines"})
		t.Equal(opts.format, "json-lines")
		t.End()
	})
}

func TestPreprocessRFlag(t *testing.T) {
	Test(t, "preprocess: -r sets report enabled", func(t *T) {
		opts, _ := preprocess([]string{"-r"})
		t.Ok(opts.report)
		t.End()
	})
}

func TestPreprocessRFlagDefaultPath(t *testing.T) {
	Test(t, "preprocess: -r defaults path to coverage.lcov", func(t *T) {
		opts, _ := preprocess([]string{"-r"})
		t.Equal(opts.path, "coverage.lcov")
		t.End()
	})
}

func TestPreprocessRFlagEqualPath(t *testing.T) {
	Test(t, "preprocess: -r=lcov.info sets path", func(t *T) {
		opts, _ := preprocess([]string{"-r=lcov.info"})
		t.Equal(opts.path, "lcov.info")
		t.End()
	})
}

func TestPreprocessReportFlagEqualPath(t *testing.T) {
	Test(t, "preprocess: --report=out.lcov sets path", func(t *T) {
		opts, _ := preprocess([]string{"--report=out.lcov"})
		t.Equal(opts.path, "out.lcov")
		t.End()
	})
}

func TestPreprocessCRCombined(t *testing.T) {
	Test(t, "preprocess: -c -r sets both enabled and report", func(t *T) {
		opts, _ := preprocess([]string{"-c", "-r"})
		t.Ok(opts.enabled && opts.report)
		t.End()
	})
}

func TestPreprocessCoverageReportLong(t *testing.T) {
	Test(t, "preprocess: --coverage --report sets both", func(t *T) {
		opts, _ := preprocess([]string{"--coverage", "--report"})
		t.Ok(opts.enabled && opts.report)
		t.End()
	})
}

func TestPreprocessCFlagDoesNotConsumeRFlag(t *testing.T) {
	Test(t, "preprocess: -c does not eat -r as format value", func(t *T) {
		opts, _ := preprocess([]string{"-c", "-r"})
		t.Equal(opts.format, "lines")
		t.End()
	})
}

func TestPreprocessRFlagDoesNotConsumePath(t *testing.T) {
	Test(t, "preprocess: bare -r does not consume next arg as path", func(t *T) {
		_, rest := preprocess([]string{"-r", "./pkg/..."})
		t.Equal(rest[0], "./pkg/...")
		t.End()
	})
}

func TestPreprocessRestPassthrough(t *testing.T) {
	Test(t, "preprocess: unknown flags pass through to rest", func(t *T) {
		_, rest := preprocess([]string{"-f", "tap", "-c", "-r"})
		t.Equal(rest[0], "-f")
		t.End()
	})
}

func TestPreprocessRestLength(t *testing.T) {
	Test(t, "preprocess: -f tap produces two rest args", func(t *T) {
		_, rest := preprocess([]string{"-f", "tap", "-c", "-r"})
		t.Equal(len(rest), 2)
		t.End()
	})
}

func TestPreprocessNoCoverageFlags(t *testing.T) {
	Test(t, "preprocess: no coverage flags leaves enabled false", func(t *T) {
		opts, _ := preprocess([]string{"-f", "tap"})
		t.NotOk(opts.enabled)
		t.End()
	})
}

func TestDirFromPatternDotDotDot(t *testing.T) {
	tape.Test(t, "preprocess: ./... maps to .", func(t *tape.T) {
		t.Equal(dirFromPattern("./..."), ".")
		t.End()
	})
}

func TestDirFromPatternDot(t *testing.T) {
	tape.Test(t, "preprocess: . maps to .", func(t *tape.T) {
		t.Equal(dirFromPattern("."), ".")
		t.End()
	})
}

func TestDirFromPatternSubpackage(t *testing.T) {
	tape.Test(t, "preprocess: ./internal/foo maps to ./internal/foo", func(t *tape.T) {
		t.Equal(dirFromPattern("./internal/foo"), "./internal/foo")
		t.End()
	})
}

func TestDirFromPatternSubpackageRecursive(t *testing.T) {
	tape.Test(t, "preprocess: ./internal/foo/... maps to ./internal/foo", func(t *tape.T) {
		t.Equal(dirFromPattern("./internal/foo/..."), "./internal/foo")
		t.End()
	})
}
