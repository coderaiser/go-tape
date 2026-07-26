package tape

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"sync"
	"testing"

	"github.com/coderaiser/go-tape/internal/config"
)

var (
	mu    sync.Mutex
	count = make(map[*testing.T]int)
)

func assertOne(t *testing.T) {
	t.Helper()
	mu.Lock()
	count[t] = 0
	mu.Unlock()
	// cleanup prevents memory leak — removes entry after test finishes
	t.Cleanup(func() {
		mu.Lock()
		delete(count, t)
		mu.Unlock()
	})
}

func hit(t *testing.T) {
	t.Helper()
	mu.Lock()
	defer mu.Unlock()
	count[t]++
	if count[t] <= 1 {
		return
	}
	if !config.CheckAssertionsCount() {
		return
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		t.Fatalf(
			"too many assertions: got %d, expected 1\nat %s:%d",
			count[t], file, line,
		)
	}
	t.Fatalf("too many assertions: got %d, expected 1", count[t])
}

func isPrimitive(v any) bool {
	switch v.(type) {
	case bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		complex64, complex128,
		string, uintptr:
		return true
	}
	t := reflect.TypeOf(v)
	return t != nil && t.Kind() == reflect.Pointer
}

func truthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case string:
		return val != ""
	default:
		return true
	}
}

func toRegexp(pattern any) (*regexp.Regexp, error) {
	switch p := pattern.(type) {
	case *regexp.Regexp:
		return p, nil
	case string:
		return regexp.Compile(p)
	default:
		return nil, fmt.Errorf("pattern must be string or *regexp.Regexp, got %T", pattern)
	}
}
