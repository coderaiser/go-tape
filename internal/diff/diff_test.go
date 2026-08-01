package diff_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/diff"
)

func TestDiffIsNonEmpty(t *testing.T) {
	tape.Test(t, "diff: Diff returns non-empty string for different values", func(t *tape.T) {
		result := diff.Diff("want", "got")
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffContainsExpected(t *testing.T) {
	tape.Test(t, "diff: Diff output contains expected value", func(t *tape.T) {
		result := diff.Diff("want", "got")
		t.Match(result, "want")
		t.End()
	})
}

func TestDiffDeepEqualEmpty(t *testing.T) {
	tape.Test(t, "diff: Diff returns empty when values deeply equal", func(t *tape.T) {
		t.Equal(diff.Diff("same", "same"), "")
		t.End()
	})
}

func TestDiffNil(t *testing.T) {
	tape.Test(t, "diff: Diff handles nil values", func(t *tape.T) {
		result := diff.Diff(nil, 5)
		t.Ok(result != "")
		t.End()
	})
}

type sampleStruct struct {
	Name string
	Age  int
}

func TestDiffStruct(t *testing.T) {
	tape.Test(t, "diff: Diff handles struct values", func(t *tape.T) {
		a := sampleStruct{Name: "aa", Age: 1}
		b := sampleStruct{Name: "bb", Age: 2}
		result := diff.Diff(a, b)
		t.Match(result, "Name")
		t.End()
	})
}

func TestDiffEmptyStruct(t *testing.T) {
	tape.Test(t, "diff: Diff hits empty struct branch", func(t *tape.T) {
		type empty struct{}
		result := diff.Diff(empty{}, 5)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffMap(t *testing.T) {
	tape.Test(t, "diff: Diff handles map values", func(t *tape.T) {
		a := map[string]int{"a": 1, "b": 2}
		b := map[string]int{"a": 1, "b": 3}
		result := diff.Diff(a, b)
		t.Match(result, "a")
		t.End()
	})
}

func TestDiffEmptyMap(t *testing.T) {
	tape.Test(t, "diff: Diff handles empty map", func(t *tape.T) {
		a := map[string]int{}
		b := map[string]int{"x": 1}
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffNilMap(t *testing.T) {
	tape.Test(t, "diff: Diff handles nil map", func(t *tape.T) {
		var a map[string]int
		b := map[string]int{"x": 1}
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffSlice(t *testing.T) {
	tape.Test(t, "diff: Diff handles slice values", func(t *tape.T) {
		a := []int{1, 2}
		b := []int{1, 3}
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffNilSlice(t *testing.T) {
	tape.Test(t, "diff: Diff handles nil slice", func(t *tape.T) {
		var a []int
		b := []int{1}
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffEmptySlice(t *testing.T) {
	tape.Test(t, "diff: Diff handles empty slice", func(t *tape.T) {
		a := []int{}
		b := []int{1}
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffArray(t *testing.T) {
	tape.Test(t, "diff: Diff handles array values", func(t *tape.T) {
		a := [2]int{1, 2}
		b := [2]int{1, 3}
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffNumbers(t *testing.T) {
	tape.Test(t, "diff: Diff handles numeric values", func(t *tape.T) {
		t.Ok(diff.Diff(1, 2) != "")
		t.End()
	})
}

func TestDiffFloats(t *testing.T) {
	tape.Test(t, "diff: Diff handles float values", func(t *tape.T) {
		t.Ok(diff.Diff(1.5, 2.5) != "")
		t.End()
	})
}

func TestDiffBools(t *testing.T) {
	tape.Test(t, "diff: Diff handles bool values", func(t *tape.T) {
		t.Ok(diff.Diff(true, false) != "")
		t.End()
	})
}

func TestDiffUint(t *testing.T) {
	tape.Test(t, "diff: Diff handles uint values", func(t *tape.T) {
		t.Ok(diff.Diff(uint(1), uint(2)) != "")
		t.End()
	})
}

func TestDiffComplex(t *testing.T) {
	tape.Test(t, "diff: Diff handles complex values", func(t *tape.T) {
		t.Ok(diff.Diff(complex(1, 1), complex(2, 2)) != "")
		t.End()
	})
}

func TestDiffPointer(t *testing.T) {
	tape.Test(t, "diff: Diff handles pointer values", func(t *tape.T) {
		a := 1
		b := 2
		result := diff.Diff(&a, &b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffNilPointer(t *testing.T) {
	tape.Test(t, "diff: Diff handles nil pointer", func(t *tape.T) {
		var a *int
		b := 1
		result := diff.Diff(a, &b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffInterfaceNil(t *testing.T) {
	tape.Test(t, "diff: Diff handles nil interface via pointer", func(t *tape.T) {
		var a any
		b := 1
		result := diff.Diff(&a, &b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffInterfaceValue(t *testing.T) {
	tape.Test(t, "diff: Diff handles interface value via pointer", func(t *tape.T) {
		var a any = 1
		var b any = 2
		result := diff.Diff(&a, &b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffChannel(t *testing.T) {
	tape.Test(t, "diff: Diff handles channel values", func(t *tape.T) {
		a := make(chan int)
		b := make(chan int)
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}

func TestDiffFunc(t *testing.T) {
	tape.Test(t, "diff: Diff handles function values", func(t *tape.T) {
		a := func() {}
		b := func() {}
		result := diff.Diff(a, b)
		t.Ok(result != "")
		t.End()
	})
}
