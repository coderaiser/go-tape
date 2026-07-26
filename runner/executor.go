package runner

import (
	"io"
	"os/exec"
)

// Executor runs a command and returns its stdout.
type Executor interface {
	Run(args ...string) (io.ReadCloser, error)
}

// OSExecutor runs commands using os/exec.
type OSExecutor struct{}

func (e *OSExecutor) Run(args ...string) (io.ReadCloser, error) {
	cmd := exec.Command("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return stdout, nil
}
