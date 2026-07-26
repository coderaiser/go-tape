//go:build no_external

package statemachine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSourceTOMLLoadsHandRolled(t *testing.T) {
	src := FileSource{Path: "testdata/runner.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) == 0 {
		t.Fatal("expected transitions from hand-rolled TOML parser")
	}
}

func TestHandRolledTOMLSectionsAndComments(t *testing.T) {
	src := FileSource{Path: "testdata/runner_full.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) != 3 {
		t.Fatalf("expected 3 transitions, got %d", len(defs))
	}
}

func TestHandRolledTOMLMalformedLine(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "malformed.toml", `[transitions.idle]
run = "running"
borked
`)
	src := FileSource{Path: dir + "/malformed.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 transition, got %d", len(defs))
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
