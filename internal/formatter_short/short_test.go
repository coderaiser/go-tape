package formatter_short_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_short"
)

func TestShortFailOmitsErrorStack(t *testing.T) {
	tape.Test(t, "formatter-short: Fail omits error stack", func(t *tape.T) {
		f := formatter_short.New()
		result := f.Fail(1, "scope: x", "Equal", "got", "want", "", "", "stack here")
		t.NotMatch(result, "stack here")
		t.End()
	})
}

func TestShortNewReturnsFormatter(t *testing.T) {
	tape.Test(t, "formatter-short: New returns a ShortFormatter", func(t *tape.T) {
		f := formatter_short.New()
		result := f.Start(1)
		t.Ok(result != "")
		t.End()
	})
}
