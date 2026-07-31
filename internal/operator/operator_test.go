package operator_test

import (
	"regexp"
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/operator"
)

func TestEqualPass(t *testing.T) {
	tape.Test(t, "operator: Equal returns ok for equal values", func(t *tape.T) {
		result := operator.Equal(1, 1)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestEqualFail(t *testing.T) {
	tape.Test(t, "operator: Equal returns not ok for different values", func(t *tape.T) {
		result := operator.Equal(1, 2)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotEqualPass(t *testing.T) {
	tape.Test(t, "operator: NotEqual returns ok for different primitives", func(t *tape.T) {
		result := operator.NotEqual(1, 2)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotEqualFail(t *testing.T) {
	tape.Test(t, "operator: NotEqual returns not ok for same primitives", func(t *tape.T) {
		result := operator.NotEqual(1, 1)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotEqualNonPrimitive(t *testing.T) {
	tape.Test(t, "operator: NotEqual returns ok for non-primitive same values", func(t *tape.T) {
		result := operator.NotEqual(struct{}{}, struct{}{})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotEqualPointer(t *testing.T) {
	tape.Test(t, "operator: NotEqual pointer same returns not ok", func(t *tape.T) {
		n := 1
		result := operator.NotEqual(&n, &n)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestDeepEqualPass(t *testing.T) {
	tape.Test(t, "operator: DeepEqual returns ok for deeply equal values", func(t *tape.T) {
		result := operator.DeepEqual([]int{1, 2}, []int{1, 2})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestDeepEqualFail(t *testing.T) {
	tape.Test(t, "operator: DeepEqual returns not ok for different values", func(t *tape.T) {
		result := operator.DeepEqual([]int{1}, []int{2})
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotDeepEqualPass(t *testing.T) {
	tape.Test(t, "operator: NotDeepEqual returns ok for different values", func(t *tape.T) {
		result := operator.NotDeepEqual([]int{1}, []int{2})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotDeepEqualFail(t *testing.T) {
	tape.Test(t, "operator: NotDeepEqual returns not ok for same values", func(t *tape.T) {
		result := operator.NotDeepEqual([]int{1, 2}, []int{1, 2})
		t.NotOk(result.Ok)
		t.End()
	})
}
func TestOkPassInt(t *testing.T) {
	tape.Test(t, "operator: Ok returns ok for non-zero int", func(t *tape.T) {
		result := operator.Ok(1)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestOkFailIntZero(t *testing.T) {
	tape.Test(t, "operator: Ok returns not ok for zero int", func(t *tape.T) {
		result := operator.Ok(0)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkFailNil(t *testing.T) {
	tape.Test(t, "operator: Ok returns not ok for nil", func(t *tape.T) {
		result := operator.Ok(nil)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkFailFalse(t *testing.T) {
	tape.Test(t, "operator: Ok returns not ok for false", func(t *tape.T) {
		result := operator.Ok(false)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkEmptyString(t *testing.T) {
	tape.Test(t, "operator: Ok returns not ok for empty string", func(t *tape.T) {
		result := operator.Ok("")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestOkNonEmptyString(t *testing.T) {
	tape.Test(t, "operator: Ok returns ok for non-empty string", func(t *tape.T) {
		result := operator.Ok("x")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestOkStruct(t *testing.T) {
	tape.Test(t, "operator: Ok returns ok for struct", func(t *tape.T) {
		result := operator.Ok(struct{}{})
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotOkPassNil(t *testing.T) {
	tape.Test(t, "operator: NotOk returns ok for nil", func(t *tape.T) {
		result := operator.NotOk(nil)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotOkPassFalse(t *testing.T) {
	tape.Test(t, "operator: NotOk returns ok for false", func(t *tape.T) {
		result := operator.NotOk(false)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotOkFailTrue(t *testing.T) {
	tape.Test(t, "operator: NotOk returns not ok for true", func(t *tape.T) {
		result := operator.NotOk(true)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestMatchPass(t *testing.T) {
	tape.Test(t, "operator: Match returns ok for matching pattern", func(t *tape.T) {
		result := operator.Match("hello", "hel")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestMatchFail(t *testing.T) {
	tape.Test(t, "operator: Match returns not ok for non-matching pattern", func(t *tape.T) {
		result := operator.Match("hello", "xyz")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestMatchInvalidPattern(t *testing.T) {
	tape.Test(t, "operator: Match returns not ok for invalid regex", func(t *tape.T) {
		result := operator.Match("hello", "[invalid")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestMatchRegexpType(t *testing.T) {
	tape.Test(t, "operator: Match works with *regexp.Regexp", func(t *tape.T) {
		re := regexp.MustCompile("x")
		result := operator.Match("x", re)
		t.Ok(result.Ok)
		t.End()
	})
}

func TestMatchInvalidType(t *testing.T) {
	tape.Test(t, "operator: Match returns not ok for invalid type", func(t *tape.T) {
		result := operator.Match("x", 42)
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotMatchPass(t *testing.T) {
	tape.Test(t, "operator: NotMatch returns ok for no match", func(t *tape.T) {
		result := operator.NotMatch("hello", "xyz")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestNotMatchFail(t *testing.T) {
	tape.Test(t, "operator: NotMatch returns not ok when matches", func(t *tape.T) {
		result := operator.NotMatch("hello", "hel")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestNotMatchInvalidPattern(t *testing.T) {
	tape.Test(t, "operator: NotMatch returns not ok for invalid regex", func(t *tape.T) {
		result := operator.NotMatch("hello", "[invalid")
		t.NotOk(result.Ok)
		t.End()
	})
}

func TestPass(t *testing.T) {
	tape.Test(t, "operator: Pass returns ok with message", func(t *tape.T) {
		result := operator.Pass("msg")
		t.Ok(result.Ok)
		t.End()
	})
}

func TestFail(t *testing.T) {
	tape.Test(t, "operator: Fail returns not ok with output", func(t *tape.T) {
		result := operator.Fail("msg")
		t.NotOk(result.Ok)
		t.End()
	})
}
