package stream

import (
	"bufio"
	"encoding/json"
	"io"
	"os/exec"
	"strings"
)

// Event type constants (kebab-case, matching json-lines protocol).
const (
	TypeTestEnd     = "test-end"
	TypeFail        = "fail"
	TypeUnknownFail = "unknown-fail"
	TypeComment     = "comment"
	TypeBuildError  = "build-error"
	TypeEnd         = "end"
)

// Event is a single typed event emitted by the stream parser.
type Event struct {
	Type string `json:"type"`

	// test-end + fail + unknown-fail + comment
	Test   string `json:"test,omitempty"`
	Count  int    `json:"count,omitempty"`
	Total  int    `json:"total,omitempty"`
	Failed int    `json:"failed,omitempty"`

	// fail fields (from TAPE: sentinel)
	Message    string `json:"message,omitempty"`
	Operator   string `json:"operator,omitempty"`
	Result     any    `json:"result,omitempty"`
	Expected   any    `json:"expected,omitempty"`
	Output     string `json:"output,omitempty"`
	At         string `json:"at,omitempty"`
	ErrorStack string `json:"error-stack,omitempty"`

	// build-error (also uses Output)
	Package string `json:"package,omitempty"`
}

// goEvent is the raw go test -json line shape.
type goEvent struct {
	Action  string
	Package string
	Test    string
	Output  string
	Elapsed float64
}

// tapeLine is the JSON payload after the TAPE: prefix.
type tapeLine struct {
	Message  string `json:"message"`
	Operator string `json:"operator"`
	Result   any    `json:"result"`
	Expected any    `json:"expected"`
	Output   string `json:"output"`
}

// execCommand is overridable in tests to inject Command failures.
var execCommand = exec.Command

// Run executes go test -json [args] and returns a channel of Events.
// The channel is closed when the process exits.
// total is the pre-counted number of tape tests (from AST scan).
func Run(total int, args ...string) (<-chan Event, error) {
	cmd := execCommand("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	ch := make(chan Event)
	go parse(stdout, total, ch)
	return ch, nil
}

// Parse reads a go test -json stream and emits typed Events.
// Pure — no exec. Used directly in tests.
func Parse(r io.Reader, total int) <-chan Event {
	ch := make(chan Event)
	go parse(r, total, ch)
	return ch
}

// parse is the shared goroutine body.
func parse(r io.Reader, total int, ch chan<- Event) {
	defer close(ch)

	// per-test output accumulation
	testOutputs := make(map[string][]string)
	// package-level output for build-error detection
	pkgOutputs := make(map[string][]string)

	count := 0
	failed := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var ge goEvent
		if err := json.Unmarshal([]byte(scanner.Text()), &ge); err != nil {
			continue
		}

		if ge.Test == "" {
			// package-level event
			switch ge.Action {
			case "output":
				pkgOutputs[ge.Package] = append(pkgOutputs[ge.Package], ge.Output)
			case "fail":
				lines := pkgOutputs[ge.Package]
				raw := strings.Join(lines, "")
				if strings.Contains(raw, "[build failed]") {
					ch <- Event{
						Type:    TypeBuildError,
						Package: ge.Package,
						Output:  raw,
					}
				}
				delete(pkgOutputs, ge.Package)
			}
			continue
		}

		// outer TestXxx wrapper — no slash means it is not a subtest
		if !strings.Contains(ge.Test, "/") {
			continue
		}

		switch ge.Action {
		case "output":
			testOutputs[ge.Test] = append(testOutputs[ge.Test], ge.Output)

		case "run":
			// initialise buffer (idempotent)
			if _, ok := testOutputs[ge.Test]; !ok {
				testOutputs[ge.Test] = nil
			}

		case "pass", "fail", "skip":
			count++
			label := decodeTestName(ge.Test)
			outputs := testOutputs[ge.Test]
			delete(testOutputs, ge.Test)

			if ge.Action == "fail" {
				tape, at := parseTapeLine(outputs)
				if tape != nil {
					failed++
					ch <- Event{
						Type:     TypeFail,
						Test:     label,
						Count:    count,
						Message:  tape.Message,
						Operator: tape.Operator,
						Result:   tape.Result,
						Expected: tape.Expected,
						Output:   tape.Output,
						At:       at,
					}
				} else {
					failed++
					ch <- Event{
						Type:   TypeUnknownFail,
						Test:   label,
						Count:  count,
						Output: strings.Join(outputs, ""),
					}
				}
			}

			ch <- Event{
				Type:   TypeTestEnd,
				Test:   label,
				Count:  count,
				Total:  total,
				Failed: failed,
			}
		}
	}
}

// decodeTestName converts "TestFoo/scope:_bar_baz" → "scope: bar baz".
func decodeTestName(test string) string {
	if i := strings.LastIndex(test, "/"); i >= 0 {
		test = test[i+1:]
	}
	return strings.ReplaceAll(test, "_", " ")
}

// parseTapeLine scans output lines for a TAPE: JSON sentinel.
// Returns the parsed tapeLine and the at field (file.go:N:) if found.
// go test -json wraps t.Log output with a file:line prefix:
//
//	"    foo_test.go:12: TAPE:{...}\n"
//
// so the prefix is stripped (and remembered as At) before the sentinel is
// matched. Lines without the prefix (plain "    TAPE:{...}") also match.
const tapePrefix = "TAPE:"

func parseTapeLine(outputs []string) (*tapeLine, string) {
	var tl tapeLine
	at := ""
	for _, line := range outputs {
		// strip leading whitespace that go test -json adds around t.Log
		trimmed := strings.TrimSpace(line)
		rest := trimmed
		if at == "" {
			if a := parseAt(trimmed); a != "" {
				at = a
				if i := strings.Index(trimmed, a); i >= 0 {
					rest = strings.TrimSpace(strings.TrimPrefix(trimmed[i+len(a):], ":"))
				}
			}
		}
		if strings.HasPrefix(rest, tapePrefix) {
			payload := rest[len(tapePrefix):]
			if err := json.Unmarshal([]byte(payload), &tl); err == nil {
				return &tl, at
			}
		}
	}
	return nil, ""
}

// parseAt extracts "file.go:N:" from a line like "    foo_test.go:12: ".
func parseAt(line string) string {
	// look for <name>.go:<digits>: pattern
	for _, part := range strings.Fields(line) {
		if strings.Contains(part, ".go:") {
			// strip trailing colon if present
			return strings.TrimRight(part, ":")
		}
	}
	return ""
}
