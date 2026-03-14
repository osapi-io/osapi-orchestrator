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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
	"github.com/stretchr/testify/suite"
)

type OrchestratorTestSuite struct {
	suite.Suite
}

func (s *OrchestratorTestSuite) TestNextOpName() {
	tests := []struct {
		name     string
		prefix   string
		expected string
	}{
		{
			name:     "First use returns prefix as-is",
			prefix:   "get-hostname",
			expected: "get-hostname",
		},
		{
			name:     "Different prefix returns as-is",
			prefix:   "get-disk",
			expected: "get-disk",
		},
		{
			name:     "Duplicate prefix gets counter suffix",
			prefix:   "get-hostname",
			expected: "get-hostname-2",
		},
		{
			name:     "Command prefix includes command name",
			prefix:   "run-uptime",
			expected: "run-uptime",
		},
		{
			name:     "Duplicate command prefix gets counter suffix",
			prefix:   "run-uptime",
			expected: "run-uptime-2",
		},
	}

	o := &Orchestrator{nameCount: make(map[string]int)}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := o.nextOpName(tc.prefix)
			s.Equal(tc.expected, got)
		})
	}
}

func (s *OrchestratorTestSuite) TestRunWithHealthCheck() {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr bool
	}{
		{
			name: "Succeeds with healthy server",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"ok"}`))
			}),
		},
		{
			name: "Fails with unhealthy server",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(`{"status":"unhealthy"}`))
			}),
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			orch := New(server.URL, "test-token")
			orch.HealthCheck()

			report, err := orch.Run()

			if tc.expectErr {
				s.Require().Error(err)
				s.Nil(report)

				return
			}

			s.Require().NoError(err)
			s.NotNil(report)
			s.Len(report.Tasks, 1)
		})
	}
}

func (s *OrchestratorTestSuite) TestRendererHooks() {
	retryErr := errors.New("transient error")

	tests := []struct {
		name         string
		setupFunc    func(hooks sdk.Hooks)
		validateFunc func(m *mockRenderer)
	}{
		{
			name: "BeforePlan calls PlanStart",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.BeforePlan(sdk.PlanSummary{
					TotalTasks: 3,
					Steps: []sdk.StepSummary{
						{Tasks: []string{"a", "b"}, Parallel: true},
						{Tasks: []string{"c"}, Parallel: false},
					},
				})
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.planStartCalled)
				s.Equal(3, m.planStartSummary.TotalTasks)
				s.Len(m.planStartSummary.Steps, 2)
			},
		},
		{
			name: "AfterPlan calls PlanDone",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.AfterPlan(&sdk.Report{
					Tasks:    []sdk.TaskResult{{Name: "a", Status: sdk.StatusChanged}},
					Duration: 5 * time.Second,
				})
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.planDoneCalled)
				s.Require().NotNil(m.planDoneReport)
				s.Len(m.planDoneReport.Tasks, 1)
				s.Equal(5*time.Second, m.planDoneReport.Duration)
			},
		},
		{
			name: "BeforeLevel sequential single task",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.BeforeLevel(0, []*sdk.Task{sdk.NewTaskFunc("task-1", nil)}, false)
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.levelStartCalled)
				s.Equal(0, m.levelStartLevel)
				s.Equal([]string{"task-1"}, m.levelStartTasks)
				s.False(m.levelStartParallel)
			},
		},
		{
			name: "BeforeLevel parallel multiple tasks",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.BeforeLevel(0, []*sdk.Task{
					sdk.NewTaskFunc("task-1", nil),
					sdk.NewTaskFunc("task-2", nil),
				}, true)
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.levelStartCalled)
				s.Equal([]string{"task-1", "task-2"}, m.levelStartTasks)
				s.True(m.levelStartParallel)
			},
		},
		{
			name: "AfterLevel no results sequential",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.BeforeLevel(1, []*sdk.Task{sdk.NewTaskFunc("t", nil)}, false)
				hooks.AfterLevel(1, nil)
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.levelDoneCalled)
				s.Equal(1, m.levelDoneLevel)
				s.Equal(0, m.levelDoneChanged)
				s.Equal(0, m.levelDoneTotal)
				s.False(m.levelDoneParallel)
			},
		},
		{
			name: "AfterLevel with changed results parallel",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.BeforeLevel(1, []*sdk.Task{sdk.NewTaskFunc("t", nil)}, true)
				hooks.AfterLevel(1, []sdk.TaskResult{
					{Name: "a", Status: sdk.StatusChanged, Changed: true},
					{Name: "b", Status: sdk.StatusUnchanged},
				})
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.levelDoneCalled)
				s.Equal(1, m.levelDoneChanged)
				s.Equal(2, m.levelDoneTotal)
				s.True(m.levelDoneParallel)
			},
		},
		{
			name: "BeforeTask calls TaskStart with name",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.BeforeTask(sdk.NewTaskFunc("fn-task", nil))
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.taskStartCalled)
				s.Equal("fn-task", m.taskStartName)
				s.Empty(m.taskStartDetail)
			},
		},
		{
			name: "AfterTask calls TaskDone",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.AfterTask(nil, sdk.TaskResult{
					Name:     "task-1",
					Status:   sdk.StatusChanged,
					Changed:  true,
					Duration: time.Second,
				})
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.taskDoneCalled)
				s.Equal("task-1", m.taskDoneResult.Name)
				s.Equal(sdk.StatusChanged, m.taskDoneResult.Status)
			},
		},
		{
			name: "OnRetry calls TaskRetry",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.OnRetry(sdk.NewTaskFunc("retry-task", nil), 2, retryErr)
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.taskRetryCalled)
				s.Equal("retry-task", m.taskRetryName)
				s.Equal(2, m.taskRetryAttempt)
				s.Equal(retryErr, m.taskRetryErr)
			},
		},
		{
			name: "OnSkip calls TaskSkip",
			setupFunc: func(hooks sdk.Hooks) {
				hooks.OnSkip(sdk.NewTaskFunc("skip-task", nil), "dependency failed")
			},
			validateFunc: func(m *mockRenderer) {
				s.True(m.taskSkipCalled)
				s.Equal("skip-task", m.taskSkipName)
				s.Equal("dependency failed", m.taskSkipReason)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			m := &mockRenderer{}
			hooks := rendererHooks(m)

			tc.setupFunc(hooks)
			tc.validateFunc(m)
		})
	}
}

// mockRenderer records renderer calls for verification.
type mockRenderer struct {
	planStartCalled  bool
	planStartSummary sdk.PlanSummary

	planDoneCalled bool
	planDoneReport *Report

	levelStartCalled   bool
	levelStartLevel    int
	levelStartTasks    []string
	levelStartParallel bool

	levelDoneCalled   bool
	levelDoneLevel    int
	levelDoneChanged  int
	levelDoneTotal    int
	levelDoneParallel bool

	taskStartCalled bool
	taskStartName   string
	taskStartDetail string

	taskDoneCalled bool
	taskDoneResult sdk.TaskResult

	taskRetryCalled  bool
	taskRetryName    string
	taskRetryAttempt int
	taskRetryErr     error

	taskSkipCalled bool
	taskSkipName   string
	taskSkipReason string
}

func (m *mockRenderer) PlanStart(
	summary sdk.PlanSummary,
) {
	m.planStartCalled = true
	m.planStartSummary = summary
}

func (m *mockRenderer) PlanDone(
	report *Report,
) {
	m.planDoneCalled = true
	m.planDoneReport = report
}

func (m *mockRenderer) LevelStart(
	level int,
	tasks []string,
	parallel bool,
) {
	m.levelStartCalled = true
	m.levelStartLevel = level
	m.levelStartTasks = tasks
	m.levelStartParallel = parallel
}

func (m *mockRenderer) LevelDone(
	level int,
	changed int,
	total int,
	parallel bool,
) {
	m.levelDoneCalled = true
	m.levelDoneLevel = level
	m.levelDoneChanged = changed
	m.levelDoneTotal = total
	m.levelDoneParallel = parallel
}

func (m *mockRenderer) TaskStart(
	name string,
	detail string,
) {
	m.taskStartCalled = true
	m.taskStartName = name
	m.taskStartDetail = detail
}

func (m *mockRenderer) TaskDone(
	result sdk.TaskResult,
) {
	m.taskDoneCalled = true
	m.taskDoneResult = result
}

func (m *mockRenderer) TaskRetry(
	name string,
	attempt int,
	err error,
) {
	m.taskRetryCalled = true
	m.taskRetryName = name
	m.taskRetryAttempt = attempt
	m.taskRetryErr = err
}

func (m *mockRenderer) TaskSkip(
	name string,
	reason string,
) {
	m.taskSkipCalled = true
	m.taskSkipName = name
	m.taskSkipReason = reason
}

func TestOrchestratorTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OrchestratorTestSuite))
}
