package tape

import (
	"testing"

	"github.com/coderaiser/go-tape/internal/operator"
)

// Operators contains the built-in assertion functions for extension packages.
// Extensions should use these functions (Equal, NotEqual, DeepEqual, etc.) to
// build named operators, then read the Result.Output field for the diff string.
// This mirrors the supertape extension API: operators are the building block,
// and the diff is encapsulated inside the Result.
type Operators struct {
	Equal        func(result, expected any) Result
	NotEqual     func(result, expected any) Result
	DeepEqual    func(result, expected any) Result
	NotDeepEqual func(result, expected any) Result
	Ok           func(result any) Result
	NotOk        func(result any) Result
	Match        func(result string, pattern any) Result
	NotMatch     func(result string, pattern any) Result
	Pass         func(message string) Result
	Fail         func(message string) Result
}

// BuiltinOperators is the canonical instance passed to extension factories.
var BuiltinOperators = Operators{
	Equal:        func(result, expected any) Result { return operator.Equal(result, expected) },
	NotEqual:     func(result, expected any) Result { return operator.NotEqual(result, expected) },
	DeepEqual:    func(result, expected any) Result { return operator.DeepEqual(result, expected) },
	NotDeepEqual: func(result, expected any) Result { return operator.NotDeepEqual(result, expected) },
	Ok:           func(result any) Result { return operator.Ok(result) },
	NotOk:        func(result any) Result { return operator.NotOk(result) },
	Match:        func(result string, pattern any) Result { return operator.Match(result, pattern) },
	NotMatch:     func(result string, pattern any) Result { return operator.NotMatch(result, pattern) },
	Pass:         func(message string) Result { return operator.Pass(message) },
	Fail:         func(message string) Result { return operator.Fail(message) },
}

// ExtendFn[XT] is the return type of Extend. Being a named type allows
// .Skip() and .Only() to be attached, mirroring TestFn for extended T.
type ExtendFn[XT any] func(t *testing.T, name string, fn func(XT))

// Skip marks the subtest as skipped without running its body.
func (f ExtendFn[XT]) Skip(_ *testing.T, _ string, _ func(XT)) {}

// Only runs one subtest; all others are skipped via the onlyName guard.
func (f ExtendFn[XT]) Only(t *testing.T, name string, fn func(XT)) {
	setOnlyName(name)
	f(t, name, fn)
}

// Extend creates a test function that passes an extended T to fn.
func Extend[XT any](factory func(*T) XT) ExtendFn[XT] {
	return ExtendFn[XT](func(t *testing.T, name string, fn func(XT)) {
		Test(t, name, func(base *T) {
			fn(factory(base))
		})
	})
}
