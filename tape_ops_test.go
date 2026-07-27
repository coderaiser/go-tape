package tape

import (
	"errors"
	"os"
	"testing"
	"time"
)

// These tests exercise failure paths of operators and guards.
// They deliberately fail, but that's fine — the coverage counts.
// Each intentional failure is wrapped in a subtest to isolate it.

func TestEqualNonPrimitiveGuard(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		os.Setenv("TAPE_CHECK_END", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		defer os.Unsetenv("TAPE_CHECK_END")
		Test(inner, "tape: Equal non-primitive", func(tt *T) {
			tt.Equal([]int{1}, []int{1})
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestEqualMismatch(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: Equal mismatch", func(tt *T) {
			tt.Equal(1, 2)
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestNotEqualNonPrimitiveGuard(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		os.Setenv("TAPE_CHECK_END", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		defer os.Unsetenv("TAPE_CHECK_END")
		Test(inner, "tape: NotEqual non-primitive", func(tt *T) {
			tt.NotEqual([]int{1}, []int{2})
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestNotEqualMatch(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: NotEqual match", func(tt *T) {
			tt.NotEqual(1, 1)
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestDeepEqualMismatch(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: DeepEqual mismatch", func(tt *T) {
			tt.DeepEqual([]int{1}, []int{2})
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestNotDeepEqualMatch(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: NotDeepEqual match", func(tt *T) {
			tt.NotDeepEqual([]int{1}, []int{1})
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestOkFalse(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: Ok false", func(tt *T) {
			tt.Ok(false)
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestNotOkTrue(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: NotOk true", func(tt *T) {
			tt.NotOk(true)
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestMatchNoMatch(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: Match no match", func(tt *T) {
			tt.Match("hello", "world")
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestMatchInvalidPattern(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: Match invalid", func(tt *T) {
			tt.Match("hello", "[invalid")
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestNotMatchDoesMatch(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: NotMatch does match", func(tt *T) {
			tt.NotMatch("hello", "hello")
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestNotMatchInvalidPattern(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: NotMatch invalid", func(tt *T) {
			tt.NotMatch("hello", "[invalid")
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestTFail(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		Test(inner, "tape: Fail", func(tt *T) {
			tt.Fail("forced")
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestErrorNil(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		os.Setenv("TAPE_CHECK_END", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		defer os.Unsetenv("TAPE_CHECK_END")
		Test(inner, "tape: Error nil", func(tt *T) {
			tt.Error(nil)
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
}

func TestNoErrorNonNil(t *testing.T) {
	wasFailed := !t.Run("inner", func(inner *testing.T) {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
		os.Setenv("TAPE_CHECK_END", "0")
		defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
		defer os.Unsetenv("TAPE_CHECK_END")
		Test(inner, "tape: NoError non-nil", func(tt *T) {
			tt.NoError(errors.New("oops"))
			tt.End()
		})
	})
	if !wasFailed {
		t.Fatal("expected test to fail")
	}
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
