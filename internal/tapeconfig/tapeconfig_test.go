package tapeconfig_test

import (
	"os"
	"path/filepath"
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
		t.DeepEqual(tapeconfig.Load("/nonexistent/path/.tape.toml"), tapeconfig.Default())
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
		t.DeepEqual(tapeconfig.Load(path), expected)
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
		t.DeepEqual(tapeconfig.Load(path), expected)
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
		t.DeepEqual(tapeconfig.Load(path), expected)
		t.End()
	})
}

func TestLoadMalformedReturnsDefaults(t *testing.T) {
	Test(t, "tapeconfig: malformed toml returns defaults with warning", func(t *T) {
		dir := t.TB().TempDir()
		path := filepath.Join(dir, ".tape.toml")
		writeTapeconfigFile(t.TB(), path, `[this is not valid toml
`)
		t.DeepEqual(tapeconfig.Load(path), tapeconfig.Default())
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
