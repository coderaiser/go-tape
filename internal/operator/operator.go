package operator

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/coderaiser/go-tape/internal/diff"
)

// Result mirrors supertape's testState object exactly.
type Result struct {
	Ok       bool
	Message  string // "should equal", "should be truthy", etc.
	Result   any    // actual value
	Expected any    // expected value
	Output   string // diff or detail — empty on pass
	At       string // source location (filled by T.Report, not operator)
}

// Equal asserts result == expected using reflect.DeepEqual.
func Equal(result, expected any) Result {
	ok := reflect.DeepEqual(result, expected)
	out := ""
	if !ok {
		out = diff.Diff(expected, result)
	}
	return Result{Ok: ok, Message: "should equal", Result: result, Expected: expected, Output: out}
}

// NotEqual asserts result != expected for primitives and pointers.
func NotEqual(result, expected any) Result {
	ok := !isPrimitive(result) || !reflect.DeepEqual(result, expected)
	out := ""
	if !ok {
		out = diff.Diff(expected, result)
	}
	return Result{Ok: ok, Message: "should not equal", Result: result, Expected: expected, Output: out}
}

// DeepEqual asserts deep equality.
func DeepEqual(result, expected any) Result {
	ok := reflect.DeepEqual(result, expected)
	out := ""
	if !ok {
		out = diff.Diff(expected, result)
	}
	return Result{Ok: ok, Message: "should deep equal", Result: result, Expected: expected, Output: out}
}

// NotDeepEqual asserts values are not deeply equal.
func NotDeepEqual(result, expected any) Result {
	ok := !reflect.DeepEqual(result, expected)
	out := ""
	if !ok {
		out = diff.Diff(expected, result)
	}
	return Result{Ok: ok, Message: "should not deep equal", Result: result, Expected: expected, Output: out}
}

// Ok asserts result is truthy.
func Ok(result any) Result {
	ok := truthy(result)
	out := ""
	if !ok {
		out = diff.Diff("truthy", result)
	}
	return Result{Ok: ok, Message: "should be truthy", Result: result, Expected: "truthy", Output: out}
}

// NotOk asserts result is falsy.
func NotOk(result any) Result {
	ok := !truthy(result)
	out := ""
	if !ok {
		out = diff.Diff("falsy", result)
	}
	return Result{Ok: ok, Message: "should be falsy", Result: result, Expected: "falsy", Output: out}
}

// Match asserts result matches pattern.
// pattern may be a string or *regexp.Regexp.
func Match(result string, pattern any) Result {
	re, err := toRegexp(pattern)
	if err != nil {
		return Result{Ok: false, Message: "should match", Result: result, Expected: pattern, Output: err.Error()}
	}
	ok := re.MatchString(result)
	out := ""
	if !ok {
		out = diff.Diff(pattern, result)
	}
	return Result{Ok: ok, Message: "should match", Result: result, Expected: pattern, Output: out}
}

// NotMatch asserts result does not match pattern.
// pattern may be a string or *regexp.Regexp.
func NotMatch(result string, pattern any) Result {
	re, err := toRegexp(pattern)
	if err != nil {
		return Result{Ok: false, Message: "should not match", Result: result, Expected: pattern, Output: err.Error()}
	}
	ok := !re.MatchString(result)
	out := ""
	if !ok {
		out = diff.Diff(pattern, result)
	}
	return Result{Ok: ok, Message: "should not match", Result: result, Expected: pattern, Output: out}
}

// Pass generates a passing result.
func Pass(message string) Result {
	return Result{Ok: true, Message: "pass", Result: message}
}

// Fail generates a failing result.
func Fail(message string) Result {
	return Result{Ok: false, Message: "fail", Result: message, Output: message}
}

// truthy checks if v is truthy.
func truthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case string:
		return val != ""
	default:
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Slice, reflect.Map, reflect.Chan,
			reflect.Func, reflect.Pointer, reflect.Interface:
			return !rv.IsNil()
		}
		return true
	}
}

// toRegexp converts a string or *regexp.Regexp to *regexp.Regexp.
func toRegexp(pattern any) (*regexp.Regexp, error) {
	switch p := pattern.(type) {
	case *regexp.Regexp:
		return p, nil
	case string:
		return regexp.Compile(regexp.QuoteMeta(p))
	default:
		return nil, fmt.Errorf("unsupported pattern type %T", pattern)
	}
}

// isPrimitive checks if v is a primitive type or pointer.
func isPrimitive(v any) bool {
	switch v.(type) {
	case bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		complex64, complex128,
		string, uintptr:
		return true
	}
	t := reflect.TypeOf(v)
	return t != nil && t.Kind() == reflect.Pointer
}
