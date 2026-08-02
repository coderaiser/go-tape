package formatter_time

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/coderaiser/go-tape/internal/formatter_progress_bar"
)

type TimeFormatter struct {
	*formatter_progress_bar.ProgressBarFormatter
	startTime time.Time
	clock     string
	w         io.Writer
}

func New(total int, w io.Writer) *TimeFormatter {
	clock := os.Getenv("TAPE_TIME_CLOCK")
	if clock == "" {
		clock = "\u23f3"
	}
	return &TimeFormatter{
		ProgressBarFormatter: formatter_progress_bar.New(total),
		clock:                clock,
		w:                    w,
	}
}

func (f *TimeFormatter) Start(total int) string {
	f.startTime = time.Now()
	return f.ProgressBarFormatter.Start(total)
}

func (f *TimeFormatter) TestEnd(count, total, failed int, name string) string {
	elapsed := time.Since(f.startTime)
	m := int(elapsed.Minutes())
	s := int(elapsed.Seconds()) % 60
	timeStr := fmt.Sprintf("%s %02d:%02d", f.clock, m, s)

	failStr := formatter_progress_bar.OkEmoji
	if failed > 0 {
		failStr = fmt.Sprintf("\033[31m%d\033[0m", failed)
	}
	bar := formatter_progress_bar.RenderBar(count, total, f.Color)
	pct := 0
	if total > 0 {
		pct = count * 100 / total
	}
	truncName := formatter_progress_bar.Truncate(name, 30)
	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s | %s",
		bar, pct, failStr, count, total, timeStr, truncName)
	fmt.Fprintf(f.w, "\r%s", line)
	return ""
}
