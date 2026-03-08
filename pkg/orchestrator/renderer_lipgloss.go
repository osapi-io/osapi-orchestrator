// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package orchestrator

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
)

// tagWidth is the visible width of the longest status tag ([unchanged]).
const tagWidth = 11

// nameWidth is the column width for task names.
const nameWidth = 25

// lipglossRenderer implements renderer with colored terminal output.
type lipglossRenderer struct {
	w       io.Writer
	verbose bool
	header  lipgloss.Style
	magenta lipgloss.Style
	cyan    lipgloss.Style
	green   lipgloss.Style
	greenB  lipgloss.Style
	red     lipgloss.Style
	yellow  lipgloss.Style
	dim     lipgloss.Style
}

// newLipglossRenderer creates a lipglossRenderer writing to stdout.
func newLipglossRenderer() *lipglossRenderer {
	return newLipglossRendererWithWriter(os.Stdout)
}

// newLipglossRendererWithWriter creates a lipglossRenderer writing to w.
func newLipglossRendererWithWriter(
	w io.Writer,
) *lipglossRenderer {
	return &lipglossRenderer{
		w:       w,
		header:  lipgloss.NewStyle().Bold(true),
		magenta: lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true),
		cyan:    lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		green:   lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		greenB:  lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true),
		red:     lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true),
		yellow:  lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		dim:     lipgloss.NewStyle().Faint(true),
	}
}

func (r *lipglossRenderer) PlanStart(
	summary sdk.PlanSummary,
) {
	r.printf("\n%s\n", r.magenta.Render("Execution Plan"))
	r.printf("%s\n", r.magenta.Render(
		fmt.Sprintf(
			"Plan: %d tasks, %d steps",
			summary.TotalTasks,
			len(summary.Steps),
		),
	))

	for i, step := range summary.Steps {
		mode := "sequential"
		style := r.magenta.Faint(true)
		if step.Parallel {
			mode = "parallel"
			style = r.magenta
		}

		r.printf(
			"  %s\n",
			style.Render(
				fmt.Sprintf(
					"Step %d (%s): %s",
					i+1,
					mode,
					strings.Join(step.Tasks, ", "),
				),
			),
		)
	}
}

func (r *lipglossRenderer) PlanDone(
	report *Report,
) {
	r.printf(
		"\n%s\n",
		r.header.Render(
			fmt.Sprintf(
				"Complete: %s in %s",
				report.Summary(),
				report.Duration,
			),
		),
	)
}

func (r *lipglossRenderer) LevelStart(
	level int,
	tasks []string,
	parallel bool,
) {
	mode := "sequential"
	style := r.dim.Bold(true)
	if parallel {
		mode = "parallel"
		style = r.cyan.Bold(true)
	}

	r.printf(
		"\n%s\n",
		style.Render(
			fmt.Sprintf(
				">>> Step %d (%s): %s",
				level+1,
				mode,
				strings.Join(tasks, ", "),
			),
		),
	)
}

func (r *lipglossRenderer) LevelDone(
	level int,
	changed int,
	total int,
	parallel bool,
) {
	style := r.dim
	if parallel {
		style = r.cyan
	}

	r.printf(
		"%s\n",
		style.Render(
			fmt.Sprintf(
				"<<< Step %d done: %d/%d changed",
				level+1,
				changed,
				total,
			),
		),
	)
}

func (r *lipglossRenderer) TaskStart(
	name string,
	detail string,
) {
	tag := padTag(r.dim.Render("[start]"), len("[start]"))
	r.printf(
		"  %s %-*s %s\n",
		tag,
		nameWidth,
		name,
		r.dim.Render(detail),
	)
}

func (r *lipglossRenderer) TaskDone(
	result sdk.TaskResult,
) {
	// Suppress the [skipped] line — the [skip] line from OnSkip
	// already shows the reason and is more useful.
	if result.Status == sdk.StatusSkipped {
		return
	}

	label := fmt.Sprintf("[%s]", result.Status)

	var tag string
	if result.Status == sdk.StatusFailed {
		tag = padTag(r.red.Render(label), len(label))
	} else {
		tag = padTag(r.green.Render(label), len(label))
	}

	changedStr := fmt.Sprintf("changed=%v", result.Changed)
	if result.Changed {
		changedStr = r.greenB.Render(changedStr)
	}

	r.printf(
		"  %s %-*s %s  %s\n",
		tag,
		nameWidth,
		result.Name,
		changedStr,
		r.dim.Render(formatDuration(result.Duration)),
	)

	// Always show error detail on failure.
	if result.Status == sdk.StatusFailed && result.Error != nil {
		r.printf(
			"  %s %s\n",
			strings.Repeat(" ", tagWidth),
			r.red.Render(result.Error.Error()),
		)
	}

	// Always show per-host results for broadcast operations.
	if len(result.HostResults) > 0 {
		r.printHostResults(result.HostResults)
	}

	// Verbose mode: show response data on success.
	if r.verbose && result.Data != nil {
		r.printResultData(result.Data)
	}
}

// printHostResults renders per-host results for broadcast operations.
func (r *lipglossRenderer) printHostResults(
	hostResults []sdk.HostResult,
) {
	indent := strings.Repeat(" ", tagWidth+2)

	for _, hr := range hostResults {
		status := r.green.Render("ok")
		if hr.Error != "" {
			status = r.red.Render("error: " + hr.Error)
		}

		changed := ""
		if hr.Changed {
			changed = r.greenB.Render(" changed")
		}

		r.printf(
			"%s%s %s%s\n",
			indent,
			r.dim.Render(hr.Hostname),
			status,
			changed,
		)
	}
}

// skipKeys are internal fields that clutter verbose output.
var skipKeys = map[string]bool{
	"duration_ms": true,
	"exit_code":   true,
	"stderr":      true,
}

// printResultData renders result data fields as indented lines.
func (r *lipglossRenderer) printResultData(
	data map[string]any,
) {
	indent := strings.Repeat(" ", tagWidth+2)

	for key, v := range data {
		if skipKeys[key] {
			continue
		}

		str := formatValue(v)
		if str != "" {
			r.printf(
				"%s%s\n",
				indent,
				r.dim.Render(fmt.Sprintf("%s: %s", key, str)),
			)
		}
	}
}

// formatValue renders a value for display, keeping simple values inline
// and omitting complex nested structures.
func formatValue(
	v any,
) string {
	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val)
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}

		return fmt.Sprintf("%.2f", val)
	case bool:
		return fmt.Sprintf("%v", val)
	case []any:
		return fmt.Sprintf("[%d items]", len(val))
	case map[string]any:
		parts := make([]string, 0, len(val))
		for k, inner := range val {
			parts = append(parts, fmt.Sprintf("%s=%v", k, inner))
		}

		return strings.Join(parts, " ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (r *lipglossRenderer) TaskRetry(
	name string,
	attempt int,
	err error,
) {
	tag := padTag(r.yellow.Render("[retry]"), len("[retry]"))
	r.printf(
		"  %s %-*s attempt=%d error=%q\n",
		tag,
		nameWidth,
		name,
		attempt,
		err,
	)
}

func (r *lipglossRenderer) TaskSkip(
	name string,
	reason string,
) {
	dimYellow := r.yellow.Faint(true)
	tag := padTag(dimYellow.Render("[skip]"), len("[skip]"))
	r.printf(
		"  %s %-*s reason=%q\n",
		tag,
		nameWidth,
		name,
		reason,
	)
}

// printf writes formatted output to the renderer's writer.
// Write errors are intentionally discarded — there is no meaningful recovery
// for a broken terminal.
func (r *lipglossRenderer) printf(
	format string,
	a ...any,
) {
	_, _ = fmt.Fprintf(r.w, format, a...)
}

// formatDuration rounds a duration to millisecond precision for cleaner output.
func formatDuration(
	d time.Duration,
) string {
	return d.Round(time.Millisecond).String()
}

// padTag right-pads a styled tag string so the visible width equals tagWidth.
func padTag(
	styled string,
	visibleLen int,
) string {
	if pad := tagWidth - visibleLen; pad > 0 {
		return styled + strings.Repeat(" ", pad)
	}

	return styled
}
