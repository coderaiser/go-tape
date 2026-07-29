# Tape
📼[**Supertape**](https://github.com/coderaiser/supertape) for Go.

## API

```go

Test "github.com/coderaiser/go-tape"

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
		t.Ok(errors.New("some error"))
		t.End()
	})
}

func TestTNoError(t *testing.T) {
	Test(t, "tape: NoError works", func(t *T) {
		t.NotOk(nil)
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
```

# License
MIT
