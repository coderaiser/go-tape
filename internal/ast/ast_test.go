package ast_test

import (
	"os"
	"sort"
	"testing"

	tape "github.com/coderaiser/go-tape"
	tapeast "github.com/coderaiser/go-tape/internal/ast"
	dedent "github.com/coderaiser/go-tape/internal/dedent"
)

// AstT extends tape.T with fixture operators.
type AstT struct{ *tape.T }

func (t *AstT) OnlyCallsInFile(file string, expected []tapeast.OnlyCall) {
	t.TB().Helper()

	src, err := os.ReadFile("fixture/" + file)
	if err != nil {
		t.TB().Fatalf("fixture not found: %v", err)
	}

	calls, err := tapeast.FindOnlyCallsInSource(string(src))
	if err != nil {
		t.TB().Fatalf("unexpected error: %v", err)
	}

	sort.Slice(calls, func(i, j int) bool {
		return calls[i].Name < calls[j].Name
	})
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Name < expected[j].Name
	})

	t.DeepEqual(calls, expected)
}

func (t *AstT) Pattern(calls []tapeast.OnlyCall, expected string) {
	t.TB().Helper()
	t.Equal(tapeast.BuildRunPattern(calls), expected)
}

// AstTest uses Extend to create the extended test convenience function.
var AstTest = tape.Extend(func(base *tape.T) *AstT {
	return &AstT{T: base}
})

func writeFile(t *testing.T, path, src string) {
	t.Helper()

	err := os.WriteFile(path, []byte(dedent.Dedent(src)), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindNoOnlyCalls(t *testing.T) {
	AstTest(t, "ast: no Only calls returns nil", func(t *AstT) {
		t.OnlyCallsInFile("no-only.go", nil)
		t.End()
	})
}

func TestFindOneOnlyCall(t *testing.T) {
	AstTest(t, "ast: one Only call returns one result", func(t *AstT) {
		t.OnlyCallsInFile("one-only.go", []tapeast.OnlyCall{
			{Parent: "TestOneOnlyParser", Name: "parser: run action"},
		})
		t.End()
	})
}

func TestFindMultipleOnlyCalls(t *testing.T) {
	AstTest(t, "ast: multiple Only calls in same func", func(t *AstT) {
		t.OnlyCallsInFile("multi-only.go", []tapeast.OnlyCall{
			{Parent: "TestMultiParser", Name: "parser: fail action"},
			{Parent: "TestMultiParser", Name: "parser: run action"},
		})
		t.End()
	})
}

func TestFindCrossFuncOnlyCalls(t *testing.T) {
	AstTest(t, "ast: Only calls across different TestXxx functions", func(t *AstT) {
		t.OnlyCallsInFile("cross-func.go", []tapeast.OnlyCall{
			{Parent: "TestCrossParser", Name: "parser: run action"},
			{Parent: "TestCrossRunner", Name: "runner: starts"},
		})
		t.End()
	})
}

func TestBuildPatternEmpty(t *testing.T) {
	AstTest(t, "ast: empty calls returns empty pattern", func(t *AstT) {
		t.Pattern(nil, "")
		t.End()
	})
}

func TestBuildPatternSingle(t *testing.T) {
	AstTest(t, "ast: single call builds pattern with parent", func(t *AstT) {
		t.Pattern([]tapeast.OnlyCall{
			{Parent: "TestParser", Name: "parser: run action"},
		}, "TestParser/parser:_run_action")
		t.End()
	})
}

func TestBuildPatternMultiple(t *testing.T) {
	AstTest(t, "ast: multiple calls joined with pipe", func(t *AstT) {
		t.Pattern([]tapeast.OnlyCall{
			{Parent: "TestParser", Name: "parser: run action"},
			{Parent: "TestRunner", Name: "runner: starts"},
		}, "TestParser/parser:_run_action|TestRunner/runner:_starts")
		t.End()
	})
}

func TestCountTests(t *testing.T) {
	AstTest(t, "ast: CountTests counts Test Only and Skip calls", func(t *AstT) {
		dir := t.TB().TempDir()

		writeFile(t.TB(), dir+"/foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: bar", func(t *Test.T) { t.Ok(true); t.End() })
				Test.Only(t, "foo: baz", func(t *Test.T) { t.Ok(true); t.End() })
			}
		`)

		n, err := tapeast.CountTests(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(n > 0)
		t.End()
	})
}

func TestFindDuplicatesFound(t *testing.T) {
	AstTest(t, "ast: FindDuplicates finds duplicate names", func(t *AstT) {
		dir := t.TB().TempDir()

		writeFile(t.TB(), dir+"/foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: bar", func(t *Test.T) { t.Ok(true); t.End() })
				Test(t, "foo: bar", func(t *Test.T) { t.Ok(true); t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(dups) > 0)
		t.End()
	})
}

func TestFindDuplicatesReturnsLocations(t *testing.T) {
	AstTest(t, "ast: FindDuplicates returns Duplicate with locations", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: bar", func(t *Test.T) { t.End() })
				Test(t, "foo: bar", func(t *Test.T) { t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(dups) == 1)
		t.End()
	})
	AstTest(t, "ast: FindDuplicates locations count equals occurrences", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: bar", func(t *Test.T) { t.End() })
				Test(t, "foo: bar", func(t *Test.T) { t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(dups[0].Locations) == 2)
		t.End()
	})
}

func TestFindDuplicatesLocationFile(t *testing.T) {
	AstTest(t, "ast: FindDuplicates Location.File is non-empty", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/a_test.go", `
			package p
			import "testing"
			var Test = func(t *testing.T, name string, fn func()) {}
			func TestA(t *testing.T) {
				Test(t, "dup name", func() {})
				Test(t, "dup name", func() {})
			}
		`)
		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Ok(dups[0].Locations[0].File != "")
		t.End()
	})
}

func TestFindDuplicatesLocationLine(t *testing.T) {
	AstTest(t, "ast: FindDuplicates Location.Line is non-zero", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/a_test.go", `
			package p
			import "testing"
			var Test = func(t *testing.T, name string, fn func()) {}
			func TestA(t *testing.T) {
				Test(t, "dup name", func() {})
				Test(t, "dup name", func() {})
			}
		`)
		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Ok(dups[0].Locations[0].Line > 0)
		t.End()
	})
}

func TestFindOnlyCallsDir(t *testing.T) {
	AstTest(t, "ast: FindOnlyCalls reads all go files in dir", func(t *AstT) {
		dir := t.TB().TempDir()

		writeFile(t.TB(), dir+"/foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test.Only(t, "foo: bar", func(t *Test.T) { t.Ok(true); t.End() })
			}
		`)

		calls, err := tapeast.FindOnlyCalls(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(calls) > 0)
		t.End()
	})
}

func TestFindOnlyCallsMissingDir(t *testing.T) {
	AstTest(t, "ast: FindOnlyCalls errors on missing dir", func(t *AstT) {
		_, err := tapeast.FindOnlyCalls("nonexistent", nil)
		t.Ok(err)
		t.End()
	})
}

func TestFindOnlyCallsInSourceInvalid(t *testing.T) {
	AstTest(t, "ast: FindOnlyCallsInSource errors on invalid Go source", func(t *AstT) {
		_, err := tapeast.FindOnlyCallsInSource("not go {{{{")
		t.Ok(err)
		t.End()
	})
}

func TestCountTestsMissingDir(t *testing.T) {
	AstTest(t, "ast: CountTests errors on missing dir", func(t *AstT) {
		_, err := tapeast.CountTests("nonexistent", nil)
		t.Ok(err)
		t.End()
	})
}

func TestFindDuplicatesMissingDir(t *testing.T) {
	AstTest(t, "ast: FindDuplicates errors on missing dir", func(t *AstT) {
		_, err := tapeast.FindDuplicates("nonexistent", nil)
		t.Ok(err)
		t.End()
	})
}

func TestFindOnlyCallsMethodSyntaxNoError(t *testing.T) {
	AstTest(t, "ast: finds Test.Only call no error", func(t *AstT) {
		src := `
			package foo

			import (
				"testing"
				Test "github.com/coderaiser/go-tape"
			)

			func TestFoo(t *testing.T) {
				Test.Only(t, "foo: bar", func(t *Test.T) {})
			}
		`

		_, err := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.NotOk(err)
		t.End()
	})
}

func TestFindOnlyCallsMethodSyntaxResult(t *testing.T) {
	AstTest(t, "ast: finds Test.Only call result", func(t *AstT) {
		src := `
			package foo

			import (
				"testing"
				Test "github.com/coderaiser/go-tape"
			)

			func TestFoo(t *testing.T) {
				Test.Only(t, "foo: bar", func(t *Test.T) {})
			}
		`

		calls, err := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		if err != nil {
			t.End()
			return
		}

		t.DeepEqual(calls, []tapeast.OnlyCall{
			{
				Parent: "TestFoo",
				Name:   "foo: bar",
			},
		})

		t.End()
	})
}

func TestFindSkipCallsMethodSyntaxResult(t *testing.T) {
	AstTest(t, "ast: finds Test.Skip call result", func(t *AstT) {
		src := `
			package foo

			import (
				"testing"
				Test "github.com/coderaiser/go-tape"
			)

			func TestFoo(t *testing.T) {
				Test.Skip(t, "foo: bar", func(t *Test.T) {})
			}
		`

		names, _ := tapeast.FindAllTestNames(t.TB().TempDir(), nil)
		_ = names
		calls, err := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		if err != nil {
			t.End()
			return
		}

		// Test.Skip is NOT an Only call — it should not appear in FindOnlyCalls
		t.Equal(len(calls), 0)
		t.End()
	})
}

func TestWalkFilesReadError(t *testing.T) {
	AstTest(t, "ast: CountTests errors on nonexistent directory", func(t *AstT) {
		_, err := tapeast.CountTests("/nonexistent/path/that/does/not/exist", nil)
		t.Ok(err)
		t.End()
	})
}

func TestCountTestsInvalidGo(t *testing.T) {
	AstTest(t, "ast: CountTests errors on invalid Go source in file", func(t *AstT) {
		dir := t.TB().TempDir()

		writeFile(t.TB(), dir+"/bad.go", `
			not go code {{{{
		`)

		_, err := tapeast.CountTests(dir, nil)
		t.Ok(err)
		t.End()
	})
}

func TestFindDuplicatesInvalidGo(t *testing.T) {
	AstTest(t, "ast: FindDuplicates errors on invalid Go source", func(t *AstT) {
		dir := t.TB().TempDir()

		writeFile(t.TB(), dir+"/bad.go", `
			not go code {{{{
		`)

		_, err := tapeast.FindDuplicates(dir, nil)
		t.Ok(err)
		t.End()
	})
}

func TestFindOnlyCallsDirInvalidGo(t *testing.T) {
	AstTest(t, "ast: FindOnlyCalls errors on invalid Go source", func(t *AstT) {
		dir := t.TB().TempDir()

		writeFile(t.TB(), dir+"/bad.go", `
			not go code {{{{
		`)

		_, err := tapeast.FindOnlyCalls(dir, nil)
		t.Ok(err)
		t.End()
	})
}

func TestWalkFilesReadFileError(t *testing.T) {
	AstTest(t, "ast: CountTests errors when file cannot be read", func(t *AstT) {
		dir := t.TB().TempDir()

		err := os.Symlink("/nonexistent", dir+"/broken.go")
		if err != nil {
			t.TB().Fatal(err)
		}

		_, err = tapeast.CountTests(dir, nil)
		t.Ok(err)
		t.End()
	})
}

func TestIsTapeCallNonCallExprNoError(t *testing.T) {
	AstTest(t, "ast: non-call expression is not a tape call no error", func(t *AstT) {
		src := `
			package foo

			import "testing"

			func TestFoo(t *testing.T) {
				(func(){})()
			}
		`

		_, err := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.NotOk(err)
		t.End()
	})
}

func TestIsTapeCallNonCallExprEmpty(t *testing.T) {
	AstTest(t, "ast: non-call expression produces no Only calls", func(t *AstT) {
		src := `
			package foo

			import "testing"

			func TestFoo(t *testing.T) {
				(func(){})()
			}
		`

		calls, _ := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.Ok(len(calls) == 0)
		t.End()
	})
}

func TestCountTestsRecursive(t *testing.T) {
	AstTest(t, "ast: CountTests counts tests in subdirectories", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test(t, "root: one", func(t *Test.T) { t.End() })
			}
		`)

		fixture("sub/sub_test.go", `
			package sub

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestSub(t *testing.T) {
				Test(t, "sub: one", func(t *Test.T) { t.End() })
			}
		`)

		n, err := tapeast.CountTests(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Equal(n, 2)
		t.End()
	})
}

func TestFindOnlyCallsRecursive(t *testing.T) {
	AstTest(t, "ast: FindOnlyCalls finds Only calls in subdirectories", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("sub/foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test.Only(t, "foo: bar", func(t *Test.T) { t.End() })
			}
		`)

		calls, err := tapeast.FindOnlyCalls(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		expected := []tapeast.OnlyCall{
			{
				Parent: "TestFoo",
				Name:   "foo: bar",
			},
		}

		t.DeepEqual(calls, expected)
		t.End()
	})
}

func TestFindDuplicatesRecursive(t *testing.T) {
	AstTest(t, "ast: FindDuplicates finds duplicates in subdirectories", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test(t, "duplicate: name", func(t *Test.T) { t.End() })
			}
		`)

		fixture("sub/sub_test.go", `
			package sub

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestSub(t *testing.T) {
				Test(t, "duplicate: name", func(t *Test.T) { t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(dups) == 1)
		t.End()
	})
	AstTest(t, "ast: FindDuplicates recursive returns correct name", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test(t, "duplicate: name", func(t *Test.T) { t.End() })
			}
		`)

		fixture("sub/sub_test.go", `
			package sub

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestSub(t *testing.T) {
				Test(t, "duplicate: name", func(t *Test.T) { t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(dups[0].Name == "duplicate: name")
		t.End()
	})
	AstTest(t, "ast: FindDuplicates recursive returns two locations", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test(t, "duplicate: name", func(t *Test.T) { t.End() })
			}
		`)

		fixture("sub/sub_test.go", `
			package sub

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestSub(t *testing.T) {
				Test(t, "duplicate: name", func(t *Test.T) { t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(dups[0].Locations) == 2)
		t.End()
	})
}

func TestCountTestsInTestFilesIgnoresNonTestFiles(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles ignores non _test.go files", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/fixture.go", `
			package foo

			func init() { tape.Test(nil, "x: y", nil) }
		`)
		n, err := tapeast.CountTestsInTestFiles(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 0)
		t.End()
	})
}

func TestCountTestsInTestFilesCountsTestFiles(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles counts tape calls in _test.go files", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: one", func(t *Test.T) { t.End() })
			}
		`)
		n, err := tapeast.CountTestsInTestFiles(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 1)
		t.End()
	})
}

func TestFindAllTestNamesReturnsNames(t *testing.T) {
	AstTest(t, "ast: FindAllTestNames returns all test names in _test.go files", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/foo_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: one", func(t *Test.T) { t.End() })
				Test(t, "foo: two", func(t *Test.T) { t.End() })
			}
		`)
		names, err := tapeast.FindAllTestNames(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(len(names), 2)
		t.End()
	})
}

func TestWalkTestFilesReadFileError(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles errors when file cannot be read", func(t *AstT) {
		dir := t.TB().TempDir()
		err := os.Symlink("/nonexistent", dir+"/broken_test.go")
		if err != nil {
			t.TB().Fatal(err)
		}
		_, err = tapeast.CountTestsInTestFiles(dir, nil)
		t.Ok(err)
		t.End()
	})
}

func TestWalkFilesIgnoresBuildIgnore(t *testing.T) {
	AstTest(t, "ast: CountTests skips files with //go:build ignore", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/ignored.go", `
			//go:build ignore

			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: bar", func(t *Test.T) { t.End() })
			}
		`)
		n, err := tapeast.CountTests(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 0)
		t.End()
	})
}

func TestWalkTestFilesIgnoresBuildIgnore(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles skips _test.go with //go:build ignore", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/ignored_test.go", `
			//go:build ignore

			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				Test(t, "foo: bar", func(t *Test.T) { t.End() })
			}
		`)
		n, err := tapeast.CountTestsInTestFiles(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 0)
		t.End()
	})
}

func TestCountTestsInTestFilesInvalidGo(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles errors on invalid Go source in _test.go", func(t *AstT) {
		dir := t.TB().TempDir()
		err := os.WriteFile(dir+"/bad_test.go", []byte("package foo\n\nnot go {{{{\n"), 0644)
		if err != nil {
			t.TB().Fatal(err)
		}
		_, err = tapeast.CountTestsInTestFiles(dir, nil)
		t.Ok(err)
		t.End()
	})
}

func TestFindAllTestNamesInvalidGo(t *testing.T) {
	AstTest(t, "ast: FindAllTestNames errors on invalid Go source in _test.go", func(t *AstT) {
		dir := t.TB().TempDir()
		err := os.WriteFile(dir+"/bad_test.go", []byte("package foo\n\nnot go {{{{\n"), 0644)
		if err != nil {
			t.TB().Fatal(err)
		}
		_, err = tapeast.FindAllTestNames(dir, nil)
		t.Ok(err)
		t.End()
	})
}

func TestWalkFilesSkipsNonGoFiles(t *testing.T) {
	AstTest(t, "ast: CountTests skips non-.go files", func(t *AstT) {
		dir := t.TB().TempDir()
		err := os.WriteFile(dir+"/readme.txt", []byte("hello"), 0644)
		if err != nil {
			t.TB().Fatal(err)
		}
		n, err := tapeast.CountTests(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 0)
		t.End()
	})
}

func TestWalkTestFilesSkipsNonTestGoFiles(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles skips non-_test.go .go files", func(t *AstT) {
		dir := t.TB().TempDir()
		err := os.WriteFile(dir+"/main.go", []byte("package main\nfunc main(){}\n"), 0644)
		if err != nil {
			t.TB().Fatal(err)
		}
		n, err := tapeast.CountTestsInTestFiles(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 0)
		t.End()
	})
}

func TestCountTestsInTestFilesMissingDir(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles errors on missing dir", func(t *AstT) {
		_, err := tapeast.CountTestsInTestFiles("nonexistent", nil)
		t.Ok(err)
		t.End()
	})
}

// -- isTestCall / isTestMethodCall branch coverage --

// TestIsTestCallFuncLiteralNoError covers line 160: isTestCall returns false
// when Fun is a function literal (not Ident or SelectorExpr) — no parse error.
func TestIsTestCallFuncLiteralNoError(t *testing.T) {
	AstTest(t, "ast: isTestCall func literal — no error", func(t *AstT) {
		src := `
			package foo

			func TestFoo(t *testing.T) {
				(func() {})()
			}
		`
		_, err := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.NotOk(err)
		t.End()
	})
}

// TestIsTestCallFuncLiteralEmpty covers line 160: isTestCall returns false
// when Fun is a function literal — produces no Only calls.
func TestIsTestCallFuncLiteralEmpty(t *testing.T) {
	AstTest(t, "ast: isTestCall func literal — no Only calls", func(t *AstT) {
		src := `
			package foo

			func TestFoo(t *testing.T) {
				(func() {})()
			}
		`
		calls, _ := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.Equal(len(calls), 0)
		t.End()
	})
}

// TestIsTestMethodCallQualifiedFormNoError covers lines 176-179: isTestMethodCall
// matches tape.Test.Only(...) — no parse error.
func TestIsTestMethodCallQualifiedFormNoError(t *testing.T) {
	AstTest(t, "ast: tape.Test.Only qualified form — no error", func(t *AstT) {
		src := `
			package foo

			import (
				"testing"
				tape "github.com/coderaiser/go-tape"
			)

			func TestFoo(t *testing.T) {
				tape.Test.Only(t, "foo: bar", func(t *tape.T) {})
			}
		`
		_, err := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.NotOk(err)
		t.End()
	})
}

// TestIsTestMethodCallQualifiedFormResult covers lines 176-179: isTestMethodCall
// matches tape.Test.Only(...) — returns the expected OnlyCall.
func TestIsTestMethodCallQualifiedFormResult(t *testing.T) {
	AstTest(t, "ast: tape.Test.Only qualified form — correct result", func(t *AstT) {
		src := `
			package foo

			import (
				"testing"
				tape "github.com/coderaiser/go-tape"
			)

			func TestFoo(t *testing.T) {
				tape.Test.Only(t, "foo: bar", func(t *tape.T) {})
			}
		`
		calls, _ := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.DeepEqual(calls, []tapeast.OnlyCall{
			{Parent: "TestFoo", Name: "foo: bar"},
		})
		t.End()
	})
}

// TestIsTestMethodCallUnknownExprFalse covers the final return false in
// isTestMethodCall: sel.X is neither *ast.Ident nor *ast.SelectorExpr.
func TestIsTestMethodCallUnknownExprFalse(t *testing.T) {
	AstTest(t, "ast: isTestMethodCall false when receiver is not Ident or SelectorExpr", func(t *AstT) {
		// An index expression like arr[0].Only(...) has a Fun whose X is *ast.IndexExpr,
		// which is neither Ident nor SelectorExpr — hits the final return false.
		src := `
			package foo

			import "testing"

			func TestFoo(t *testing.T) {
				fns := []struct{ Only func(*testing.T, string, func()) }{}
				fns[0].Only(t, "foo: bar", func() {})
			}
		`
		calls, _ := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.Equal(len(calls), 0)
		t.End()
	})
}

// TestCountTestsQualifiedOnlyForm ensures CountTests counts tape.Test.Only
// (the qualified SelectorExpr form) alongside plain Test(...) calls.
func TestCountTestsQualifiedOnlyForm(t *testing.T) {
	AstTest(t, "ast: CountTests counts tape.Test.Only qualified form", func(t *AstT) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir+"/foo_test.go", `
			package foo

			import (
				"testing"
				tape "github.com/coderaiser/go-tape"
			)

			func TestFoo(t *testing.T) {
				tape.Test(t, "foo: one", func(t *tape.T) { t.End() })
				tape.Test.Only(t, "foo: two", func(t *tape.T) { t.End() })
			}
		`)
		n, err := tapeast.CountTests(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 2)
		t.End()
	})
}

func TestCountTestsInTestFilesSkipsExcludedDirs(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles skips dirs in exclude list", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test(t, "root: one", func(t *Test.T) { t.End() })
			}
		`)

		fixture("fixture/fixture_test.go", `
			package fixture

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFixture(t *testing.T) {
				Test(t, "fixture: one", func(t *Test.T) { t.End() })
			}
		`)

		n, err := tapeast.CountTestsInTestFiles(dir, []string{"fixture"})
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 1)
		t.End()
	})
}

func TestCountTestsInTestFilesDoesNotSkipNonExcludedDirs(t *testing.T) {
	AstTest(t, "ast: CountTestsInTestFiles does not skip non-excluded dirs", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test(t, "root: one", func(t *Test.T) { t.End() })
			}
		`)

		fixture("fixture/fixture_test.go", `
			package fixture

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFixture(t *testing.T) {
				Test(t, "fixture: one", func(t *Test.T) { t.End() })
			}
		`)

		n, err := tapeast.CountTestsInTestFiles(dir, nil)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(n, 2)
		t.End()
	})
}

func TestFindDuplicatesSkipsExcludedDirs(t *testing.T) {
	AstTest(t, "ast: FindDuplicates skips excluded dirs", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test(t, "same: name", func(t *Test.T) { t.End() })
			}
		`)

		fixture("fixture/fixture_test.go", `
			package fixture

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFixture(t *testing.T) {
				Test(t, "same: name", func(t *Test.T) { t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir, []string{"fixture"})
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(len(dups), 0)
		t.End()
	})
}

func TestFindOnlyCallsSkipsExcludedDirs(t *testing.T) {
	AstTest(t, "ast: FindOnlyCalls skips excluded dirs", func(t *AstT) {
		dir, fixture := Fixture(t.TB())

		fixture("root_test.go", `
			package foo

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				Test.Only(t, "root: only", func(t *Test.T) { t.End() })
			}
		`)

		fixture("fixture/fixture_test.go", `
			package fixture

			import Test "github.com/coderaiser/go-tape"
			import "testing"

			func TestFixture(t *testing.T) {
				Test.Only(t, "fixture: only", func(t *Test.T) { t.End() })
			}
		`)

		calls, err := tapeast.FindOnlyCalls(dir, []string{"fixture"})
		if err != nil {
			t.TB().Fatal(err)
		}
		t.DeepEqual(calls, []tapeast.OnlyCall{{Parent: "TestRoot", Name: "root: only"}})
		t.End()
	})
}

func TestIsExcludedDirMatchesExactName(t *testing.T) {
	AstTest(t, "ast: isExcludedDir matches exact name", func(t *AstT) {
		t.Ok(tapeast.IsExcludedDir("fixture", []string{"fixture"}))
		t.End()
	})
}

func TestIsExcludedDirMatchesGlobPattern(t *testing.T) {
	AstTest(t, "ast: isExcludedDir matches glob pattern", func(t *AstT) {
		t.Ok(tapeast.IsExcludedDir("cmd/fixture", []string{"cmd/f*"}))
		t.End()
	})
}

func TestIsExcludedDirDoesNotMatchUnrelatedDir(t *testing.T) {
	AstTest(t, "ast: isExcludedDir does not match unrelated dir", func(t *AstT) {
		t.NotOk(tapeast.IsExcludedDir("somedir", []string{"fixture"}))
		t.End()
	})
}
