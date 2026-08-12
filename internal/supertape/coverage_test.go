package supertape_test

import (
	"testing"
	"time"

	. "github.com/coderaiser/go-tape"
	st "github.com/coderaiser/go-tape/internal/supertape"
)

// TestOptions covers the config Option constructors.
func TestOptions(t *testing.T) {
	Test(t, "supertape: options set runner fields", func(t *T) {
		r := st.New(
			st.WithCheckAssertionsCount(false),
			st.WithCheckEnd(false),
			st.WithCheckScopes(false),
			st.WithTimeout(time.Second),
		)
		r.Test("scope: ok", func(tt st.T) {
			tt.Ok(true)
			tt.Ok(true)
			tt.End()
		})
		passed, failed, _, _ := r.Counts()
		t.Equal(passed == 1 && failed == 0, true)
		t.End()
	})
}

// TestSkip covers the Skip method.
func TestSkip(t *testing.T) {
	Test(t, "supertape: Skip records skipped and emits test-end", func(t *T) {
		r, evs := newCollector()
		r.Skip("scope: skipped")
		_, _, skipped, total := r.Counts()
		t.Equal(skipped == 1 && total == 1 && len(*evs) == 1, true)
		t.End()
	})
}

// TestInvalidScope covers the scope guard.
func TestInvalidScope(t *testing.T) {
	Test(t, "supertape: invalid scope name records fail", func(t *T) {
		r, evs := newCollector()
		r.Test("not-a-scope", func(tt st.T) {
			tt.Ok(true)
			tt.End()
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed == 1 && (*evs)[0].Operator == "scope", true)
		t.End()
	})
}

// TestEndNotCalled covers the End() requirement guard.
func TestEndNotCalled(t *testing.T) {
	Test(t, "supertape: missing End records fail", func(t *T) {
		r, evs := newCollector()
		r.Test("scope: noend", func(tt st.T) {
			tt.Ok(true)
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed == 1 && (*evs)[0].Operator == "end", true)
		t.End()
	})
}

// TestTooManyAssertions covers the assertion-count guard.
func TestTooManyAssertions(t *testing.T) {
	Test(t, "supertape: more than one assertion records fail", func(t *T) {
		r, evs := newCollector()
		r.Test("scope: toomany", func(tt st.T) {
			tt.Ok(true)
			tt.Ok(true)
			tt.End()
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed == 1 && (*evs)[0].Operator == "ok", true)
		t.End()
	})
}

// TestTimeout covers the timeout guard.
func TestTimeout(t *testing.T) {
	Test(t, "supertape: timeout records fail", func(t *T) {
		evs := &[]st.Event{}
		r := st.New(
			st.WithTimeout(time.Millisecond*5),
			st.WithCheckEnd(false),
			st.WithCheckAssertionsCount(false),
			st.WithHandler(func(ev st.Event) { *evs = append(*evs, ev) }),
		)
		r.Test("scope: slow", func(tt st.T) {
			time.Sleep(time.Second)
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed == 1 && (*evs)[0].Operator == "timeout", true)
		t.End()
	})
}

// failErr is an error value used to exercise the Fail(error) branch.
type failErr string

func (e failErr) Error() string { return string(e) }

// TestPackageLevelTest covers the package-level Test function (default runner).
func TestPackageLevelTest(t *testing.T) {
	Test(t, "supertape: package-level Test runs via default runner", func(t *T) {
		st.Test("scope: package", func(tt st.T) {
			tt.Ok(true)
			tt.End()
		})
		t.End()
	})
}

// TestFailBranches covers the Fail string/error/default branches.
func TestFailBranches(t *testing.T) {
	Test(t, "supertape: Fail string, error and default branches", func(t *T) {
		r := st.New()
		r.Test("scope: failstr", func(tt st.T) { tt.Fail("boom"); tt.End() })
		r.Test("scope: failerr", func(tt st.T) { tt.Fail(failErr("custom")); tt.End() })
		r.Test("scope: faildef", func(tt st.T) { tt.Fail(42); tt.End() })
		_, failed, _, _ := r.Counts()
		t.Equal(failed == 3, true)
		t.End()
	})
}

// TestCallerLocation covers the caller helper via a failing assertion.
func TestCallerLocation(t *testing.T) {
	Test(t, "supertape: assertion stamps caller location", func(t *T) {
		r, evs := newCollector()
		r.Test("scope: loc", func(tt st.T) {
			tt.Equal(1, 2)
			tt.End()
		})
		t.Ok((*evs)[0].At != "")
		t.End()
	})
}

// TestNotEqualFail covers a failing NotEqual for the report !Ok path variety.
func TestNotEqualFail(t *testing.T) {
	Test(t, "supertape: failing NotEqual records fail", func(t *T) {
		r := st.New()
		r.Test("scope: noteqfail", func(tt st.T) {
			tt.NotEqual(1, 1)
			tt.End()
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed, 1)
		t.End()
	})
}

// TestNotDeepEqualFail covers a failing NotDeepEqual.
func TestNotDeepEqualFail(t *testing.T) {
	Test(t, "supertape: failing NotDeepEqual records fail", func(t *T) {
		r := st.New()
		r.Test("scope: notdeepeq", func(tt st.T) {
			tt.NotDeepEqual([]int{1}, []int{1})
			tt.End()
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed, 1)
		t.End()
	})
}