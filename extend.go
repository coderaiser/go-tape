package tape

import (
	"testing"

	"github.com/coderaiser/go-tape/internal/operator"
)

// Operators contains the built-in assertion functions for extension packages.
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
	Equal:        func(result, expected any) Result { return toResult(operator.Equal(result, expected)) },
	NotEqual:     func(result, expected any) Result { return toResult(operator.NotEqual(result, expected)) },
	DeepEqual:    func(result, expected any) Result { return toResult(operator.DeepEqual(result, expected)) },
	NotDeepEqual: func(result, expected any) Result { return toResult(operator.NotDeepEqual(result, expected)) },
	Ok:           func(result any) Result { return toResult(operator.Ok(result)) },
	NotOk:        func(result any) Result { return toResult(operator.NotOk(result)) },
	Match:        func(result string, pattern any) Result { return toResult(operator.Match(result, pattern)) },
	NotMatch:     func(result string, pattern any) Result { return toResult(operator.NotMatch(result, pattern)) },
	Pass:         func(message string) Result { return toResult(operator.Pass(message)) },
	Fail:         func(message string) Result { return toResult(operator.Fail(message)) },
}

// Extend[XT any] creates a test function that passes an extended T to fn.
func Extend[XT any](factory func(*T) XT) func(*testing.T, string, func(XT)) {
	return func(t *testing.T, name string, fn func(XT)) {
		Test(t, name, func(base *T) {
			fn(factory(base))
		})
	}
}
