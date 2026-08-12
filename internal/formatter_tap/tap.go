package formatter_tap

import (
	"fmt"
	"strings"

	"github.com/coderaiser/go-tape/internal/diff"
	"github.com/coderaiser/go-tape/internal/stream"
)

// TAPFormatter outputs TAP version 13.
type TAPFormatter struct{}

func New() *TAPFormatter { return &TAPFormatter{} }

// Event dispatches on e.Type and returns a TAP line or block.
func (f *TAPFormatter) Event(e stream.Event) string {
	switch e.Type {
	case stream.TypeTestEnd:
		return fmt.Sprintf("ok %d %s\n", e.Count, e.Test)
	case stream.TypeFail:
		return f.failBlock(e.Count, e.Test, e.Operator, e.Result, e.Expected, e.Output, e.At, e.ErrorStack)
	case stream.TypeUnknownFail:
		return fmt.Sprintf("not ok %d %s\n%s\n", e.Count, e.Test, e.Output)
	case stream.TypeBuildError:
		return fmt.Sprintf("# build-error: %s\n%s\n", e.Package, e.Output)
	case stream.TypePackageError:
		return fmt.Sprintf("# package-error: %s\n%s\n", e.Package, e.Output)
	case stream.TypeComment:
		return fmt.Sprintf("# %s\n", e.Message)
	}
	return ""
}

func (f *TAPFormatter) failBlock(count int, name, operator string, result, expected any, output, at, errorStack string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "not ok %d %s\n", count, name)

	if output != "" {
		sb.WriteString(output)
	} else {
		sb.WriteString("  ---\n")
		if operator != "" {
			fmt.Fprintf(&sb, "    operator: %s\n", operator)
		}
		if d := diff.Diff(expected, result); d != "" {
			fmt.Fprintf(&sb, "      diff: |-\n")
			for _, line := range strings.Split(strings.TrimRight(d, "\n"), "\n") {
				fmt.Fprintf(&sb, "      %s\n", line)
			}
		} else {
			if expected != nil {
				fmt.Fprintf(&sb, "    expected: |-\n      %v\n", expected)
			}
			if result != nil {
				fmt.Fprintf(&sb, "    result: |-\n      %v\n", result)
			}
		}
		if at != "" {
			fmt.Fprintf(&sb, "    %s\n", at)
		}
		if errorStack != "" {
			fmt.Fprintf(&sb, "    stack: |-\n%s\n", errorStack)
		}
		sb.WriteString("  ...\n")
	}

	sb.WriteString("\n")
	return sb.String()
}

// End emits the TAP summary block.
func (f *TAPFormatter) End(passed, failed, skipped int) string {
	total := passed + failed + skipped
	var sb strings.Builder
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "1..%d\n", total)
	fmt.Fprintf(&sb, "# tests %d\n", total)
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
