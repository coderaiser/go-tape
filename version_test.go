package tape

import (
	"testing"
)

func TestTapeVersion(t *testing.T) {
	Test(t, "version: TapeVersionFromJSON returns version string", func(t *T) {
		result := TapeVersionFromJSON([]byte(`{"version":"1.2.3"}`))
		t.Equal(result, "1.2.3")
		t.End()
	})

	Test(t, "version: TapeVersionFromJSON returns unknown on invalid JSON", func(t *T) {
		result := TapeVersionFromJSON([]byte(`{invalid`))
		t.Equal(result, "unknown")
		t.End()
	})

	Test(t, "version: TapeVersionFromJSON returns unknown on empty version", func(t *T) {
		result := TapeVersionFromJSON([]byte(`{"version":""}`))
		t.Equal(result, "unknown")
		t.End()
	})

	Test(t, "version: TapeVersionLine contains v", func(t *T) {
		result := TapeVersionLine()
		t.Match(result, "v")
		t.End()
	})
}
