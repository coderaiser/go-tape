package tapeconfig_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/tapeconfig"
)

func writeTapeconfigFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	Test(t, "tapeconfig: missing file returns defaults", func(t *T) {
		t.DeepEqual(tapeconfig.Load("/nonexistent/path"), tapeconfig.Default())
		t.End()
	})
}

func TestLoadMissingFileNoWarning(t *testing.T) {
	Test(t, "tapeconfig: missing file does not print warning to stderr", func(t *T) {
		// Capture stderr by redirecting os.Stderr temporarily.
		r, w, _ := os.Pipe()
		old := os.Stderr
		os.Stderr = w

		tapeconfig.Load("/nonexistent/path")

		w.Close()
		os.Stderr = old

		var buf strings.Builder
		io.Copy(&buf, r)

		t.Equal(buf.String(), "")
		t.End()
	})
}

func TestLoadValidFile(t *testing.T) {
	Test(t, "tapeconfig: valid file is parsed correctly", func(t *T) {
		dir := t.TB().TempDir()
		path := filepath.Join(dir, ".tape.toml")
		writeTapeconfigFile(t.TB(), path, `[test]
formatter = "tap"
exclude = ["vendor"]

[coverage]
exclude = ["docs"]
`)
		var expected tapeconfig.Config
		expected.Test.Formatter = "tap"
		expected.Test.Exclude = []string{"vendor"}
		expected.Coverage.Exclude = []string{"docs"}
		t.DeepEqual(tapeconfig.Load(dir), expected)
		t.End()
	})
}

func TestLoadPartialFileKeepsDefaults(t *testing.T) {
	Test(t, "tapeconfig: partial file keeps defaults for absent keys", func(t *T) {
		dir := t.TB().TempDir()
		path := filepath.Join(dir, ".tape.toml")
		writeTapeconfigFile(t.TB(), path, `[coverage]
exclude = ["docs"]
`)
		expected := tapeconfig.Default()
		expected.Coverage.Exclude = []string{"docs"}
		t.DeepEqual(tapeconfig.Load(dir), expected)
		t.End()
	})
}

func TestLoadUnknownKeysIgnored(t *testing.T) {
	Test(t, "tapeconfig: unknown keys are ignored (no panic)", func(t *T) {
		dir := t.TB().TempDir()
		path := filepath.Join(dir, ".tape.toml")
		writeTapeconfigFile(t.TB(), path, `[test]
formatter = "tap"
bogus = "value"
`)
		expected := tapeconfig.Default()
		expected.Test.Formatter = "tap"
		t.DeepEqual(tapeconfig.Load(dir), expected)
		t.End()
	})
}

func TestLoadMalformedReturnsDefaults(t *testing.T) {
	Test(t, "tapeconfig: malformed toml returns defaults with warning", func(t *T) {
		dir := t.TB().TempDir()
		path := filepath.Join(dir, ".tape.toml")
		writeTapeconfigFile(t.TB(), path, `[this is not valid toml
`)
		t.DeepEqual(tapeconfig.Load(dir), tapeconfig.Default())
		t.End()
	})
}

func TestDefaultReturnsExpected(t *testing.T) {
	Test(t, "tapeconfig: Default() returns expected values", func(t *T) {
		var expected tapeconfig.Config
		expected.Test.Formatter = "progress-bar"
		expected.Test.Exclude = []string{"fixture"}
		expected.Coverage.Exclude = []string{"node_modules"}
		t.DeepEqual(tapeconfig.Default(), expected)
		t.End()
	})
}

func TestLoadDotPrefix(t *testing.T) {
	Test(t, "tapeconfig: Load reads .tape.toml when present", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".tape.toml"), []byte("[test]\nformatter = \"json\"\n"), 0644)
		cfg := tapeconfig.Load(dir)
		t.Equal(cfg.Test.Formatter, "json")
		t.End()
	})
}

func TestLoadNoDotFallback(t *testing.T) {
	Test(t, "tapeconfig: Load falls back to tape.toml when no dot-prefix file", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, "tape.toml"), []byte("[test]\nformatter = \"json\"\n"), 0644)
		cfg := tapeconfig.Load(dir)
		t.Equal(cfg.Test.Formatter, "json")
		t.End()
	})
}

func TestLoadMissingUsesDefault(t *testing.T) {
	Test(t, "tapeconfig: Load returns Default when no config file present", func(t *T) {
		cfg := tapeconfig.Load(t.TB().TempDir())
		t.Equal(cfg.Test.Formatter, "progress-bar")
		t.End()
	})
}

func TestLoadDotPrefixTakesPriority(t *testing.T) {
	Test(t, "tapeconfig: .tape.toml takes priority over tape.toml", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".tape.toml"), []byte("[test]\nformatter = \"tap\"\n"), 0644)
		os.WriteFile(filepath.Join(dir, "tape.toml"), []byte("[test]\nformatter = \"json\"\n"), 0644)
		cfg := tapeconfig.Load(dir)
		t.Equal(cfg.Test.Formatter, "tap")
		t.End()
	})
}
