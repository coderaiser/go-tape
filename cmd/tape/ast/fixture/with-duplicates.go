package fixture

import (
	"testing"
	tape "github.com/coderaiser/go-tape"
)

func TestDuplicates(t *testing.T) {
	tape.Test(t, "parser: run action", func(t *tape.T) { t.Ok(true); t.End() })
	tape.Test(t, "parser: run action", func(t *tape.T) { t.Ok(true); t.End() })
}
