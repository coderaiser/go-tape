//go:build linux || darwin

package formatter_progress_bar

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func termWidth() int {
	if v := os.Getenv("TAPE_TERM_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	type winsize struct {
		Row, Col       uint16
		Xpixel, Ypixel uint16
	}
	ws := &winsize{}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stderr.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(ws)))
	if errno != 0 || ws.Col == 0 {
		return 80
	}
	return int(ws.Col)
}
