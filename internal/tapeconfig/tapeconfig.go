// Package tapeconfig loads `.tape.toml` files that configure test discovery
// (formatter, directory excludes) and coverage excludes for the tape CLI.
package tapeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds tape configuration loaded from .tape.toml.
type Config struct {
	Test struct {
		Formatter string   `toml:"formatter"`
		Exclude   []string `toml:"exclude"`
	} `toml:"test"`
	Coverage struct {
		Exclude []string `toml:"exclude"`
	} `toml:"coverage"`
}

// Default returns the built-in configuration shipped with tape.
func Default() Config {
	cfg := Config{}
	cfg.Test.Formatter = "progress-bar"
	cfg.Test.Exclude = []string{"fixture"}
	cfg.Coverage.Exclude = []string{"node_modules"}
	return cfg
}

// Load reads tape config from dir, trying .tape.toml then tape.toml.
// Missing file returns Default(); malformed file prints a warning and
// returns Default().
func Load(dir string) Config {
	for _, name := range []string{".tape.toml", "tape.toml"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		cfg := Default()
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not load %s: %v\n", path, err)
			return Default()
		}
		return cfg
	}
	return Default()
}
