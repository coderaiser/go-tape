package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coderaiser/go-tape/internal/formatter"
	"github.com/coderaiser/go-tape/internal/interp"
	"github.com/coderaiser/go-tape/internal/stream"
)

// supertapeSuffix is the file-name convention for supertape sources.
const supertapeSuffix = "_tape.go"

// isSupertapeFile reports whether src is a supertape-style package-main source
// that imports the virtual tapeapi package.
// fileSource pairs a discovered supertape file with its (already-read) source
// text, so each file is read exactly once.
type fileSource struct {
	path string
	src  string
}

func isSupertapeFile(src string) bool {
	return strings.Contains(src, `import . "tapeapi"`) ||
		strings.Contains(src, `import "tapeapi"`)
}

// findSupertapeFiles walks dir (which may be a directory or a single file) and
// returns the supertape source files it finds, along with their source text.
// Only *_tape.go files are considered supertape sources; the tapeapi import
// confirms it.
func findSupertapeFiles(dir string, exclude []string) ([]fileSource, error) {
	var files []fileSource
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, supertapeSuffix) {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if isSupertapeFile(string(src)) {
			files = append(files, fileSource{path: path, src: string(src)})
		}
		return nil
	})
	_ = exclude
	return files, err
}

// interpState tracks the current test while converting interp events to stream
// events.
type interpState struct {
	current string
	failed  bool
}

// interpToStream converts a single interp.Event into zero or more stream
// events. A failing assertion emits a stream.TypeFail; the final end event for
// a test emits a stream.TypeTestEnd.
func interpToStream(ev interp.Event, st *interpState) []stream.Event {
	switch ev.Kind {
	case "test":
		st.current = ev.Name
		st.failed = false
		return nil
	case "assert":
		if !ev.Ok {
			st.failed = true
			return []stream.Event{{
				Type:     stream.TypeFail,
				Test:     st.current,
				Message:  ev.Msg,
				Operator: ev.Name,
				Result:   ev.Got,
				Expected: ev.Wanted,
			}}
		}
		return nil
	case "end":
		// The runTest end event has a Name; the t.End() event does not.
		if ev.Name != "" {
			e := stream.Event{Type: stream.TypeTestEnd, Test: ev.Name}
			if !ev.Ok {
				e.Type = stream.TypeFail
			}
			return []stream.Event{e}
		}
		return nil
	}
	return nil
}

// runInterpMode runs all supertape sources under dir through the interpreter,
// feeding the resulting stream events to formatter d. It returns the process
// exit code (0 = pass, 1 = any failure or error).
func runInterpMode(dir string, d *formatter.Dispatcher, stdout io.Writer) int {
	files, err := findSupertapeFiles(dir, nil)
	if err != nil {
		_, _ = io.WriteString(stdout, "tape: scan supertape files: "+err.Error()+"\n")
		return 1
	}
	if len(files) == 0 {
		_, _ = io.WriteString(stdout, "tape: no supertape files found\n")
		return 1
	}

	passed, failed := 0, 0
	for _, f := range files {
		var st interpState
		_, err := interp.Run(f.src, nil, func(ev interp.Event) {
			for _, se := range interpToStream(ev, &st) {
				if se.Type == stream.TypeFail {
					failed++
				} else {
					passed++
				}
				d.Emit(se)
			}
		})
		if err != nil {
			_, _ = io.WriteString(stdout, "tape: interpret "+f.path+": "+err.Error()+"\n")
			return 1
		}
	}

	d.End(passed, failed, 0)
	if failed > 0 {
		return 1
	}
	return 0
}
