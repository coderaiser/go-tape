package statemachine_test

import (
	"testing"

	"github.com/coderaiser/go-tape/internal/statemachine"
)

func TestValidateEmptyFrom(t *testing.T) {
	MachineTest(t, "statemachine: Validate errors on empty From", func(t *MachineT) {
		error := statemachine.Validate([]statemachine.TransitionDef{{From: "", Event: "run", To: "running"}})
		t.Equal(error.Error(), "transition has empty From")
		t.End()
	})
}

func TestValidateEmptyEvent(t *testing.T) {
	MachineTest(t, "statemachine: Validate errors on empty Event", func(t *MachineT) {
		error := statemachine.Validate([]statemachine.TransitionDef{{From: "idle", Event: "", To: "running"}})
		t.Equal(error.Error(), "transition has empty Event")
		t.End()
	})
}

func TestValidateEmptyTo(t *testing.T) {
	MachineTest(t, "statemachine: Validate errors on empty To", func(t *MachineT) {
		error := statemachine.Validate([]statemachine.TransitionDef{{From: "idle", Event: "run", To: ""}})
		t.Equal(error.Error(), "transition has empty To")
		t.End()
	})
}

func TestValidatePass(t *testing.T) {
	MachineTest(t, "statemachine: Validate passes for complete definitions", func(t *MachineT) {
		error := statemachine.Validate([]statemachine.TransitionDef{{From: "idle", Event: "run", To: "running"}})
		t.NotOk(error)
		t.End()
	})
}
