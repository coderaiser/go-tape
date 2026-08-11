package tape

import "github.com/coderaiser/go-tape/internal/operator"

// Result is aliased from operator.Result so the public API is stable
// without a separate type definition. Extension packages build named
// operators on top of it, the same way supertape extensions register
// custom operator names.
//
// The diff string lives inside Output — extension packages should call
// BuiltinOperators.Equal(got, expected) and read the Result.Output field
// instead of calling a separate Diff function.
type Result = operator.Result
