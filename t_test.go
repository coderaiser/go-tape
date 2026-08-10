package tape

import (
	"os"
	"testing"
)

func TestTB(t *testing.T) {
	Test(t, "tape: TB returns underlying testing.T", func(t *T) {
		t.Ok(t.TB() != nil)
		t.End()
	})
}

func TestSetenv(t *testing.T) {
	Test(t, "tape: Setenv sets env variable", func(t *T) {
		t.Setenv("TAPE_TEST_VAR", "hello")
		t.Equal(os.Getenv("TAPE_TEST_VAR"), "hello")
		t.End()
	})
}

func TestReportOkDoesNotFail(t *testing.T) {
	Test(t, "tape: Report with Ok:true does not fail test", func(t *T) {
		inner := &testing.T{}
		tt := newT(inner)
		tt.Report(Result{Ok: true, Message: "ok"})
		t.NotOk(inner.Failed())
		t.End()
	})
}

func TestReportFailEmitsTAPELine(t *testing.T) {
	Test(t, "tape: Report with Ok:false marks test failed", func(t *T) {
		// The TAPE: log line is captured by go test -json in real runs.
		// A zero-value *testing.T lets us observe Fail() without failing the
		// test runner (the outer tape.Test subtest stays green).
		inner := &testing.T{}
		tt := newT(inner)
		tt.Report(Result{Ok: false, Message: "should equal", Result: 2, Expected: 1})
		t.Ok(inner.Failed())
		t.End()
	})
}
