# Tape

📼[**Supertape**](https://github.com/coderaiser/supertape) for Go.

## API

```go
import (
    Test "github.com/coderaiser/go-tape"
)

func TestTEqual(t *testing.T) {
	Test(t, "tape: Equal works", func(t *T) {
		result := 42
		expected := 42
		
		t.Equal(result, expected)
		t.End()
	})
}

func TestTNotEqual(t *testing.T) {
	Test(t, "tape: NotEqual works", func(t *T) {
		result := 1
		expected := 2
		
		t.NotEqual(result, exptected)
		t.End()
	})
}

func TestTOk(t *testing.T) {
	Test(t, "tape: Ok works", func(t *T) {
		result := true
		
		t.Ok(result)
		t.End()
	})
}

func TestTNotOk(t *testing.T) {
	Test(t, "tape: NotOk works", func(t *T) {
		result = false
		
		t.NotOk(result)
		t.End()
	})
}

func TestTDeepEqual(t *testing.T) {
	Test(t, "tape: DeepEqual works", func(t *T) {
		result := []int{1, 2}
		expected := []int{1, 2}
		
		t.DeepEqual(result, expected)
		t.End()
	})
}

func TestTNotDeepEqual(t *testing.T) {
	Test(t, "tape: NotDeepEqual works", func(t *T) {
		result := []int{1}
		expected := []int{2}
		
		t.NotDeepEqual(result, expected)
		t.End()
	})
}

func TestTError(t *testing.T) {
	Test(t, "tape: Error works", func(t *T) {
		result := errors.New("some error")
		
		t.Ok(result)
		t.End()
	})
}

func TestTNoError(t *testing.T) {
	Test(t, "tape: NoError works", func(t *T) {
		result := nil
		
		t.NotOk(result)
		t.End()
	})
}

func TestTMatch(t *testing.T) {
	Test(t, "tape: Match works with string pattern", func(t *T) {
		result := "hello 123"
		expected := `hello \d+`
		
		t.Match(result, expected)
		t.End()
	})
}

func TestTMatchRegexp(t *testing.T) {
	Test(t, "tape: Match works with regexp pattern", func(t *T) {
		result := "hello 123"
		expected := regexp.MustCompile(`hello \d+`)
		
		t.Match(result, expected)
		t.End()
	})
}

func TestTNotMatch(t *testing.T) {
	Test(t, "tape: NotMatch works", func(t *T) {
		result :=	"hello"
		expected := `\d+`
		t.NotMatch(result, expected)
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
		const result = true;
		
		t.Comment("just a note")
		t.Ok(result)
		t.End()
	})
}
```

# License

MIT
