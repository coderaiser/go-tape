package tape

import (
	"testing"

	"github.com/coderaiser/go-tape/internal/operator"
)

// Operators contains the built-in assertion functions for extension packages.
type Operators struct {
	Equal        func(result, expected any) operator.Result
	NotEqual     func(result, expected any) operator.Result
	DeepEqual    func(result, expected any) operator.Result
	NotDeepEqual func(result, expected any) operator.Result
	Ok           func(result any) operator.Result
	NotOk        func(result any) operator.Result
	Match        func(result string, pattern any) operator.Result
	NotMatch     func(result string, pattern any) operator.Result
	Pass         func(message string) operator.Result
	Fail         func(message string) operator.Result
}

// BuiltinOperators is the canonical instance passed to extension factories.
var BuiltinOperators = Operators{
	Equal:        operator.Equal,
	NotEqual:     operator.NotEqual,
	DeepEqual:    operator.DeepEqual,
	NotDeepEqual: operator.NotDeepEqual,
	Ok:           operator.Ok,
	NotOk:        operator.NotOk,
	Match:        operator.Match,
	NotMatch:     operator.NotMatch,
	Pass:         operator.Pass,
	Fail:         operator.Fail,
}

// Extend[XT any] creates a test function that passes an extended T to fn.
func Extend[XT any](factory func(*T) XT) func(*testing.T, string, func(XT)) {
	return func(t *testing.T, name string, fn func(XT)) {
		Test(t, name, func(base *T) {
			fn(factory(base))
		})
	}
}
