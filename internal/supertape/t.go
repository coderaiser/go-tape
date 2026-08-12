package supertape

import (
	"runtime"
	"sync"

	"github.com/coderaiser/go-tape/internal/operator"
)

// t is the concrete implementation of T. It holds no *testing.T field; it
// reports results to the Runner, which records counts and emits events.
type t struct {
	runner *Runner
	name   string

	mu        sync.Mutex
	endedFlag bool
}

// setEnded records whether End() was called, guarded for the timeout path.
func (tt *t) setEnded(v bool) {
	tt.mu.Lock()
	tt.endedFlag = v
	tt.mu.Unlock()
}

// ended reports whether End() was called.
func (tt *t) ended() bool {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	return tt.endedFlag
}

// assert computes a Result via fn, stamps At from the caller, and reports it.
// depth is the number of frames from this helper to the test-body caller.
func (tt *t) assert(res operator.Result) {
	res.At = caller(2)
	tt.runner.report(tt, res)
}

// caller returns the "file.go:line" of the assertion call site.
func caller(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return ""
	}
	return file + ":" + itoa(line)
}

// itoa converts an int to a decimal string without fmt (keeps it light).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// Equal asserts result == expected.
func (tt *t) Equal(result, expected any) {
	tt.assert(operator.Equal(result, expected))
}

// NotEqual asserts result != expected.
func (tt *t) NotEqual(result, expected any) {
	tt.assert(operator.NotEqual(result, expected))
}

// DeepEqual asserts deep equality.
func (tt *t) DeepEqual(result, expected any) {
	tt.assert(operator.DeepEqual(result, expected))
}

// NotDeepEqual asserts values are not deeply equal.
func (tt *t) NotDeepEqual(result, expected any) {
	tt.assert(operator.NotDeepEqual(result, expected))
}

// Ok asserts result is truthy.
func (tt *t) Ok(result any) {
	tt.assert(operator.Ok(result))
}

// NotOk asserts result is falsy.
func (tt *t) NotOk(result any) {
	tt.assert(operator.NotOk(result))
}

// Match asserts result matches pattern.
func (tt *t) Match(result string, pattern any) {
	tt.assert(operator.Match(result, pattern))
}

// NotMatch asserts result does not match pattern.
func (tt *t) NotMatch(result string, pattern any) {
	tt.assert(operator.NotMatch(result, pattern))
}

// Pass generates an unconditional passing assertion.
func (tt *t) Pass(message ...string) {
	msg := "(unnamed assert)"
	if len(message) > 0 {
		msg = message[0]
	}
	tt.assert(operator.Pass(msg))
}

// Fail generates an unconditional failing assertion.
// message may be a string, error, or any other value.
func (tt *t) Fail(message any) {
	switch msg := message.(type) {
	case string:
		tt.assert(operator.Fail(msg))
	case error:
		tt.assert(operator.Fail(msg.Error()))
	default:
		tt.assert(operator.Fail(failDefault(message)))
	}
}

// failDefault converts an arbitrary value to a fail message string.
func failDefault(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return "fail"
}

// Comment is a no-op annotation; it does not count as an assertion.
func (tt *t) Comment(message string) {
	_ = message
}

// End marks the test as intentionally complete.
func (tt *t) End() {
	tt.setEnded(true)
}
