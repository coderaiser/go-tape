package formatter_short

import (
	"github.com/coderaiser/go-tape/internal/formatter_tap"
)

type ShortFormatter struct {
	*formatter_tap.TAPFormatter
}

func New() *ShortFormatter {
	return &ShortFormatter{TAPFormatter: formatter_tap.New()}
}

func (f *ShortFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	return f.TAPFormatter.Fail(count, message, operator, result, expected, output, at, "")
}
