// Package supertape provides a supertape-style assertion API that does NOT
// depend on *testing.T. It is a normal Go library and a clean interpretation
// target for the ixgo interpreter in --super mode.
//
// The concrete assertion type T is a plain interface, so test bodies written
// against it can be interpreted without a Go test binary. Assertions reuse
// internal/operator, and each emitted event mirrors internal/stream.Event.
package supertape

import (
	"sync"
	"time"

	"github.com/coderaiser/go-tape/internal/config"
	"github.com/coderaiser/go-tape/internal/operator"
	"github.com/coderaiser/go-tape/internal/scope"
)

// Event type constants, mirroring stream.Event's typed event protocol.
const (
	TypeTestEnd = "test-end"
	TypeFail    = "fail"
)

// T is the assertion interface. It intentionally holds no *testing.T field so
// it can be implemented and interpreted independently of the Go test runner.
type T interface {
	Equal(result, expected any)
	NotEqual(result, expected any)
	DeepEqual(result, expected any)
	NotDeepEqual(result, expected any)
	Ok(result any)
	NotOk(result any)
	Match(result string, pattern any)
	NotMatch(result string, pattern any)
	Pass(message ...string)
	Fail(message any)
	Comment(message string)
	End()
}

// Event is a single structured event emitted per test, mirroring the fields of
// stream.Event used for test-end and fail events.
type Event struct {
	Type     string
	Test     string
	Message  string
	Operator string
	Result   any
	Expected any
	Output   string
	At       string
}

// Option configures a Runner.
type Option func(*Runner)

// WithCheckAssertionsCount enables/disables the one-assertion-per-test guard.
func WithCheckAssertionsCount(b bool) Option {
	return func(r *Runner) { r.checkAssertionsCount = b }
}

// WithCheckEnd enables/disables the requirement that t.End() be called.
func WithCheckEnd(b bool) Option {
	return func(r *Runner) { r.checkEnd = b }
}

// WithCheckScopes enables/disables the 'scope: message' name check.
func WithCheckScopes(b bool) Option {
	return func(r *Runner) { r.checkScopes = b }
}

// WithTimeout sets the per-test timeout.
func WithTimeout(d time.Duration) Option {
	return func(r *Runner) { r.timeout = d }
}

// Runner records pass/fail/skip counts and emits one structured event per test.
type Runner struct {
	mu      sync.Mutex
	passed  int
	failed  int
	skipped int
	total   int

	checkAssertionsCount bool
	checkEnd             bool
	checkScopes          bool
	timeout              time.Duration
	handler              func(Event)

	currentHits   int
	currentFailed bool
}

// New returns a Runner using the internal/config defaults, overridable by opts.
func New(opts ...Option) *Runner {
	r := &Runner{
		checkAssertionsCount: config.CheckAssertionsCount(),
		checkEnd:             config.CheckEnd(),
		checkScopes:          config.CheckScopes(),
		timeout:              config.Timeout(),
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Counts returns (passed, failed, skipped, total).
func (r *Runner) Counts() (int, int, int, int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.passed, r.failed, r.skipped, r.total
}

// Skip records a skipped test without running a body.
func (r *Runner) Skip(name string) {
	r.mu.Lock()
	r.total++
	r.skipped++
	r.mu.Unlock()
	r.emit(Event{Type: TypeTestEnd, Test: name})
}

// Test runs a single test under the runner's guards: scope check, timeout,
// assertion count, and End() requirement. fn runs in a goroutine so the
// timeout can be enforced with select; the call blocks until the test settles.
func (r *Runner) Test(name string, fn func(T)) {
	r.mu.Lock()
	r.total++
	r.currentHits = 0
	r.currentFailed = false
	r.mu.Unlock()

	// guard 1: scope
	if r.checkScopes && !scope.Valid(name) {
		r.recordFail(Event{
			Type:     TypeFail,
			Test:     name,
			Message:  "invalid scope name",
			Operator: "scope",
			Output:   "expected 'scope: message'",
		})
		r.emitTestEnd(name)
		return
	}

	tt := &t{runner: r, name: name}

	done := make(chan struct{})
	go func() {
		defer close(done)
		fn(tt)
	}()

	// guard 2: timeout
	select {
	case <-done:
	case <-time.After(r.timeout):
		tt.setEnded(true)
		r.recordFail(Event{
			Type:     TypeFail,
			Test:     name,
			Message:  "test timed out",
			Operator: "timeout",
			Output:   "test timed out after " + r.timeout.String(),
		})
	}

	// guard 3: End() requirement
	if r.checkEnd && !tt.ended() {
		r.recordFail(Event{
			Type:     TypeFail,
			Test:     name,
			Message:  "t.End() not called",
			Operator: "end",
			Output:   "t.End() not called",
		})
	}

	r.mu.Lock()
	if !r.currentFailed {
		r.passed++
	}
	r.mu.Unlock()

	r.emitTestEnd(name)
}

// WithHandler sets the event sink that receives every emitted Event.
func WithHandler(h func(Event)) Option {
	return func(r *Runner) { r.handler = h }
}


// emitTestEnd emits the closing test-end event for the given test name.
func (r *Runner) emitTestEnd(name string) {
	r.emit(Event{Type: TypeTestEnd, Test: name})
}

// recordFail marks the current test as failed, bumps the failed count and
// emits the fail event.
func (r *Runner) recordFail(ev Event) {
	r.mu.Lock()
	r.failed++
	r.currentFailed = true
	r.mu.Unlock()
	r.emit(ev)
}

// emit forwards ev to the configured handler, if any.
func (r *Runner) emit(ev Event) {
	if r.handler != nil {
		r.handler(ev)
	}
}

// report processes a single assertion Result from the concrete T.
func (r *Runner) report(tt *t, res operator.Result) {
	r.mu.Lock()
	r.currentHits++
	hits := r.currentHits
	r.mu.Unlock()

	// guard: one assertion per test
	if hits > 1 && r.checkAssertionsCount {
		r.recordFail(Event{
			Type:     TypeFail,
			Test:     tt.name,
			Message:  "too many assertions: got more than 1",
			Operator: res.Operator,
			Result:   res.Result,
			Expected: res.Expected,
			Output:   res.Output,
			At:       res.At,
		})
		return
	}

	if !res.Ok {
		r.recordFail(Event{
			Type:     TypeFail,
			Test:     tt.name,
			Message:  res.Message,
			Operator: res.Operator,
			Result:   res.Result,
			Expected: res.Expected,
			Output:   res.Output,
			At:       res.At,
		})
	}
}

// defaultRunner backs the package-level Test function.
var defaultRunner = New()

// Test is the top-level entry point: it runs fn against a fresh T via the
// default runner. It has no *testing.T parameter, so it is a clean ixgo
// interpretation target.
func Test(name string, fn func(T)) {
	defaultRunner.Test(name, fn)
}
