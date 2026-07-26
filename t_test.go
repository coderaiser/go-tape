package tape

import (
	"os"
	"testing"
)


func TestT(t *testing.T) {
    tape.Test(t, "tape: TB returns underlying testing.T", func(t *tape.T) {
        t.Ok(t.TB() != nil)
        t.End()
    })

    tape.Test(t, "tape: Setenv sets env variable", func(t *tape.T) {
        t.Setenv("TAPE_TEST_VAR", "hello")
        t.Equal(os.Getenv("TAPE_TEST_VAR"), "hello")
        t.End()
    })
}
