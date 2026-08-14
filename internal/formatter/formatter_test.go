package formatter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestNewTap(t *testing.T) {
	Test(t, "formatter: tap format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("tap", &buf, 3) != nil)
		t.End()
	})
}

func TestNewShort(t *testing.T) {
	Test(t, "formatter: short format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("short", &buf, 3) != nil)
		t.End()
	})
}

func TestNewFail(t *testing.T) {
	Test(t, "formatter: fail format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("fail", &buf, 3) != nil)
		t.End()
	})
}

func TestNewTime(t *testing.T) {
	Test(t, "formatter: time format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("time", &buf, 3) != nil)
		t.End()
	})
}

func TestNewJSONLines(t *testing.T) {
	Test(t, "formatter: json-lines format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("json-lines", &buf, 3) != nil)
		t.End()
	})
}

func TestNewProgressBar(t *testing.T) {
	Test(t, "formatter: progress-bar format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("progress-bar", &buf, 3) != nil)
		t.End()
	})
}

func TestNewDebug(t *testing.T) {
	Test(t, "formatter: debug format constructs", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("debug", &buf, 3) != nil)
		t.End()
	})
}

func TestNewUnknownDefaultsToProgressBar(t *testing.T) {
	Test(t, "formatter: unknown format falls back to progress-bar", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		t.Ok(New("unknown-format", &buf, 3) != nil)
		t.End()
	})
}

func TestNewEmptyFormatDefaultsToProgressBar(t *testing.T) {
	Test(t, "formatter: empty format defaults to progress-bar", func(t *T) {
		t.TB().Setenv("CI", "0")
		var buf strings.Builder
		t.Ok(New("", &buf, 3) != nil)
		t.End()
	})
}

// TestNewCIForcesFail verifies that CI=1 forces the fail formatter even
// when "tap" is requested.
func TestNewCIForcesFail(t *testing.T) {
	Test(t, "formatter: CI=1 forces fail format", func(t *T) {
		t.TB().Setenv("CI", "1")
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		var buf strings.Builder
		d := New("tap", &buf, 1)
		d.Emit(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1, Total: 1})
		d.Emit(stream.Event{Type: stream.TypeFail, Test: "scope: x", Message: "m", Operator: "op"})
		t.Match(buf.String(), "# scope: x")
		t.End()
	})
}

func TestNewCITrueForcesFail(t *testing.T) {
	Test(t, "formatter: CI=true forces fail format", func(t *T) {
		t.TB().Setenv("CI", "true")
		t.TB().Setenv("TAPE_PROGRESS_BAR", "0")
		var buf strings.Builder
		d := New("tap", &buf, 1)
		d.Emit(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1, Total: 1})
		d.Emit(stream.Event{Type: stream.TypeFail, Test: "scope: x", Message: "m", Operator: "op"})
		t.Match(buf.String(), "# scope: x")
		t.End()
	})
}

func TestEmitTestEnd(t *testing.T) {
	Test(t, "formatter: Emit test-end writes output", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		d := New("tap", &buf, 1)
		d.Emit(stream.Event{Type: stream.TypeTestEnd, Test: "scope: x", Count: 1, Total: 1})
		t.Match(buf.String(), "ok 1 scope: x")
		t.End()
	})
}

func TestEmitFail(t *testing.T) {
	Test(t, "formatter: Emit fail writes not-ok line", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		d := New("tap", &buf, 1)
		d.Emit(stream.Event{
			Type: stream.TypeFail, Test: "scope: x", Count: 1,
			Message: "should equal", Operator: "should equal",
			Result: "got", Expected: "want",
		})
		t.Match(buf.String(), "not ok 1 scope: x")
		t.End()
	})
}

func TestEmitBuildError(t *testing.T) {
	Test(t, "formatter: Emit build-error writes output", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		d := New("tap", &buf, 0)
		d.Emit(stream.Event{
			Type:    stream.TypeBuildError,
			Package: "example.com/foo",
			Output:  "foo.go:1:1: undefined: x\n",
		})
		t.Match(buf.String(), "build-error")
		t.End()
	})
}

func TestEmitUnknownFail(t *testing.T) {
	Test(t, "formatter: Emit unknown-fail writes output", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		d := New("tap", &buf, 1)
		d.Emit(stream.Event{
			Type:   stream.TypeUnknownFail,
			Test:   "scope: x",
			Count:  1,
			Output: "panic: something\n",
		})
		t.Match(buf.String(), "not ok")
		t.End()
	})
}

func TestEmitResolvesAtToURI(t *testing.T) {
	Test(t, "formatter: Emit resolves At field to file URI", func(t *T) {
		t.TB().Setenv("CI", "false")
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, "go.mod"),
			[]byte("module example.com/mymod\n\ngo 1.25\n"), 0644)
		var buf strings.Builder
		d := newWithDir("tap", &buf, 1, dir)
		d.Emit(stream.Event{
			Type:     stream.TypeFail,
			Test:     "scope: x",
			Count:    1,
			At:       "foo_test.go:10",
			Package:  "example.com/mymod/internal/pkg",
			Message:  "m",
			Operator: "op",
		})
		t.Match(buf.String(), "internal/pkg/foo_test.go:10")
		t.End()
	})
}

func TestEndWritesSummary(t *testing.T) {
	Test(t, "formatter: End writes summary line", func(t *T) {
		t.TB().Setenv("CI", "false")
		var buf strings.Builder
		d := New("tap", &buf, 1)
		d.End(1, 0, 0)
		t.Match(buf.String(), "# tests 1")
		t.End()
	})
}

// --- helper function coverage ---

func TestTestLabelWithSlash(t *testing.T) {
	Test(t, "formatter: testLabel strips prefix before slash", func(t *T) {
		t.Equal(testLabel("TestFoo/scope:_bar_baz"), "scope: bar baz")
		t.End()
	})
}

func TestTestLabelWithoutSlash(t *testing.T) {
	Test(t, "formatter: testLabel returns name unchanged when no slash", func(t *T) {
		t.Equal(testLabel("scope:_bar"), "scope: bar")
		t.End()
	})
}

func TestFileLinkEmpty(t *testing.T) {
	Test(t, "formatter: fileLink returns empty string for empty at", func(t *T) {
		t.Equal(fileLink("", "/dir"), "")
		t.End()
	})
}

func TestFileLinkWithDir(t *testing.T) {
	Test(t, "formatter: fileLink prepends dir to relative at", func(t *T) {
		t.Equal(fileLink("file.go:10:", "/proj"), "at file:///proj/file.go:10")
		t.End()
	})
}

func TestFileLinkAbsolute(t *testing.T) {
	Test(t, "formatter: fileLink does not prepend dir to absolute at", func(t *T) {
		t.Equal(fileLink("/abs/file.go:10:", "/proj"), "at file:///abs/file.go:10")
		t.End()
	})
}

func TestFileLinkNoDir(t *testing.T) {
	Test(t, "formatter: fileLink with no dir uses at as-is", func(t *T) {
		t.Equal(fileLink("file.go:5:", ""), "at file://file.go:5")
		t.End()
	})
}

func TestPkgDirRoot(t *testing.T) {
	Test(t, "formatter: pkgDir returns module root for root package", func(t *T) {
		t.Equal(pkgDir("example.com/mymod", "example.com/mymod", "/abs"), "/abs")
		t.End()
	})
}

func TestPkgDirSub(t *testing.T) {
	Test(t, "formatter: pkgDir appends sub-package path", func(t *T) {
		t.Equal(pkgDir("example.com/mymod/internal/foo", "example.com/mymod", "/abs"), "/abs/internal/foo")
		t.End()
	})
}

func TestPkgDirEmptyPkg(t *testing.T) {
	Test(t, "formatter: pkgDir returns dir when pkg is empty", func(t *T) {
		t.Equal(pkgDir("", "example.com/mymod", "/abs"), "/abs")
		t.End()
	})
}

func TestReadModuleName(t *testing.T) {
	Test(t, "formatter: readModuleName reads module from go.mod", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/mymod\n\ngo 1.21\n"), 0644)
		t.Equal(readModuleName(dir), "example.com/mymod")
		t.End()
	})
}

func TestReadModuleNameMissing(t *testing.T) {
	Test(t, "formatter: readModuleName returns empty when no go.mod", func(t *T) {
		t.Equal(readModuleName(t.TB().TempDir()), "")
		t.End()
	})
}

func TestReadModuleNameWithoutModuleLine(t *testing.T) {
	Test(t, "formatter: readModuleName returns empty when go.mod has no module line", func(t *T) {
		dir := t.TB().TempDir()
		err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("go 1.21\n"), 0644)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Equal(readModuleName(dir), "")
		t.End()
	})
}

func TestWriteEmpty(t *testing.T) {
	Test(t, "formatter: write ignores empty string", func(t *T) {
		var buf strings.Builder
		write(&buf, "")
		t.Equal(buf.String(), "")
		t.End()
	})
}

func TestWriteNonEmpty(t *testing.T) {
	Test(t, "formatter: write emits non-empty string", func(t *T) {
		var buf strings.Builder
		write(&buf, "hello\n")
		t.Equal(buf.String(), "hello\n")
		t.End()
	})
}
