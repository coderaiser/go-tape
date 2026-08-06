package tape

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed package.json
var packageJSONBytes []byte

// TapeVersion returns the version string embedded from package.json at build time.
func TapeVersionFromJSON(packageJSONBytes []byte) string {
	var pkg struct {
		TapeVersion string `json:"version"`
	}
	if err := json.Unmarshal(packageJSONBytes, &pkg); err != nil {
		return "unknown"
	}
	if pkg.TapeVersion == "" {
		return "unknown"
	}
	return pkg.TapeVersion
}

func TapeVersion() string {
	return TapeVersionFromJSON(packageJSONBytes)
}

// TapeVersionLine returns "go-coverage x.y.z" for -v output.
func TapeVersionLine() string {
	return fmt.Sprintf("v%s", TapeVersion())
}
