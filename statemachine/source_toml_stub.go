//go:build no_external

package statemachine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// loadTOML is a hand-rolled parser for the subset of TOML used by go-tape.
// Only handles [transitions.X] sections with key = "value" pairs.
// Used when building with -tags no_external.
func loadTOML(path string) ([]TransitionDef, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("FileSource: open %s: %w", path, err)
	}
	defer f.Close()

	var defs []TransitionDef
	var currentFrom string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// section header: [transitions.idle]
		if strings.HasPrefix(line, "[transitions.") && strings.HasSuffix(line, "]") {
			currentFrom = line[len("[transitions.") : len(line)-1]
			continue
		}

		// skip other sections: [states] etc.
		if strings.HasPrefix(line, "[") {
			currentFrom = ""
			continue
		}

		if currentFrom == "" {
			continue
		}

		// key = "value"
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		event := strings.TrimSpace(parts[0])
		toVal := strings.TrimSpace(strings.Trim(strings.TrimSpace(parts[1]), "\""))

		defs = append(defs, TransitionDef{From: currentFrom, Event: event, To: toVal})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("FileSource: scan %s: %w", path, err)
	}

	return defs, nil
}
