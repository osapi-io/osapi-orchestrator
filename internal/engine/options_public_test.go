package engine_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	osapiclient "github.com/retr0h/osapi/pkg/sdk/client"
	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

type OptionsPublicTestSuite struct {
	suite.Suite
}

func TestOptionsPublicTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsPublicTestSuite))
}

func (s *OptionsPublicTestSuite) TestErrorStrategy() {
	tests := []struct {
		name     string
		strategy engine.ErrorStrategy
		wantStr  string
	}{
		{
			name:     "stop all",
			strategy: engine.StopAll,
			wantStr:  "stop_all",
		},
		{
			name:     "continue",
			strategy: engine.Continue,
			wantStr:  "continue",
		},
		{
			name:     "retry",
			strategy: engine.Retry(3),
			wantStr:  "retry(3)",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.wantStr, tt.strategy.String())
		})
	}
}

func (s *OptionsPublicTestSuite) TestRetryCount() {
	tests := []struct {
		name     string
		strategy engine.ErrorStrategy
		want     int
	}{
		{
			name:     "stop all has zero retries",
			strategy: engine.StopAll,
			want:     0,
		},
		{
			name:     "continue has zero retries",
			strategy: engine.Continue,
			want:     0,
		},
		{
			name:     "retry has n retries",
			strategy: engine.Retry(5),
			want:     5,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.want, tt.strategy.RetryCount())
		})
	}
}

func (s *OptionsPublicTestSuite) TestWithHooks() {
	called := false
	hooks := engine.Hooks{
		BeforeTask: func(_ *engine.Task) {
			called = true
		},
	}

	cfg := engine.PlanConfig{}
	opt := engine.WithHooks(hooks)
	opt(&cfg)

	s.NotNil(cfg.Hooks)
	s.NotNil(cfg.Hooks.BeforeTask)

	// Create a task to pass to the callback.
	t := engine.NewTaskFunc(
		"test",
		func(
			_ context.Context,
			_ *osapiclient.Client,
		) (*engine.Result, error) {
			return &engine.Result{}, nil
		},
	)
	cfg.Hooks.BeforeTask(t)
	s.True(called)
}

func (s *OptionsPublicTestSuite) TestHooksDefaults() {
	h := engine.Hooks{}

	// Nil callbacks should be safe — no panic.
	s.Nil(h.BeforePlan)
	s.Nil(h.AfterPlan)
	s.Nil(h.BeforeLevel)
	s.Nil(h.AfterLevel)
	s.Nil(h.BeforeTask)
	s.Nil(h.AfterTask)
	s.Nil(h.OnRetry)
	s.Nil(h.OnSkip)
}

func (s *OptionsPublicTestSuite) TestPlanOption() {
	tests := []struct {
		name        string
		option      engine.PlanOption
		wantOnError engine.ErrorStrategy
	}{
		{
			name:        "on error sets strategy",
			option:      engine.OnError(engine.Continue),
			wantOnError: engine.Continue,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			cfg := &engine.PlanConfig{}
			tt.option(cfg)
			s.Equal(tt.wantOnError.String(), cfg.OnErrorStrategy.String())
		})
	}
}
