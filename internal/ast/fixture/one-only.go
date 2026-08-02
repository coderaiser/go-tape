//go:build ignore

package fixture

import (
	"testing"
	Test "github.com/coderaiser/go-tape"
)

func TestOneOnlyParser(t *testing.T) {
	Test.Only(t, "parser: run action", func(t *Test.T) { t.Ok(true); t.End() })
	Test(t, "parser: other test", func(t *Test.T) { t.Ok(true); t.End() })
}
