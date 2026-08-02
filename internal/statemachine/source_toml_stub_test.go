//go:build no_external

package statemachine_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/coderaiser/go-tape/internal/statemachine"
)

func TestFileSourceTOMLLoadsHandRolled(t *testing.T) {
	MachineTest(t, "statemachine: FileSource loads TOML via hand-rolled parser", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/runner.toml"}
		defs, error := src.Load()
		if error != nil {
			t.TB().Fatalf("Load: %v", error)
		}
		t.Ok(len(defs) > 0)
		t.End()
	})
}

func TestHandRolledTOMLSectionsAndComments(t *testing.T) {
	MachineTest(t, "statemachine: hand-rolled TOML handles sections and comments", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/runner_full.toml"}
		defs, error := src.Load()
		if error != nil {
			t.TB().Fatalf("Load: %v", error)
		}
		t.Equal(len(defs), 3)
		t.End()
	})
}

func TestHandRolledTOMLMalformedLine(t *testing.T) {
	MachineTest(t, "statemachine: hand-rolled TOML skips malformed lines", func(t *MachineT) {
		dir := t.TB().TempDir()
		writeFixture(t.TB(), dir, "malformed.toml", "[transitions.idle]\nrun = \"running\"\nborked\n")
		src := statemachine.FileSource{Path: dir + "/malformed.toml"}
		defs, error := src.Load()
		if error != nil {
			t.TB().Fatalf("Load: %v", error)
		}
		t.Equal(len(defs), 1)
		t.End()
	})
}

func writeFixture(t *testing.T, dir, name, content string) {
	t.Helper()
	if error := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); error != nil {
		t.Fatal(error)
	}
}

// errReader errors on the second Read call.
type errReader struct{ calls int }

func (r *errReader) Read(p []byte) (int, error) {
	r.calls++
	if r.calls > 1 {
		return 0, errors.New("read error")
	}
	content := "[transitions.idle]\n"
	n := copy(p, content)
	return n, nil
}

func TestParseTOMLReaderScannerError(t *testing.T) {
	MachineTest(t, "statemachine: parseTOMLReader reports scanner error", func(t *MachineT) {
		_, error := statemachine.ParseTOMLReader(&errReader{}, "test")
		t.Ok(error)
		t.End()
	})
}
