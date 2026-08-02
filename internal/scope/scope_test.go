package scope_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/scope"
)

func TestValidScopeMessage(t *testing.T) {
	tape.Test(t, "scope: valid scope message passes", func(t *tape.T) {
		t.Ok(scope.Valid("config: loads defaults"))
		t.End()
	})
}

func TestValidSingleWordScope(t *testing.T) {
	tape.Test(t, "scope: single word scope passes", func(t *tape.T) {
		t.Ok(scope.Valid("a: message"))
		t.End()
	})
}

func TestInvalidNoColon(t *testing.T) {
	tape.Test(t, "scope: missing colon is invalid", func(t *tape.T) {
		t.NotOk(scope.Valid("just a message"))
		t.End()
	})
}

func TestInvalidEmptyScope(t *testing.T) {
	tape.Test(t, "scope: empty scope is invalid", func(t *tape.T) {
		t.NotOk(scope.Valid(": message"))
		t.End()
	})
}

func TestInvalidEmptyMessage(t *testing.T) {
	tape.Test(t, "scope: empty message is invalid", func(t *tape.T) {
		t.NotOk(scope.Valid("scope: "))
		t.End()
	})
}

func TestInvalidEmptyString(t *testing.T) {
	tape.Test(t, "scope: empty string is invalid", func(t *tape.T) {
		t.NotOk(scope.Valid(""))
		t.End()
	})
}

func TestInvalidNoSpaceAfterColon(t *testing.T) {
	tape.Test(t, "scope: missing space after colon is invalid", func(t *tape.T) {
		t.NotOk(scope.Valid("scope:message"))
		t.End()
	})
}
