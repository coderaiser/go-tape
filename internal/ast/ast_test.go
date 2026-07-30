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

func AstTest(tb *testing.T, name string, fn func(*AstT)) {
	tape.Test(tb, name, func(base *tape.T) {
		fn(&AstT{T: base})
	})
}

func writeFile(t *testing.T, path, src string) {
	t.Helper()

	err := os.WriteFile(path, []byte(dedent.Dedent(src)), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func mkdir(t *testing.T, path string) {
	t.Helper()

	err := os.MkdirAll(path, 0755)
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

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				tape.Test(t, "foo: bar", func(t *tape.T) { t.Ok(true); t.End() })
				tape.Only(t, "foo: baz", func(t *tape.T) { t.Ok(true); t.End() })
			}
		`)

		n, err := tapeast.CountTests(dir)
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

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				tape.Test(t, "foo: bar", func(t *tape.T) { t.Ok(true); t.End() })
				tape.Test(t, "foo: bar", func(t *tape.T) { t.Ok(true); t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(dups) > 0)
		t.End()
	})
}

func TestFindOnlyCallsDir(t *testing.T) {
	AstTest(t, "ast: FindOnlyCalls reads all go files in dir", func(t *AstT) {
		dir := t.TB().TempDir()

		writeFile(t.TB(), dir+"/foo_test.go", `
			package foo

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				tape.Only(t, "foo: bar", func(t *tape.T) { t.Ok(true); t.End() })
			}
		`)

		calls, err := tapeast.FindOnlyCalls(dir)
		if err != nil {
			t.TB().Fatal(err)
		}

		t.Ok(len(calls) > 0)
		t.End()
	})
}

func TestFindOnlyCallsMissingDir(t *testing.T) {
	AstTest(t, "ast: FindOnlyCalls errors on missing dir", func(t *AstT) {
		_, err := tapeast.FindOnlyCalls("nonexistent")
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
		_, err := tapeast.CountTests("nonexistent")
		t.Ok(err)
		t.End()
	})
}

func TestFindDuplicatesMissingDir(t *testing.T) {
	AstTest(t, "ast: FindDuplicates errors on missing dir", func(t *AstT) {
		_, err := tapeast.FindDuplicates("nonexistent")
		t.Ok(err)
		t.End()
	})
}

func TestFindOnlyCallsUnqualifiedNoError(t *testing.T) {
	AstTest(t, "ast: finds Only call without package qualifier no error", func(t *AstT) {
		src := `
			package foo

			import "testing"

			func TestFoo(t *testing.T) {
				Only(t, "foo: bar", func(t *T) {})
			}
		`

		_, err := tapeast.FindOnlyCallsInSource(dedent.Dedent(src))
		t.NotOk(err)
		t.End()
	})
}

func TestFindOnlyCallsUnqualifiedResult(t *testing.T) {
	AstTest(t, "ast: finds Only call without package qualifier result", func(t *AstT) {
		src := `
			package foo

			import "testing"

			func TestFoo(t *testing.T) {
				Only(t, "foo: bar", func(t *T) {})
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

func TestWalkFilesReadError(t *testing.T) {
	AstTest(t, "ast: CountTests errors on nonexistent directory", func(t *AstT) {
		_, err := tapeast.CountTests("/nonexistent/path/that/does/not/exist")
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

		_, err := tapeast.CountTests(dir)
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

		_, err := tapeast.FindDuplicates(dir)
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

		_, err := tapeast.FindOnlyCalls(dir)
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

		_, err = tapeast.CountTests(dir)
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

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				tape.Test(t, "root: one", func(t *tape.T) { t.End() })
			}
		`)

		fixture("sub/sub_test.go", `
			package sub

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestSub(t *testing.T) {
				tape.Test(t, "sub: one", func(t *tape.T) { t.End() })
			}
		`)

		n, err := tapeast.CountTests(dir)
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

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestFoo(t *testing.T) {
				tape.Only(t, "foo: bar", func(t *tape.T) { t.End() })
			}
		`)

		calls, err := tapeast.FindOnlyCalls(dir)
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

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestRoot(t *testing.T) {
				tape.Test(t, "duplicate: name", func(t *tape.T) { t.End() })
			}
		`)

		fixture("sub/sub_test.go", `
			package sub

			import tape "github.com/coderaiser/go-tape"
			import "testing"

			func TestSub(t *testing.T) {
				tape.Test(t, "duplicate: name", func(t *tape.T) { t.End() })
			}
		`)

		dups, err := tapeast.FindDuplicates(dir)
		if err != nil {
			t.TB().Fatal(err)
		}

		expected := []string{"duplicate: name"}

		t.DeepEqual(dups, expected)
		t.End()
	})
}
