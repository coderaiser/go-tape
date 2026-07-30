package formatter_fail

import (
	"fmt"

	"github.com/coderaiser/go-tape/internal/formatter_short"
)

type FailFormatter struct {
	*formatter_short.ShortFormatter
	currentTest string
}

func New() *FailFormatter {
	return &FailFormatter{ShortFormatter: formatter_short.New()}
}

func (f *FailFormatter) Test(name string) string {
	f.currentTest = name
	return ""
}

func (f *FailFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	prefix := ""
	if f.currentTest != "" {
		prefix = fmt.Sprintf("# %s\n", f.currentTest)
	}
	return prefix + f.ShortFormatter.Fail(count, message, operator, result, expected, output, at, errorStack)
}
