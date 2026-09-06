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

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapiclient "github.com/osapi-io/osapi/pkg/sdk/client"
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
	_ *osapiclient.Client,
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
	task := engine.NewTaskFunc("t", noop)
	dep := engine.NewTaskFunc("dep", noop)
	task.DependsOn(dep).OnlyIfChanged()

	s.True(task.RequiresChange())
}

func (s *TaskPublicTestSuite) TestWhen() {
	task := engine.NewTaskFunc("t", noop)
	called := false
	task.When(func(_ engine.Results) bool {
		called = true

		return true
	})

	guard := task.Guard()
	s.NotNil(guard)
	s.True(guard(engine.Results{}))
	s.True(called)
}

func (s *TaskPublicTestSuite) TestTaskFunc() {
	fn := func(
		_ context.Context,
		_ *osapiclient.Client,
	) (*engine.Result, error) {
		return &engine.Result{Changed: true}, nil
	}

	task := engine.NewTaskFunc("custom", fn)
	s.Equal("custom", task.Name())
	s.True(task.IsFunc())
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

func (s *TaskPublicTestSuite) TestWhenWithReason() {
	tests := []struct {
		name         string
		guardResult  bool
		reason       string
		validateFunc func(engine.GuardFn)
	}{
		{
			name:        "sets guard and reason when guard returns false",
			guardResult: false,
			reason:      "host is unreachable",
			validateFunc: func(guard engine.GuardFn) {
				s.Require().NotNil(guard)
				s.False(guard(engine.Results{}))
			},
		},
		{
			name:        "sets guard and reason when guard returns true",
			guardResult: true,
			reason:      "custom reason",
			validateFunc: func(guard engine.GuardFn) {
				s.Require().NotNil(guard)
				s.True(guard(engine.Results{}))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			task.WhenWithReason(func(_ engine.Results) bool {
				return tt.guardResult
			}, tt.reason)

			tt.validateFunc(task.Guard())
		})
	}
}

func (s *TaskPublicTestSuite) TestSetGuardReason() {
	tests := []struct {
		name         string
		initial      string
		updated      string
		validateFunc func(task *engine.Task)
	}{
		{
			name:    "updates reason from initial value",
			initial: "initial reason",
			updated: "dynamic reason",
			validateFunc: func(task *engine.Task) {
				s.Equal("dynamic reason", task.GuardReason())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			task.WhenWithReason(func(_ engine.Results) bool {
				return false
			}, tt.initial)

			task.SetGuardReason(tt.updated)
			tt.validateFunc(task)
		})
	}
}

func (s *TaskPublicTestSuite) TestOnErrorOverride() {
	task := engine.NewTaskFunc("t", noop)
	task.OnError(engine.Continue)

	s.NotNil(task.ErrorStrategy())
	s.Equal("continue", task.ErrorStrategy().String())
}

func (s *TaskPublicTestSuite) TestFn() {
	fnTask := engine.NewTaskFunc("fn", func(
		_ context.Context,
		_ *osapiclient.Client,
	) (*engine.Result, error) {
		return nil, nil
	})
	s.NotNil(fnTask.Fn())
}
