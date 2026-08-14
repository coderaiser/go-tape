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
	File   string // absolute path to the file containing the Only call
}

// Location is a file:line position for a test-name string literal.
type Location struct {
	File string // absolute path
	Line int
}

// Duplicate holds a test name that appears more than once and the file/line
// of each occurrence.
type Duplicate struct {
	Name      string
	Locations []Location
}

// CountTests returns total tape.Test + tape.Only calls in dir.
// tape.Skip calls are excluded — they never emit test-end events,
// so including them causes skipped count inflation.
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
func FindDuplicates(dir string, exclude []string) ([]Duplicate, error) {
	seen := make(map[string][]Location)
	err := walkFilesWithPath(dir, exclude, func(path, src string) error {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			return err
		}
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
					name := strings.Trim(lit.Value, `"`)
					pos := fset.Position(lit.Pos())
					seen[name] = append(seen[name], Location{
						File: pos.Filename,
						Line: pos.Line,
					})
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	var dups []Duplicate
	for name, locs := range seen {
		if len(locs) > 1 {
			dups = append(dups, Duplicate{Name: name, Locations: locs})
		}
	}
	return dups, nil
}

// FindOnlyCalls returns all tape.Only() calls with enclosing TestXxx name.
func FindOnlyCalls(dir string, exclude []string) ([]OnlyCall, error) {
	var all []OnlyCall
	err := walkFilesWithPath(dir, exclude, func(path, src string) error {
		calls, err := FindOnlyCallsInSource(src)
		if err != nil {
			return err
		}
		for i := range calls {
			calls[i].File = path // absolute path
		}
		all = append(all, calls...)
		return nil
	})
	return all, err
}

// funcSpan is the position range of a top-level TestXxx function declaration.
type funcSpan struct {
	name  string
	start token.Pos
	end   token.Pos
}

// FindOnlyCallsInSource parses Go source and returns Only calls.
// Pure — no I/O. Used directly in tests.
//
// Parent lookup uses position spans rather than traversal order, so Only calls
// inside anonymous function literals are attributed to the correct enclosing
// TestXxx function.
func FindOnlyCallsInSource(src string) ([]OnlyCall, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return nil, err
	}

	// First pass: collect position spans of all top-level TestXxx functions.
	var spans []funcSpan
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Body == nil {
			continue
		}
		if strings.HasPrefix(fn.Name.Name, "Test") {
			spans = append(spans, funcSpan{
				name:  fn.Name.Name,
				start: fn.Pos(),
				end:   fn.End(),
			})
		}
	}

	// parentOf returns the TestXxx function name whose span contains pos,
	// or "" if none does.
	parentOf := func(pos token.Pos) string {
		for _, s := range spans {
			if pos >= s.start && pos <= s.end {
				return s.name
			}
		}
		return ""
	}

	// Second pass: find all Only calls and look up their parent by position.
	var calls []OnlyCall
	ast.Inspect(f, func(n ast.Node) bool {
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
				parent := parentOf(call.Pos())
				calls = append(calls, OnlyCall{Parent: parent, Name: name})
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

// findCallNames extracts the second string-literal argument from tape
// Test(...) and Test.Only(...) calls in src.
// tape.Skip is intentionally excluded — Skip calls never emit test-end
// events, so counting them inflates the total and produces phantom skips.
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
		if _, ok := fn.X.(*ast.Ident); ok {
			return fn.Sel.Name == "Test"
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

// walkFilesWithPath is like walkFiles but passes the absolute file path as the
// first argument to fn.
func walkFilesWithPath(dir string, exclude []string, fn func(string, string) error) error {
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

		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if hasBuildIgnore(string(src)) {
			return nil
		}

		return fn(abs, string(src))
	})
}

// hasBuildIgnore returns true only if //go:build ignore appears as a real
// file-level build constraint, not inside a string literal or function body.
func hasBuildIgnore(src string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		// unparseable — fall back to string search on first 512 bytes
		if len(src) > 512 {
			src = src[:512]
		}
		return strings.Contains(src, "//go:build ignore")
	}
	// Only look at comments before the package declaration position.
	pkgPos := f.Package
	for _, cg := range f.Comments {
		if cg.Pos() >= pkgPos {
			break
		}
		for _, c := range cg.List {
			if c.Text == "//go:build ignore" {
				return true
			}
		}
	}
	return false
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

// walkTestFilesWithPath is like walkTestFiles but passes (path, src) to fn.
func walkTestFilesWithPath(dir string, exclude []string, fn func(path, src string) error) error {
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
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		return fn(abs, string(src))
	})
}

// TestCall is a single tape.Test / tape.Only call found during AST scan.
type TestCall struct {
	Kind string   // "Test" or "Only"
	Name string   // second string-literal argument
	File string   // absolute path
	Line int
}

// FindTestsWithLocations returns every tape.Test and tape.Only call in
// *_test.go files under dir, with file and line information.
// Used by formatter-debug to show exactly what the AST scan counted.
func FindTestsWithLocations(dir string, exclude []string) ([]TestCall, error) {
	var calls []TestCall
	err := walkTestFilesWithPath(dir, exclude, func(path, src string) error {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			kind := ""
			switch {
			case isTestCall(call):
				kind = "Test"
			case isTestMethodCall(call, "Only"):
				kind = "Only"
			}
			if kind == "" {
				return true
			}
			if len(call.Args) >= 2 {
				if lit, ok := call.Args[1].(*ast.BasicLit); ok {
					pos := fset.Position(lit.Pos())
					calls = append(calls, TestCall{
						Kind: kind,
						Name: strings.Trim(lit.Value, `"`),
						File: path,
						Line: pos.Line,
					})
				}
			}
			return true
		})
		return nil
	})
	return calls, err
}


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
