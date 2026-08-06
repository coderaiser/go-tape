package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// OnlyCall represents a tape.Only() call with its parent TestXxx function.
type OnlyCall struct {
	Parent string // enclosing TestXxx function name
	Name   string // name string argument
}

// CountTests returns total tape.Test + tape.Only + tape.Skip calls in dir.
func CountTests(dir string, exclude []string) (int, error) {
	total := 0
	err := walkFiles(dir, exclude, func(src string) error {
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
func FindDuplicates(dir string, exclude []string) ([]string, error) {
	seen := make(map[string]int)
	err := walkFiles(dir, exclude, func(src string) error {
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
func FindOnlyCalls(dir string, exclude []string) ([]OnlyCall, error) {
	var all []OnlyCall
	err := walkFiles(dir, exclude, func(src string) error {
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
		if !isTestMethodCall(call, "Only") {
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

// findCallNames extracts the second string-literal argument from all tape
// Test(...), Test.Only(...), and Test.Skip(...) calls in src.
// The fnNames parameter is kept for API compatibility but is now ignored —
// all three call forms are always matched.
func findCallNames(src string, _ ...string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return nil, err
	}
	var names []string
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		matched := isTestCall(call) ||
			isTestMethodCall(call, "Only") ||
			isTestMethodCall(call, "Skip")
		if matched && len(call.Args) >= 2 {
			if lit, ok := call.Args[1].(*ast.BasicLit); ok {
				names = append(names, strings.Trim(lit.Value, `"`))
			}
		}
		return true
	})
	return names, nil
}

// isTestCall reports whether call is a plain Test(...) invocation —
// either a bare Ident ("Test") or a package-qualified SelectorExpr ("tape.Test").
// Used for counting all tape tests (Test, Test.Only, Test.Skip).
func isTestCall(call *ast.CallExpr) bool {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		return fn.Name == "Test"
	case *ast.SelectorExpr:
		if id, ok := fn.X.(*ast.Ident); ok {
			return fn.Sel.Name == "Test" && id.Name != "Test"
		}
	}
	return false
}

// isTestMethodCall reports whether call is a Test.<method>(...) invocation.
// Matches both the aliased form (Test.Only) and the qualified form (tape.Test.Only).
// Does NOT match bare Only(...) or tape.Only(...) — those are the old API.
func isTestMethodCall(call *ast.CallExpr, method string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != method {
		return false
	}
	// Test.Only — X is Ident "Test" (import alias form)
	if x, ok := sel.X.(*ast.Ident); ok {
		return x.Name == "Test"
	}
	// tape.Test.Only — X is SelectorExpr with Sel "Test"
	if x, ok := sel.X.(*ast.SelectorExpr); ok {
		return x.Sel.Name == "Test"
	}
	return false
}

// isExcludedDir reports whether path matches any of the exclusion patterns.
// Matching is done against both the full relative path and the directory's
// base name using doublestar.Match. For example, "fixture" skips any
// fixture/ subtree; "cmd/fixture" skips only that specific path.
func isExcludedDir(path string, patterns []string) bool {
	name := filepath.Base(path)
	for _, pattern := range patterns {
		if m, _ := doublestar.Match(pattern, path); m {
			return true
		}
		if m, _ := doublestar.Match(pattern, name); m {
			return true
		}
	}
	return false
}

func walkFiles(dir string, exclude []string, fn func(string) error) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if isExcludedDir(path, exclude) {
				return fs.SkipDir
			}
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
func walkTestFiles(dir string, exclude []string, fn func(string) error) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if isExcludedDir(path, exclude) {
				return fs.SkipDir
			}
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
func CountTestsInTestFiles(dir string, exclude []string) (int, error) {
	total := 0
	err := walkTestFiles(dir, exclude, func(src string) error {
		names, err := findCallNames(src, "Test", "Only", "Skip")
		if err != nil {
			return err
		}
		total += len(names)
		return nil
	})
	return total, err
}

// FindAllTestNames returns all tape.Test/Only/Skip name strings in *_test.go files.
func FindAllTestNames(dir string, exclude []string) ([]string, error) {
	var all []string
	err := walkTestFiles(dir, exclude, func(src string) error {
		names, err := findCallNames(src, "Test", "Only", "Skip")
		if err != nil {
			return err
		}
		all = append(all, names...)
		return nil
	})
	return all, err
}
