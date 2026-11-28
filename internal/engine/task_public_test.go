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
	"testing"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

type TaskPublicTestSuite struct {
	suite.Suite
}

func TestTaskPublicTestSuite(t *testing.T) {
	suite.Run(t, new(TaskPublicTestSuite))
}

// noop is a no-op TaskFn for tests that only need a valid task.
func noop(
	_ context.Context,
	_ *osapi.Client,
) (*engine.Result, error) {
	return &engine.Result{}, nil
}

func (s *TaskPublicTestSuite) TestDependsOn() {
	tests := []struct {
		name         string
		setupDeps    func(a, b, c *engine.Task)
		checkTask    string
		validateFunc func(*engine.Task)
	}{
		{
			name: "single dependency",
			setupDeps: func(a, b, _ *engine.Task) {
				b.DependsOn(a)
			},
			checkTask: "b",
			validateFunc: func(task *engine.Task) {
				s.Len(task.Dependencies(), 1)
			},
		},
		{
			name: "multiple dependencies",
			setupDeps: func(a, b, c *engine.Task) {
				c.DependsOn(a, b)
			},
			checkTask: "c",
			validateFunc: func(task *engine.Task) {
				s.Len(task.Dependencies(), 2)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			a := engine.NewTaskFunc("a", noop)
			b := engine.NewTaskFunc("b", noop)
			c := engine.NewTaskFunc("c", noop)
			tt.setupDeps(a, b, c)

			tasks := map[string]*engine.Task{"a": a, "b": b, "c": c}

			tt.validateFunc(tasks[tt.checkTask])
		})
	}
}

func (s *TaskPublicTestSuite) TestOnlyIfChanged() {
	tests := []struct {
		name          string
		onlyIfChanged bool
		validateFunc  func(*engine.Task)
	}{
		{
			name: "a new task requires no change",
			validateFunc: func(task *engine.Task) {
				s.False(task.RequiresChange())
			},
		},
		{
			name:          "a task after OnlyIfChanged requires a change",
			onlyIfChanged: true,
			validateFunc: func(task *engine.Task) {
				s.True(task.RequiresChange())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			task.DependsOn(engine.NewTaskFunc("dep", noop))
			if tt.onlyIfChanged {
				task.OnlyIfChanged()
			}

			tt.validateFunc(task)
		})
	}
}

func (s *TaskPublicTestSuite) TestWhen() {
	tests := []struct {
		name         string
		guard        bool
		validateFunc func(task *engine.Task, called *bool)
	}{
		{
			name: "a new task has no guard, and is evaluated to run",
			validateFunc: func(task *engine.Task, called *bool) {
				s.False(task.HasGuard())

				run, reason := task.Evaluate(engine.Results{})
				s.True(run)
				s.Empty(reason)
				s.False(*called)
			},
		},
		{
			name:  "a guard is called when the task is evaluated",
			guard: true,
			validateFunc: func(task *engine.Task, called *bool) {
				s.True(task.HasGuard())

				run, reason := task.Evaluate(engine.Results{})
				s.True(run)
				s.Empty(reason)
				s.True(*called)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			called := false
			if tt.guard {
				task.When(func(_ engine.Results) bool {
					called = true

					return true
				})
			}

			tt.validateFunc(task, &called)
		})
	}
}

func (s *TaskPublicTestSuite) TestTaskFunc() {
	tests := []struct {
		name         string
		fn           engine.TaskFn
		validateFunc func(*engine.Task)
	}{
		{
			name: "names the task and keeps its function",
			fn:   noop,
			validateFunc: func(task *engine.Task) {
				s.Equal("custom", task.Name())
				s.True(task.IsFunc())
				s.NotNil(task.Fn())
			},
		},
		{
			name: "a task with no function is not a function",
			validateFunc: func(task *engine.Task) {
				s.Equal("custom", task.Name())
				s.False(task.IsFunc())
				s.Nil(task.Fn())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(engine.NewTaskFunc("custom", tt.fn))
		})
	}
}

func (s *TaskPublicTestSuite) TestSetName() {
	tests := []struct {
		name         string
		initial      string
		renamed      string
		validateFunc func(string)
	}{
		{
			name:    "changes task name",
			initial: "original",
			renamed: "renamed",
			validateFunc: func(got string) {
				s.Equal("renamed", got)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc(tt.initial, noop)
			task.SetName(tt.renamed)

			tt.validateFunc(task.Name())
		})
	}
}

func (s *TaskPublicTestSuite) TestEvaluate() {
	yes := func(_ engine.Results) bool { return true }
	no := func(_ engine.Results) bool { return false }

	tests := []struct {
		name       string
		setup      func(task *engine.Task)
		wantRun    bool
		wantReason string
	}{
		{
			name:    "a task with no guard runs",
			setup:   func(_ *engine.Task) {},
			wantRun: true,
		},
		{
			name: "When gives its default reason",
			setup: func(task *engine.Task) {
				task.When(no)
			},
			wantReason: "guard returned false",
		},
		{
			name: "WhenWithReason gives its reason",
			setup: func(task *engine.Task) {
				task.WhenWithReason(no, "host is unreachable")
			},
			wantReason: "host is unreachable",
		},
		{
			name: "Check gives its own reason",
			setup: func(task *engine.Task) {
				task.Check(func(_ engine.Results) (bool, string) {
					return false, "dynamic reason"
				})
			},
			wantReason: "dynamic reason",
		},
		{
			name: "guards join with AND: a later guard that fails skips the task",
			setup: func(task *engine.Task) {
				task.WhenWithReason(yes, "first")
				task.WhenWithReason(no, "second")
			},
			wantReason: "second",
		},
		{
			name: "the first guard that fails gives the reason",
			setup: func(task *engine.Task) {
				task.WhenWithReason(no, "first")
				task.WhenWithReason(no, "second")
			},
			wantReason: "first",
		},
		{
			name: "each guard holds",
			setup: func(task *engine.Task) {
				task.WhenWithReason(yes, "first")
				task.WhenWithReason(yes, "second")
			},
			wantRun: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			tt.setup(task)

			run, reason := task.Evaluate(engine.Results{})
			s.Equal(tt.wantRun, run)
			s.Equal(tt.wantReason, reason)
		})
	}
}

func (s *TaskPublicTestSuite) TestRetry() {
	tests := []struct {
		name         string
		count        int
		validateFunc func(*engine.Task)
	}{
		{
			name: "a new task has no retry policy",
			validateFunc: func(task *engine.Task) {
				s.Nil(task.RetryPolicy())
			},
		},
		{
			name:  "Retry sets the policy",
			count: 3,
			validateFunc: func(task *engine.Task) {
				s.Require().NotNil(task.RetryPolicy())
				s.Equal(3, task.RetryPolicy().Count())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			if tt.count > 0 {
				task.Retry(engine.Retry(tt.count))
			}

			tt.validateFunc(task)
		})
	}
}

func (s *TaskPublicTestSuite) TestOnErrorOverride() {
	tests := []struct {
		name         string
		strategy     *engine.ErrorStrategy
		validateFunc func(*engine.Task)
	}{
		{
			name: "a new task has no error strategy of its own",
			validateFunc: func(task *engine.Task) {
				s.Nil(task.ErrorStrategy())
			},
		},
		{
			name:     "OnError sets the task's error strategy",
			strategy: &engine.Continue,
			validateFunc: func(task *engine.Task) {
				s.Require().NotNil(task.ErrorStrategy())
				s.Equal("continue", task.ErrorStrategy().String())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			if tt.strategy != nil {
				task.OnError(*tt.strategy)
			}

			tt.validateFunc(task)
		})
	}
}
