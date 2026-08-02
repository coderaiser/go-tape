//go:build ignore

package fixture

import (
	"testing"
	Test "github.com/coderaiser/go-tape"
)

func TestCrossParser(t *testing.T) {
	Test.Only(t, "parser: run action", func(t *Test.T) { t.Ok(true); t.End() })
}

func TestCrossRunner(t *testing.T) {
	Test.Only(t, "runner: starts", func(t *Test.T) { t.Ok(true); t.End() })
}
