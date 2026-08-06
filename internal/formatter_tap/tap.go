package formatter_tap

import (
	"fmt"
	"strings"
)

// TAPFormatter outputs TAP version 13.
type TAPFormatter struct{}

func New() *TAPFormatter { return &TAPFormatter{} }

func (f *TAPFormatter) Start(total int) string {
	return "TAP version 13\n"
}

func (f *TAPFormatter) Test(name string) string {
	return fmt.Sprintf("# %s\n", name)
}

func (f *TAPFormatter) TestEnd(count, total, failed int, name string) string {
	return ""
}

func (f *TAPFormatter) Success(count int, message string) string {
	return fmt.Sprintf("ok %d %s\n", count, message)
}

// Fail emits a TAP13 not-ok line followed by diagnostic detail.
// When diff is non-empty it is written verbatim on its own; otherwise the
// operator/expected/result/at/stack fields are written as indented lines.
func (f *TAPFormatter) Fail(count int, message, operator string, result, expected any, diff, at, errorStack string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "not ok %d %s\n", count, message)

	if diff != "" {
		sb.WriteString(diff)
	} else {
		if operator != "" {
			fmt.Fprintf(&sb, "    operator: %s\n", operator)
		}
		if expected != nil {
			fmt.Fprintf(&sb, "    expected: %v\n", expected)
		}
		if result != nil {
			fmt.Fprintf(&sb, "    result: %v\n", result)
		}
		if at != "" {
			fmt.Fprintf(&sb, "    at %s\n", at)
		}
		if errorStack != "" {
			fmt.Fprintf(&sb, "    stack: %s\n", errorStack)
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

func (f *TAPFormatter) Comment(message string) string {
	return fmt.Sprintf("# %s\n", message)
}

func (f *TAPFormatter) End(count, passed, failed, skipped int) string {
	var sb strings.Builder
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "1..%d\n", count)
	fmt.Fprintf(&sb, "# tests %d\n", count)
	fmt.Fprintf(&sb, "# pass %d\n", passed)
	if skipped > 0 {
		fmt.Fprintf(&sb, "# skip %d\n", skipped)
	}
	if failed > 0 {
		fmt.Fprintf(&sb, "# fail %d\n", failed)
	}
	sb.WriteString("\n")
	if failed == 0 {
		sb.WriteString("# ok\n")
	}
	sb.WriteString("\n")
	return sb.String()
}
