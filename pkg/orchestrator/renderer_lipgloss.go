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
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

// tagWidth is the visible width of the longest status tag ([unchanged]).
const tagWidth = 11

// nameWidth is the column width for task names.
const nameWidth = 25

// lipglossRenderer implements renderer with colored terminal output.
type lipglossRenderer struct {
	mu      sync.Mutex
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
	summary engine.PlanSummary,
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
	mode levelMode,
) {
	style := r.dim.Bold(true)
	if mode == modeParallel {
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
	mode levelMode,
) {
	style := r.dim
	if mode == modeParallel {
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
	r.mu.Lock()
	defer r.mu.Unlock()

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
	result engine.TaskResult,
) {
	// Suppress the [skipped] line — the [skip] line from OnSkip
	// already shows the reason and is more useful.
	if result.Status == engine.StatusSkipped {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.printf(
		"  %s %-*s %s  %s\n",
		r.taskTag(result.Status),
		nameWidth,
		result.Name,
		r.changedLabel(result),
		r.dim.Render(r.durationLabel(result)),
	)

	r.printTaskError(result)

	// A broadcast operation reports per host, and those lines stand in
	// for the aggregated data printTaskDetail would otherwise write.
	if len(result.HostResults) > 0 {
		r.printHostResults(result.HostResults)
	}

	r.printTaskDetail(result)
}

// printTaskError writes the error detail under a failed task. It is
// shown whether or not the run is verbose.
func (r *lipglossRenderer) printTaskError(result engine.TaskResult) {
	if result.Status != engine.StatusFailed || result.Error == nil {
		return
	}

	r.printf(
		"  %s %s\n",
		tagIndent,
		r.red.Render(result.Error.Error()),
	)
}

// printTaskDetail writes the job id and response data a verbose run
// asks for. The aggregated data is left out of a broadcast, where the
// per-host lines have already said it.
func (r *lipglossRenderer) printTaskDetail(result engine.TaskResult) {
	if !r.verbose {
		return
	}

	if result.JobID != "" {
		r.printf(
			"%s%s\n",
			detailIndent,
			r.dim.Render(fmt.Sprintf("job_id: %s", result.JobID)),
		)
	}

	if result.Data != nil && len(result.HostResults) == 0 {
		r.printResultData(result.Data)
	}
}

// taskTag renders the bracketed status tag, red on failure.
func (r *lipglossRenderer) taskTag(status engine.Status) string {
	label := fmt.Sprintf("[%s]", status)

	style := r.green
	if status == engine.StatusFailed {
		style = r.red
	}

	return padTag(style.Render(label), len(label))
}

// changedLabel renders the changed flag, highlighted when it is set.
func (r *lipglossRenderer) changedLabel(result engine.TaskResult) string {
	label := fmt.Sprintf("changed=%v", result.Changed)
	if !result.Changed {
		return label
	}

	return r.greenB.Render(label)
}

// durationLabel renders how long a task took, with the job's own time
// alongside it in verbose mode.
func (r *lipglossRenderer) durationLabel(result engine.TaskResult) string {
	label := formatDuration(result.Duration)
	if r.verbose && result.JobDuration > 0 {
		label += fmt.Sprintf(" (job: %s)", formatDuration(result.JobDuration))
	}

	return label
}

// printHostResults renders per-host results for broadcast operations.
// Each host gets a bracketed header line with status and optional timing,
// followed by indented response data in verbose mode.
func (r *lipglossRenderer) printHostResults(
	hostResults []engine.HostResult,
) {
	for _, hr := range hostResults {
		r.printf(
			"%s%s %s%s%s\n",
			detailIndent,
			r.cyan.Render(fmt.Sprintf("[%s]", hr.Hostname)),
			r.hostStatus(hr),
			r.hostChanged(hr),
			r.hostDuration(hr),
		)

		if r.verbose && hr.Data != nil {
			r.printHostData(detailIndent+"  ", hr.Data)
		}
	}
}

// hostStatus renders a host's outcome, carrying its error message when
// it has one. A skipped host is yellow rather than red: the operation
// was unavailable there, which is not a failure.
func (r *lipglossRenderer) hostStatus(hr engine.HostResult) string {
	switch hr.Status {
	case string(engine.StatusSkipped):
		return r.yellow.Render(withDetail("skipped", hr.Error))
	case string(engine.StatusFailed):
		return r.red.Render(withDetail("failed", hr.Error))
	default:
		if hr.Error != "" {
			return r.red.Render("error: " + hr.Error)
		}

		return r.green.Render("ok")
	}
}

// withDetail appends an error to a status word when there is one.
func withDetail(status string, err string) string {
	if err == "" {
		return status
	}

	return status + ": " + err
}

// hostChanged marks a host that reported a change.
func (r *lipglossRenderer) hostChanged(hr engine.HostResult) string {
	if !hr.Changed {
		return ""
	}

	return r.greenB.Render(" changed")
}

// hostDuration reports how long the job took on a host, in verbose mode.
func (r *lipglossRenderer) hostDuration(hr engine.HostResult) string {
	if !r.verbose || hr.JobDuration <= 0 {
		return ""
	}

	return r.dim.Render(
		fmt.Sprintf(" (job: %s)", formatDuration(hr.JobDuration)),
	)
}

// printHostData renders a host's response fields, less the noisy ones.
func (r *lipglossRenderer) printHostData(
	indent string,
	data map[string]any,
) {
	for key, v := range data {
		if skipKeys[key] {
			continue
		}

		if str := formatValue(v); str != "" {
			r.printf(
				"%s%s\n",
				indent,
				r.dim.Render(fmt.Sprintf("%s: %s", key, str)),
			)
		}
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
	for key, v := range data {
		if skipKeys[key] {
			continue
		}

		str := formatValue(v)
		if str != "" {
			r.printf(
				"%s%s\n",
				detailIndent,
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
	r.mu.Lock()
	defer r.mu.Unlock()

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
	r.mu.Lock()
	defer r.mu.Unlock()

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

// The indents continuation lines are written at: tagIndent lines an
// error up under the status tag, detailIndent lines everything else up
// under the task name.
var (
	tagIndent    = strings.Repeat(" ", tagWidth)
	detailIndent = strings.Repeat(" ", tagWidth+2)
)

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
