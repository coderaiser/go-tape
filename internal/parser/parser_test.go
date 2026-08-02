package parser_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/model"
	"github.com/coderaiser/go-tape/internal/parser"
)

func TestParsePassEvent(t *testing.T) {
	tape.Test(t, "parser: parse pass event", func(t *tape.T) {
		e, error := parser.Parse(`{"Action":"pass","Package":"mypkg","Test":"TestFoo","Elapsed":0.1}`)
		if error != nil {
			t.TB().Fatalf("Parse: %v", error)
		}
		t.DeepEqual(e, model.Event{Action: "pass", Package: "mypkg", Test: "TestFoo", Elapsed: 0.1})
		t.End()
	})
}

func TestParseFailEvent(t *testing.T) {
	tape.Test(t, "parser: parse fail event", func(t *tape.T) {
		e, error := parser.Parse(`{"Action":"fail","Package":"mypkg","Test":"TestBar"}`)
		if error != nil {
			t.TB().Fatalf("Parse: %v", error)
		}
		t.DeepEqual(e, model.Event{Action: "fail", Package: "mypkg", Test: "TestBar"})
		t.End()
	})
}

func TestParseSkipEvent(t *testing.T) {
	tape.Test(t, "parser: parse skip event", func(t *tape.T) {
		e, error := parser.Parse(`{"Action":"skip","Package":"mypkg","Test":"TestBaz"}`)
		if error != nil {
			t.TB().Fatalf("Parse: %v", error)
		}
		t.DeepEqual(e, model.Event{Action: "skip", Package: "mypkg", Test: "TestBaz"})
		t.End()
	})
}

func TestParseOutputEvent(t *testing.T) {
	tape.Test(t, "parser: parse output event", func(t *tape.T) {
		e, error := parser.Parse(`{"Action":"output","Package":"mypkg","Test":"TestFoo","Output":"ok\\n"}`)
		if error != nil {
			t.TB().Fatalf("Parse: %v", error)
		}
		t.DeepEqual(e, model.Event{Action: "output", Package: "mypkg", Test: "TestFoo", Output: "ok\\n"})
		t.End()
	})
}

func TestParseRunEvent(t *testing.T) {
	tape.Test(t, "parser: parse run event", func(t *tape.T) {
		e, error := parser.Parse(`{"Action":"run","Package":"mypkg","Test":"TestFoo"}`)
		if error != nil {
			t.TB().Fatalf("Parse: %v", error)
		}
		t.DeepEqual(e, model.Event{Action: "run", Package: "mypkg", Test: "TestFoo"})
		t.End()
	})
}

func TestParsePauseEvent(t *testing.T) {
	tape.Test(t, "parser: parse pause event", func(t *tape.T) {
		e, error := parser.Parse(`{"Action":"pause","Package":"mypkg"}`)
		if error != nil {
			t.TB().Fatalf("Parse: %v", error)
		}
		t.DeepEqual(e, model.Event{Action: "pause", Package: "mypkg"})
		t.End()
	})
}

func TestParseContEvent(t *testing.T) {
	tape.Test(t, "parser: parse cont event", func(t *tape.T) {
		e, error := parser.Parse(`{"Action":"cont","Package":"mypkg"}`)
		if error != nil {
			t.TB().Fatalf("Parse: %v", error)
		}
		t.DeepEqual(e, model.Event{Action: "cont", Package: "mypkg"})
		t.End()
	})
}

func TestParseInvalidJSON(t *testing.T) {
	tape.Test(t, "parser: parse invalid json errors", func(t *tape.T) {
		_, error := parser.Parse("not json")
		t.Equal(error.Error(), "parse event: invalid character 'o' in literal null (expecting 'u')")
		t.End()
	})
}

func TestParseEmptyLine(t *testing.T) {
	tape.Test(t, "parser: parse empty line errors", func(t *tape.T) {
		_, error := parser.Parse("")
		t.Equal(error.Error(), "parse event: unexpected end of JSON input")
		t.End()
	})
}
