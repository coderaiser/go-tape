package formatter

import (
	"fmt"
	"strings"

	"github.com/coderaiser/go-tape/state"
)

// TAPFormatter outputs TAP version 14.
type TAPFormatter struct {
	passed  []string
	failed  []string
	skipped []string
}

func NewTAP() *TAPFormatter {
	return &TAPFormatter{}
}

func (f *TAPFormatter) Add(passed, failed, skipped []string) {
	f.passed = passed
	f.failed = failed
	f.skipped = skipped
}

func (f *TAPFormatter) Format() string {
	var sb strings.Builder
	total := len(f.passed) + len(f.failed) + len(f.skipped)
	fmt.Fprintf(&sb, "TAP version 14\n")
	fmt.Fprintf(&sb, "1..%d\n", total)

	n := 0
	for _, test := range f.passed {
		n++
		fmt.Fprintf(&sb, "ok %d - %s\n", n, test)
	}
	for _, test := range f.failed {
		n++
		fmt.Fprintf(&sb, "not ok %d - %s\n", n, test)
	}
	for _, test := range f.skipped {
		n++
		fmt.Fprintf(&sb, "ok %d - # SKIP %s\n", n, test)
	}

	fmt.Fprintf(&sb, "# pass %d\n", len(f.passed))
	fmt.Fprintf(&sb, "# fail %d\n", len(f.failed))
	fmt.Fprintf(&sb, "# skip %d\n", len(f.skipped))

	return sb.String()
}

// FormatFromStore builds TAP output from a store.
func FormatFromStore(s *state.Store) string {
	passed, failed, skipped := s.Summary()
	f := NewTAP()
	f.Add(passed, failed, skipped)
	return f.Format()
}
