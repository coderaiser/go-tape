package uncovered

// Covered is exercised by tests.
func Covered() string {
	return "covered"
}

// Uncovered is never called by any test, leaving its block uncovered.
func Uncovered() string {
	return "uncovered"
}
