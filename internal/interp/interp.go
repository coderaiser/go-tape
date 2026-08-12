package interp

import (
	"bytes"
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
	cb   func(Event)
	out  *bytes.Buffer
	mu   sync.Mutex
}

func (h *handler) emit(e Event) {
	// h.cb is set before Run starts; the interpreter runs synchronously so no
	// concurrent access beyond the mutex guard used by RegisterExternal.
	if h.cb != nil {
		h.cb(e)
	}
}

// nativeT implements host T; its methods emit assertion events.
type nativeT struct {
	h *handler
}

func (tt nativeT) Equal(a, b any) {
	tt.h.emit(Event{Kind: "assert", Name: "equal", Ok: eq(a, b), Got: a, Wanted: b})
}

func (tt nativeT) EqualText(a, b any, msg string) {
	tt.h.emit(Event{Kind: "assert", Name: "equal", Ok: eq(a, b), Got: a, Wanted: b, Msg: msg})
}

func (tt nativeT) DeepEqual(a, b any) {
	tt.h.emit(Event{Kind: "assert", Name: "deepEqual", Ok: deepEq(a, b), Got: a, Wanted: b})
}

func (tt nativeT) Ok(a any) {
	tt.h.emit(Event{Kind: "assert", Name: "ok", Ok: truthy(a), Got: a})
}

func (tt nativeT) NotOk(a any) {
	tt.h.emit(Event{Kind: "assert", Name: "notOk", Ok: !truthy(a), Got: a})
}

func (tt nativeT) End() {
	tt.h.emit(Event{Kind: "end"})
}

// Run interprets a plain `package main` supertape-style source that uses the
// tape API (Test/T/Report) via `import . "tapeapi"`, emitting structured events
// to the given callback. It returns the interpreter exit code. No *testing.T,
// no go test process, and none of ixgo's testing machinery is used.
func Run(src string, w *bytes.Buffer, emit func(Event)) (code int, err error) {
	h := &handler{cb: emit, out: w}

	ctx := ixgo.NewContext(0)
	if w != nil {
		ctx.SetPrintOutput(w)
	}

	// Point the tapeapi vars at this run's handler (ctx override wins over the
	// globally-registered Vars used for type resolution).
	ctx.RegisterExternal(apiPath+".Test", &testFn{h})
	ctx.RegisterExternal(apiPath+".Report", &reportFn{h})

	return ctx.RunFile("main.go", src, nil)
}
