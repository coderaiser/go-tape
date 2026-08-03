package main

import "strings"

type coverageOpts struct {
	enabled bool
	format  string // "lines" | "code-frame" | "json-lines"
	report  bool
	path    string // default "coverage.lcov"
}

var coverageFormats = map[string]bool{
	"lines": true, "code-frame": true, "json-lines": true,
}

func preprocess(args []string) (coverageOpts, []string) {
	opts := coverageOpts{format: "lines", path: "coverage.lcov"}
	var rest []string
	i := 0
	for i < len(args) {
		a := args[i]
		switch {
		case a == "-c" || a == "--coverage":
			opts.enabled = true
			if i+1 < len(args) && coverageFormats[args[i+1]] {
				opts.format = args[i+1]
				i += 2
			} else {
				i++
			}
		case strings.HasPrefix(a, "-c="):
			opts.enabled = true
			opts.format = strings.TrimPrefix(a, "-c=")
			i++
		case strings.HasPrefix(a, "--coverage="):
			opts.enabled = true
			opts.format = strings.TrimPrefix(a, "--coverage=")
			i++
		case a == "-r" || a == "--report":
			opts.report = true
			i++
		case strings.HasPrefix(a, "-r="):
			opts.report = true
			opts.path = strings.TrimPrefix(a, "-r=")
			i++
		case strings.HasPrefix(a, "--report="):
			opts.report = true
			opts.path = strings.TrimPrefix(a, "--report=")
			i++
		default:
			rest = append(rest, a)
			i++
		}
	}
	return opts, rest
}
