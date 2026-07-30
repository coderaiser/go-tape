//go:build windows

package formatter_progress_bar

import (
	"os"
	"strconv"
)

func termWidth() int {
	if v := os.Getenv("TAPE_TERM_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 80
}
