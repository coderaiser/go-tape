package tape

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/goplus/ixgo"
)

// testT is the host T interface backing the tape-style API.
type testT interface {
	Equal(a, b any)
	End()
}

type nativeT struct{}

func (nativeT) Equal(a, b any) {
	fmt.Printf("EVENT Equal(%v,%v)\n", a, b)
}

func (nativeT) End() {
	fmt.Printf("EVENT End\n")
}

var hostTest = func(name string, fn func(testT)) {
	fmt.Printf("EVENT test %q\n", name)
	fn(nativeT{})
}

type result struct {
	Ok     bool
	Name   string
	Msg    string
}

var hostReport = func(r result) {
	fmt.Printf("EVENT report ok=%v name=%q msg=%q\n", r.Ok, r.Name, r.Msg)
}

// TestIxgoSourcePkg registers a package with Source + Interfaces + Vars.
func TestIxgoSourcePkg(t *testing.T) {
	tapeapiSrc := `package tapeapi

type T interface {
	Equal(a, b any)
	End()
}

type Result struct {
	Ok     bool
	Name   string
	Msg    string
}

var Test func(name string, fn func(T))
var Report func(r Result)
`
	ixgo.RegisterPackage(&ixgo.Package{
		Name:   "tapeapi",
		Path:   "tapeapi",
		Source: tapeapiSrc,
		Interfaces: map[string]reflect.Type{
			"T": reflect.TypeOf((*testT)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"Result": reflect.TypeOf(result{}),
		},
		Vars: map[string]reflect.Value{
			"Test":   reflect.ValueOf(&hostTest),
			"Report": reflect.ValueOf(&hostReport),
		},
	})

	src := `package main

import . "tapeapi"

func main() {
	Test("scope: works", func(t T) {
		t.Equal(1, 1)
		t.End()
	})
	Report(Result{Ok: true, Name: "x", Msg: "y"})
}
`
	ctx := ixgo.NewContext(0)
	code, err := ctx.RunFile("main.go", src, nil)
	if err != nil {
		t.Fatalf("RunFile error: %v (code=%d)", err, code)
	}
	t.Logf("run ok code=%d", code)
}



