package runner

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/model"
)

func TestRunnerReturnsEventsFromFakeExecutor(t *testing.T) {
	tape.Test(t, "runner: returns run and pass events", func(t *tape.T) {
		lines := []string{
			`{"Action":"run","Package":"mypkg","Test":"TestFoo"}`,
			`{"Action":"pass","Package":"mypkg","Test":"TestFoo","Elapsed":0.1}`,
		}
		executor := &FakeExecutor{Lines: lines}
		r := New(executor)

		ch, err := r.Run("test", "-json")

		var events []model.Event
		for e := range ch {
			events = append(events, e)
		}

		t.Ok(err == nil && len(events) == 2 &&
			events[0].Action == "run" && events[1].Action == "pass")
		t.End()
	})
}

func TestRunnerSkipsInvalidJSON(t *testing.T) {
	tape.Test(t, "runner: skips invalid json lines", func(t *tape.T) {
		executor := &FakeExecutor{Lines: []string{"not json"}}
		r := New(executor)

		ch, err := r.Run("test", "-json")

		count := 0
		for range ch {
			count++
		}

		t.Ok(err == nil && count == 0)
		t.End()
	})
}

func TestFakeExecutorReturnsLines(t *testing.T) {
	tape.Test(t, "runner: fake executor yields queued lines", func(t *tape.T) {
		e := &FakeExecutor{Lines: []string{"line1", "line2"}}
		rc, _ := e.Run()
		buf := make([]byte, 100)
		n, _ := rc.Read(buf)
		t.Equal(string(buf[:n]), "line1\n")
		t.End()
	})
}

func TestOSExecutorRun(t *testing.T) {
	tape.Test(t, "runner: OS executor reads output", func(t *tape.T) {
		// inject echo as the command — no real go test spawned
		e := &OSExecutor{
			Command: func(name string, args ...string) *exec.Cmd {
				return exec.Command("echo", `{"Action":"run","Test":"TestFoo"}`)
			},
		}
		rc, err := e.Run("test", "-json")
		if err == nil {
			defer rc.Close()
		}
		buf := make([]byte, 256)
		n := 0
		if rc != nil {
			n, _ = rc.Read(buf)
		}
		t.Ok(err == nil && rc != nil && n > 0)
		t.End()
	})
}

func TestRunnerClosesChannelOnEmptyInput(t *testing.T) {
	tape.Test(t, "runner: closes channel on empty input", func(t *tape.T) {
		executor := &FakeExecutor{Lines: []string{}}
		r := New(executor)
		ch, err := r.Run("test", "-json")
		count := 0
		for range ch {
			count++
		}
		t.Ok(err == nil && count == 0)
		t.End()
	})
}

func TestOSExecutorStartError(t *testing.T) {
	tape.Test(t, "runner: reports Start error", func(t *tape.T) {
		e := &OSExecutor{
			Command: func(name string, args ...string) *exec.Cmd {
				// nonexistent command — Start() will fail
				return exec.Command("nonexistent-command-xyz")
			},
		}
		_, err := e.Run("test")
		t.Ok(err != nil)
		t.End()
	})
}

func TestOSExecutorConformsToInterface(t *testing.T) {
	tape.Test(t, "runner: OS executor implements Executor", func(t *tape.T) {
		var e Executor = NewOSExecutor()
		t.Ok(e != nil)
		t.End()
	})
}

func TestOSExecutorStdoutPipeError(t *testing.T) {
	tape.Test(t, "runner: reports StdoutPipe error", func(t *tape.T) {
		e := &OSExecutor{
			Command: func(name string, args ...string) *exec.Cmd {
				cmd := exec.Command("echo", "hello")
				// pre-assign stdout — StdoutPipe will fail
				cmd.Stdout = os.Stdout
				return cmd
			},
		}
		_, err := e.Run("test")
		t.Ok(err != nil)
		t.End()
	})
}

// errReader returns an error after the first read
type errReader struct{ done bool }

func (r *errReader) Read(p []byte) (int, error) {
	if r.done {
		return 0, errors.New("read error")
	}
	r.done = true
	copy(p, []byte("not json\n"))
	return 9, nil
}

func (r *errReader) Close() error { return nil }

// errExecutor returns an errReader
type errExecutor struct{}

func (e errExecutor) Run(args ...string) (io.ReadCloser, error) {
	return &errReader{}, nil
}

func TestRunnerHandlesScannerError(t *testing.T) {
	tape.Test(t, "runner: survives scanner error", func(t *tape.T) {
		r := New(errExecutor{})
		ch, err := r.Run("test", "-json")
		if err == nil {
			// drain channel — scanner error causes goroutine to exit
			for range ch {
			}
		}
		t.Ok(err == nil)
		t.End()
	})
}

// errExecutorRunError always fails
type errExecutorRunError struct{}

func (e errExecutorRunError) Run(args ...string) (io.ReadCloser, error) {
	return nil, errors.New("executor failed")
}

func TestRunnerExecutorError(t *testing.T) {
	tape.Test(t, "runner: reports executor error", func(t *tape.T) {
		r := New(errExecutorRunError{})
		_, err := r.Run("test")
		t.Ok(err != nil)
		t.End()
	})
}
