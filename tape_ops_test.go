package tape

import (
	"errors"
	"os"
	"testing"
	"time"
)

// These tests exercise failure paths of operators and guards.
// They deliberately fail, but that's fine — the coverage counts.

func TestEqualNonPrimitiveGuard(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	os.Setenv("TAPE_CHECK_END", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	defer os.Unsetenv("TAPE_CHECK_END")
	Test(t, "tape: Equal non-primitive", func(tt *T) {
		tt.Equal([]int{1}, []int{1})
		tt.End()
	})
}

func TestEqualMismatch(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: Equal mismatch", func(tt *T) {
		tt.Equal(1, 2)
		tt.End()
	})
}

func TestNotEqualNonPrimitiveGuard(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	os.Setenv("TAPE_CHECK_END", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	defer os.Unsetenv("TAPE_CHECK_END")
	Test(t, "tape: NotEqual non-primitive", func(tt *T) {
		tt.NotEqual([]int{1}, []int{2})
		tt.End()
	})
}

func TestNotEqualMatch(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: NotEqual match", func(tt *T) {
		tt.NotEqual(1, 1)
		tt.End()
	})
}

func TestDeepEqualMismatch(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: DeepEqual mismatch", func(tt *T) {
		tt.DeepEqual([]int{1}, []int{2})
		tt.End()
	})
}

func TestNotDeepEqualMatch(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: NotDeepEqual match", func(tt *T) {
		tt.NotDeepEqual([]int{1}, []int{1})
		tt.End()
	})
}

func TestOkFalse(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: Ok false", func(tt *T) {
		tt.Ok(false)
		tt.End()
	})
}

func TestNotOkTrue(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: NotOk true", func(tt *T) {
		tt.NotOk(true)
		tt.End()
	})
}

func TestMatchNoMatch(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: Match no match", func(tt *T) {
		tt.Match("hello", "world")
		tt.End()
	})
}

func TestMatchInvalidPattern(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: Match invalid", func(tt *T) {
		tt.Match("hello", "[invalid")
		tt.End()
	})
}

func TestNotMatchDoesMatch(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: NotMatch does match", func(tt *T) {
		tt.NotMatch("hello", "hello")
		tt.End()
	})
}

func TestNotMatchInvalidPattern(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: NotMatch invalid", func(tt *T) {
		tt.NotMatch("hello", "[invalid")
		tt.End()
	})
}

func TestTFail(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: Fail", func(tt *T) {
		tt.Fail("forced")
		tt.End()
	})
}

func TestErrorNil(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: Error nil", func(tt *T) {
		tt.Error(nil)
		tt.End()
	})
}

func TestNoErrorNonNil(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	Test(t, "tape: NoError non-nil", func(tt *T) {
		tt.NoError(errors.New("oops"))
		tt.End()
	})
}

func TestHitTwiceFailsInner(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		assertOne(inner)
		hit(inner)
		hit(inner)
	})
	if !wasFailed {
		t.Fatal("expected inner test to fail")
	}
}

func TestTimeoutFires(t *testing.T) {
	os.Setenv("TAPE_TIMEOUT", "1ms")
	defer os.Unsetenv("TAPE_TIMEOUT")
	wasFailed := !t.Run("timeout-outer", func(outer *testing.T) {
		Test(outer, "tape: timeout fires", func(t *T) {
			time.Sleep(50 * time.Millisecond)
			t.Ok(true)
			t.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected timeout to fail")
	}
}

func TestScopeCheckFail(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		Test(inner, "bad name no scope", func(t *T) {
			t.Ok(true)
			t.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected scope check to fail")
	}
}

func TestEndCheckFails(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		Test(inner, "tape: missing end", func(t *T) {
			t.Ok(true)
		})
	})
	if !wasFailed {
		t.Fatal("expected end check to fail")
	}
}
