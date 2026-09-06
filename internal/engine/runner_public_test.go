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

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapiclient "github.com/osapi-io/osapi/pkg/sdk/client"
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
		name         string
		setup        func() []*engine.Task
		validateFunc func([][]*engine.Task)
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
			validateFunc: func(levels [][]*engine.Task) {
				s.Len(levels, 3)
			},
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
			validateFunc: func(levels [][]*engine.Task) {
				s.Len(levels, 3)
			},
		},
		{
			name: "independent tasks in 1 level",
			setup: func() []*engine.Task {
				a := engine.NewTaskFunc("a", noop)
				b := engine.NewTaskFunc("b", noop)

				return []*engine.Task{a, b}
			},
			validateFunc: func(levels [][]*engine.Task) {
				s.Len(levels, 1)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(engine.ExportLevelize(tt.setup()))
		})
	}
}

func (s *RunnerPublicTestSuite) TestRunTaskStoresResultForAllPaths() {
	tests := []struct {
		name         string
		setup        func() *engine.Plan
		taskName     string
		validateFunc func(*engine.Result)
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
			taskName: "child",
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result)
				s.Equal(engine.StatusSkipped, result.Status)
			},
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
			taskName: "failing",
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result)
				s.Equal(engine.StatusFailed, result.Status)
			},
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
			taskName: "guarded",
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result)
				s.Equal(engine.StatusSkipped, result.Status)
			},
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
			taskName: "child",
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result)
				s.Equal(engine.StatusSkipped, result.Status)
			},
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
			taskName: "ok",
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result)
				s.Equal(engine.StatusChanged, result.Status)
			},
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
			taskName: "ok",
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result)
				s.Equal(engine.StatusUnchanged, result.Status)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			er := engine.ExportNewRunner(tt.setup())

			_, _ = er.Run(context.Background())

			tt.validateFunc(er.GetResult(tt.taskName))
		})
	}
}

func (s *RunnerPublicTestSuite) TestDownstreamGuardInspectsSkippedStatus() {
	tests := []struct {
		name         string
		setup        func() (*engine.Plan, *bool)
		observerName string
		validateFunc func(bool, *engine.Result)
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
			observerName: "observer",
			validateFunc: func(guardCalled bool, result *engine.Result) {
				s.True(guardCalled)
				s.Require().NotNil(result)
				s.Equal(engine.StatusUnchanged, result.Status)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan, guardCalled := tt.setup()
			er := engine.ExportNewRunner(plan)

			_, _ = er.Run(context.Background())

			tt.validateFunc(*guardCalled, er.GetResult(tt.observerName))
		})
	}
}

func (s *RunnerPublicTestSuite) TestTaskFuncWithResultsReceivesResults() {
	tests := []struct {
		name         string
		setup        func() (*engine.Plan, *string)
		validateFunc func(string, error)
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
			validateFunc: func(captured string, err error) {
				s.Require().NoError(err)
				s.Equal("web-01", captured)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan, captured := tt.setup()

			_, err := plan.Run(context.Background())

			tt.validateFunc(*captured, err)
		})
	}
}

func (s *RunnerPublicTestSuite) TestTaskResultCarriesData() {
	tests := []struct {
		name         string
		setup        func() *engine.Plan
		taskName     string
		validateFunc func(map[string]any)
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
			validateFunc: func(data map[string]any) {
				s.Equal("hello", data["stdout"])
			},
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

					tt.validateFunc(tr.Data)
				}
			}

			s.True(found, "task %q should be in report", tt.taskName)
		})
	}
}

func (s *RunnerPublicTestSuite) TestBackoffDelay() {
	tests := []struct {
		name         string
		initial      time.Duration
		max          time.Duration
		attempt      int
		validateFunc func(time.Duration)
	}{
		{
			name:    "first attempt uses initial interval",
			initial: 100 * time.Millisecond,
			max:     10 * time.Second,
			attempt: 0,
			validateFunc: func(got time.Duration) {
				s.Equal(100*time.Millisecond, got)
			},
		},
		{
			name:    "second attempt doubles",
			initial: 100 * time.Millisecond,
			max:     10 * time.Second,
			attempt: 1,
			validateFunc: func(got time.Duration) {
				s.Equal(200*time.Millisecond, got)
			},
		},
		{
			name:    "clamped to max interval",
			initial: 100 * time.Millisecond,
			max:     300 * time.Millisecond,
			attempt: 5,
			validateFunc: func(got time.Duration) {
				s.Equal(300*time.Millisecond, got)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(
				engine.ExportBackoffDelay(tt.initial, tt.max, tt.attempt),
			)
		})
	}
}

func (s *RunnerPublicTestSuite) TestRunTaskPreservesResultOnError() {
	tests := []struct {
		name         string
		setup        func() *engine.Plan
		taskName     string
		validateFunc func(*engine.Result)
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
			taskName: "failing",
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result)
				s.Equal(engine.StatusFailed, result.Status)
				s.True(result.Changed)
				s.Len(result.HostResults, 2)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			er := engine.ExportNewRunner(tt.setup())

			_, _ = er.Run(context.Background())

			tt.validateFunc(er.GetResult(tt.taskName))
		})
	}
}
