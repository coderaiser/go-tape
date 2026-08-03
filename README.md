# Tape [![License][LicenseIMGURL]][LicenseURL] [![Build Status][BuildStatusIMGURL]][BuildStatusURL] [![Coverage Status][CoverageIMGURL]][CoverageURL]

[BuildStatusURL]: https://github.com/coderaiser/go-tape/actions/workflows/test.yml
[BuildStatusIMGURL]: https://github.com/coderaiser/go-tape/actions/workflows/test.yml/badge.svg
[LicenseURL]: https://tldrlegal.com/license/mit-license "MIT License"
[LicenseIMGURL]: https://img.shields.io/badge/license-MIT-317BF9.svg?style=flat
[CoverageURL]: https://coveralls.io/github/coderaiser/go-tape?branch=master
[CoverageIMGURL]: https://coveralls.io/repos/coderaiser/go-tape/badge.svg?branch=master&service=github

📼 [**Supertape**](https://github.com/coderaiser/supertape) for Go — same assertions, same philosophy.

## Install

```sh
go install github.com/coderaiser/go-tape/cmd/tape@latest
```

## Usage

```sh
tape [options] [path]
tape ./...
tape ./internal/...
```

## API

Import the package aliased as `Test` to get the `Test(...)`, `Test.Only(...)`, and `Test.Skip(...)` syntax that mirrors supertape:

```go
import (
    "testing"
    Test "github.com/coderaiser/go-tape"
)

func TestEqual(t *testing.T) {
    Test(t, "tape: Equal works", func(t *Test.T) {
        t.Equal(42, 42)
        t.End()
    })
}

func TestNotEqual(t *testing.T) {
    Test(t, "tape: NotEqual works", func(t *Test.T) {
        t.NotEqual(1, 2)
        t.End()
    })
}

func TestOk(t *testing.T) {
    Test(t, "tape: Ok works", func(t *Test.T) {
        t.Ok(true)
        t.End()
    })
}

func TestNotOk(t *testing.T) {
    Test(t, "tape: NotOk works", func(t *Test.T) {
        t.NotOk(false)
        t.End()
    })
}

func TestDeepEqual(t *testing.T) {
    Test(t, "tape: DeepEqual works", func(t *Test.T) {
        t.DeepEqual([]int{1, 2}, []int{1, 2})
        t.End()
    })
}

func TestMatch(t *testing.T) {
    Test(t, "tape: Match works", func(t *Test.T) {
        t.Match("hello 123", `hello \d+`)
        t.End()
    })
    func TestMatch(t *testing.T) {
    Test(t, "tape: Match: RegExp", func(t *Test.T) {
        t.Match("hello 123", regexp.MustCompile(`hello \d+`))
        t.End()
    })
    Test(t, "tape: Match: string", func(t *Test.T) {
        t.Match("hello 123", `hello`)
        t.End()
    })
}

}

func TestComment(t *testing.T) {
    Test(t, "tape: Comment does not count as assertion", func(t *Test.T) {
        t.Comment("just a note")
        t.Ok(true)
        t.End()
    })
}
```

### Test.Only

Run a single test, skipping all others — identical to supertape's `test.only`:

```go
func TestParser(t *testing.T) {
    Test.Only(t, "parser: run action", func(t *Test.T) {
        t.Ok(true)
        t.End()
    })
    Test(t, "parser: other test", func(t *Test.T) {
        // skipped — only the above runs
        t.Ok(true)
        t.End()
    })
}
```

`tape` detects `Test.Only` calls via AST scan and passes a `-run` filter to `go test` automatically. No changes to the test command needed.

### Test.Skip

Skip a single test without removing it:

```go
func TestParser(t *testing.T) {
    Test.Skip(t, "parser: known broken", func(t *Test.T) {
        t.Ok(false) // never runs
        t.End()
    })
    Test(t, "parser: works fine", func(t *Test.T) {
        t.Ok(true)
        t.End()
    })
}
```

### Extend

Create a custom assertion type that wraps `*Test.T`:

```go
type MyT struct{ *Test.T }

func (t *MyT) FileExists(path string) {
    t.TB().Helper()
    _, err := os.Stat(path)
    t.NotOk(err)
}

var MyTest = Test.Extend(func(base *Test.T) *MyT {
    return &MyT{T: base}
})

func TestFiles(t *testing.T) {
    MyTest(t, "files: config exists", func(t *MyT) {
        t.FileExists("config.toml")
        t.End()
    })
}
```

## CLI

```
Usage: go-tape [options] [path]

Options:
  -h                           display this help and exit
  -v                           output version information and exit
  -f, --format format          use a specific output format
                               default: progress-bar (tap on CI)
                               values: tap|progress-bar|short|fail|time|json-lines
  --no-check-scopes            do not check that messages contain scope: 'scope: message'
  --no-check-assertions-count  do not check that assertion count is no more than 1
  --no-check-duplicates        do not check messages for duplicates
```

## Environment variables

| Variable | Description |
|---|---|
| `TAPE_PROGRESS_BAR` | Force progress bar on (`1`) or off (`0`) |
| `TAPE_PROGRESS_BAR_MIN` | Minimum test count to show progress bar (default: `10`) |
| `TAPE_PROGRESS_BAR_COLOR` | Progress bar color — hex (`#f9d472`) or ANSI escape |
| `TAPE_PROGRESS_BAR_STACK` | Set to `0` to hide error stack traces |
| `TAPE_TERM_WIDTH` | Override terminal width for the progress bar |
| `TAPE_TIMEOUT` | Per-test timeout (default: `5s`, e.g. `TAPE_TIMEOUT=30s`) |
| `TAPE_CHECK_SCOPES` | Set to `0` to disable scope format check |
| `TAPE_CHECK_ASSERTIONS_COUNT` | Set to `0` to allow more than one assertion per test |
| `TAPE_CHECK_END` | Set to `0` to allow omitting `t.End()` |
| `TAPE_CHECK_SKIPPED` | Set to `1` to exit with error when skipped tests exist |
| `TAPE_TIME_CLOCK` | Override the clock emoji used by the `time` formatter |

## License

MIT
