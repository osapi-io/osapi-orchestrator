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

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type RendererLipglossTestSuite struct {
	suite.Suite
}

func (s *RendererLipglossTestSuite) TestNewLipglossRenderer() {
	r := newLipglossRenderer()

	s.NotNil(r)
	s.NotNil(r.w)
}

func (s *RendererLipglossTestSuite) TestPlanStart() {
	var buf bytes.Buffer
	r := newLipglossRendererWithWriter(&buf)

	summary := sdk.PlanSummary{
		TotalTasks: 3,
		Steps: []sdk.StepSummary{
			{Tasks: []string{"health-check", "dns-update"}, Parallel: true},
			{Tasks: []string{"get-hostname"}, Parallel: false},
		},
	}
	r.PlanStart(summary)

	output := buf.String()
	s.Contains(output, "Execution Plan")
	s.NotContains(output, "===")
	s.Contains(output, "Plan: 3 tasks, 2 steps")
	s.Contains(output, "Step 1 (parallel): health-check, dns-update")
	s.Contains(output, "Step 2 (sequential): get-hostname")
}

func (s *RendererLipglossTestSuite) TestPlanDone() {
	var buf bytes.Buffer
	r := newLipglossRendererWithWriter(&buf)

	report := &Report{
		Tasks: []sdk.TaskResult{
			{Name: "a", Status: sdk.StatusChanged, Changed: true},
			{Name: "b", Status: sdk.StatusUnchanged},
		},
		Duration: 287 * time.Millisecond,
	}
	r.PlanDone(report)

	output := buf.String()
	s.Contains(output, "Complete:")
	s.NotContains(output, "===")
	s.Contains(output, "2 tasks")
	s.Contains(output, "1 changed")
	s.Contains(output, "1 unchanged")
	s.Contains(output, "287ms")
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
	var buf bytes.Buffer
	r := newLipglossRendererWithWriter(&buf)

	r.TaskStart("health-check", "target=_any")

	output := buf.String()
	s.Contains(output, "[start]")
	s.Contains(output, "health-check")
	s.Contains(output, "target=_any")
}

func (s *RendererLipglossTestSuite) TestTaskDone() {
	tests := []struct {
		name     string
		result   sdk.TaskResult
		contains []string
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
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)

			r.TaskDone(tc.result)

			output := buf.String()
			for _, c := range tc.contains {
				s.Contains(output, c)
			}
		})
	}
}

func (s *RendererLipglossTestSuite) TestTaskRetry() {
	var buf bytes.Buffer
	r := newLipglossRendererWithWriter(&buf)

	r.TaskRetry("run-command", 1, errors.New("timeout"))

	output := buf.String()
	s.Contains(output, "[retry]")
	s.Contains(output, "run-command")
	s.Contains(output, "attempt=1")
	s.Contains(output, "timeout")
}

func (s *RendererLipglossTestSuite) TestTaskSkip() {
	var buf bytes.Buffer
	r := newLipglossRendererWithWriter(&buf)

	r.TaskSkip("whoami", "dependency failed")

	output := buf.String()
	s.Contains(output, "[skip]")
	s.Contains(output, "whoami")
	s.Contains(output, "dependency failed")
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

func (s *RendererLipglossTestSuite) TestTaskDoneShowsErrorOnFailure() {
	tests := []struct {
		name     string
		result   sdk.TaskResult
		contains []string
	}{
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
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)

			r.TaskDone(tc.result)

			output := buf.String()
			for _, c := range tc.contains {
				s.Contains(output, c)
			}
		})
	}
}

func (s *RendererLipglossTestSuite) TestTaskDoneVerbose() {
	tests := []struct {
		name        string
		verbose     bool
		result      sdk.TaskResult
		contains    []string
		notContains []string
	}{
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
			name:    "Normal mode hides stdout",
			verbose: false,
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
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var buf bytes.Buffer
			r := newLipglossRendererWithWriter(&buf)
			r.verbose = tc.verbose

			r.TaskDone(tc.result)

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

func TestRendererLipglossTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(RendererLipglossTestSuite))
}
