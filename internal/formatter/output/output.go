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
)

func ParseOutput(lines []string) OutputFields {
	raw := strings.Join(lines, "")
	fields := OutputFields{Raw: raw}

	var cutLines []string
	for _, line := range lines {
		if m := reAt.FindStringSubmatch(line); m != nil {
			fields.At = m[1]
			// Strip the "    file.go:N: " prefix, keep the rest as the first cut line.
			rest := reAtPrefix.ReplaceAllString(line, "")
			if rest != "" && rest != "\n" {
				cutLines = append(cutLines, rest)
			}
			continue
		}
		if m := reResult.FindStringSubmatch(line); m != nil {
			fields.Result = strings.TrimSpace(m[1])
			cutLines = append(cutLines, line)
			continue
		}
		if m := reExpected.FindStringSubmatch(line); m != nil {
			fields.Expected = strings.TrimSpace(m[1])
			cutLines = append(cutLines, line)
			continue
		}
		if m := reOperator.FindStringSubmatch(line); m != nil {
			fields.Operator = strings.TrimSpace(m[1])
			cutLines = append(cutLines, line)
			continue
		}
		// Skip Go test runner noise lines.
		if reNoise.MatchString(strings.TrimSpace(line)) {
			continue
		}
		cutLines = append(cutLines, line)
	}

	fields.Cut = strings.Join(cutLines, "")
	return fields
}
