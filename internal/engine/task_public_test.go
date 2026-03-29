package engine_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapiclient "github.com/retr0h/osapi/pkg/sdk/client"
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
		name       string
		setupDeps  func(a, b, c *engine.Task)
		checkTask  string
		wantDepLen int
	}{
		{
			name: "single dependency",
			setupDeps: func(a, b, _ *engine.Task) {
				b.DependsOn(a)
			},
			checkTask:  "b",
			wantDepLen: 1,
		},
		{
			name: "multiple dependencies",
			setupDeps: func(a, b, c *engine.Task) {
				c.DependsOn(a, b)
			},
			checkTask:  "c",
			wantDepLen: 2,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			a := engine.NewTaskFunc("a", noop)
			b := engine.NewTaskFunc("b", noop)
			c := engine.NewTaskFunc("c", noop)
			tt.setupDeps(a, b, c)

			tasks := map[string]*engine.Task{"a": a, "b": b, "c": c}
			s.Len(tasks[tt.checkTask].Dependencies(), tt.wantDepLen)
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
		name     string
		initial  string
		renamed  string
		wantName string
	}{
		{
			name:     "changes task name",
			initial:  "original",
			renamed:  "renamed",
			wantName: "renamed",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc(tt.initial, noop)
			task.SetName(tt.renamed)
			s.Equal(tt.wantName, task.Name())
		})
	}
}

func (s *TaskPublicTestSuite) TestWhenWithReason() {
	tests := []struct {
		name        string
		guardResult bool
		reason      string
	}{
		{
			name:        "sets guard and reason when guard returns false",
			guardResult: false,
			reason:      "host is unreachable",
		},
		{
			name:        "sets guard and reason when guard returns true",
			guardResult: true,
			reason:      "custom reason",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			task := engine.NewTaskFunc("t", noop)
			task.WhenWithReason(func(_ engine.Results) bool {
				return tt.guardResult
			}, tt.reason)

			guard := task.Guard()
			s.NotNil(guard)
			s.Equal(tt.guardResult, guard(engine.Results{}))
		})
	}
}

func (s *TaskPublicTestSuite) TestSetGuardReason() {
	tests := []struct {
		name       string
		initial    string
		updated    string
		validateFn func(task *engine.Task)
	}{
		{
			name:    "updates reason from initial value",
			initial: "initial reason",
			updated: "dynamic reason",
			validateFn: func(task *engine.Task) {
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
			tt.validateFn(task)
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
