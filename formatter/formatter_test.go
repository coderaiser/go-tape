package formatter

import (
	"strings"
	"testing"

	"github.com/coderaiser/go-tape/model"
	"github.com/coderaiser/go-tape/state"
)

func TestTAPOkLine(t *testing.T) {
	f := NewTAP()
	f.Add([]string{"pkg: TestFoo"}, nil, nil)
	output := f.Format()
	if !strings.Contains(output, "ok 1 - pkg: TestFoo") {
		t.Errorf("expected ok line, got: %s", output)
	}
}

func TestTAPNotOkLine(t *testing.T) {
	f := NewTAP()
	f.Add(nil, []string{"pkg: TestFoo"}, nil)
	output := f.Format()
	if !strings.Contains(output, "not ok 1 - pkg: TestFoo") {
		t.Errorf("expected not ok line, got: %s", output)
	}
}

func TestTAPSkipLine(t *testing.T) {
	f := NewTAP()
	f.Add(nil, nil, []string{"pkg: TestFoo"})
	output := f.Format()
	if !strings.Contains(output, "ok 1 - # SKIP pkg: TestFoo") {
		t.Errorf("expected skip line, got: %s", output)
	}
}

func TestTAPSummaryCounts(t *testing.T) {
	f := NewTAP()
	f.Add([]string{"a", "b"}, []string{"c"}, []string{"d"})
	output := f.Format()
	if !strings.Contains(output, "# pass 2") {
		t.Errorf("expected pass 2, got: %s", output)
	}
	if !strings.Contains(output, "# fail 1") {
		t.Errorf("expected fail 1, got: %s", output)
	}
	if !strings.Contains(output, "# skip 1") {
		t.Errorf("expected skip 1, got: %s", output)
	}
}

func TestTAPPlanLine(t *testing.T) {
	f := NewTAP()
	f.Add([]string{"a", "b"}, nil, nil)
	output := f.Format()
	if !strings.Contains(output, "1..2") {
		t.Errorf("expected plan line, got: %s", output)
	}
}

func TestShortPass(t *testing.T) {
	var sf ShortFormatter
	sf.Add([]string{"a", "b"}, nil, nil)
	output := sf.Format()
	if !strings.Contains(output, "PASS: 2/2 tests passed") {
		t.Errorf("expected PASS line, got: %s", output)
	}
}

func TestShortFail(t *testing.T) {
	var sf ShortFormatter
	sf.Add(nil, []string{"a"}, nil)
	output := sf.Format()
	if !strings.Contains(output, "FAIL: 1/1 tests failed") {
		t.Errorf("expected FAIL line, got: %s", output)
	}
}

func TestShortSkip(t *testing.T) {
	var sf ShortFormatter
	sf.Add([]string{"a"}, nil, []string{"b"})
	output := sf.Format()
	if !strings.Contains(output, "PASS: 1/2 tests passed") {
		t.Errorf("expected PASS with skip, got: %s", output)
	}
	if !strings.Contains(output, "(1 skipped)") {
		t.Errorf("expected skip count, got: %s", output)
	}
}

func TestProgressPass(t *testing.T) {
	var pf ProgressFormatter
	pf.Add([]string{"a"}, nil, nil)
	output := pf.Format()
	if !strings.Contains(output, "##[section]Tests passed: 1/1") {
		t.Errorf("expected section line, got: %s", output)
	}
}

func TestProgressFail(t *testing.T) {
	var pf ProgressFormatter
	pf.Add(nil, []string{"a"}, nil)
	output := pf.Format()
	if !strings.Contains(output, "##[error]Tests failed: 1/1") {
		t.Errorf("expected error line, got: %s", output)
	}
}

func TestTAPVersionLine(t *testing.T) {
	f := NewTAP()
	f.Add(nil, nil, nil)
	output := f.Format()
	if !strings.Contains(output, "TAP version 14") {
		t.Errorf("expected TAP version 14, got: %s", output)
	}
}

func TestFormatFromStore(t *testing.T) {
	s := state.New()
	s.Apply(model.Event{Action: "run", Test: "TestA"})
	s.Apply(model.Event{Action: "pass", Test: "TestA"})
	s.Apply(model.Event{Action: "output", Test: "TestA", Output: "ok"})
	output := FormatFromStore(s)
	if !strings.Contains(output, "TAP version 14") {
		t.Errorf("expected TAP header, got: %s", output)
	}
}
