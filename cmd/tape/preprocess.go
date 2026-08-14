package main

import (
	"os"
	"path/filepath"
	"strings"

	tapeast "github.com/coderaiser/go-tape/internal/ast"
)

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

// dirFromPattern derives the filesystem directory to scan from a go test
// package pattern. "./..." and "." both map to ".".
// "./internal/foo/..." maps to "./internal/foo".
func dirFromPattern(path string) string {
	dir := strings.TrimSuffix(path, "/...")
	if dir == "" || dir == "." || dir == "./..." {
		return "."
	}
	return dir
}

// uniquePkgDirs returns one "./rel/path/to/pkg" entry per unique package
// directory found in the Only calls, relative to the working directory.
// When all Only calls are in one package, returns a single entry.
func uniquePkgDirs(calls []tapeast.OnlyCall) []string {
	seen := make(map[string]bool)
	var dirs []string
	wd, _ := os.Getwd()
	for _, c := range calls {
		dir := filepath.Dir(c.File)
		rel, err := filepath.Rel(wd, dir)
		if err != nil {
			rel = dir
		}
		pkg := "./" + rel
		if !seen[pkg] {
			seen[pkg] = true
			dirs = append(dirs, pkg)
		}
	}
	return dirs
}
