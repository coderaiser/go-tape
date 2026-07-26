package formatter

import (
	"fmt"
	"strings"
)

// ShortFormatter outputs compact format.
type ShortFormatter struct {
	total   int
	passed  int
	failed  int
	skipped int
}

func (f *ShortFormatter) Add(passed, failed, skipped []string) {
	f.passed = len(passed)
	f.failed = len(failed)
	f.skipped = len(skipped)
	f.total = f.passed + f.failed + f.skipped
}

func (f *ShortFormatter) Format() string {
	var sb strings.Builder
	if f.failed > 0 {
		fmt.Fprintf(&sb, "FAIL: %d/%d tests failed", f.failed, f.total)
	} else {
		fmt.Fprintf(&sb, "PASS: %d/%d tests passed", f.passed, f.total)
	}
	if f.skipped > 0 {
		fmt.Fprintf(&sb, " (%d skipped)", f.skipped)
	}
	sb.WriteString("\n")
	return sb.String()
}
