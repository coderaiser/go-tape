package ast_test

import (
	"os"
	"path/filepath"
	"testing"

	dedent "github.com/lithammer/dedent"
)

// Fixture creates an isolated temporary directory and returns:
// - the directory path
// - a helper for creating fixture files inside it.
//
// The directory is automatically removed by testing.T after the test.
func Fixture(t *testing.T) (string, func(name, src string)) {
	t.Helper()

	dir := t.TempDir()

	write := func(name, src string) {
		t.Helper()

		path := filepath.Join(dir, name)

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(
			path,
			[]byte(dedent.Dedent(src)),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	return dir, write
}
