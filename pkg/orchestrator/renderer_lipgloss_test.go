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
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
	"github.com/stretchr/testify/suite"
)

type RendererLipglossTestSuite struct {
	suite.Suite
}

func (s *RendererLipglossTestSuite) TestPlanLifecycle() {
	tests := []struct {
		name        string
		setupFunc   func(r *lipglossRenderer)
		contains    []string
		notContains []string
	}{
		{
			name: "NewLipglossRenderer creates valid renderer",
			setupFunc: func(_ *lipglossRenderer) {
				r := newLipglossRenderer()
				s.NotNil(r)
				s.NotNil(r.w)
			},
		},
		{
			name: "PlanStart renders execution plan",
			setupFunc: func(r *lipglossRenderer) {
				r.PlanStart(sdk.PlanSummary{
					TotalTasks: 3,
					Steps: []sdk.StepSummary{
						{Tasks: []string{"health-check", "dns-update"}, Parallel: true},
						{Tasks: []string{"get-hostname"}, Parallel: false},
					},
				})
			},
			contains: []string{
				"Execution Plan",
				"Plan: 3 tasks, 2 steps",
				"Step 1 (parallel): health-check, dns-update",
				"Step 2 (sequential): get-hostname",
			},
			notContains: []string{"==="},
		},
		{
			name: "PlanDone renders completion summary",
			setupFunc: func(r *lipglossRenderer) {
				r.PlanDone(&Report{
					Tasks: []sdk.TaskResult{
						{Name: "a", Status: sdk.StatusChanged, Changed: true},
						{Name: "b", Status: sdk.StatusUnchanged},
					},
					Duration: 287 * time.Millisecond,
				})
			},
			contains: []string{
				"Complete:",
				"2 tasks",
				"1 changed",
				"1 unchanged",
				"287ms",
			},
			notContains: []string{"==="},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)

			tc.setupFunc(r)

			output := buf.String()
			for _, c := range tc.contains {
				s.Contains(output, c)
			}
			for _, c := range tc.notContains {
				s.NotContains(output, c)
			}
		})
	}
}

func (s *RendererLipglossTestSuite) TestLevelStart() {
	tests := []struct {
		name     string
		level    int
		tasks    []string
		parallel bool
		contains string
	}{
		{
			name:     "Sequential level",
			level:    0,
			tasks:    []string{"health-check"},
			parallel: false,
			contains: "sequential",
		},
		{
			name:     "Parallel level",
			level:    1,
			tasks:    []string{"get-hostname", "get-disk"},
			parallel: true,
			contains: "parallel",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)

			r.LevelStart(tc.level, tc.tasks, tc.parallel)

			output := buf.String()
			s.Contains(output, ">>> Step")
			s.Contains(output, tc.contains)
			for _, task := range tc.tasks {
				s.Contains(output, task)
			}
		})
	}
}

func (s *RendererLipglossTestSuite) TestLevelDone() {
	tests := []struct {
		name     string
		level    int
		changed  int
		total    int
		parallel bool
	}{
		{
			name:     "No changes sequential",
			level:    0,
			changed:  0,
			total:    2,
			parallel: false,
		},
		{
			name:     "Some changes parallel",
			level:    1,
			changed:  1,
			total:    3,
			parallel: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)

			r.LevelDone(tc.level, tc.changed, tc.total, tc.parallel)

			output := buf.String()
			s.Contains(output, "<<< Step")
			s.Contains(output, "changed")
		})
	}
}

func (s *RendererLipglossTestSuite) TestTaskEvents() {
	tests := []struct {
		name     string
		callFunc func(r *lipglossRenderer)
		contains []string
	}{
		{
			name: "TaskStart renders start tag and details",
			callFunc: func(r *lipglossRenderer) {
				r.TaskStart("health-check", "target=_any")
			},
			contains: []string{"[start]", "health-check", "target=_any"},
		},
		{
			name: "TaskRetry renders retry tag with attempt",
			callFunc: func(r *lipglossRenderer) {
				r.TaskRetry("run-command", 1, errors.New("timeout"))
			},
			contains: []string{"[retry]", "run-command", "attempt=1", "timeout"},
		},
		{
			name: "TaskSkip renders skip tag with reason",
			callFunc: func(r *lipglossRenderer) {
				r.TaskSkip("whoami", "dependency failed")
			},
			contains: []string{"[skip]", "whoami", "dependency failed"},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)

			tc.callFunc(r)

			output := buf.String()
			for _, c := range tc.contains {
				s.Contains(output, c)
			}
		})
	}
}

func (s *RendererLipglossTestSuite) TestTaskDone() {
	tests := []struct {
		name        string
		verbose     bool
		result      sdk.TaskResult
		contains    []string
		notContains []string
		expectEmpty bool
	}{
		{
			name: "Unchanged task",
			result: sdk.TaskResult{
				Name:     "health-check",
				Status:   sdk.StatusUnchanged,
				Changed:  false,
				Duration: 12 * time.Millisecond,
			},
			contains: []string{
				"health-check",
				"changed=false",
				"12ms",
			},
		},
		{
			name: "Changed task",
			result: sdk.TaskResult{
				Name:     "run-command",
				Status:   sdk.StatusChanged,
				Changed:  true,
				Duration: 230 * time.Millisecond,
			},
			contains: []string{
				"run-command",
				"changed=true",
				"230ms",
			},
		},
		{
			name: "Failed task",
			result: sdk.TaskResult{
				Name:     "bad-task",
				Status:   sdk.StatusFailed,
				Changed:  false,
				Duration: 5 * time.Millisecond,
			},
			contains: []string{
				"bad-task",
				"failed",
			},
		},
		{
			name: "Skipped task suppressed",
			result: sdk.TaskResult{
				Name:     "skipped-task",
				Status:   sdk.StatusSkipped,
				Duration: 1 * time.Millisecond,
			},
			expectEmpty: true,
		},
		{
			name: "Shows error message on failure",
			result: sdk.TaskResult{
				Name:     "deploy",
				Status:   sdk.StatusFailed,
				Duration: 45 * time.Millisecond,
				Error:    fmt.Errorf("command exited with code 1"),
			},
			contains: []string{
				"[failed]",
				"deploy",
				"command exited with code 1",
			},
		},
		{
			name: "No error line when error is nil",
			result: sdk.TaskResult{
				Name:     "deploy",
				Status:   sdk.StatusFailed,
				Duration: 45 * time.Millisecond,
			},
			contains: []string{
				"[failed]",
				"deploy",
			},
		},
		{
			name:    "Verbose shows stdout",
			verbose: true,
			result: sdk.TaskResult{
				Name:     "run-cmd",
				Status:   sdk.StatusChanged,
				Changed:  true,
				Duration: 100 * time.Millisecond,
				Data: map[string]any{
					"stdout":    "hello world",
					"exit_code": float64(0),
				},
			},
			contains: []string{
				"[changed]",
				"stdout: hello world",
			},
		},
		{
			name: "Normal mode hides stdout",
			result: sdk.TaskResult{
				Name:     "run-cmd",
				Status:   sdk.StatusChanged,
				Changed:  true,
				Duration: 100 * time.Millisecond,
				Data: map[string]any{
					"stdout":    "hello world",
					"exit_code": float64(0),
				},
			},
			contains: []string{
				"[changed]",
			},
			notContains: []string{
				"stdout:",
			},
		},
		{
			name:    "Verbose skips empty values",
			verbose: true,
			result: sdk.TaskResult{
				Name:     "run-cmd",
				Status:   sdk.StatusChanged,
				Changed:  true,
				Duration: 100 * time.Millisecond,
				Data: map[string]any{
					"stdout": "output",
					"stderr": "",
				},
			},
			contains: []string{
				"stdout: output",
			},
			notContains: []string{
				"stderr:",
			},
		},
		{
			name: "Shows per-host results",
			result: sdk.TaskResult{
				Name:    "deploy-all",
				Status:  sdk.StatusChanged,
				Changed: true,
				HostResults: []sdk.HostResult{
					{Hostname: "web-01", Changed: true},
					{Hostname: "web-02", Changed: false, Error: "timeout"},
				},
			},
			contains: []string{
				"web-01",
				"web-02",
				"error: timeout",
			},
		},
		{
			name: "No host results for non-broadcast",
			result: sdk.TaskResult{
				Name:   "get-hostname",
				Status: sdk.StatusUnchanged,
			},
			notContains: []string{
				"web-",
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)
			r.verbose = tc.verbose

			r.TaskDone(tc.result)

			output := buf.String()

			if tc.expectEmpty {
				s.Empty(output)
				return
			}

			for _, c := range tc.contains {
				s.Contains(output, c)
			}
			for _, c := range tc.notContains {
				s.NotContains(output, c)
			}
		})
	}
}

func (s *RendererLipglossTestSuite) TestFormatValue() {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "String value trimmed",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "Empty string returns empty",
			input:    "",
			expected: "",
		},
		{
			name:     "Float64 whole number as integer",
			input:    float64(42),
			expected: "42",
		},
		{
			name:     "Float64 decimal to two places",
			input:    float64(3.14159),
			expected: "3.14",
		},
		{
			name:     "Bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "Bool false",
			input:    false,
			expected: "false",
		},
		{
			name:     "Slice shows item count",
			input:    []any{"a", "b", "c"},
			expected: "[3 items]",
		},
		{
			name:     "Empty slice shows zero items",
			input:    []any{},
			expected: "[0 items]",
		},
		{
			name:     "Map shows key=value pairs",
			input:    map[string]any{"env": "prod"},
			expected: "env=prod",
		},
		{
			name:     "Default type uses Sprintf",
			input:    42,
			expected: "42",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := formatValue(tc.input)
			s.Equal(tc.expected, got)
		})
	}
}

func (s *RendererLipglossTestSuite) TestFormatDuration() {
	tests := []struct {
		name     string
		d        time.Duration
		expected string
	}{
		{
			name:     "Pure milliseconds",
			d:        942191541 * time.Nanosecond,
			expected: "942ms",
		},
		{
			name:     "Seconds with fractional ms",
			d:        1035578833 * time.Nanosecond,
			expected: "1.036s",
		},
		{
			name:     "Clean milliseconds unchanged",
			d:        230 * time.Millisecond,
			expected: "230ms",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := formatDuration(tc.d)
			s.Equal(tc.expected, got)
		})
	}
}

func (s *RendererLipglossTestSuite) TestPadTag() {
	tests := []struct {
		name       string
		styled     string
		visibleLen int
		wantSuffix string
	}{
		{
			name:       "Short tag gets padded",
			styled:     "[ok]",
			visibleLen: 4,
			wantSuffix: strings.Repeat(" ", tagWidth-4),
		},
		{
			name:       "Exact width tag not padded",
			styled:     "[unchanged]",
			visibleLen: tagWidth,
			wantSuffix: "",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := padTag(tc.styled, tc.visibleLen)
			s.True(
				strings.HasPrefix(got, tc.styled),
				"should start with styled text",
			)
			suffix := got[len(tc.styled):]
			s.Equal(tc.wantSuffix, suffix)
		})
	}
}

func TestRendererLipglossTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(RendererLipglossTestSuite))
}
