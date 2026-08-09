package tape

import "github.com/coderaiser/go-tape/internal/operator"

// Result is the public shape of an operator outcome. Extension packages build
// named operators on top of it, the same way supertape extensions register
// custom operator names.
//
// The diff string lives inside Output — extension packages should call
// BuiltinOperators.Equal(got, expected) and read the Result.Output field
// instead of calling a separate Diff function.
type Result struct {
	Ok       bool
	Message  string
	Result   any
	Expected any
	Output   string
}

// toResult converts an internal operator.Result into the public Result shape.
func toResult(r operator.Result) Result {
	return Result{
		Ok:       r.Ok,
		Message:  r.Message,
		Result:   r.Result,
		Expected: r.Expected,
		Output:   r.Output,
	}
}
