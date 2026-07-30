//go:build no_external

package diff

import "fmt"

// Diff returns a formatted diff string between expected and result.
func Diff(expected, result any) string {
	return fmt.Sprintf("expected: %v\ngot:      %v", expected, result)
}
