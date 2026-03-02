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

package orchestrator_test

import (
	"testing"
	"time"

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

type ResultPublicTestSuite struct {
	suite.Suite
}

func (s *ResultPublicTestSuite) TestDecode() {
	tests := []struct {
		name         string
		results      sdk.Results
		lookupName   string
		target       any
		expectErr    bool
		errContains  string
		validateFunc func()
	}{
		{
			name: "Decodes hostname result",
			results: sdk.Results{
				"get-hostname": &sdk.Result{
					Changed: false,
					Data: map[string]any{
						"hostname": "web-01",
						"labels":   map[string]any{"env": "prod"},
					},
				},
			},
			lookupName: "get-hostname",
			target:     &orchestrator.HostnameResult{},
			validateFunc: func() {
				r := orchestrator.NewResults(sdk.Results{
					"get-hostname": &sdk.Result{
						Data: map[string]any{
							"hostname": "web-01",
							"labels":   map[string]any{"env": "prod"},
						},
					},
				})
				var h orchestrator.HostnameResult
				err := r.Decode("get-hostname", &h)
				s.Require().NoError(err)
				s.Equal("web-01", h.Hostname)
				s.Equal(map[string]string{"env": "prod"}, h.Labels)
			},
		},
		{
			name: "Decodes command result",
			results: sdk.Results{
				"run-uptime": &sdk.Result{
					Data: map[string]any{
						"stdout":      "12:34:56 up 10 days",
						"stderr":      "",
						"exit_code":   float64(0),
						"duration_ms": float64(42),
					},
				},
			},
			lookupName: "run-uptime",
			target:     &orchestrator.CommandResult{},
			validateFunc: func() {
				r := orchestrator.NewResults(sdk.Results{
					"run-uptime": &sdk.Result{
						Data: map[string]any{
							"stdout":      "12:34:56 up 10 days",
							"exit_code":   float64(0),
							"duration_ms": float64(42),
						},
					},
				})
				var cmd orchestrator.CommandResult
				err := r.Decode("run-uptime", &cmd)
				s.Require().NoError(err)
				s.Equal("12:34:56 up 10 days", cmd.Stdout)
				s.Equal(0, cmd.ExitCode)
				s.Equal(int64(42), cmd.DurationMs)
			},
		},
		{
			name: "Decodes nil data without error",
			results: sdk.Results{
				"empty-task": &sdk.Result{Data: nil},
			},
			lookupName: "empty-task",
			target:     &orchestrator.HostnameResult{},
		},
		{
			name:        "Returns error for missing result",
			results:     sdk.Results{},
			lookupName:  "nonexistent",
			target:      &orchestrator.HostnameResult{},
			expectErr:   true,
			errContains: "no result for",
		},
		{
			name: "Returns error when unmarshal fails",
			results: sdk.Results{
				"bad-task": &sdk.Result{
					Data: map[string]any{"hostname": 12345},
				},
			},
			lookupName:  "bad-task",
			target:      &struct{ Hostname chan int }{},
			expectErr:   true,
			errContains: "decode result data",
		},
		{
			name: "Returns error when marshal fails",
			results: sdk.Results{
				"bad-marshal": &sdk.Result{
					Data: map[string]any{"fn": func() {}},
				},
			},
			lookupName:  "bad-marshal",
			target:      &orchestrator.HostnameResult{},
			expectErr:   true,
			errContains: "marshal result data",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			results := orchestrator.NewResults(tc.results)
			err := results.Decode(tc.lookupName, tc.target)

			if tc.expectErr {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)

				return
			}

			s.NoError(err)

			if tc.validateFunc != nil {
				tc.validateFunc()
			}
		})
	}
}

func (s *ResultPublicTestSuite) TestReportSummary() {
	tests := []struct {
		name     string
		tasks    []sdk.TaskResult
		expected string
	}{
		{
			name:     "Empty report",
			tasks:    nil,
			expected: "0 tasks",
		},
		{
			name: "All changed",
			tasks: []sdk.TaskResult{
				{Name: "a", Status: sdk.StatusChanged},
				{Name: "b", Status: sdk.StatusChanged},
			},
			expected: "2 tasks, 2 changed",
		},
		{
			name: "All unchanged",
			tasks: []sdk.TaskResult{
				{Name: "a", Status: sdk.StatusUnchanged},
			},
			expected: "1 tasks, 1 unchanged",
		},
		{
			name: "Mixed statuses",
			tasks: []sdk.TaskResult{
				{Name: "a", Status: sdk.StatusChanged},
				{Name: "b", Status: sdk.StatusUnchanged},
				{Name: "c", Status: sdk.StatusSkipped},
				{Name: "d", Status: sdk.StatusFailed},
			},
			expected: "4 tasks, 1 changed, 1 unchanged, 1 skipped, 1 failed",
		},
		{
			name: "Skipped only",
			tasks: []sdk.TaskResult{
				{Name: "a", Status: sdk.StatusSkipped},
			},
			expected: "1 tasks, 1 skipped",
		},
		{
			name: "Failed only",
			tasks: []sdk.TaskResult{
				{Name: "a", Status: sdk.StatusFailed},
			},
			expected: "1 tasks, 1 failed",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			report := &orchestrator.Report{
				Tasks:    tc.tasks,
				Duration: 5 * time.Second,
			}

			s.Equal(tc.expected, report.Summary())
		})
	}
}

func (s *ResultPublicTestSuite) TestStatus() {
	tests := []struct {
		name       string
		results    sdk.Results
		lookupName string
		wantStatus orchestrator.TaskStatus
	}{
		{
			name: "Returns TaskStatusChanged for changed result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Changed: true,
					Status:  sdk.StatusChanged,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusChanged,
		},
		{
			name: "Returns TaskStatusUnchanged for unchanged result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Changed: false,
					Status:  sdk.StatusUnchanged,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusUnchanged,
		},
		{
			name: "Returns TaskStatusSkipped for skipped result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Status: sdk.StatusSkipped,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusSkipped,
		},
		{
			name: "Returns TaskStatusFailed for failed result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Status: sdk.StatusFailed,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusFailed,
		},
		{
			name:       "Returns TaskStatusUnknown for missing result",
			results:    sdk.Results{},
			lookupName: "nonexistent",
			wantStatus: orchestrator.TaskStatusUnknown,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			r := orchestrator.NewResults(tc.results)
			s.Equal(tc.wantStatus, r.Status(tc.lookupName))
		})
	}
}

func TestResultPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(ResultPublicTestSuite))
}
