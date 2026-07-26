package statemachine

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// TransitionSource loads transition definitions once at startup.
// Load() is called once in New() — never per Apply().
type TransitionSource interface {
	Load() ([]TransitionDef, error)
}

// MemorySource holds transitions in memory.
// Primary use: tests — no file I/O needed.
type MemorySource struct {
	Defs []TransitionDef
}

func (s *MemorySource) Load() ([]TransitionDef, error) {
	return s.Defs, nil
}

// FileSource loads transitions from a TOML file.
//
// Expected TOML format:
//
//   [transitions.idle]
//   run = "running"
//
//   [transitions.running]
//   pass = "passed"
//   fail = "failed"
type FileSource struct {
	Path string
}

func (s FileSource) Load() ([]TransitionDef, error) {
	var cfg struct {
		Transitions map[string]map[string]string `toml:"transitions"`
	}
	if _, err := toml.DecodeFile(s.Path, &cfg); err != nil {
		return nil, fmt.Errorf("FileSource: load %s: %w", s.Path, err)
	}
	var defs []TransitionDef
	for from, events := range cfg.Transitions {
		for event, to := range events {
			defs = append(defs, TransitionDef{
				From:  from,
				Event: event,
				To:    to,
			})
		}
	}
	return defs, nil
}
