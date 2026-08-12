// Package interp proves that ixgo can interpret a plain-Go supertape-style
// test source (package main with top-level Test(name, fn) calls and a plain T
// interface with assertion methods), where the test API is supplied as ixgo
// externals via the package-export mechanism - no *testing.T, no go test, no
// testing.MainStart/RunTests.
package interp

import (
	"reflect"

	"github.com/goplus/ixgo"
)

// apiPath is the import path under which the tape-style API is registered.
const apiPath = "tapeapi"

// T is the host interface backing the interpreted test body's T parameter.
// It mirrors the supertape API: assertions plus End.
type T interface {
	Equal(a, b any)
	EqualText(a, b any, msg string)
	DeepEqual(a, b any)
	Ok(a any)
	NotOk(a any)
	End()
}

// Result mirrors a tape Report value (structured assertion result).
type Result struct {
	Ok       bool
	Operator string
	Got      any
	Expected any
	Message  string
}

// Register exports the tape-style API as an ixgo package so that a plain
// `package main` source can `import . "tapeapi"` and use Test/T/Report with no
// package qualifier. The Source provides the Go types the interpreter needs to
// type-check the test body; Interfaces/NamedTypes map those source types to the
// host types; Vars supplies the runtime values (overridden per-run via the
// context so each Run can route events to its own handler).
func init() {
	tapeapiSrc := `package tapeapi

type T interface {
	Equal(a, b any)
	EqualText(a, b any, msg string)
	DeepEqual(a, b any)
	Ok(a any)
	NotOk(a any)
	End()
}

type Result struct {
	Ok       bool
	Operator string
	Got      any
	Expected any
	Message  string
}

var Test func(name string, fn func(T))

var Report func(r Result)
`
	ixgo.RegisterPackage(&ixgo.Package{
		Name:   "tapeapi",
		Path:   apiPath,
		Source: tapeapiSrc,
		Interfaces: map[string]reflect.Type{
			"T": reflect.TypeOf((*T)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"Result": reflect.TypeOf(Result{}),
		},
		Vars: map[string]reflect.Value{
			"Test":   reflect.ValueOf(&defaultTest),
			"Report": reflect.ValueOf(&defaultReport),
		},
	})
}
