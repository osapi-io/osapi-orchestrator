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

	engine "github.com/osapi-io/osapi-orchestrator/internal/engine"
	"github.com/stretchr/testify/suite"
)

type RendererLipglossTestSuite struct {
	suite.Suite
}

func (s *RendererLipglossTestSuite) TestNewLipglossRenderer() {
	tests := []struct {
		name         string
		validateFunc func(r *lipglossRenderer)
	}{
		{
			name: "Creates valid renderer with non-nil writer",
			validateFunc: func(r *lipglossRenderer) {
				s.NotNil(r)
				s.NotNil(r.w)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			r := newLipglossRenderer()
			tc.validateFunc(r)
		})
	}
}

func (s *RendererLipglossTestSuite) TestPlanStart() {
	tests := []struct {
		name        string
		setupFunc   func(r *lipglossRenderer)
		contains    []string
		notContains []string
	}{
		{
			name: "Renders execution plan",
			setupFunc: func(r *lipglossRenderer) {
				r.PlanStart(engine.PlanSummary{
					TotalTasks: 3,
					Steps: []engine.StepSummary{
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

func (s *RendererLipglossTestSuite) TestPlanDone() {
	tests := []struct {
		name        string
		setupFunc   func(r *lipglossRenderer)
		contains    []string
		notContains []string
	}{
		{
			name: "Renders completion summary",
			setupFunc: func(r *lipglossRenderer) {
				r.PlanDone(&Report{
					Tasks: []engine.TaskResult{
						{Name: "a", Status: engine.StatusChanged, Changed: true},
						{Name: "b", Status: engine.StatusUnchanged},
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

func (s *RendererLipglossTestSuite) TestTaskStart() {
	tests := []struct {
		name     string
		callFunc func(r *lipglossRenderer)
		contains []string
	}{
		{
			name: "Renders start tag and details",
			callFunc: func(r *lipglossRenderer) {
				r.TaskStart("health-check", "target=_any")
			},
			contains: []string{"[start]", "health-check", "target=_any"},
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

func (s *RendererLipglossTestSuite) TestTaskRetry() {
	tests := []struct {
		name     string
		callFunc func(r *lipglossRenderer)
		contains []string
	}{
		{
			name: "Renders retry tag with attempt",
			callFunc: func(r *lipglossRenderer) {
				r.TaskRetry("run-command", 1, errors.New("timeout"))
			},
			contains: []string{"[retry]", "run-command", "attempt=1", "timeout"},
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

func (s *RendererLipglossTestSuite) TestTaskSkip() {
	tests := []struct {
		name     string
		callFunc func(r *lipglossRenderer)
		contains []string
	}{
		{
			name: "Renders skip tag with reason",
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
		result      engine.TaskResult
		contains    []string
		notContains []string
		expectEmpty bool
	}{
		{
			name: "Unchanged task",
			result: engine.TaskResult{
				Name:     "health-check",
				Status:   engine.StatusUnchanged,
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
			result: engine.TaskResult{
				Name:     "run-command",
				Status:   engine.StatusChanged,
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
			result: engine.TaskResult{
				Name:     "bad-task",
				Status:   engine.StatusFailed,
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
			result: engine.TaskResult{
				Name:     "skipped-task",
				Status:   engine.StatusSkipped,
				Duration: 1 * time.Millisecond,
			},
			expectEmpty: true,
		},
		{
			name: "Shows error message on failure",
			result: engine.TaskResult{
				Name:     "deploy",
				Status:   engine.StatusFailed,
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
			result: engine.TaskResult{
				Name:     "deploy",
				Status:   engine.StatusFailed,
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
			result: engine.TaskResult{
				Name:     "run-cmd",
				Status:   engine.StatusChanged,
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
			result: engine.TaskResult{
				Name:     "run-cmd",
				Status:   engine.StatusChanged,
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
			result: engine.TaskResult{
				Name:     "run-cmd",
				Status:   engine.StatusChanged,
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
			name: "Shows per-host results with bracketed hostnames",
			result: engine.TaskResult{
				Name:    "deploy-all",
				Status:  engine.StatusChanged,
				Changed: true,
				HostResults: []engine.HostResult{
					{Hostname: "web-01", Changed: true},
					{Hostname: "web-02", Changed: false, Error: "timeout"},
				},
			},
			contains: []string{
				"[web-01]",
				"[web-02]",
				"error: timeout",
			},
		},
		{
			name: "Host with status ok renders as ok",
			result: engine.TaskResult{
				Name:   "broadcast-ok",
				Status: engine.StatusUnchanged,
				HostResults: []engine.HostResult{
					{Hostname: "web-01", Status: "ok"},
				},
			},
			contains: []string{
				"[web-01]",
				"ok",
			},
			notContains: []string{
				"skipped",
				"failed",
				"error:",
			},
		},
		{
			name: "Host with status skipped and error renders with skipped prefix",
			result: engine.TaskResult{
				Name:   "broadcast-skipped",
				Status: engine.StatusUnchanged,
				HostResults: []engine.HostResult{
					{Hostname: "darwin-01", Status: "skipped", Error: "unsupported"},
				},
			},
			contains: []string{
				"[darwin-01]",
				"skipped: unsupported",
			},
			notContains: []string{
				"error: unsupported",
			},
		},
		{
			name: "Host with status skipped and no error renders as skipped",
			result: engine.TaskResult{
				Name:   "broadcast-skipped-no-err",
				Status: engine.StatusUnchanged,
				HostResults: []engine.HostResult{
					{Hostname: "darwin-02", Status: "skipped"},
				},
			},
			contains: []string{
				"[darwin-02]",
				"skipped",
			},
			notContains: []string{
				"error:",
			},
		},
		{
			name: "Host with status failed and error renders with failed prefix",
			result: engine.TaskResult{
				Name:   "broadcast-failed",
				Status: engine.StatusFailed,
				HostResults: []engine.HostResult{
					{Hostname: "db-01", Status: "failed", Error: "permission denied"},
				},
			},
			contains: []string{
				"[db-01]",
				"failed: permission denied",
			},
			notContains: []string{
				"error: permission denied",
			},
		},
		{
			name: "Host with empty status but error falls back to error prefix",
			result: engine.TaskResult{
				Name:   "broadcast-legacy",
				Status: engine.StatusFailed,
				HostResults: []engine.HostResult{
					{Hostname: "legacy-01", Error: "connection refused"},
				},
			},
			contains: []string{
				"[legacy-01]",
				"error: connection refused",
			},
			notContains: []string{
				"skipped",
				"failed:",
			},
		},
		{
			name:    "Verbose shows per-host data and suppresses inline data",
			verbose: true,
			result: engine.TaskResult{
				Name:     "get-hostname-all",
				Status:   engine.StatusUnchanged,
				Duration: 1500 * time.Millisecond,
				Data: map[string]any{
					"hostname": "nerd",
				},
				HostResults: []engine.HostResult{
					{
						Hostname:    "nerd",
						JobDuration: 2 * time.Millisecond,
						Data: map[string]any{
							"hostname": "nerd",
						},
					},
				},
			},
			contains: []string{
				"[nerd]",
				"(job: 2ms)",
			},
			notContains: []string{
				// Inline data is suppressed when per-host results
				// are present — the per-host section already shows it.
				"             hostname: nerd\n             hostname: nerd",
			},
		},
		{
			name: "Normal mode hides per-host data",
			result: engine.TaskResult{
				Name:   "get-hostname-all",
				Status: engine.StatusUnchanged,
				HostResults: []engine.HostResult{
					{
						Hostname: "nerd",
						Data: map[string]any{
							"hostname": "nerd",
						},
					},
				},
			},
			contains: []string{
				"[nerd]",
			},
			notContains: []string{
				"hostname: nerd",
			},
		},
		{
			name: "No host results for non-broadcast",
			result: engine.TaskResult{
				Name:   "get-hostname",
				Status: engine.StatusUnchanged,
			},
			notContains: []string{
				"web-",
			},
		},
		{
			name:    "Verbose shows job ID when present",
			verbose: true,
			result: engine.TaskResult{
				JobID:    "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				Name:     "set-dns",
				Status:   engine.StatusChanged,
				Changed:  true,
				Duration: 50 * time.Millisecond,
			},
			contains: []string{
				"job_id: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			},
		},
		{
			name: "Normal mode hides job ID",
			result: engine.TaskResult{
				JobID:    "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				Name:     "set-dns",
				Status:   engine.StatusChanged,
				Changed:  true,
				Duration: 50 * time.Millisecond,
			},
			notContains: []string{
				"job_id:",
			},
		},
		{
			name:    "Verbose shows job duration when present",
			verbose: true,
			result: engine.TaskResult{
				Name:        "set-dns",
				Status:      engine.StatusChanged,
				Changed:     true,
				Duration:    50 * time.Millisecond,
				JobDuration: 30 * time.Millisecond,
			},
			contains: []string{
				"(job: 30ms)",
			},
		},
		{
			name:    "Verbose shows host data but skips internal keys",
			verbose: true,
			result: engine.TaskResult{
				Name:     "broadcast-cmd",
				Status:   engine.StatusChanged,
				Changed:  true,
				Duration: 100 * time.Millisecond,
				HostResults: []engine.HostResult{
					{
						Hostname: "web-01",
						Changed:  true,
						Data: map[string]any{
							"stdout":    "deployed",
							"exit_code": float64(0),
						},
					},
				},
			},
			contains: []string{
				"[web-01]",
				"stdout: deployed",
			},
			notContains: []string{
				"exit_code:",
			},
		},
		{
			name:    "Verbose hides job ID when empty",
			verbose: true,
			result: engine.TaskResult{
				Name:     "health-check",
				Status:   engine.StatusUnchanged,
				Duration: 10 * time.Millisecond,
			},
			notContains: []string{
				"job_id:",
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
