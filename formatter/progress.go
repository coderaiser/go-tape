package formatter

import "fmt"

// ProgressFormatter outputs CI-friendly format.
type ProgressFormatter struct {
	passed  int
	failed  int
	skipped int
	total   int
}

func (f *ProgressFormatter) Add(passed, failed, skipped []string) {
	f.passed = len(passed)
	f.failed = len(failed)
	f.skipped = len(skipped)
	f.total = f.passed + f.failed + f.skipped
}

func (f *ProgressFormatter) Format() string {
	if f.failed > 0 {
		return fmt.Sprintf("##[error]Tests failed: %d/%d\n", f.failed, f.total)
	}
	return fmt.Sprintf("##[section]Tests passed: %d/%d\n", f.passed, f.total)
}
