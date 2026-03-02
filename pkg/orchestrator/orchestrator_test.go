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

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type OrchestratorTestSuite struct {
	suite.Suite
}

func (s *OrchestratorTestSuite) TestNextName() {
	o := &Orchestrator{}

	tests := []struct {
		name      string
		operation string
		expected  string
	}{
		{
			name:      "First call increments to 1",
			operation: "node.hostname.get",
			expected:  "node.hostname.get-1",
		},
		{
			name:      "Second call increments to 2",
			operation: "node.disk.get",
			expected:  "node.disk.get-2",
		},
		{
			name:      "Third call increments to 3",
			operation: "node.hostname.get",
			expected:  "node.hostname.get-3",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := o.nextName(tc.operation)
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
			orch.HealthCheck("_any")

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
	s.Run("BeforePlan calls PlanStart", func() {
		m := &mockRenderer{}
		hooks := rendererHooks(m)

		summary := sdk.PlanSummary{
			TotalTasks: 3,
			Steps: []sdk.StepSummary{
				{Tasks: []string{"a", "b"}, Parallel: true},
				{Tasks: []string{"c"}, Parallel: false},
			},
		}
		hooks.BeforePlan(summary)

		s.True(m.planStartCalled)
		s.Equal(summary, m.planStartSummary)
	})

	s.Run("AfterPlan calls PlanDone", func() {
		m := &mockRenderer{}
		hooks := rendererHooks(m)

		hooks.AfterPlan(&sdk.Report{
			Tasks:    []sdk.TaskResult{{Name: "a", Status: sdk.StatusChanged}},
			Duration: 5 * time.Second,
		})

		s.True(m.planDoneCalled)
		s.Require().NotNil(m.planDoneReport)
		s.Len(m.planDoneReport.Tasks, 1)
		s.Equal(5*time.Second, m.planDoneReport.Duration)
	})

	s.Run("BeforeLevel calls LevelStart", func() {
		tests := []struct {
			name          string
			tasks         []*sdk.Task
			parallel      bool
			expectedNames []string
		}{
			{
				name:          "Sequential single task",
				tasks:         []*sdk.Task{sdk.NewTask("task-1", nil)},
				parallel:      false,
				expectedNames: []string{"task-1"},
			},
			{
				name: "Parallel multiple tasks",
				tasks: []*sdk.Task{
					sdk.NewTask("task-1", nil),
					sdk.NewTask("task-2", nil),
				},
				parallel:      true,
				expectedNames: []string{"task-1", "task-2"},
			},
		}

		for _, tc := range tests {
			s.Run(tc.name, func() {
				m := &mockRenderer{}
				hooks := rendererHooks(m)

				hooks.BeforeLevel(0, tc.tasks, tc.parallel)

				s.True(m.levelStartCalled)
				s.Equal(0, m.levelStartLevel)
				s.Equal(tc.expectedNames, m.levelStartTasks)
				s.Equal(tc.parallel, m.levelStartParallel)
			})
		}
	})

	s.Run("AfterLevel calls LevelDone", func() {
		tests := []struct {
			name             string
			parallel         bool
			results          []sdk.TaskResult
			expectedChanged  int
			expectedTotal    int
			expectedParallel bool
		}{
			{
				name:             "No results sequential",
				parallel:         false,
				results:          nil,
				expectedChanged:  0,
				expectedTotal:    0,
				expectedParallel: false,
			},
			{
				name:     "With changed results parallel",
				parallel: true,
				results: []sdk.TaskResult{
					{Name: "a", Status: sdk.StatusChanged, Changed: true},
					{Name: "b", Status: sdk.StatusUnchanged},
				},
				expectedChanged:  1,
				expectedTotal:    2,
				expectedParallel: true,
			},
		}

		for _, tc := range tests {
			s.Run(tc.name, func() {
				m := &mockRenderer{}
				hooks := rendererHooks(m)

				// BeforeLevel sets the parallel state for the level.
				hooks.BeforeLevel(1, []*sdk.Task{sdk.NewTask("t", nil)}, tc.parallel)
				hooks.AfterLevel(1, tc.results)

				s.True(m.levelDoneCalled)
				s.Equal(1, m.levelDoneLevel)
				s.Equal(tc.expectedChanged, m.levelDoneChanged)
				s.Equal(tc.expectedTotal, m.levelDoneTotal)
				s.Equal(tc.expectedParallel, m.levelDoneParallel)
			})
		}
	})

	s.Run("BeforeTask calls TaskStart", func() {
		tests := []struct {
			name           string
			task           *sdk.Task
			expectedName   string
			expectedDetail string
		}{
			{
				name: "Op task",
				task: sdk.NewTask("op-task", &sdk.Op{
					Operation: "node.hostname.get",
					Target:    "_any",
				}),
				expectedName:   "op-task",
				expectedDetail: "operation=node.hostname.get target=_any",
			},
			{
				name:           "Func task",
				task:           sdk.NewTaskFunc("fn-task", nil),
				expectedName:   "fn-task",
				expectedDetail: "(custom function)",
			},
		}

		for _, tc := range tests {
			s.Run(tc.name, func() {
				m := &mockRenderer{}
				hooks := rendererHooks(m)

				hooks.BeforeTask(tc.task)

				s.True(m.taskStartCalled)
				s.Equal(tc.expectedName, m.taskStartName)
				s.Equal(tc.expectedDetail, m.taskStartDetail)
			})
		}
	})

	s.Run("AfterTask calls TaskDone", func() {
		m := &mockRenderer{}
		hooks := rendererHooks(m)

		result := sdk.TaskResult{
			Name:     "task-1",
			Status:   sdk.StatusChanged,
			Changed:  true,
			Duration: time.Second,
		}
		hooks.AfterTask(nil, result)

		s.True(m.taskDoneCalled)
		s.Equal(result, m.taskDoneResult)
	})

	s.Run("OnRetry calls TaskRetry", func() {
		m := &mockRenderer{}
		hooks := rendererHooks(m)
		retryErr := errors.New("transient error")

		hooks.OnRetry(
			sdk.NewTask("retry-task", nil),
			2,
			retryErr,
		)

		s.True(m.taskRetryCalled)
		s.Equal("retry-task", m.taskRetryName)
		s.Equal(2, m.taskRetryAttempt)
		s.Equal(retryErr, m.taskRetryErr)
	})

	s.Run("OnSkip calls TaskSkip", func() {
		m := &mockRenderer{}
		hooks := rendererHooks(m)

		hooks.OnSkip(
			sdk.NewTask("skip-task", nil),
			"dependency failed",
		)

		s.True(m.taskSkipCalled)
		s.Equal("skip-task", m.taskSkipName)
		s.Equal("dependency failed", m.taskSkipReason)
	})
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
