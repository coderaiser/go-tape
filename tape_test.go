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

func TestTError(t *testing.T) {
	Test(t, "tape: Error works", func(t *T) {
		t.Error(errors.New("some error"))
		t.End()
	})
}

func TestTNoError(t *testing.T) {
	Test(t, "tape: NoError works", func(t *T) {
		t.NoError(nil)
		t.End()
	})
}

func TestTMatch(t *testing.T) {
	Test(t, "tape: Match works with string pattern", func(t *T) {
		t.Match("hello 123", `hello \d+`)
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
	Test(t, "tape: truthy true for true", func(t *T) {
		t.Ok(truthy(true))
		t.End()
	})
}

func TestHelperFalsy(t *testing.T) {
	Test(t, "tape: truthy false for false", func(t *T) {
		t.NotOk(truthy(false))
		t.End()
	})
}

// -- Only / Skip --

func TestOnly(t *testing.T) {
	Only(t, "tape: Only runs the test", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestSkip(t *testing.T) {
	Skip(t, "tape: this test is skipped", func(t *T) {
		t.Ok(false)
		t.End()
	})
}

// -- Extend --

func TestExtend(t *testing.T) {
	called := false
	ext := Extensions{func(t *T) { called = true }}
	Extend(ext)(t, "tape: extend works", func(t *T) {
		t.Ok(called)
		t.End()
	})
}

// -- internal helpers coverage --

func TestHitWithDisabledCount(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	assertOne(t)
	hit(t)
	hit(t) // should not fail — count check disabled
}

func TestEndCalledFlag(t *testing.T) {
	tt := newT(t)
	if tt.ended {
		t.Error("ended should be false initially")
	}
	tt.End()
	if !tt.ended {
		t.Error("ended should be true after End()")
	}
}

func TestNewT(t *testing.T) {
	tt := newT(t)
	if tt == nil {
		t.Fatal("expected non-nil T")
	}
}
