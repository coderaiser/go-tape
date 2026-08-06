package output

import (
	"regexp"
	"strings"
)

type OutputFields struct {
	Operator   string
	Result     string
	Expected   string
	At         string
	Diff       string
	ErrorStack string
	Raw        string
	Cut        string
}

var (
	reAt       = regexp.MustCompile(`^\s+(\S+\.go:\d+:)\s*`)
	reResult   = regexp.MustCompile(`^\s+result:\s+(.+)`)
	reExpected = regexp.MustCompile(`^\s+expected:\s+(.+)`)
	reOperator = regexp.MustCompile(`^\s+operator:\s+(.+)`)
	reNoise    = regexp.MustCompile(`^(=== RUN|--- FAIL:|--- PASS:|FAIL|PASS|ok\s)`)
	reAtPrefix = regexp.MustCompile(`^\s+\S+\.go:\d+:\s*`)
	reANSI     = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	// diff header lines emitted by diff.Diff
	reDiffHeader = regexp.MustCompile(`^\s*[-+] (?:expected|received)\s*$`)
	// diff content lines: start with "- " or "+ " (after stripping leading space+ANSI)
	reDiffLine = regexp.MustCompile(`^\s*[-+] `)
)

func stripANSI(s string) string {
	return reANSI.ReplaceAllString(s, "")
}

func ParseOutput(lines []string) OutputFields {
	raw := strings.Join(lines, "")
	fields := OutputFields{Raw: raw}

	var cutLines []string
	var diffLines []string
	inDiff := false

	for _, line := range lines {
		clean := stripANSI(line)
		stripped := strings.TrimRight(strings.TrimSpace(clean), "\n")

		// diff header: "- expected" or "+ received" (possibly with ANSI)
		if reDiffHeader.MatchString(clean) {
			inDiff = true
			continue
		}

		// blank line between diff header and diff body — skip
		if inDiff && strings.TrimSpace(clean) == "" {
			continue
		}

		// diff body lines
		if inDiff {
			if reDiffLine.MatchString(clean) {
				diffLines = append(diffLines, strings.TrimSpace(stripped))
				continue
			}
			// not a diff line any more — fall through to normal parsing
			inDiff = false
		}

		if m := reAt.FindStringSubmatch(clean); m != nil {
			fields.At = m[1]
			rest := reAtPrefix.ReplaceAllString(clean, "")
			if rest != "" && rest != "\n" {
				cutLines = append(cutLines, rest)
			}
			continue
		}
		if m := reResult.FindStringSubmatch(clean); m != nil {
			fields.Result = strings.TrimSpace(m[1])
			cutLines = append(cutLines, clean)
			continue
		}
		if m := reExpected.FindStringSubmatch(clean); m != nil {
			fields.Expected = strings.TrimSpace(m[1])
			cutLines = append(cutLines, clean)
			continue
		}
		if m := reOperator.FindStringSubmatch(clean); m != nil {
			fields.Operator = strings.TrimSpace(m[1])
			cutLines = append(cutLines, clean)
			continue
		}
		// Skip Go test runner noise lines.
		if reNoise.MatchString(strings.TrimSpace(clean)) {
			continue
		}
		cutLines = append(cutLines, clean)
	}

	if len(diffLines) > 0 {
		fields.Diff = strings.Join(diffLines, "\n")
	}
	fields.Cut = strings.Join(cutLines, "")
	return fields
}
