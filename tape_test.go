package tape

import (
	"errors"
	"os"
	"regexp"
	"testing"
)

// -- scope guard --

func TestValidScopePasses(t *testing.T) {
	Test(t, "tape: valid scope runs fine", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestScopesDisabled(t *testing.T) {
	os.Setenv("TAPE_CHECK_SCOPES", "0")
	defer os.Unsetenv("TAPE_CHECK_SCOPES")
	Test(t, "no scope format", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

// -- assertion count guard --

func TestAssertionsCountDisabled(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: two assertions", func(t *T) {
		t.Ok(true)
		t.Ok(true)
		t.End()
	})
}

// -- end guard --

func TestEndCheckDisabled(t *testing.T) {
	os.Setenv("TAPE_CHECK_END", "0")
	defer os.Unsetenv("TAPE_CHECK_END")
	Test(t, "tape: no end check", func(t *T) {
		t.Ok(true)
	})
}

// -- operators happy path --

func TestTEqual(t *testing.T) {
	Test(t, "tape: Equal works", func(t *T) {
		t.Equal(42, 42)
		t.End()
	})
}

func TestTNotEqual(t *testing.T) {
	Test(t, "tape: NotEqual works", func(t *T) {
		t.NotEqual(1, 2)
		t.End()
	})
}

func TestTOk(t *testing.T) {
	Test(t, "tape: Ok works", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestTNotOk(t *testing.T) {
	Test(t, "tape: NotOk works", func(t *T) {
		t.NotOk(false)
		t.End()
	})
}

func TestTDeepEqual(t *testing.T) {
	Test(t, "tape: DeepEqual works", func(t *T) {
		t.DeepEqual([]int{1, 2}, []int{1, 2})
		t.End()
	})
}

func TestTNotDeepEqual(t *testing.T) {
	Test(t, "tape: NotDeepEqual works", func(t *T) {
		t.NotDeepEqual([]int{1}, []int{2})
		t.End()
	})
}

func TestTMatch(t *testing.T) {
	Test(t, "tape: Match works with string pattern", func(t *T) {
		t.Match("hello 123", regexp.MustCompile(`hello \d+`))
		t.End()
	})
}

func TestTMatchRegexp(t *testing.T) {
	Test(t, "tape: Match works with regexp pattern", func(t *T) {
		t.Match("hello 123", regexp.MustCompile(`hello \d+`))
		t.End()
	})
}

func TestTNotMatch(t *testing.T) {
	Test(t, "tape: NotMatch works", func(t *T) {
		t.NotMatch("hello", `\d+`)
		t.End()
	})
}

func TestTPass(t *testing.T) {
	Test(t, "tape: Pass works", func(t *T) {
		t.Pass("all good")
		t.End()
	})
}

func TestTComment(t *testing.T) {
	Test(t, "tape: Comment does not count as assertion", func(t *T) {
		t.Comment("just a note")
		t.Ok(true)
		t.End()
	})
}

// -- ReportCustom / public Result --

func TestTReportCustomPass(t *testing.T) {
	Test(t, "tape: ReportCustom pass counts assertion", func(t *T) {
		t.ReportCustom(true, "transform", "", "got", "expected")
		t.End()
	})
}

func TestTReportCustomFail(t *testing.T) {
	tt := &T{t: &testing.T{}}
	tt.ReportCustom(false, "transform", "diff", "got", "expected")
	if !tt.t.Failed() {
		t.Fatal("expected failure recorded")
	}
}

func TestBuiltinOperatorsPublic(t *testing.T) {
	Test(t, "tape: BuiltinOperators return public results", func(t *T) {
		equal := BuiltinOperators.Equal(1, 1)
		notEqual := BuiltinOperators.NotEqual(1, 2)
		deepEqual := BuiltinOperators.DeepEqual([]int{1}, []int{1})
		notDeepEqual := BuiltinOperators.NotDeepEqual([]int{1}, []int{2})
		ok := BuiltinOperators.Ok(true)
		notOk := BuiltinOperators.NotOk(false)
		match := BuiltinOperators.Match("hello", "hello")
		notMatch := BuiltinOperators.NotMatch("hello", `\d+`)
		pass := BuiltinOperators.Pass("all good")
		fail := BuiltinOperators.Fail("boom")
		t.Ok(equal.Ok && notEqual.Ok && deepEqual.Ok && notDeepEqual.Ok && ok.Ok &&
			notOk.Ok && match.Ok && notMatch.Ok && pass.Ok && !fail.Ok)
		t.End()
	})
}
func TestBuiltinOperatorsOutputContainsDiff(t *testing.T) {
	Test(t, "tape: BuiltinOperators.Equal returns Result with Ok false", func(t *T) {
		r := BuiltinOperators.Equal("got", "expected")
		t.Ok(r.Ok == false)
		t.End()
	})
	Test(t, "tape: BuiltinOperators.Equal returns Result with non-empty Output", func(t *T) {
		r := BuiltinOperators.Equal("got", "expected")
		t.Ok(r.Output != "")
		t.End()
	})
}


func TestTestSkipNoop(t *testing.T) {
	ran := false
	Test.Skip(t, "tape: skipped block", func(t *T) { ran = true })
	if ran {
		t.Fatal("Skip must not run the function")
	}
}

func TestTestOnlyMethod(t *testing.T) {
	ran := false
	Test.Only(t, "tape: only block", func(t *T) {
		ran = true
		t.End()
	})
	if !ran {
		t.Fatal("Only must run the function")
	}
}

// -- helpers --

func TestHelperIsPrimitive(t *testing.T) {
	Test(t, "tape: isPrimitive true for int", func(t *T) {
		t.Ok(isPrimitive(42))
		t.End()
	})
}

func TestHelperIsNotPrimitive(t *testing.T) {
	Test(t, "tape: isPrimitive false for slice", func(t *T) {
		t.NotOk(isPrimitive([]int{1}))
		t.End()
	})
}

func TestHelperTruthy(t *testing.T) {
	Test(t, "tape: Ok true for true", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestHelperFalsy(t *testing.T) {
	Test(t, "tape: NotOk false for false", func(t *T) {
		t.NotOk(false)
		t.End()
	})
}

// -- Test / Test.Skip --

func TestOnly(t *testing.T) {
	Test(t, "tape: Test runs the test", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestSkip(t *testing.T) {
	Test.Skip(t, "tape: Test.Skip skips the test", func(t *T) {
		t.Ok(false)
		t.End()
	})
}

func TestSkipDoesNotSkipParent(t *testing.T) {
	passed := false
	Test.Skip(t, "tape: Test.Skip does not skip parent", func(t *T) {
		t.Ok(false)
		t.End()
	})
	passed = true
	if !passed {
		t.Fatal("parent should not have been skipped")
	}
}

func TestSkipDoesNotRunFn(t *testing.T) {
	ran := false
	Test.Skip(t, "tape: Test.Skip fn must not run", func(t *T) {
		ran = true
		t.End()
	})
	if ran {
		t.Fatal("Test.Skip must not execute fn")
	}
}

func TestSkipParentIsNotMarkedSkipped(t *testing.T) {
	Test.Skip(t, "tape: Test.Skip parent stays clean", func(t *T) {
		t.End()
	})
	if t.Skipped() {
		t.Fatal("parent test must not be marked skipped by Test.Skip()")
	}
}

// -- Extend --

func TestExtend(t *testing.T) {
	called := false
	type myT struct{ *T }
	factory := func(base *T) *myT { called = true; return &myT{T: base} }
	Extend(factory)(t, "tape: extend works", func(t *myT) {
		t.Ok(called)
		t.End()
	})
}

func TestExtendFnSkipIsNoop(t *testing.T) {
	Test(t, "tape: ExtendFn.Skip does not run the body", func(t *T) {
		ran := false
		ef := Extend(func(base *T) *T { return base })
		ef.Skip(t.TB(), "tape: skip test", func(_ *T) { ran = true })
		t.Equal(ran, false)
		t.End()
	})
}

func TestExtendFnOnlyRuns(t *testing.T) {
	Test(t, "tape: ExtendFn.Only runs the body", func(t *T) {
		ran := false
		ef := Extend(func(base *T) *T { return base })
		ef.Only(t.TB(), "tape: only test", func(inner *T) {
			ran = true
			inner.End()
		})
		t.Ok(ran)
		t.End()
	})
}

func TestOnlyGuardSkipsOthers(t *testing.T) {
	ran := false
	Test.Only(t, "tape: only guard: target", func(t *T) {
		t.Ok(true)
		t.End()
	})
	Test(t, "tape: only guard: other", func(t *T) {
		ran = true
		t.End()
	})
	if ran {
		t.Fatal("Only guard: non-matching test ran its body")
	}
}

// -- internal helpers coverage --

func TestAssertOneCleanup(t *testing.T) {
	inner := &testing.T{}
	assertOne(inner)
	mu.Lock()
	_, exists := count[inner]
	mu.Unlock()
	if !exists {
		t.Fatal("expected count entry after assertOne")
	}
	mu.Lock()
	delete(count, inner)
	mu.Unlock()
	mu.Lock()
	_, exists = count[inner]
	mu.Unlock()
	if exists {
		t.Fatal("expected count entry deleted after cleanup")
	}
}

func TestHitWithDisabledCount(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	assertOne(t)
	hit(t)
	hit(t)
}

func TestEndCalledFlag(t *testing.T) {
	tt := newT(t)
	if tt.ended {
		t.Error("ended should be false initially")
	}
	tt.End()
	if !tt.ended {
		t.Fatal("ended should be true after End()")
	}
}

func TestNewT(t *testing.T) {
	tt := newT(t)
	if tt == nil {
		t.Fatal("expected non-nil T")
	}
}

func TestTEqualPointer(t *testing.T) {
	x := 42
	Test(t, "tape: Equal pointer", func(tt *T) {
		tt.Equal(&x, &x)
		tt.End()
	})
}

// -- isPrimitive coverage --

func TestIsPrimitivePointer(t *testing.T) {
	x := 42
	Test(t, "tape: isPrimitive true for pointer", func(tt *T) {
		tt.Ok(isPrimitive(&x))
		tt.End()
	})
}

func TestIsPrimitiveNil(t *testing.T) {
	Test(t, "tape: isPrimitive false for nil", func(tt *T) {
		tt.NotOk(isPrimitive(nil))
		tt.End()
	})
}

// -- truthy coverage via t.Ok / t.NotOk --

func TestTruthyNil(t *testing.T) {
	Test(t, "tape: NotOk false for nil", func(tt *T) {
		tt.NotOk(nil)
		tt.End()
	})
}

func TestTruthyIntZero(t *testing.T) {
	Test(t, "tape: NotOk false for int 0", func(tt *T) {
		tt.NotOk(0)
		tt.End()
	})
}

func TestTruthyStringEmpty(t *testing.T) {
	Test(t, "tape: NotOk false for empty string", func(tt *T) {
		tt.NotOk("")
		tt.End()
	})
}

func TestTruthyDefault(t *testing.T) {
	Test(t, "tape: Ok true for struct", func(tt *T) {
		tt.Ok(struct{}{})
		tt.End()
	})
}

// -- toRegexp coverage --

func TestToRegexpInvalidRegex(t *testing.T) {
	_, err := toRegexp("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestToRegexpInvalidType(t *testing.T) {
	_, err := toRegexp(42)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestOnlyRuns(t *testing.T) {
	ran := false
	Test(t, "tape: Test delegates to Test", func(t *T) {
		ran = true
		t.Ok(ran)
		t.End()
	})
	if !ran {
		t.Fatal("Test did not run the function")
	}
}

func TestTFailString(t *testing.T) {
	tt := &T{t: &testing.T{}}
	tt.Fail("forced failure")
}

func TestTFailError(t *testing.T) {
	tt := &T{t: &testing.T{}}
	tt.Fail(errors.New("error message"))
}

func TestTFailDefault(t *testing.T) {
	tt := &T{t: &testing.T{}}
	tt.Fail(42)
}

func TestTPassNoArgs(t *testing.T) {
	Test(t, "tape: Pass with no args", func(tt *T) {
		tt.Pass()
		tt.End()
	})
}

func TestTestOnlyCallsFnDirectly(t *testing.T) {
	t.Setenv("TAPE_CHECK_SCOPES", "0")
	t.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	t.Setenv("TAPE_CHECK_END", "0")
	ran := false
	Test(t, "tape: Test calls fn", func(t *T) {
		ran = true
		t.End()
	})
	if !ran {
		t.Fatal("Test did not call fn")
	}
}

func TestToRegexpWithRegexpType(t *testing.T) {
	re := regexp.MustCompile("hello")
	got, err := toRegexp(re)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if got != re {
		t.Fatal("expected same regexp")
	}
}
