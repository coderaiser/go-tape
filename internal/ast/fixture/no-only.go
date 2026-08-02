//go:build ignore

package fixture

import (
	"testing"
	Test "github.com/coderaiser/go-tape"
)

func TestNoOnlyParser(t *testing.T) {
	Test(t, "parser: run action", func(t *Test.T) { t.Ok(true); t.End() })
	Test(t, "parser: fail action", func(t *Test.T) { t.Ok(true); t.End() })
}
