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

package engine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	osapiclient "github.com/retr0h/osapi/pkg/sdk/client"
	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

type RunnerPublicTestSuite struct {
	suite.Suite
}

func TestRunnerPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerPublicTestSuite))
}

func (s *RunnerPublicTestSuite) TestLevelize() {
	noop := func(
		_ context.Context,
		_ *osapiclient.Client,
	) (*engine.Result, error) {
		return &engine.Result{}, nil
	}

	tests := []struct {
		name       string
		setup      func() []*engine.Task
		wantLevels int
	}{
		{
			name: "linear chain has 3 levels",
			setup: func() []*engine.Task {
				a := engine.NewTaskFunc("a", noop)
				b := engine.NewTaskFunc("b", noop)
				c := engine.NewTaskFunc("c", noop)
				b.DependsOn(a)
				c.DependsOn(b)

				return []*engine.Task{a, b, c}
			},
			wantLevels: 3,
		},
		{
			name: "diamond has 3 levels",
			setup: func() []*engine.Task {
				a := engine.NewTaskFunc("a", noop)
				b := engine.NewTaskFunc("b", noop)
				c := engine.NewTaskFunc("c", noop)
				d := engine.NewTaskFunc("d", noop)
				b.DependsOn(a)
				c.DependsOn(a)
				d.DependsOn(b, c)

				return []*engine.Task{a, b, c, d}
			},
			wantLevels: 3,
		},
		{
			name: "independent tasks in 1 level",
			setup: func() []*engine.Task {
				a := engine.NewTaskFunc("a", noop)
				b := engine.NewTaskFunc("b", noop)

				return []*engine.Task{a, b}
			},
			wantLevels: 1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tasks := tt.setup()
			levels := engine.ExportLevelize(tasks)
			s.Len(levels, tt.wantLevels)
		})
	}
}

func (s *RunnerPublicTestSuite) TestRunTaskStoresResultForAllPaths() {
	tests := []struct {
		name       string
		setup      func() *engine.Plan
		taskName   string
		wantStatus engine.Status
	}{
		{
			name: "OnlyIfChanged skip stores StatusSkipped",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))

				dep := plan.TaskFunc("dep", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: false}, nil
				})

				child := plan.TaskFunc("child", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: true}, nil
				})
				child.DependsOn(dep)
				child.OnlyIfChanged()

				return plan
			},
			taskName:   "child",
			wantStatus: engine.StatusSkipped,
		},
		{
			name: "failed task stores StatusFailed",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))

				plan.TaskFunc("failing", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return nil, fmt.Errorf("deliberate error")
				})

				return plan
			},
			taskName:   "failing",
			wantStatus: engine.StatusFailed,
		},
		{
			name: "guard-false skip stores StatusSkipped",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))

				plan.TaskFunc("guarded", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: true}, nil
				}).When(func(_ engine.Results) bool {
					return false
				})

				return plan
			},
			taskName:   "guarded",
			wantStatus: engine.StatusSkipped,
		},
		{
			name: "dependency-failed skip stores StatusSkipped",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))

				dep := plan.TaskFunc("dep", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return nil, fmt.Errorf("deliberate error")
				})

				child := plan.TaskFunc("child", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: true}, nil
				})
				child.DependsOn(dep)

				return plan
			},
			taskName:   "child",
			wantStatus: engine.StatusSkipped,
		},
		{
			name: "successful changed task stores StatusChanged",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))

				plan.TaskFunc("ok", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: true}, nil
				})

				return plan
			},
			taskName:   "ok",
			wantStatus: engine.StatusChanged,
		},
		{
			name: "successful unchanged task stores StatusUnchanged",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))

				plan.TaskFunc("ok", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: false}, nil
				})

				return plan
			},
			taskName:   "ok",
			wantStatus: engine.StatusUnchanged,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan := tt.setup()
			er := engine.ExportNewRunner(plan)

			_, err := er.Run(context.Background())
			_ = err

			result := er.GetResult(tt.taskName)
			s.NotNil(
				result,
				"results map should contain entry for %q",
				tt.taskName,
			)
			s.Equal(
				tt.wantStatus,
				result.Status,
				"result status for %q",
				tt.taskName,
			)
		})
	}
}

func (s *RunnerPublicTestSuite) TestDownstreamGuardInspectsSkippedStatus() {
	tests := []struct {
		name            string
		setup           func() (*engine.Plan, *bool)
		observerName    string
		wantGuardCalled bool
		wantTaskStatus  engine.Status
	}{
		{
			name: "guard can see guard-skipped task status",
			setup: func() (*engine.Plan, *bool) {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))
				guardCalled := false

				guarded := plan.TaskFunc("guarded", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: true}, nil
				})
				guarded.When(func(_ engine.Results) bool {
					return false
				})

				observer := plan.TaskFunc("observer", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{Changed: false}, nil
				})
				observer.DependsOn(guarded)
				observer.When(func(r engine.Results) bool {
					guardCalled = true
					res := r.Get("guarded")

					return res != nil && res.Status == engine.StatusSkipped
				})

				return plan, &guardCalled
			},
			observerName:    "observer",
			wantGuardCalled: true,
			wantTaskStatus:  engine.StatusUnchanged,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan, guardCalled := tt.setup()
			er := engine.ExportNewRunner(plan)

			_, err := er.Run(context.Background())
			_ = err

			s.Equal(
				tt.wantGuardCalled,
				*guardCalled,
				"guard should have been called",
			)

			result := er.GetResult(tt.observerName)
			s.NotNil(
				result,
				"observer should have a result entry",
			)
			s.Equal(
				tt.wantTaskStatus,
				result.Status,
				"observer task status",
			)
		})
	}
}

func (s *RunnerPublicTestSuite) TestTaskFuncWithResultsReceivesResults() {
	tests := []struct {
		name        string
		setup       func() (*engine.Plan, *string)
		wantCapture string
	}{
		{
			name: "receives upstream result data",
			setup: func() (*engine.Plan, *string) {
				plan := engine.NewPlan(nil, engine.OnError(engine.StopAll))
				var captured string

				a := plan.TaskFunc("a", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{
						Changed: true,
						Data:    map[string]any{"hostname": "web-01"},
					}, nil
				})

				b := plan.TaskFuncWithResults("b", func(
					_ context.Context,
					_ *osapiclient.Client,
					results engine.Results,
				) (*engine.Result, error) {
					r := results.Get("a")
					if r != nil {
						if h, ok := r.Data["hostname"].(string); ok {
							captured = h
						}
					}

					return &engine.Result{Changed: false}, nil
				})
				b.DependsOn(a)

				return plan, &captured
			},
			wantCapture: "web-01",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan, captured := tt.setup()

			_, err := plan.Run(context.Background())

			s.Require().NoError(err)
			s.Equal(tt.wantCapture, *captured)
		})
	}
}

func (s *RunnerPublicTestSuite) TestTaskResultCarriesData() {
	tests := []struct {
		name     string
		setup    func() *engine.Plan
		taskName string
		wantKey  string
		wantVal  any
	}{
		{
			name: "success result includes data",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.StopAll))

				plan.TaskFunc("a", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{
						Changed: true,
						Data:    map[string]any{"stdout": "hello"},
					}, nil
				})

				return plan
			},
			taskName: "a",
			wantKey:  "stdout",
			wantVal:  "hello",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan := tt.setup()

			report, err := plan.Run(context.Background())

			s.Require().NoError(err)

			var found bool
			for _, tr := range report.Tasks {
				if tr.Name == tt.taskName {
					found = true
					s.Equal(tt.wantVal, tr.Data[tt.wantKey])
				}
			}

			s.True(found, "task %q should be in report", tt.taskName)
		})
	}
}

func (s *RunnerPublicTestSuite) TestBackoffDelay() {
	tests := []struct {
		name    string
		initial time.Duration
		max     time.Duration
		attempt int
		want    time.Duration
	}{
		{
			name:    "first attempt uses initial interval",
			initial: 100 * time.Millisecond,
			max:     10 * time.Second,
			attempt: 0,
			want:    100 * time.Millisecond,
		},
		{
			name:    "second attempt doubles",
			initial: 100 * time.Millisecond,
			max:     10 * time.Second,
			attempt: 1,
			want:    200 * time.Millisecond,
		},
		{
			name:    "clamped to max interval",
			initial: 100 * time.Millisecond,
			max:     300 * time.Millisecond,
			attempt: 5,
			want:    300 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := engine.ExportBackoffDelay(tt.initial, tt.max, tt.attempt)
			s.Equal(tt.want, got)
		})
	}
}

func (s *RunnerPublicTestSuite) TestRunTaskPreservesResultOnError() {
	tests := []struct {
		name            string
		setup           func() *engine.Plan
		taskName        string
		wantChanged     bool
		wantHostResults int
	}{
		{
			name: "TaskFunc error preserves Changed and HostResults",
			setup: func() *engine.Plan {
				plan := engine.NewPlan(nil, engine.OnError(engine.Continue))

				plan.TaskFunc("failing", func(
					_ context.Context,
					_ *osapiclient.Client,
				) (*engine.Result, error) {
					return &engine.Result{
						Changed: true,
						HostResults: []engine.HostResult{
							{Hostname: "web-01", Changed: true},
							{Hostname: "web-02", Error: "timeout"},
						},
					}, fmt.Errorf("partial failure")
				})

				return plan
			},
			taskName:        "failing",
			wantChanged:     true,
			wantHostResults: 2,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan := tt.setup()
			er := engine.ExportNewRunner(plan)

			_, err := er.Run(context.Background())
			_ = err

			result := er.GetResult(tt.taskName)
			s.Require().NotNil(result)
			s.Equal(engine.StatusFailed, result.Status)
			s.Equal(tt.wantChanged, result.Changed)
			s.Len(result.HostResults, tt.wantHostResults)
		})
	}
}
