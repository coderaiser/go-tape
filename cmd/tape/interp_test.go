package main

import (
	"os"
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/formatter"
	"github.com/coderaiser/go-tape/internal/interp"
	"github.com/coderaiser/go-tape/internal/stream"
)

func TestIsSupertapeFile(t *testing.T) {
	Test(t, "interp: isSupertapeFile detects tapeapi imports", func(t *T) {
		t.Ok(isSupertapeFile(`import . "tapeapi"`) &&
			isSupertapeFile(`import "tapeapi"`) &&
			!isSupertapeFile(`package foo`))
		t.End()
	})
}

func TestInterpToStream(t *testing.T) {
	Test(t, "interp: interpToStream maps events", func(t *T) {
		var st interpState

		// test start produces nothing
		if n := interpToStream(interp.Event{Kind: "test", Name: "scope: x"}, &st); len(n) != 0 {
			t.TB().Fatal("test event should emit nothing")
		}

		// passing assert produces nothing
		if n := interpToStream(interp.Event{Kind: "assert", Name: "equal", Ok: true}, &st); len(n) != 0 {
			t.TB().Fatal("passing assert should emit nothing")
		}

		// failing assert produces one TypeFail
		evs := interpToStream(interp.Event{Kind: "assert", Name: "equal", Ok: false, Got: 1, Wanted: 2}, &st)
		if len(evs) != 1 || evs[0].Type != stream.TypeFail || evs[0].Test != "scope: x" {
			t.TB().Fatalf("expected fail event, got %+v", evs)
		}

		// t.End() (no Name) emits nothing
		if n := interpToStream(interp.Event{Kind: "end"}, &st); len(n) != 0 {
			t.TB().Fatal("t.End() event should emit nothing")
		}

		// report event is not a transition; emits nothing (covers default return nil)
		if n := interpToStream(interp.Event{Kind: "report", Name: "ok", Ok: true}, &st); len(n) != 0 {
			t.TB().Fatal("report event should emit nothing")
		}

		// final end with Name + Ok emits TypeTestEnd
		evs = interpToStream(interp.Event{Kind: "end", Name: "scope: x", Ok: true}, &st)
		if len(evs) != 1 || evs[0].Type != stream.TypeTestEnd {
			t.TB().Fatalf("expected test-end event, got %+v", evs)
		}

		// final end with Name + !Ok emits TypeFail
		evs = interpToStream(interp.Event{Kind: "end", Name: "scope: x", Ok: false}, &st)
		if len(evs) != 1 || evs[0].Type != stream.TypeFail {
			t.TB().Fatalf("expected fail end event, got %+v", evs)
		}

		t.End()
	})
}

func TestFindSupertapeFiles(t *testing.T) {
	Test(t, "interp: findSupertapeFiles finds tape files only", func(t *T) {
		files, err := findSupertapeFiles("testdata/interp", nil)
		ok := err == nil && len(files) == 2
		for _, f := range files {
			ok = ok && strings.HasSuffix(f.path, "_tape.go") && f.src != ""
		}
		t.Ok(ok)
		t.End()
	})
}

func TestFindSupertapeFilesMissingDir(t *testing.T) {
	Test(t, "interp: findSupertapeFiles errors on missing dir", func(t *T) {
		_, err := findSupertapeFiles("testdata/does-not-exist", nil)
		t.Ok(err != nil)
		t.End()
	})
}

func TestFindSupertapeFilesUnreadable(t *testing.T) {
	Test(t, "interp: findSupertapeFiles errors on unreadable tape file", func(t *T) {
		dir := t.TB().TempDir()
		path := dir + "/bad_tape.go"
		if err := os.WriteFile(path, []byte(`package main
import . "tapeapi"
`), 0o644); err != nil {
			t.TB().Fatal(err)
		}
		if err := os.Chmod(path, 0o000); err != nil {
			t.TB().Fatal(err)
		}
		defer os.Chmod(path, 0o644)
		_, err := findSupertapeFiles(dir, nil)
		t.Ok(err != nil)
		t.End()
	})
}

func TestRunInterpModePass(t *testing.T) {
	Test(t, "interp: runInterpMode passes for passing source", func(t *T) {
		var out strings.Builder
		d := formatter.New("tap", &out, 1)
		code := runInterpMode("testdata/interp/pass_tape.go", d, &out)
		t.Ok(code == 0 && strings.Contains(out.String(), "ok"))
		t.End()
	})
}

func TestRunInterpModeFail(t *testing.T) {
	Test(t, "interp: runInterpMode fails for failing source", func(t *T) {
		var out strings.Builder
		d := formatter.New("tap", &out, 1)
		code := runInterpMode("testdata/interp/fail_tape.go", d, &out)
		t.Ok(code == 1 && strings.Contains(out.String(), "not ok"))
		t.End()
	})
}

func TestRunInterpModeNoFiles(t *testing.T) {
	Test(t, "interp: runInterpMode errors when no files found", func(t *T) {
		var out strings.Builder
		d := formatter.New("tap", &out, 1)
		code := runInterpMode("testdata/interp/plain.go", d, &out)
		t.Ok(code == 1 && strings.Contains(out.String(), "no supertape files found"))
		t.End()
	})
}

func TestRunInterpModeMissingDir(t *testing.T) {
	Test(t, "interp: runInterpMode errors on missing dir", func(t *T) {
		var out strings.Builder
		d := formatter.New("tap", &out, 1)
		code := runInterpMode("testdata/does-not-exist", d, &out)
		t.Equal(code, 1)
		t.End()
	})
}

func TestRunInterpModeInvalidSource(t *testing.T) {
	Test(t, "interp: runInterpMode exits 1 on interpret error", func(t *T) {
		dir := t.TB().TempDir()
		src := `package main

import . "tapeapi"

func main() { Test("bad", func(t T) { t this is not valid go }) }`
		if err := os.WriteFile(dir+"/bad_tape.go", []byte(src), 0o644); err != nil {
			t.TB().Fatal(err)
		}
		var out strings.Builder
		d := formatter.New("tap", &out, 1)
		code := runInterpMode(dir, d, &out)
		t.Ok(code == 1 && strings.Contains(out.String(), "tape: interpret"))
		t.End()
	})
}
