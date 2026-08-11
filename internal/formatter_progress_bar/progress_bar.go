package formatter_progress_bar

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/coderaiser/go-tape/internal/stream"
)

const (
	barWidth      = 40
	barComplete   = '\u2588'
	barEmpty      = '\u2591' // ░ LIGHT SHADE — eaw=N, always 1 column
	OkEmoji       = "\U0001f44c"
	FailEmoji     = "\u274c"
	SkipEmoji     = "\u26a0\ufe0f"
	OkMark        = "\u2705"
	DEFAULT_COLOR = "#f9d472"
)

var okEmoji = OkEmoji
var failEmoji = FailEmoji
var skipEmoji = SkipEmoji

type ProgressBarFormatter struct {
	total    int
	show     bool
	Color    string
	stackEnv string
	out      strings.Builder
}

func New(total int) *ProgressBarFormatter {
	color := os.Getenv("TAPE_PROGRESS_BAR_COLOR")
	if color == "" {
		color = DEFAULT_COLOR
	}

	min := 100
	if v := os.Getenv("TAPE_PROGRESS_BAR_MIN"); v != "" {
		if n, err := fmt.Sscanf(v, "%d", &min); n != 1 || err != nil {
			min = 100
		}
	}

	show := total >= min
	if v := os.Getenv("TAPE_PROGRESS_BAR"); v == "1" {
		show = true
	} else if v == "0" {
		show = false
	}

	return &ProgressBarFormatter{
		total:    total,
		Color:    color,
		stackEnv: os.Getenv("TAPE_PROGRESS_BAR_STACK"),
		show:     show,
	}
}

// Event dispatches on e.Type. Fail blocks are buffered; other types are no-ops
// for the progress-bar formatter (it shows progress on stderr, TAP on stdout).
func (f *ProgressBarFormatter) Event(e stream.Event) string {
	switch e.Type {
	case stream.TypeTestEnd:
		// drive the progress bar on stderr
		if f.show {
			failStr := okEmoji
			if e.Failed > 0 {
				failStr = fmt.Sprintf("\033[31m%d\033[0m", e.Failed)
			}
			bar := renderBar(e.Count, e.Total, f.Color)
			pct := 0
			if e.Total > 0 {
				pct = e.Count * 100 / e.Total
				if pct > 100 {
					pct = 100
				}
			}
			displayTotal := e.Total
			if e.Count > displayTotal {
				displayTotal = e.Count
			}
			truncName := truncate(e.Test, 40)
			line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s", bar, pct, failStr, e.Count, displayTotal, truncName)
			width := termWidth()
			if visibleLen(line) > width {
				line = truncateANSI(line, width)
			}
			fmt.Fprintf(os.Stderr, "\r%s", line)
		}
		return ""

	case stream.TypeFail, stream.TypeUnknownFail:
		// buffer for flush in End
		var sb strings.Builder
		fmt.Fprintf(&sb, "\n# %s\n", e.Test)
		fmt.Fprintf(&sb, "%s not ok %d %s\n", failEmoji, e.Count, e.Test)
		sb.WriteString("  ---\n")
		fmt.Fprintf(&sb, "    operator: %s\n", e.Operator)
		if e.Output != "" {
			sb.WriteString(e.Output)
		} else {
			fmt.Fprintf(&sb, "    expected: |-\n      %v\n", e.Expected)
			fmt.Fprintf(&sb, "    result: |-\n      %v\n", e.Result)
		}
		if e.At != "" {
			fmt.Fprintf(&sb, "    %s\n", e.At)
		}
		if f.stackEnv != "0" && e.ErrorStack != "" {
			fmt.Fprintf(&sb, "    stack: |-\n%s\n", e.ErrorStack)
		}
		sb.WriteString("  ...\n")
		sb.WriteString("\n")
		f.out.WriteString(sb.String())
		return ""

	case stream.TypeComment:
		return fmt.Sprintf("# %s\n", e.Message)
	}
	return ""
}

// End clears the progress bar and flushes all buffered output + summary.
func (f *ProgressBarFormatter) End(passed, failed, skipped int) string {
	if f.show {
		fmt.Fprintf(os.Stderr, "\r\033[2K\033[?25h")
	}

	lines := []string{}
	// flush buffered fail output
	if s := f.out.String(); s != "" {
		parts := strings.Split(strings.TrimRight(s, "\n"), "\n")
		lines = append(lines, parts...)
	}

	total := passed + failed + skipped
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("1..%d", total))
	lines = append(lines, fmt.Sprintf("# tests %d", total))
	lines = append(lines, fmt.Sprintf("# pass %d", passed))
	if skipped > 0 {
		lines = append(lines, fmt.Sprintf("# %s skip %d", skipEmoji, skipped))
	}
	lines = append(lines, "")
	if failed > 0 {
		lines = append(lines, fmt.Sprintf("# %s fail %d", failEmoji, failed))
	} else {
		lines = append(lines, strings.TrimRight(okLine(), "\n"))
	}
	lines = append(lines, "")
	lines = append(lines, "")

	result := strings.Join(lines, "\n")
	if f.show {
		return "\r" + result
	}
	return result
}

// okLine returns the success marker line. JetBrains' JediTerm misaligns the
// emoji, so an extra space is added when running under that terminal.
func okLine() string {
	spaces := ""
	if os.Getenv("TERMINAL_EMULATOR") == "JetBrains-JediTerm" {
		spaces = " "
	}
	return fmt.Sprintf("# %s%s ok\n", OkMark, spaces)
}

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

func visibleLen(s string) int {
	return len([]rune(ansiEscape.ReplaceAllString(s, "")))
}

func truncateANSI(s string, n int) string {
	visible := 0
	var out []byte
	i := 0
	b := []byte(s)
	for i < len(b) {
		if b[i] == 0x1b && i+1 < len(b) && b[i+1] == '[' {
			j := i + 2
			for j < len(b) && b[j] != 'm' {
				j++
			}
			if j < len(b) {
				j++
			}
			out = append(out, b[i:j]...)
			i = j
			continue
		}
		if visible >= n {
			break
		}
		r, size := decodeRuneAt(b, i)
		_ = r
		out = append(out, b[i:i+size]...)
		visible++
		i += size
	}
	return string(out)
}

func decodeRuneAt(b []byte, i int) (rune, int) {
	if b[i] < 0x80 {
		return rune(b[i]), 1
	}
	var size int
	switch {
	case b[i]&0xE0 == 0xC0:
		size = 2
	case b[i]&0xF0 == 0xE0:
		size = 3
	case b[i]&0xF8 == 0xF0:
		size = 4
	default:
		return '\uFFFD', 1
	}
	if i+size > len(b) {
		return '\uFFFD', 1
	}
	r := rune(b[i] & (0xFF >> size))
	for k := 1; k < size; k++ {
		if b[i+k]&0xC0 != 0x80 {
			return '\uFFFD', 1
		}
		r = r<<6 | rune(b[i+k]&0x3F)
	}
	return r, size
}

func hexToANSI(color string) string {
	if len(color) != 7 || color[0] != '#' {
		return color
	}
	parse := func(s string) int {
		n := 0
		for _, c := range s {
			n <<= 4
			switch {
			case c >= '0' && c <= '9':
				n |= int(c - '0')
			case c >= 'a' && c <= 'f':
				n |= int(c-'a') + 10
			case c >= 'A' && c <= 'F':
				n |= int(c-'A') + 10
			}
		}
		return n
	}
	r := parse(color[1:3])
	g := parse(color[3:5])
	b := parse(color[5:7])
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func renderBar(done, total int, color string) string {
	ansi := hexToANSI(color)
	if total == 0 {
		return fmt.Sprintf("%s%s\033[0m", ansi, strings.Repeat(string(barEmpty), barWidth))
	}
	filled := done * barWidth / total
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat(string(barComplete), filled) + strings.Repeat(string(barEmpty), barWidth-filled)
	return fmt.Sprintf("%s%s\033[0m", ansi, bar)
}

func truncate(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	return string([]rune(s)[:n-3]) + "..."
}

// RenderBar exported for time formatter
var RenderBar = renderBar

// Truncate exported for time formatter
var Truncate = truncate
