package formatter_fail

import (
	"fmt"

	"github.com/coderaiser/go-tape/internal/formatter_tap"
)

type FailFormatter struct {
	*formatter_tap.TAPFormatter
	currentTest string
}

func New() *FailFormatter {
	return &FailFormatter{TAPFormatter: formatter_tap.New()}
}

func (f *FailFormatter) Test(name string) string {
	f.currentTest = name
	return ""
}

// Success suppresses passing test output — formatter-fail only shows failures.
func (f *FailFormatter) Success(count int, message string) string {
	return ""
}

// Fail always prepends the current test name and delegates to TAPFormatter
// (keeping the full error stack, unlike formatter-short which strips it).
func (f *FailFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	return fmt.Sprintf("# %s\n", f.currentTest) +
		f.TAPFormatter.Fail(count, message, operator, result, expected, output, at, errorStack)
}
