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

	"github.com/charmbracelet/lipgloss"
	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
)

// tagWidth is the visible width of the longest status tag ([unchanged]).
const tagWidth = 11

// lipglossRenderer implements renderer with colored terminal output.
type lipglossRenderer struct {
	w       io.Writer
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
		"  %s %-20s %s\n",
		tag,
		name,
		r.dim.Render(detail),
	)
}

func (r *lipglossRenderer) TaskDone(
	result sdk.TaskResult,
) {
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
		"  %s %-20s %s  %s\n",
		tag,
		result.Name,
		changedStr,
		r.dim.Render(result.Duration.String()),
	)
}

func (r *lipglossRenderer) TaskRetry(
	name string,
	attempt int,
	err error,
) {
	tag := padTag(r.yellow.Render("[retry]"), len("[retry]"))
	r.printf(
		"  %s %-20s attempt=%d error=%q\n",
		tag,
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
		"  %s %-20s reason=%q\n",
		tag,
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
