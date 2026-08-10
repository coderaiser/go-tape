package stream

import (
	"os"
	"os/exec"
	"testing"
)

// TestStreamRunStdoutPipeError covers the StdoutPipe failure branch of Run,
// which is unreachable through the public API. It swaps the execCommand seam
// (declared in stream.go) for a command whose Stdout is already set.
func TestStreamRunStdoutPipeError(t *testing.T) {
	old := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		cmd := exec.Command(name, args...)
		cmd.Stdout = os.Stderr // pre-set → StdoutPipe must fail
		return cmd
	}
	defer func() { execCommand = old }()

	_, err := Run(0, "version")
	if err == nil {
		t.Fatal("expected StdoutPipe error")
	}
}
