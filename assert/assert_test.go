package assert

import (
	"errors"
	"testing"
)

func TestEqualMatch(t *testing.T) {
	One(t)
	Equal(t, 42, 42)
}

func TestContainsMatch(t *testing.T) {
	One(t)
	Contains(t, "hello world", "world")
}

func TestNoErrorNil(t *testing.T) {
	One(t)
	NoError(t, nil)
}

func TestErrorNonNil(t *testing.T) {
	One(t)
	Error(t, errors.New("some error"))
}

func TestOkTrue(t *testing.T) {
	One(t)
	Ok(t, true)
}

func TestNotOkFalse(t *testing.T) {
	One(t)
	NotOk(t, false)
}
