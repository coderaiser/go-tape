package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// OnlyCall represents a tape.Only() call with its parent TestXxx function.
type OnlyCall struct {
	Parent string // enclosing TestXxx function name
	Name   string // name string argument
}

// CountTests returns total tape.Test + tape.Only + tape.Skip calls in dir.
func CountTests(dir string) (int, error) {
	total := 0
	err := walkFiles(dir, func(src string) error {
		names, err := findCallNames(src, "Test", "Only", "Skip")
		if err != nil {
			return err
		}
		total += len(names)
		return nil
	})
	return total, err
}

// FindDuplicates returns test names appearing more than once in dir.
func FindDuplicates(dir string) ([]string, error) {
	seen := make(map[string]int)
	err := walkFiles(dir, func(src string) error {
		names, err := findCallNames(src, "Test", "Only", "Skip")
		if err != nil {
			return err
		}
		for _, n := range names {
			seen[n]++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var dups []string
	for name, count := range seen {
		if count > 1 {
			dups = append(dups, name)
		}
	}
	return dups, nil
}

// FindOnlyCalls returns all tape.Only() calls with enclosing TestXxx name.
func FindOnlyCalls(dir string) ([]OnlyCall, error) {
	var all []OnlyCall
	err := walkFiles(dir, func(src string) error {
		calls, err := FindOnlyCallsInSource(src)
		if err != nil {
			return err
		}
		all = append(all, calls...)
		return nil
	})
	return all, err
}

// FindOnlyCallsInSource parses Go source and returns Only calls.
// Pure — no I/O. Used directly in tests.
func FindOnlyCallsInSource(src string) ([]OnlyCall, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return nil, err
	}
	var calls []OnlyCall
	var currentFunc string
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name != nil && strings.HasPrefix(fn.Name.Name, "Test") {
				currentFunc = fn.Name.Name
			}
			return true
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if !isTapeCall(call, "Only") {
			return true
		}
		if len(call.Args) >= 2 {
			if lit, ok := call.Args[1].(*ast.BasicLit); ok {
				name := strings.Trim(lit.Value, `"`)
				calls = append(calls, OnlyCall{Parent: currentFunc, Name: name})
			}
		}
		return true
	})
	return calls, nil
}

// BuildRunPattern converts OnlyCalls to a go test -run pattern.
func BuildRunPattern(calls []OnlyCall) string {
	if len(calls) == 0 {
		return ""
	}
	patterns := make([]string, len(calls))
	for i, c := range calls {
		patterns[i] = c.Parent + "/" + strings.ReplaceAll(c.Name, " ", "_")
	}
	return strings.Join(patterns, "|")
}

func findCallNames(src string, fnNames ...string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return nil, err
	}
	nameSet := make(map[string]bool)
	for _, n := range fnNames {
		nameSet[n] = true
	}
	var names []string
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		for fname := range nameSet {
			if isTapeCall(call, fname) && len(call.Args) >= 2 {
				if lit, ok := call.Args[1].(*ast.BasicLit); ok {
					names = append(names, strings.Trim(lit.Value, `"`))
				}
			}
		}
		return true
	})
	return names, nil
}

func isTapeCall(call *ast.CallExpr, name string) bool {
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		return fn.Sel.Name == name
	case *ast.Ident:
		return fn.Name == name
	}
	return false
}

func walkFiles(dir string, fn func(string) error) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if hasBuildIgnore(string(src)) {
			return nil
		}

		return fn(string(src))
	})
}

// hasBuildIgnore returns true if the source contains //go:build ignore.
func hasBuildIgnore(src string) bool {
	return strings.Contains(src, "//go:build ignore")
}

// walkTestFiles is like walkFiles but restricted to _test.go files.
func walkTestFiles(dir string, fn func(string) error) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if hasBuildIgnore(string(src)) {
			return nil
		}
		return fn(string(src))
	})
}

// CountTestsInTestFiles counts tape.Test/Only/Skip calls only in *_test.go files.
// Use this from the CLI to avoid counting fixture files.
func CountTestsInTestFiles(dir string) (int, error) {
	total := 0
	err := walkTestFiles(dir, func(src string) error {
		names, err := findCallNames(src, "Test", "Only", "Skip")
		if err != nil {
			return err
		}
		total += len(names)
		return nil
	})
	return total, err
}

// CountSkipCallsInTestFiles counts only tape.Skip calls in *_test.go files.
// Skip calls never run — they are counted separately like supertape's model.
func CountSkipCallsInTestFiles(dir string) (int, error) {
	total := 0
	err := walkTestFiles(dir, func(src string) error {
		names, err := findCallNames(src, "Skip")
		if err != nil {
			return err
		}
		total += len(names)
		return nil
	})
	return total, err
}

// FindAllTestNames returns all tape.Test/Only/Skip name strings in *_test.go files.
func FindAllTestNames(dir string) ([]string, error) {
	var all []string
	err := walkTestFiles(dir, func(src string) error {
		names, err := findCallNames(src, "Test", "Only", "Skip")
		if err != nil {
			return err
		}
		all = append(all, names...)
		return nil
	})
	return all, err
}
