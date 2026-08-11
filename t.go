package tape

import (
	"encoding/json"
	"testing"

	"github.com/coderaiser/go-tape/internal/operator"
)

// T wraps *testing.T with supertape-compatible assertions.
type T struct {
	t     *testing.T
	ended bool
}

func newT(t *testing.T) *T { return &T{t: t} }

// TB returns the underlying *testing.T for use with helpers
// that require it directly (e.g. t.TB().TempDir(), t.TB().Fatal(...)).
func (tt *T) TB() *testing.T {
	return tt.t
}

// Setenv sets an environment variable for the duration of the test.
// It is automatically restored when the test ends.
func (tt *T) Setenv(key, value string) {
	tt.t.Helper()
	tt.t.Setenv(key, value)
}

// Report records a Result, marking pass or fail on the underlying test.
func (tt *T) Report(r Result) {
	tt.t.Helper()
	if !r.Ok {
		b, _ := json.Marshal(struct {
			Message  string `json:"message"`
			Operator string `json:"operator"`
			Result   any    `json:"result"`
			Expected any    `json:"expected"`
			Output   string `json:"output"`
		}{
			Message:  r.Message,
			Operator: r.Operator,
			Result:   r.Result,
			Expected: r.Expected,
			Output:   r.Output,
		})
		tt.t.Log("TAPE:" + string(b)) // structured sentinel — for stream parser
		tt.t.Fail()                   // marks test failed — no freeform text
	}
}

// ReportCustom records a custom operator result and counts it against the
// one-assertion-per-block guard. Extension packages use it to report named
// operators such as transform, noTransform, report and noReport.
func (tt *T) ReportCustom(ok bool, operatorName, output string, got, expected any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(Result{
		Ok:       ok,
		Message:  operatorName,
		Result:   got,
		Expected: expected,
		Output:   output,
	})
}

// Equal asserts result == expected using strict equality.
// For primitives and pointers only.
// Use DeepEqual for structs, slices, and maps.
func (tt *T) Equal(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.Equal(result, expected))
}

// NotEqual asserts result != expected.
// For primitives and pointers only.
func (tt *T) NotEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.NotEqual(result, expected))
}

// DeepEqual asserts deep equality using reflect.DeepEqual.
// Use for structs, slices, and maps.
func (tt *T) DeepEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.DeepEqual(result, expected))
}

// NotDeepEqual asserts values are not deeply equal.
func (tt *T) NotDeepEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.NotDeepEqual(result, expected))
}

// Ok asserts result is truthy.
func (tt *T) Ok(result any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.Ok(result))
}

// NotOk asserts result is falsy.
func (tt *T) NotOk(result any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.NotOk(result))
}

// Match asserts result matches pattern.
// pattern may be a string or *regexp.Regexp.
func (tt *T) Match(result string, pattern any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.Match(result, pattern))
}

// NotMatch asserts result does not match pattern.
// pattern may be a string or *regexp.Regexp.
func (tt *T) NotMatch(result string, pattern any) {
	tt.t.Helper()
	hit(tt.t)
	tt.Report(operator.NotMatch(result, pattern))
}

// Pass generates an unconditional passing assertion.
func (tt *T) Pass(message ...string) {
	tt.t.Helper()
	hit(tt.t)
	msg := "(unnamed assert)"
	if len(message) > 0 {
		msg = message[0]
	}
	tt.Report(operator.Pass(msg))
}

// Fail generates an unconditional failing assertion.
// message may be a string or error.
func (tt *T) Fail(message any) {
	tt.t.Helper()
	hit(tt.t)
	switch msg := message.(type) {
	case string:
		tt.Report(operator.Fail(msg))
	case error:
		tt.Report(operator.Fail(msg.Error()))
	default:
		tt.Report(operator.Fail("fail"))
	}
}

// Comment prints a TAP comment without counting as an assertion.
func (tt *T) Comment(message string) {
	tt.t.Helper()
	tt.t.Log("#", message)
}

// End marks the test as intentionally complete.
// Required when TAPE_CHECK_END is enabled (default: on).
func (tt *T) End() {
	tt.ended = true
}
