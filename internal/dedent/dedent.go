//go:build !no_external

package dedent

import "github.com/lithammer/dedent"

// Dedent removes any common leading whitespace from every line in text.
func Dedent(text string) string {
	return dedent.Dedent(text)
}
