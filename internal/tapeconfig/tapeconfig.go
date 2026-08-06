// Package tapeconfig loads `.tape.toml` files that configure test discovery
// (formatter, directory excludes) and coverage excludes for the tape CLI.
package tapeconfig

import (
	"fmt"
	"os"

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

// Load reads the TOML config at path. A missing file returns Default();
// a malformed file prints a warning to stderr and returns Default().
func Load(path string) Config {
	cfg := Default()
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: could not load %s: %v\n", path, err)
		}
	}
	return cfg
}
