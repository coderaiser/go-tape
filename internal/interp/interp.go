package interp

import (
	"bytes"
	"reflect"
	"sync"

	"github.com/goplus/ixgo"
)

// Event is a structured event emitted while interpreting a test source.
type Event struct {
	Kind   string // "test" | "assert" | "end" | "report"
	Name   string // test name (Kind == "test") or assertion operator
	Ok     bool
	Got    any
	Wanted any
	Msg    string
}

// Handler receives events produced during interpretation.
type Handler interface {
	OnEvent(Event)
}

// handler holds the per-run event sink plus print output.
type handler struct {
	cb  func(Event)
	out *bytes.Buffer
	mu  sync.Mutex
}

func (h *handler) emit(e Event) {
	// h.cb is set before Run starts; the interpreter runs synchronously so no
	// concurrent access beyond the mutex guard used by RegisterExternal.
	if h.cb != nil {
		h.cb(e)
	}
}

// runTest is the host implementation of the tapeapi `Test` var for a single
// run. It emits a "test" event, runs the interpreted body against a fresh
// nativeT (whose assertion methods emit "assert" events), then emits an "end"
// event. The "end" event carries the accumulated pass/fail state.
func (h *handler) runTest(name string, fn func(T)) {
	h.emit(Event{Kind: "test", Name: name})
	tt := nativeT{h: h}
	failed := false
	tt.fail = &failed
	fn(tt)
	h.emit(Event{Kind: "end", Name: name, Ok: !failed})
}

// runReport is the host implementation of the tapeapi `Report` var. It emits a
// "report" event carrying the structured assertion result.
func (h *handler) runReport(r Result) {
	h.emit(Event{Kind: "report", Name: r.Operator, Ok: r.Ok, Got: r.Got, Wanted: r.Expected, Msg: r.Message})
}

// nativeT implements host T; its methods emit assertion events.
type nativeT struct {
	h    *handler
	fail *bool
}

func (tt nativeT) Equal(a, b any) {
	ok := eq(a, b)
	if !ok && tt.fail != nil {
		*tt.fail = true
	}
	tt.h.emit(Event{Kind: "assert", Name: "equal", Ok: ok, Got: a, Wanted: b})
}

func (tt nativeT) EqualText(a, b any, msg string) {
	ok := eq(a, b)
	if !ok && tt.fail != nil {
		*tt.fail = true
	}
	tt.h.emit(Event{Kind: "assert", Name: "equal", Ok: ok, Got: a, Wanted: b, Msg: msg})
}

func (tt nativeT) DeepEqual(a, b any) {
	ok := deepEq(a, b)
	if !ok && tt.fail != nil {
		*tt.fail = true
	}
	tt.h.emit(Event{Kind: "assert", Name: "deepEqual", Ok: ok, Got: a, Wanted: b})
}

func (tt nativeT) Ok(a any) {
	ok := truthy(a)
	if !ok && tt.fail != nil {
		*tt.fail = true
	}
	tt.h.emit(Event{Kind: "assert", Name: "ok", Ok: ok, Got: a})
}

func (tt nativeT) NotOk(a any) {
	ok := !truthy(a)
	if !ok && tt.fail != nil {
		*tt.fail = true
	}
	tt.h.emit(Event{Kind: "assert", Name: "notOk", Ok: ok, Got: a})
}

func (tt nativeT) End() {
	tt.h.emit(Event{Kind: "end"})
}

// eq reports strict shallow equality.
func eq(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// deepEq reports deep equality.
func deepEq(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// truthy mirrors operator.truthy semantics for the simple native assertions.
func truthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != ""
	case int:
		return val != 0
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan,
		reflect.Func, reflect.Pointer, reflect.Interface:
		return !rv.IsNil()
	}
	return true
}

// Run interprets a plain `package main` supertape-style source that uses the
// tape API (Test/T/Report) via `import . "tapeapi"`, emitting structured events
// to the given callback. It returns the interpreter exit code. No *testing.T,
// no go test process, and none of ixgo's testing machinery is used.
//
// ixgo resolves tapeapi global vars to the package-level defaultTest/
// defaultReport functions registered in the Package.Vars map, so per-run event
// routing is done through the package-level `current` handler (Runs are
// sequential by design).
func Run(src string, w *bytes.Buffer, emit func(Event)) (code int, err error) {
	h := &handler{cb: emit, out: w}

	ctx := ixgo.NewContext(0)
	if w != nil {
		ctx.SetPrintOutput(w)
	}

	current = h
	defer func() { current = nil }()

	return ctx.RunFile("main.go", src, nil)
}
