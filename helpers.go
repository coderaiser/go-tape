package tape

import (
	"reflect"
	"regexp"
	"runtime"
	"strings"
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
			count[t],
			file,
			line,
		)
	}
	t.Fatalf(
		"too many assertions: got %d, expected 1",
		count[t],
	)
}

func deepEqual(want, got any) bool {
	return reflect.DeepEqual(want, got)
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

func match(s, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}
