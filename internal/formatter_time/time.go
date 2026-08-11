package formatter_time

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/coderaiser/go-tape/internal/formatter_progress_bar"
	"github.com/coderaiser/go-tape/internal/stream"
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

// Event handles test-end to write the time-enhanced progress line;
// all other event types are delegated to ProgressBarFormatter.Event.
func (f *TimeFormatter) Event(e stream.Event) string {
	if f.startTime.IsZero() {
		f.startTime = time.Now()
	}
	if e.Type != stream.TypeTestEnd {
		return f.ProgressBarFormatter.Event(e)
	}
	elapsed := time.Since(f.startTime)
	m := int(elapsed.Minutes())
	s := int(elapsed.Seconds()) % 60
	timeStr := fmt.Sprintf("%s %02d:%02d", f.clock, m, s)

	failStr := formatter_progress_bar.OkEmoji
	if e.Failed > 0 {
		failStr = fmt.Sprintf("\033[31m%d\033[0m", e.Failed)
	}
	bar := formatter_progress_bar.RenderBar(e.Count, e.Total, f.Color)
	pct := 0
	if e.Total > 0 {
		pct = e.Count * 100 / e.Total
	}
	truncName := formatter_progress_bar.Truncate(e.Test, 30)
	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s | %s",
		bar, pct, failStr, e.Count, e.Total, timeStr, truncName)
	fmt.Fprintf(f.w, "\r%s", line)
	return ""
}
