//go:build ignore

package fixture

import (
	"testing"
	Test "github.com/coderaiser/go-tape"
)

func TestMultiParser(t *testing.T) {
	Test.Only(t, "parser: run action", func(t *Test.T) { t.Ok(true); t.End() })
	Test.Only(t, "parser: fail action", func(t *Test.T) { t.Ok(true); t.End() })
}
