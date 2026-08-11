package operator_test

import (
	"regexp"
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/operator"
)

func TestEqualPass(t *testing.T) {
	Test(t, "operator: Equal returns ok for equal values", func(t *T) {
		result := operator.Equal(1, 1)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestEqualFail(t *testing.T) {
	Test(t, "operator: Equal returns not ok for different values", func(t *T) {
		result := operator.Equal(1, 2)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestEqualSlicesSameContent(t *testing.T) {
	Test(t, "operator: Equal fails for slices with same content", func(t *T) {
		result := operator.Equal([]int{1}, []int{1})
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestEqualSlicesSameContentOutput(t *testing.T) {
	Test(t, "operator: Equal emits 'values not equal, but deepEqual' for matching slices", func(t *T) {
		result := operator.Equal([]int{1}, []int{1})
		t.Equal(result.Output, "values not equal, but deepEqual")
		t.End()
	})
}

func TestEqualEmptySlices(t *testing.T) {
	Test(t, "operator: Equal fails for two empty slices", func(t *T) {
		result := operator.Equal([]int{}, []int{})
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestEqualEmptySlicesOutput(t *testing.T) {
	Test(t, "operator: Equal emits 'values not equal, but deepEqual' for two empty slices", func(t *T) {
		result := operator.Equal([]int{}, []int{})
		t.Equal(result.Output, "values not equal, but deepEqual")
		t.End()
	})
}

func TestEqualNilNil(t *testing.T) {
	Test(t, "operator: Equal passes for nil == nil", func(t *T) {
		result := operator.Equal(nil, nil)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestEqualStringsDiff(t *testing.T) {
	Test(t, "operator: Equal fails with diff for different strings", func(t *T) {
		result := operator.Equal("a", "b")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotEqualPass(t *testing.T) {
	Test(t, "operator: NotEqual returns ok for different primitives", func(t *T) {
		result := operator.NotEqual(1, 2)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotEqualFail(t *testing.T) {
	Test(t, "operator: NotEqual returns not ok for same primitives", func(t *T) {
		result := operator.NotEqual(1, 1)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotEqualNonPrimitive(t *testing.T) {
	Test(t, "operator: NotEqual returns ok for non-primitive same values", func(t *T) {
		result := operator.NotEqual(struct{}{}, struct{}{})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotEqualPointer(t *testing.T) {
	Test(t, "operator: NotEqual pointer same returns not ok", func(t *T) {
		n := 1
		result := operator.NotEqual(&n, &n)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestDeepEqualPass(t *testing.T) {
	Test(t, "operator: DeepEqual returns ok for deeply equal values", func(t *T) {
		result := operator.DeepEqual([]int{1, 2}, []int{1, 2})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestDeepEqualFail(t *testing.T) {
	Test(t, "operator: DeepEqual returns not ok for different values", func(t *T) {
		result := operator.DeepEqual([]int{1}, []int{2})
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotDeepEqualPass(t *testing.T) {
	Test(t, "operator: NotDeepEqual returns ok for different values", func(t *T) {
		result := operator.NotDeepEqual([]int{1}, []int{2})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotDeepEqualFail(t *testing.T) {
	Test(t, "operator: NotDeepEqual returns not ok for same values", func(t *T) {
		result := operator.NotDeepEqual([]int{1, 2}, []int{1, 2})
		t.NotOk(result.Ok)
		t.End()
	})
}
func TestOkPassInt(t *testing.T) {
	Test(t, "operator: Ok returns ok for non-zero int", func(t *T) {
		result := operator.Ok(1)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestOkFailIntZero(t *testing.T) {
	Test(t, "operator: Ok returns not ok for zero int", func(t *T) {
		result := operator.Ok(0)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkFailNil(t *testing.T) {
	Test(t, "operator: Ok returns not ok for nil", func(t *T) {
		result := operator.Ok(nil)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkFailFalse(t *testing.T) {
	Test(t, "operator: Ok returns not ok for false", func(t *T) {
		result := operator.Ok(false)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkEmptyString(t *testing.T) {
	Test(t, "operator: Ok returns not ok for empty string", func(t *T) {
		result := operator.Ok("")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkNonEmptyString(t *testing.T) {
	Test(t, "operator: Ok returns ok for non-empty string", func(t *T) {
		result := operator.Ok("x")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestOkStruct(t *testing.T) {
	Test(t, "operator: Ok returns ok for struct", func(t *T) {
		result := operator.Ok(struct{}{})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotOkPassNil(t *testing.T) {
	Test(t, "operator: NotOk returns ok for nil", func(t *T) {
		result := operator.NotOk(nil)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotOkPassFalse(t *testing.T) {
	Test(t, "operator: NotOk returns ok for false", func(t *T) {
		result := operator.NotOk(false)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotOkFailTrue(t *testing.T) {
	Test(t, "operator: NotOk returns not ok for true", func(t *T) {
		result := operator.NotOk(true)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestMatchPass(t *testing.T) {
	Test(t, "operator: Match returns ok for matching pattern", func(t *T) {
		result := operator.Match("hello", "hel")
		t.Ok(result.Ok)
		t.End()
	})

	Test(t, "operator: Match: []", func(t *T) {
		result := operator.Match("usage: coverage [options]", "usage: coverage [options]")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestMatchFail(t *testing.T) {
	Test(t, "operator: Match returns not ok for non-matching pattern", func(t *T) {
		result := operator.Match("hello", "xyz")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestMatchInvalidPattern(t *testing.T) {
	Test(t, "operator: Match returns not ok for invalid regex", func(t *T) {
		result := operator.Match("hello", "[invalid")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestMatchRegexpType(t *testing.T) {
	Test(t, "operator: Match works with *regexp.Regexp", func(t *T) {
		re := regexp.MustCompile("x")
		result := operator.Match("x", re)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestMatchInvalidType(t *testing.T) {
	Test(t, "operator: Match returns not ok for invalid type", func(t *T) {
		result := operator.Match("x", 42)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotMatchPass(t *testing.T) {
	Test(t, "operator: NotMatch returns ok for no match", func(t *T) {
		result := operator.NotMatch("hello", "xyz")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotMatchFail(t *testing.T) {
	Test(t, "operator: NotMatch returns not ok when matches", func(t *T) {
		result := operator.NotMatch("hello", "hel")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotMatchInvalidPattern(t *testing.T) {
	Test(t, "operator: NotMatch returns ok for invalid regex", func(t *T) {
		result := operator.NotMatch("hello", "[invalid")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotMatchInvalidType(t *testing.T) {
	Test(t, "operator: NotMatch returns not ok for invalid type", func(t *T) {
		result := operator.NotMatch("x", 42)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestPass(t *testing.T) {
	Test(t, "operator: Pass returns ok with message", func(t *T) {
		result := operator.Pass("msg")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestFail(t *testing.T) {
	Test(t, "operator: Fail returns not ok with output", func(t *T) {
		result := operator.Fail("msg")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkTypedNilSlice(t *testing.T) {
	Test(t, "operator: Ok returns not ok for typed-nil slice", func(t *T) {
		var s []int
		result := operator.Ok(s)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotOkTypedNilSlice(t *testing.T) {
	Test(t, "operator: NotOk returns ok for typed-nil slice", func(t *T) {
		var s []int
		result := operator.NotOk(s)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestOkTypedNilMap(t *testing.T) {
	Test(t, "operator: Ok returns not ok for typed-nil map", func(t *T) {
		var m map[string]int
		result := operator.Ok(m)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkTypedNilPointer(t *testing.T) {
	Test(t, "operator: Ok returns not ok for typed-nil pointer", func(t *T) {
		var p *int
		result := operator.Ok(p)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkTypedNilFunc(t *testing.T) {
	Test(t, "operator: Ok returns not ok for typed-nil func", func(t *T) {
		var f func()
		result := operator.Ok(f)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkNonNilSlice(t *testing.T) {
	Test(t, "operator: Ok returns ok for non-nil empty slice", func(t *T) {
		s := []int{}
		result := operator.Ok(s)
		t.Ok(result.Ok)
		t.End()
	})
}
