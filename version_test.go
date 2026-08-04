package tape

import (
	"testing"
)

func TestVersion(t *testing.T) {
	Test(t, "version: VersionFromJSON returns version string", func(t *T) {
		result := VersionFromJSON([]byte(`{"version":"1.2.3"}`))
		t.Equal(result, "1.2.3")
		t.End()
	})

	Test(t, "version: VersionFromJSON returns unknown on invalid JSON", func(t *T) {
		result := VersionFromJSON([]byte(`{invalid`))
		t.Equal(result, "unknown")
		t.End()
	})

	Test(t, "version: VersionFromJSON returns unknown on empty version", func(t *T) {
		result := VersionFromJSON([]byte(`{"version":""}`))
		t.Equal(result, "unknown")
		t.End()
	})

	Test(t, "version: VersionLine contains v", func(t *T) {
		result := VersionLine()
		t.Match(result, "v")
		t.End()
	})
}
