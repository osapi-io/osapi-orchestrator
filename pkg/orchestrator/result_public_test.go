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

	engine "github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

type ResultPublicTestSuite struct {
	suite.Suite
}

func (s *ResultPublicTestSuite) TestDecode() {
	tests := []struct {
		name         string
		results      engine.Results
		lookupName   string
		target       any
		expectErr    bool
		errContains  string
		validateFunc func()
	}{
		{
			name: "Decodes hostname from host results",
			results: engine.Results{
				"get-hostname": &engine.Result{
					Changed: false,
					Data:    map[string]any{"results": []any{}},
					HostResults: []engine.HostResult{
						{
							Hostname: "web-01",
							Data: map[string]any{
								"hostname": "web-01",
								"labels":   map[string]any{"env": "prod"},
							},
						},
					},
				},
			},
			lookupName: "get-hostname",
			target:     &osapi.HostnameResult{},
			validateFunc: func() {
				r := orchestrator.NewResults(engine.Results{
					"get-hostname": &engine.Result{
						Data: map[string]any{"results": []any{}},
						HostResults: []engine.HostResult{
							{
								Hostname: "web-01",
								Data: map[string]any{
									"hostname": "web-01",
									"labels":   map[string]any{"env": "prod"},
								},
							},
						},
					},
				})
				var h osapi.HostnameResult
				err := r.Decode("get-hostname", &h)
				s.Require().NoError(err)
				s.Equal("web-01", h.Hostname)
				s.Equal(map[string]string{"env": "prod"}, h.Labels)
			},
		},
		{
			name: "Decodes command from host results",
			results: engine.Results{
				"run-uptime": &engine.Result{
					Data: map[string]any{"results": []any{}},
					HostResults: []engine.HostResult{
						{
							Hostname: "web-01",
							Data: map[string]any{
								"stdout":      "12:34:56 up 10 days",
								"exit_code":   float64(0),
								"duration_ms": float64(42),
							},
						},
					},
				},
			},
			lookupName: "run-uptime",
			target:     &osapi.CommandResult{},
			validateFunc: func() {
				r := orchestrator.NewResults(engine.Results{
					"run-uptime": &engine.Result{
						Data: map[string]any{"results": []any{}},
						HostResults: []engine.HostResult{
							{
								Hostname: "web-01",
								Data: map[string]any{
									"stdout":      "12:34:56 up 10 days",
									"exit_code":   float64(0),
									"duration_ms": float64(42),
								},
							},
						},
					},
				})
				var cmd osapi.CommandResult
				err := r.Decode("run-uptime", &cmd)
				s.Require().NoError(err)
				s.Equal("12:34:56 up 10 days", cmd.Stdout)
				s.Equal(0, cmd.ExitCode)
				s.Equal(int64(42), cmd.DurationMs)
			},
		},
		{
			name: "Falls back to Data when no host results",
			results: engine.Results{
				"summarize": &engine.Result{
					Changed: true,
					Data:    map[string]any{"host": "web-01"},
				},
			},
			lookupName: "summarize",
			target:     &map[string]any{},
			validateFunc: func() {
				r := orchestrator.NewResults(engine.Results{
					"summarize": &engine.Result{
						Changed: true,
						Data:    map[string]any{"host": "web-01"},
					},
				})
				var m map[string]any
				err := r.Decode("summarize", &m)
				s.Require().NoError(err)
				s.Equal("web-01", m["host"])
			},
		},
		{
			name: "Decodes nil data without error",
			results: engine.Results{
				"empty-task": &engine.Result{Data: nil},
			},
			lookupName: "empty-task",
			target:     &osapi.HostnameResult{},
		},
		{
			name:        "Returns error for missing result",
			results:     engine.Results{},
			lookupName:  "nonexistent",
			target:      &osapi.HostnameResult{},
			expectErr:   true,
			errContains: "no result for",
		},
		{
			name: "Decodes command error from host results",
			results: engine.Results{
				"run-cmd": &engine.Result{
					Changed: true,
					Data:    map[string]any{"results": []any{}},
					HostResults: []engine.HostResult{
						{
							Hostname: "web-01",
							Error:    "exec failed",
							Data: map[string]any{
								"stdout":      "partial output",
								"stderr":      "command not found",
								"exit_code":   float64(127),
								"duration_ms": float64(50),
								"error":       "exec failed",
							},
						},
					},
				},
			},
			lookupName: "run-cmd",
			target:     &osapi.CommandResult{},
			validateFunc: func() {
				r := orchestrator.NewResults(engine.Results{
					"run-cmd": &engine.Result{
						Changed: true,
						Data:    map[string]any{"results": []any{}},
						HostResults: []engine.HostResult{
							{
								Hostname: "web-01",
								Error:    "exec failed",
								Data: map[string]any{
									"stdout":      "partial output",
									"stderr":      "command not found",
									"exit_code":   float64(127),
									"duration_ms": float64(50),
									"error":       "exec failed",
								},
							},
						},
					},
				})
				var cmd osapi.CommandResult
				err := r.Decode("run-cmd", &cmd)
				s.Require().NoError(err)
				s.Equal("exec failed", cmd.Error)
				s.Equal(127, cmd.ExitCode)
				s.Equal("command not found", cmd.Stderr)
			},
		},
		{
			name: "Returns error when unmarshal fails",
			results: engine.Results{
				"bad-task": &engine.Result{
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
			results: engine.Results{
				"bad-marshal": &engine.Result{
					Data: map[string]any{"fn": func() {}},
				},
			},
			lookupName:  "bad-marshal",
			target:      &osapi.HostnameResult{},
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
		tasks    []engine.TaskResult
		expected string
	}{
		{
			name:     "Empty report",
			tasks:    nil,
			expected: "0 tasks",
		},
		{
			name: "All changed",
			tasks: []engine.TaskResult{
				{Name: "a", Status: engine.StatusChanged},
				{Name: "b", Status: engine.StatusChanged},
			},
			expected: "2 tasks, 2 changed",
		},
		{
			name: "All unchanged",
			tasks: []engine.TaskResult{
				{Name: "a", Status: engine.StatusUnchanged},
			},
			expected: "1 tasks, 1 unchanged",
		},
		{
			name: "Mixed statuses",
			tasks: []engine.TaskResult{
				{Name: "a", Status: engine.StatusChanged},
				{Name: "b", Status: engine.StatusUnchanged},
				{Name: "c", Status: engine.StatusSkipped},
				{Name: "d", Status: engine.StatusFailed},
			},
			expected: "4 tasks, 1 changed, 1 unchanged, 1 skipped, 1 failed",
		},
		{
			name: "Skipped only",
			tasks: []engine.TaskResult{
				{Name: "a", Status: engine.StatusSkipped},
			},
			expected: "1 tasks, 1 skipped",
		},
		{
			name: "Failed only",
			tasks: []engine.TaskResult{
				{Name: "a", Status: engine.StatusFailed},
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
		results    engine.Results
		lookupName string
		wantStatus orchestrator.TaskStatus
	}{
		{
			name: "Returns TaskStatusChanged for changed result",
			results: engine.Results{
				"step-a": &engine.Result{
					Changed: true,
					Status:  engine.StatusChanged,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusChanged,
		},
		{
			name: "Returns TaskStatusUnchanged for unchanged result",
			results: engine.Results{
				"step-a": &engine.Result{
					Changed: false,
					Status:  engine.StatusUnchanged,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusUnchanged,
		},
		{
			name: "Returns TaskStatusSkipped for skipped result",
			results: engine.Results{
				"step-a": &engine.Result{
					Status: engine.StatusSkipped,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusSkipped,
		},
		{
			name: "Returns TaskStatusFailed for failed result",
			results: engine.Results{
				"step-a": &engine.Result{
					Status: engine.StatusFailed,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusFailed,
		},
		{
			name:       "Returns TaskStatusUnknown for missing result",
			results:    engine.Results{},
			lookupName: "nonexistent",
			wantStatus: orchestrator.TaskStatusUnknown,
		},
		{
			name: "Returns TaskStatusUnknown for unrecognized SDK status",
			results: engine.Results{
				"step-a": &engine.Result{
					Status: engine.Status("unknown-status"),
				},
			},
			lookupName: "step-a",
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

func (s *ResultPublicTestSuite) TestChanged() {
	tests := []struct {
		name       string
		results    engine.Results
		lookupName string
		want       bool
	}{
		{
			name: "Returns true when step reported changes",
			results: engine.Results{
				"deploy": &engine.Result{Changed: true, Status: engine.StatusChanged},
			},
			lookupName: "deploy",
			want:       true,
		},
		{
			name: "Returns false when step did not report changes",
			results: engine.Results{
				"deploy": &engine.Result{Changed: false, Status: engine.StatusUnchanged},
			},
			lookupName: "deploy",
			want:       false,
		},
		{
			name:       "Returns false when step not found",
			results:    engine.Results{},
			lookupName: "nonexistent",
			want:       false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			r := orchestrator.NewResults(tc.results)
			s.Equal(tc.want, r.Changed(tc.lookupName))
		})
	}
}

func (s *ResultPublicTestSuite) TestHostResults() {
	tests := []struct {
		name       string
		results    engine.Results
		lookupName string
		wantNil    bool
		wantLen    int
		validateFn func(hrs []orchestrator.HostResult)
	}{
		{
			name: "Returns per-host results",
			results: engine.Results{
				"deploy": &engine.Result{
					Changed: true,
					Status:  engine.StatusChanged,
					HostResults: []engine.HostResult{
						{
							Hostname: "web-01",
							Changed:  true,
							Data: map[string]any{
								"stdout": "deployed",
							},
						},
						{
							Hostname: "web-02",
							Changed:  false,
							Error:    "timeout",
						},
					},
				},
			},
			lookupName: "deploy",
			wantLen:    2,
			validateFn: func(hrs []orchestrator.HostResult) {
				s.Equal("web-01", hrs[0].Hostname)
				s.True(hrs[0].Changed)
				s.Empty(hrs[0].Error)
				s.Equal("web-02", hrs[1].Hostname)
				s.Equal("timeout", hrs[1].Error)
			},
		},
		{
			name:       "Returns nil for missing step",
			results:    engine.Results{},
			lookupName: "nonexistent",
			wantNil:    true,
		},
		{
			name: "Returns nil for unicast result",
			results: engine.Results{
				"get-host": &engine.Result{
					Changed: false,
					Status:  engine.StatusUnchanged,
				},
			},
			lookupName: "get-host",
			wantNil:    true,
		},
		{
			name: "with status fields",
			results: engine.Results{
				"deploy": &engine.Result{
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "ok", Changed: true},
						{Hostname: "web-02", Status: "skipped", Error: "unsupported"},
						{Hostname: "web-03", Status: "failed", Error: "permission denied"},
					},
				},
			},
			lookupName: "deploy",
			wantLen:    3,
			validateFn: func(hrs []orchestrator.HostResult) {
				s.Require().Len(hrs, 3)
				s.Equal("ok", hrs[0].Status)
				s.True(hrs[0].Changed)
				s.Equal("skipped", hrs[1].Status)
				s.Equal("unsupported", hrs[1].Error)
				s.Equal("failed", hrs[2].Status)
				s.Equal("permission denied", hrs[2].Error)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			r := orchestrator.NewResults(tc.results)
			hrs := r.HostResults(tc.lookupName)

			if tc.wantNil {
				s.Nil(hrs)
				return
			}

			s.Len(hrs, tc.wantLen)

			if tc.validateFn != nil {
				tc.validateFn(hrs)
			}
		})
	}
}

func (s *ResultPublicTestSuite) TestHostResultDecode() {
	tests := []struct {
		name        string
		hostResult  orchestrator.HostResult
		target      any
		expectErr   bool
		errContains string
		validateFn  func(cmd osapi.CommandResult)
	}{
		{
			name: "Decodes host result data into typed struct",
			hostResult: orchestrator.HostResult{
				Hostname: "web-01",
				Changed:  true,
				Data: map[string]any{
					"stdout":    "hello",
					"stderr":    "",
					"exit_code": float64(0),
				},
			},
			validateFn: func(cmd osapi.CommandResult) {
				s.Equal("hello", cmd.Stdout)
				s.Equal(0, cmd.ExitCode)
			},
		},
		{
			name: "Returns error for unmarshalable data",
			hostResult: orchestrator.HostResult{
				Data: map[string]any{"fn": func() {}},
			},
			expectErr:   true,
			errContains: "marshal host result data",
		},
		{
			name: "Returns error when decode target is invalid",
			hostResult: orchestrator.HostResult{
				Data: map[string]any{"stdout": "hello"},
			},
			target:      &struct{ Stdout chan int }{},
			expectErr:   true,
			errContains: "decode host result data",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			if tc.target != nil {
				err := tc.hostResult.Decode(tc.target)
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)

				return
			}

			var cmd osapi.CommandResult
			err := tc.hostResult.Decode(&cmd)

			if tc.expectErr {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)

				return
			}

			s.Require().NoError(err)

			if tc.validateFn != nil {
				tc.validateFn(cmd)
			}
		})
	}
}

func (s *ResultPublicTestSuite) TestReportDecode() {
	tests := []struct {
		name        string
		tasks       []engine.TaskResult
		lookupName  string
		target      any
		expectErr   bool
		errContains string
		validateFn  func(cmd osapi.CommandResult)
	}{
		{
			name: "Decodes task result from report",
			tasks: []engine.TaskResult{
				{
					Name:    "run-cmd",
					Status:  engine.StatusChanged,
					Changed: true,
					Data: map[string]any{
						"stdout":    "hello",
						"exit_code": float64(0),
					},
				},
			},
			lookupName: "run-cmd",
			validateFn: func(cmd osapi.CommandResult) {
				s.Equal("hello", cmd.Stdout)
				s.Equal(0, cmd.ExitCode)
			},
		},
		{
			name: "Decodes from host results when present",
			tasks: []engine.TaskResult{
				{
					Name:    "run-cmd",
					Status:  engine.StatusChanged,
					Changed: true,
					Data:    map[string]any{"results": []any{}},
					HostResults: []engine.HostResult{
						{
							Hostname: "web-01",
							Data: map[string]any{
								"stdout":    "hello from host",
								"exit_code": float64(0),
							},
						},
					},
				},
			},
			lookupName: "run-cmd",
			validateFn: func(cmd osapi.CommandResult) {
				s.Equal("hello from host", cmd.Stdout)
				s.Equal(0, cmd.ExitCode)
			},
		},
		{
			name:        "Returns error for missing task",
			tasks:       []engine.TaskResult{},
			lookupName:  "nonexistent",
			expectErr:   true,
			errContains: "no result for",
		},
		{
			name: "Returns error for nil data",
			tasks: []engine.TaskResult{
				{
					Name:   "no-data",
					Status: engine.StatusSkipped,
				},
			},
			lookupName:  "no-data",
			expectErr:   true,
			errContains: "no result data for",
		},
		{
			name: "Returns error when marshal fails",
			tasks: []engine.TaskResult{
				{
					Name:   "bad-marshal",
					Status: engine.StatusChanged,
					Data:   map[string]any{"fn": func() {}},
				},
			},
			lookupName:  "bad-marshal",
			expectErr:   true,
			errContains: "marshal result data",
		},
		{
			name: "Returns error when decode target is invalid",
			tasks: []engine.TaskResult{
				{
					Name:   "bad-decode",
					Status: engine.StatusChanged,
					Data:   map[string]any{"stdout": "hello"},
				},
			},
			lookupName:  "bad-decode",
			target:      &struct{ Stdout chan int }{},
			expectErr:   true,
			errContains: "decode result data",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			report := &orchestrator.Report{
				Tasks: tc.tasks,
			}

			if tc.target != nil {
				err := report.Decode(tc.lookupName, tc.target)
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)

				return
			}

			var cmd osapi.CommandResult
			err := report.Decode(tc.lookupName, &cmd)

			if tc.expectErr {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)

				return
			}

			s.Require().NoError(err)

			if tc.validateFn != nil {
				tc.validateFn(cmd)
			}
		})
	}
}

func TestResultPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(ResultPublicTestSuite))
}
