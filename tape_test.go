package tape

import (
	"errors"
	"os"
	"testing"
)

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

func TestAssertionsCountDisabled(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: two assertions", func(t *T) {
		t.Ok(true)
		t.Ok(true)
		t.End()
	})
}

func TestEndCheckDisabled(t *testing.T) {
	os.Setenv("TAPE_CHECK_END", "0")
	defer os.Unsetenv("TAPE_CHECK_END")
	Test(t, "tape: no end check", func(t *T) {
		t.Ok(true)
	})
}

func TestTEqual(t *testing.T) {
	Test(t, "tape: Equal works", func(t *T) {
		t.Equal(42, 42)
		t.End()
	})
}

func TestTEqualMismatch(t *testing.T) {
	// Errorf doesn't stop execution, this runs fine
	Test(t, "tape: Equal mismatch", func(t *T) {
		t.Equal(1, 2)
		t.End()
	})
}

func TestTOk(t *testing.T) {
	Test(t, "tape: Ok works", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestTOkFalse(t *testing.T) {
	Test(t, "tape: Ok false", func(t *T) {
		t.Ok(false)
		t.End()
	})
}

func TestTNotOk(t *testing.T) {
	Test(t, "tape: NotOk works", func(t *T) {
		t.NotOk(false)
		t.End()
	})
}

func TestTNotOkTrue(t *testing.T) {
	Test(t, "tape: NotOk true", func(t *T) {
		t.NotOk(true)
		t.End()
	})
}

func TestTDeepEqual(t *testing.T) {
	Test(t, "tape: DeepEqual works", func(t *T) {
		t.DeepEqual([]int{1, 2}, []int{1, 2})
		t.End()
	})
}

func TestTDeepEqualMismatch(t *testing.T) {
	Test(t, "tape: DeepEqual mismatch", func(t *T) {
		t.DeepEqual([]int{1}, []int{2})
		t.End()
	})
}

func TestTContains(t *testing.T) {
	Test(t, "tape: Contains works", func(t *T) {
		t.Contains("hello world", "world")
		t.End()
	})
}

func TestTContainsNoMatch(t *testing.T) {
	Test(t, "tape: Contains no match", func(t *T) {
		t.Contains("hello", "xyz")
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
	Test(t, "tape: Match works", func(t *T) {
		t.Match("hello 123", `hello \d+`)
		t.End()
	})
}

func TestTMatchNoMatch(t *testing.T) {
	Test(t, "tape: Match no match", func(t *T) {
		t.Match("hello", `\d+`)
		t.End()
	})
}

func TestExtend(t *testing.T) {
	called := false
	ext := Extensions{func(t *T) { called = true }}
	Extend(ext)(t, "tape: extend works", func(t *T) {
		t.Ok(called)
		t.End()
	})
}

func TestHelperFunctions(t *testing.T) {
	if !deepEqual(1, 1) { t.Error("fail") }
	if deepEqual(1, 2) { t.Error("fail") }
	if !contains("hello", "hell") { t.Error("fail") }
	if contains("hello", "xyz") { t.Error("fail") }
	if !match("abc123", `abc\d+`) { t.Error("fail") }
	if match("abc", `\d+`) { t.Error("fail") }
	if match("hello", `[invalid`) { t.Error("fail") }
}

func TestOnly(t *testing.T) {
	Test(t, "tape: Only test", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestSkip(t *testing.T) {
	Skip(t, "tape: skip test", func(t *T) {
		t.Ok(true)
		t.End()
	})
}

func TestHitWithDisabledCount(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	
	assertOne(t)
	hit(t)
	hit(t) // Should not fail because count check disabled
}

func TestEndCalledFlag(t *testing.T) {
	tt := newT(t)
	if tt.ended {
		t.Error("ended should be false initially")
	}
	tt.End()
	if !tt.ended {
		t.Error("ended should be true after call")
	}
}

func TestNewT(t *testing.T) {
	tt := newT(t)
	if tt == nil {
		t.Fatal("expected non-nil T")
	}
}
