package formatter_progress_bar

import (
	"fmt"
	"os"
	"regexp"
	"strings"
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

func (f *ProgressBarFormatter) Start(total int) string {
	return "TAP version 13\n"
}

func (f *ProgressBarFormatter) Test(name string) string { return "" }

func (f *ProgressBarFormatter) TestEnd(count, total, failed int, name string) string {
	if !f.show {
		return ""
	}

	failStr := okEmoji
	if failed > 0 {
		failStr = fmt.Sprintf("\033[31m%d\033[0m", failed)
	}
	bar := renderBar(count, total, f.Color)
	pct := 0
	if total > 0 {
		pct = count * 100 / total
		if pct > 100 {
			pct = 100
		}
	}
	displayTotal := total
	if count > displayTotal {
		displayTotal = count
	}
	truncName := truncate(name, 40)
	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s", bar, pct, failStr, count, displayTotal, truncName)
	width := termWidth()
	if visibleLen(line) > width {
		line = truncateANSI(line, width)
	}
	fmt.Fprintf(os.Stderr, "\r%s", line)
	return ""
}

func (f *ProgressBarFormatter) Success(count int, message string) string { return "" }

func (f *ProgressBarFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "\n# %s\n", message)
	fmt.Fprintf(&sb, "%s not ok %d %s\n", failEmoji, count, message)
	sb.WriteString("  ---\n")
	if operator != "" {
		fmt.Fprintf(&sb, "    operator: %s\n", operator)
	}
	if output != "" && operator == "" {
		sb.WriteString(output)
	} else {
		fmt.Fprintf(&sb, "    expected: |-\n      %v\n", expected)
		fmt.Fprintf(&sb, "    result: |-\n      %v\n", result)
	}
	if at != "" {
		fmt.Fprintf(&sb, "    %s\n", at)
	}
	if f.stackEnv != "0" && errorStack != "" {
		fmt.Fprintf(&sb, "    stack: |-\n%s\n", errorStack)
	}
	sb.WriteString("  ...\n\n")
	f.out.WriteString(sb.String())
	return ""
}

func (f *ProgressBarFormatter) Comment(message string) string {
	return fmt.Sprintf("# %s\n", message)
}

func (f *ProgressBarFormatter) End(count, passed, failed, skipped int) string {
	var sb strings.Builder
	if f.show {
		fmt.Fprintf(os.Stderr, "\r\033[2K")
	}
	sb.WriteString(f.out.String())
	if f.show {
		sb.WriteString("\n")
	}
	fmt.Fprintf(&sb, "1..%d\n", count)
	fmt.Fprintf(&sb, "# tests %d\n", count)
	fmt.Fprintf(&sb, "# pass %d\n", passed)
	if skipped > 0 {
		fmt.Fprintf(&sb, "# %s skip %d\n", skipEmoji, skipped)
	}
	sb.WriteString("\n")
	if failed > 0 {
		fmt.Fprintf(&sb, "# %s fail %d\n", failEmoji, failed)
	} else {
		sb.WriteString(okLine())
	}
	sb.WriteString("\n")
	return sb.String()
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
