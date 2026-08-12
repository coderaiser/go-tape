package interp

import (
	"bytes"
	"strings"
	"testing"
)

// run collects all events emitted while interpreting src.
func run(src string) ([]Event, error) {
	var evs []Event
	_, err := Run(src, nil, func(e Event) {
		evs = append(evs, e)
	})
	return evs, err
}

// kinds returns the sequence of event kinds emitted.
func kinds(evs []Event) []string {
	out := make([]string, 0, len(evs))
	for _, e := range evs {
		out = append(out, e.Kind)
	}
	return out
}

func TestRunPass(t *testing.T) {
	evs, err := run(`package main

import . "tapeapi"

func main() {
	Test("scope: works", func(t T) {
		t.Equal(1, 1)
		t.End()
	})
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 4 {
		t.Fatalf("expected 4 events, got %d: %v", len(evs), evs)
	}
	if evs[0].Kind != "test" || evs[0].Name != "scope: works" {
		t.Fatalf("expected test event, got %+v", evs[0])
	}
	if evs[1].Kind != "assert" || evs[1].Name != "equal" || !evs[1].Ok {
		t.Fatalf("expected passing equal assert, got %+v", evs[1])
	}
	if evs[2].Kind != "end" {
		t.Fatalf("expected end event from t.End(), got %+v", evs[2])
	}
	if evs[3].Kind != "end" || !evs[3].Ok {
		t.Fatalf("expected final end event with ok=true, got %+v", evs[3])
	}
}

func TestRunFail(t *testing.T) {
	evs, err := run(`package main

import . "tapeapi"

func main() {
	Test("scope: fails", func(t T) {
		t.Equal(1, 2)
		t.End()
	})
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 4 {
		t.Fatalf("expected 4 events, got %d: %v", len(evs), evs)
	}
	if evs[1].Ok {
		t.Fatalf("expected failing assert, got %+v", evs[1])
	}
	if evs[3].Ok {
		t.Fatalf("expected final end ok=false for failing test, got %+v", evs[3])
	}
}

func TestRunAssertions(t *testing.T) {
	evs, err := run(`package main

import . "tapeapi"

func main() {
	Test("scope: assertions", func(t T) {
		t.DeepEqual([]int{1, 2}, []int{1, 2})
		t.Ok(true)
		t.NotOk(false)
		t.EqualText("a", "a", "msg")
		t.End()
	})
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 7 {
		t.Fatalf("expected 7 events, got %d: %v", len(evs), evs)
	}
	ops := []string{evs[1].Name, evs[2].Name, evs[3].Name, evs[4].Name}
	want := []string{"deepEqual", "ok", "notOk", "equal"}
	for i := range want {
		if ops[i] != want[i] {
			t.Fatalf("event %d operator = %q, want %q", i+1, ops[i], want[i])
		}
		if !evs[i+1].Ok {
			t.Fatalf("event %d should pass, got %+v", i+1, evs[i+1])
		}
	}
	if evs[4].Msg != "msg" {
		t.Fatalf("EqualText msg = %q, want %q", evs[4].Msg, "msg")
	}
}

func TestRunReport(t *testing.T) {
	evs, err := run(`package main

import . "tapeapi"

func main() {
	Report(Result{Ok: true, Operator: "equal", Got: 1, Expected: 1, Message: "should equal"})
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 1 {
		t.Fatalf("expected 1 event, got %d: %v", len(evs), evs)
	}
	e := evs[0]
	if e.Kind != "report" || e.Name != "equal" || !e.Ok || e.Got != 1 || e.Wanted != 1 || e.Msg != "should equal" {
		t.Fatalf("unexpected report event: %+v", e)
	}
}

func TestRunPrintOutput(t *testing.T) {
	var buf bytes.Buffer
	_, err := Run(`package main

import . "tapeapi"

func main() {
	println("hello from interp")
	Test("scope: works", func(t T) {
		t.Ok(true)
		t.End()
	})
}
`, &buf, nil)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello from interp") {
		t.Fatalf("expected print output captured, got %q", buf.String())
	}
}

func TestRunInvalidSource(t *testing.T) {
	_, err := run(`package main

import . "tapeapi"

func main() {
	this is not valid go
}
`)
	if err == nil {
		t.Fatal("expected error for invalid source")
	}
}

func TestRunNoTests(t *testing.T) {
	evs, err := run(`package main

func main() {
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 0 {
		t.Fatalf("expected no events, got %d: %v", len(evs), evs)
	}
}

func TestRunMultipleTests(t *testing.T) {
	evs, err := run(`package main

import . "tapeapi"

func main() {
	Test("scope: first", func(t T) {
		t.Ok(true)
		t.End()
	})
	Test("scope: second", func(t T) {
		t.NotOk(false)
		t.End()
	})
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 8 {
		t.Fatalf("expected 8 events, got %d: %v", len(evs), evs)
	}
	if evs[0].Name != "scope: first" || evs[4].Name != "scope: second" {
		t.Fatalf("unexpected test names: %q, %q", evs[0].Name, evs[4].Name)
	}
}

func TestTruthy(t *testing.T) {
	type strct struct{ X int }
	cases := []struct {
		in   any
		want bool
	}{
		{nil, false},
		{true, true},
		{false, false},
		{"", false},
		{"x", true},
		{0, false},
		{1, true},
		{[]int{}, true},
		{[]int(nil), false},
		{map[string]int{}, true},
		{map[string]int(nil), false},
		{make(chan int), true},
		{(*int)(nil), false},
		{new(int), true},
		{varIface(nil), false},
		{varIface(1), true},
		{strct{}, true},
	}
	for i, c := range cases {
		if got := truthy(c.in); got != c.want {
			t.Fatalf("truthy(%v) = %v, want %v", c.in, got, c.want)
		}
		_ = i
	}
}

// varIface boxes v into an interface{} to exercise the interface branch.
func varIface(v any) any { return v }

func TestEqHelpers(t *testing.T) {
	if !eq(1, 1) {
		t.Fatal("eq(1,1) should be true")
	}
	if eq(1, 2) {
		t.Fatal("eq(1,2) should be false")
	}
	if !deepEq([]int{1}, []int{1}) {
		t.Fatal("deepEq should hold for equal slices")
	}
	if deepEq([]int{1}, []int{2}) {
		t.Fatal("deepEq should fail for different slices")
	}
}

func TestFailingAssertionsMarkEnd(t *testing.T) {
	evs, err := run(`package main

import . "tapeapi"

func main() {
	Test("scope: all fail", func(t T) {
		t.DeepEqual([]int{1}, []int{2})
		t.Ok(false)
		t.NotOk(true)
		t.End()
	})
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 6 {
		t.Fatalf("expected 6 events, got %d: %v", len(evs), evs)
	}
	for i := 1; i <= 3; i++ {
		if evs[i].Ok {
			t.Fatalf("event %d should fail, got %+v", i, evs[i])
		}
	}
	if evs[5].Ok {
		t.Fatalf("final end should be ok=false, got %+v", evs[5])
	}
}

// TestFailingEqualText exercises the EqualText failing branch.
func TestFailingEqualText(t *testing.T) {
	evs, err := run(`package main

import . "tapeapi"

func main() {
	Test("scope: eqtext", func(t T) {
		t.EqualText("a", "b", "boom")
		t.End()
	})
}
`)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(evs) != 4 {
		t.Fatalf("expected 4 events, got %d: %v", len(evs), evs)
	}
	if evs[1].Ok {
		t.Fatalf("EqualText should fail, got %+v", evs[1])
	}
	if evs[1].Msg != "boom" {
		t.Fatalf("EqualText msg = %q, want boom", evs[1].Msg)
	}
}