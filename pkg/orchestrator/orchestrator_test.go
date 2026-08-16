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
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	engine "github.com/osapi-io/osapi-orchestrator/internal/engine"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
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

func (s *OrchestratorTestSuite) TestRun() {
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

			report, err := orch.Run(context.Background())

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
		name       string
		expectFunc func(m *Mockrenderer)
		setupFunc  func(hooks engine.Hooks)
	}{
		{
			name: "BeforePlan calls PlanStart",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().PlanStart(gomock.Cond(func(x engine.PlanSummary) bool {
					return x.TotalTasks == 3 && len(x.Steps) == 2
				}))
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.BeforePlan(engine.PlanSummary{
					TotalTasks: 3,
					Steps: []engine.StepSummary{
						{Tasks: []string{"a", "b"}, Parallel: true},
						{Tasks: []string{"c"}, Parallel: false},
					},
				})
			},
		},
		{
			name: "AfterPlan calls PlanDone",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().PlanDone(gomock.Cond(func(r *Report) bool {
					return r != nil && len(r.Tasks) == 1 &&
						r.Duration == 5*time.Second
				}))
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.AfterPlan(&engine.Report{
					Tasks:    []engine.TaskResult{{Name: "a", Status: engine.StatusChanged}},
					Duration: 5 * time.Second,
				})
			},
		},
		{
			name: "BeforeLevel sequential single task",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().LevelStart(0, []string{"task-1"}, false)
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.BeforeLevel(0, []*engine.Task{engine.NewTaskFunc("task-1", nil)}, false)
			},
		},
		{
			name: "BeforeLevel parallel multiple tasks",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().LevelStart(0, []string{"task-1", "task-2"}, true)
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.BeforeLevel(0, []*engine.Task{
					engine.NewTaskFunc("task-1", nil),
					engine.NewTaskFunc("task-2", nil),
				}, true)
			},
		},
		{
			name: "AfterLevel no results sequential",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().LevelStart(1, []string{"t"}, false)
				m.EXPECT().LevelDone(1, 0, 0, false)
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.BeforeLevel(1, []*engine.Task{engine.NewTaskFunc("t", nil)}, false)
				hooks.AfterLevel(1, nil)
			},
		},
		{
			name: "AfterLevel with changed results parallel",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().LevelStart(1, []string{"t"}, true)
				m.EXPECT().LevelDone(1, 1, 2, true)
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.BeforeLevel(1, []*engine.Task{engine.NewTaskFunc("t", nil)}, true)
				hooks.AfterLevel(1, []engine.TaskResult{
					{Name: "a", Status: engine.StatusChanged, Changed: true},
					{Name: "b", Status: engine.StatusUnchanged},
				})
			},
		},
		{
			name: "BeforeTask calls TaskStart with name",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().TaskStart("fn-task", "")
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.BeforeTask(engine.NewTaskFunc("fn-task", nil))
			},
		},
		{
			name: "AfterTask calls TaskDone",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().TaskDone(gomock.Cond(func(r engine.TaskResult) bool {
					return r.Name == "task-1" && r.Status == engine.StatusChanged
				}))
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.AfterTask(nil, engine.TaskResult{
					Name:     "task-1",
					Status:   engine.StatusChanged,
					Changed:  true,
					Duration: time.Second,
				})
			},
		},
		{
			name: "OnRetry calls TaskRetry",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().TaskRetry("retry-task", 2, retryErr)
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.OnRetry(engine.NewTaskFunc("retry-task", nil), 2, retryErr)
			},
		},
		{
			name: "OnSkip calls TaskSkip",
			expectFunc: func(m *Mockrenderer) {
				m.EXPECT().TaskSkip("skip-task", "dependency failed")
			},
			setupFunc: func(hooks engine.Hooks) {
				hooks.OnSkip(engine.NewTaskFunc("skip-task", nil), "dependency failed")
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			m := NewMockrenderer(gomock.NewController(s.T()))

			tc.expectFunc(m)
			tc.setupFunc(rendererHooks(m))
		})
	}
}

func TestOrchestratorTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OrchestratorTestSuite))
}
