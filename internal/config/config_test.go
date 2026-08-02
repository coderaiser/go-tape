package config_test

import (
	"testing"
	"time"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-tape/internal/config"
)

func TestCheckScopesDefault(t *testing.T) {
	tape.Test(t, "config: CheckScopes defaults to true", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_SCOPES", "")
		t.Ok(config.CheckScopes())
		t.End()
	})
}

func TestCheckScopesDisabled(t *testing.T) {
	tape.Test(t, "config: CheckScopes disabled with 0", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_SCOPES", "0")
		t.NotOk(config.CheckScopes())
		t.End()
	})
}

func TestCheckAssertionsCountDefault(t *testing.T) {
	tape.Test(t, "config: CheckAssertionsCount defaults to true", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "")
		t.Ok(config.CheckAssertionsCount())
		t.End()
	})
}

func TestCheckAssertionsCountDisabled(t *testing.T) {
	tape.Test(t, "config: CheckAssertionsCount disabled with false", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "false")
		t.NotOk(config.CheckAssertionsCount())
		t.End()
	})
}

func TestCheckEndDefault(t *testing.T) {
	tape.Test(t, "config: CheckEnd defaults to true", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_END", "")
		t.Ok(config.CheckEnd())
		t.End()
	})
}

func TestCheckEndDisabled(t *testing.T) {
	tape.Test(t, "config: CheckEnd disabled with 0", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_END", "0")
		t.NotOk(config.CheckEnd())
		t.End()
	})
}

func TestTimeoutDefault(t *testing.T) {
	tape.Test(t, "config: Timeout defaults to 3s", func(t *tape.T) {
		t.Setenv("TAPE_TIMEOUT", "")
		t.Equal(config.Timeout(), 3*time.Second)
		t.End()
	})
}

func TestTimeoutCustom(t *testing.T) {
	tape.Test(t, "config: Timeout honors TAPE_TIMEOUT", func(t *tape.T) {
		t.Setenv("TAPE_TIMEOUT", "5s")
		t.Equal(config.Timeout(), 5*time.Second)
		t.End()
	})
}

func TestTimeoutInvalid(t *testing.T) {
	tape.Test(t, "config: Timeout falls back on invalid value", func(t *tape.T) {
		t.Setenv("TAPE_TIMEOUT", "invalid")
		t.Equal(config.Timeout(), 3*time.Second)
		t.End()
	})
}

func TestStrictTransitionsDefault(t *testing.T) {
	tape.Test(t, "config: StrictTransitions defaults to true", func(t *tape.T) {
		t.Setenv("TAPE_STRICT_TRANSITIONS", "")
		t.Ok(config.StrictTransitions())
		t.End()
	})
}

func TestStrictTransitionsDisabled(t *testing.T) {
	tape.Test(t, "config: StrictTransitions disabled with 0", func(t *tape.T) {
		t.Setenv("TAPE_STRICT_TRANSITIONS", "0")
		t.NotOk(config.StrictTransitions())
		t.End()
	})
}

func TestCheckSkippedDefault(t *testing.T) {
	tape.Test(t, "config: CheckSkipped defaults to false", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_SKIPPED", "")
		t.NotOk(config.CheckSkipped())
		t.End()
	})
}

func TestCheckSkippedEnabled(t *testing.T) {
	tape.Test(t, "config: CheckSkipped enabled with 1", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_SKIPPED", "1")
		t.Ok(config.CheckSkipped())
		t.End()
	})
}

func TestEnvBoolTrue(t *testing.T) {
	tape.Test(t, "config: envBool parses true", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_SCOPES", "true")
		t.Ok(config.CheckScopes())
		t.End()
	})
}

func TestEnvBoolInvalid(t *testing.T) {
	tape.Test(t, "config: envBool falls back on invalid value", func(t *tape.T) {
		t.Setenv("TAPE_CHECK_SCOPES", "invalid")
		t.Ok(config.CheckScopes())
		t.End()
	})
}
