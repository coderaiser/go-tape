package supertape_test

import (
	"regexp"
	"testing"

	. "github.com/coderaiser/go-tape"
	st "github.com/coderaiser/go-tape/internal/supertape"
)

// newCollector returns a st.Runner plus a slice that receives emitted events.
func newCollector() (*st.Runner, *[]st.Event) {
	events := &[]st.Event{}
	r := st.New(st.WithHandler(func(ev st.Event) {
		*events = append(*events, ev)
	}))
	return r, events
}

// -- basic happy path --

func TestEqualPass(t *testing.T) {
	Test(t, "supertape: Equal passes for equal values", func(t *T) {
		r := st.New()
		r.Test("supertape: inner equal", func(tt st.T) {
			tt.Equal(1, 1)
			tt.End()
		})
		passed, failed, _, _ := r.Counts()
		t.Equal(passed, 1)
		_ = failed
		t.End()
	})
}

func TestNotEqualPass(t *testing.T) {
	Test(t, "supertape: NotEqual passes for different values", func(t *T) {
		r := st.New()
		r.Test("supertape: inner notEqual", func(tt st.T) {
			tt.NotEqual(1, 2)
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestDeepEqualPass(t *testing.T) {
	Test(t, "supertape: DeepEqual passes for deeply equal values", func(t *T) {
		r := st.New()
		r.Test("supertape: inner deepEqual", func(tt st.T) {
			tt.DeepEqual([]int{1, 2}, []int{1, 2})
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestNotDeepEqualPass(t *testing.T) {
	Test(t, "supertape: NotDeepEqual passes for non-equal values", func(t *T) {
		r := st.New()
		r.Test("supertape: inner notDeepEqual", func(tt st.T) {
			tt.NotDeepEqual([]int{1}, []int{2})
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestOkPass(t *testing.T) {
	Test(t, "supertape: Ok passes for truthy value", func(t *T) {
		r := st.New()
		r.Test("supertape: inner ok", func(tt st.T) {
			tt.Ok(true)
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestNotOkPass(t *testing.T) {
	Test(t, "supertape: NotOk passes for falsy value", func(t *T) {
		r := st.New()
		r.Test("supertape: inner notOk", func(tt st.T) {
			tt.NotOk(false)
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestMatchPass(t *testing.T) {
	Test(t, "supertape: Match passes with string pattern", func(t *T) {
		r := st.New()
		r.Test("supertape: inner match", func(tt st.T) {
			tt.Match("hello 123", "hello 123")
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestMatchRegexpPass(t *testing.T) {
	Test(t, "supertape: Match passes with *regexp.Regexp", func(t *T) {
		r := st.New()
		r.Test("supertape: inner match regexp", func(tt st.T) {
			tt.Match("hello 123", regexp.MustCompile(`hello \d+`))
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestNotMatchPass(t *testing.T) {
	Test(t, "supertape: NotMatch passes when no match", func(t *T) {
		r := st.New()
		r.Test("supertape: inner notMatch", func(tt st.T) {
			tt.NotMatch("hello", `\d+`)
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestPassPass(t *testing.T) {
	Test(t, "supertape: Pass passes with message", func(t *T) {
		r := st.New()
		r.Test("supertape: inner pass", func(tt st.T) {
			tt.Pass("all good")
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestPassNoArgs(t *testing.T) {
	Test(t, "supertape: Pass passes with no args", func(t *T) {
		r := st.New()
		r.Test("supertape: inner pass noargs", func(tt st.T) {
			tt.Pass()
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

func TestCommentDoesNotCount(t *testing.T) {
	Test(t, "supertape: Comment does not count as assertion", func(t *T) {
		r := st.New()
		r.Test("supertape: inner comment", func(tt st.T) {
			tt.Comment("note")
			tt.Ok(true)
			tt.End()
		})
		passed, _, _, _ := r.Counts()
		t.Equal(passed, 1)
		t.End()
	})
}

// -- failures --

func TestEqualFailRecordsFail(t *testing.T) {
	Test(t, "supertape: failing Equal records a failed test", func(t *T) {
		r := st.New()
		r.Test("supertape: inner fail", func(tt st.T) {
			tt.Equal(1, 2)
			tt.End()
		})
		_, failed, _, total := r.Counts()
		t.Equal(failed == 1 && total == 1, true)
		t.End()
	})
}

func TestFailString(t *testing.T) {
	Test(t, "supertape: Fail with string message", func(t *T) {
		r := st.New()
		r.Test("supertape: inner failstr", func(tt st.T) {
			tt.Fail("boom")
			tt.End()
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed, 1)
		t.End()
	})
}

func TestFailError(t *testing.T) {
	Test(t, "supertape: Fail with error", func(t *T) {
		r := st.New()
		r.Test("supertape: inner failerr", func(tt st.T) {
			tt.Fail("boom error")
			tt.End()
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed, 1)
		t.End()
	})
}

func TestFailDefault(t *testing.T) {
	Test(t, "supertape: Fail with arbitrary value", func(t *T) {
		r := st.New()
		r.Test("supertape: inner faildefault", func(tt st.T) {
			tt.Fail(42)
			tt.End()
		})
		_, failed, _, _ := r.Counts()
		t.Equal(failed, 1)
		t.End()
	})
}

// -- events --

func TestEmitsTestEnd(t *testing.T) {
	Test(t, "supertape: emits test-end event", func(t *T) {
		r, evs := newCollector()
		r.Test("scope: inner test", func(tt st.T) {
			tt.Ok(true)
			tt.End()
		})
		t.Equal(len(*evs), 1)
		t.End()
	})
}

