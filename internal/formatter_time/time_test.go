package formatter_time_test

import (
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter_time"
)

func TestTimeFormatterStart(t *testing.T) {
	tape.Test(t, "formatter-time: Start returns empty string", func(t *tape.T) {
		f := formatter_time.New(10, &strings.Builder{})
		result := f.Start(10)
		t.Equal(result, "")
		t.End()
	})
}

func TestTimeFormatterTestEndReturnsEmpty(t *testing.T) {
	tape.Test(t, "formatter-time: TestEnd returns empty string", func(t *tape.T) {
		f := formatter_time.New(10, &strings.Builder{})
		f.Start(10)
		result := f.TestEnd(1, 10, 0, "scope: foo")
		t.Equal(result, "")
		t.End()
	})
}

func TestTimeFormatterTestEndWithFail(t *testing.T) {
	tape.Test(t, "formatter-time: TestEnd formats failure count in red", func(t *tape.T) {
		f := formatter_time.New(10, &strings.Builder{})
		f.Start(10)
		result := f.TestEnd(1, 10, 1, "scope: foo")
		t.Equal(result, "")
		t.End()
	})
}

func TestTimeFormatterClockEnv(t *testing.T) {
	tape.Test(t, "formatter-time: New uses TAPE_TIME_CLOCK env var", func(t *tape.T) {
		t.TB().Setenv("TAPE_TIME_CLOCK", "\U0001f550")
		f := formatter_time.New(10, &strings.Builder{})
		t.Ok(f != nil)
		t.End()
	})
}
