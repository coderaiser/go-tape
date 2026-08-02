//go:build !no_external

package statemachine_test

import (
	"testing"

	"github.com/coderaiser/go-tape/internal/statemachine"
)

func TestFileSourceTOMLLoadsBurntSushi(t *testing.T) {
	MachineTest(t, "statemachine: FileSource loads TOML via BurntSushi", func(t *MachineT) {
		src := statemachine.FileSource{Path: "testdata/runner.toml"}
		defs, error := src.Load()
		if error != nil {
			t.TB().Fatalf("Load: %v", error)
		}
		t.Ok(len(defs) > 0)
		t.End()
	})
}
