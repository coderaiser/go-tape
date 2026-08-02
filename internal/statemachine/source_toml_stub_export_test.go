//go:build no_external

package statemachine

import "io"

// ParseTOMLReader exposes the hand-rolled TOML parser for tests built with
// -tags no_external.
func ParseTOMLReader(r io.Reader, name string) ([]TransitionDef, error) {
	return parseTOMLReader(r, name)
}
