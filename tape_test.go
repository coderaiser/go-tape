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

	Test(t, "tape: two assertions allowed", func(t *T) {
		t.Ok(true)
		t.Ok(true)
		t.End()
	})
}

func TestEndCheckDisabled(t *testing.T) {
	os.Setenv("TAPE_CHECK_END", "0")
	defer os.Unsetenv("TAPE_CHECK_END")

	Test(t, "tape: no end with check disabled", func(t *T) {
		t.Ok(true)
	})
}

func TestTEqual(t *testing.T) {
	Test(t, "tape: Equal works", func(t *T) {
		t.Equal(42, 42)
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

func TestTContains(t *testing.T) {
	Test(t, "tape: Contains works", func(t *T) {
		t.Contains("hello world", "world")
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

func TestExtend(t *testing.T) {
	called := false
	ext := Extensions{
		func(t *T) {
			called = true
		},
	}
	Extend(ext)(t, "tape: extend works", func(t *T) {
		t.Ok(called)
		t.End()
	})
}

func TestDeepEqualHelper(t *testing.T) {
	if !deepEqual(1, 1) {
		t.Error("expected equal")
	}
}

func TestDeepEqualMismatch(t *testing.T) {
	if deepEqual(1, 2) {
		t.Error("expected not equal")
	}
}

func TestContainsHelper(t *testing.T) {
	if !contains("hello", "hell") {
		t.Error("expected contains")
	}
}

func TestContainsNoMatch(t *testing.T) {
	if contains("hello", "xyz") {
		t.Error("expected not contains")
	}
}

func TestMatchHelper(t *testing.T) {
	if !match("abc123", `abc\d+`) {
		t.Error("expected match")
	}
}

func TestMatchNoMatch(t *testing.T) {
	if match("abc", `\d+`) {
		t.Error("expected no match")
	}
}

func TestMatchInvalidPattern(t *testing.T) {
	if match("hello", `[invalid`) {
		t.Error("expected false for invalid pattern")
	}
}

func TestHitOneCall(t *testing.T) {
	assertOne(t)
	hit(t)
}

func TestOnly(t *testing.T) {
	Test(t, "tape: Only test", func(t *T) {
		t.Ok(true)
		t.End()
	})
}
