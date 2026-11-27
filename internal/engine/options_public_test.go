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
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapiclient "github.com/osapi-io/osapi/pkg/sdk/client"
)

type OptionsPublicTestSuite struct {
	suite.Suite
}

func TestOptionsPublicTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsPublicTestSuite))
}

func (s *OptionsPublicTestSuite) TestErrorStrategy() {
	tests := []struct {
		name         string
		strategy     engine.ErrorStrategy
		validateFunc func(string)
	}{
		{
			name:     "stop all",
			strategy: engine.StopAll,
			validateFunc: func(got string) {
				s.Equal("stop_all", got)
			},
		},
		{
			name:     "continue",
			strategy: engine.Continue,
			validateFunc: func(got string) {
				s.Equal("continue", got)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(tt.strategy.String())
		})
	}
}

func (s *OptionsPublicTestSuite) TestRetryPolicy() {
	tests := []struct {
		name         string
		policy       engine.RetryPolicy
		validateFunc func(engine.RetryPolicy)
	}{
		{
			name:   "the zero policy has no retries",
			policy: engine.RetryPolicy{},
			validateFunc: func(got engine.RetryPolicy) {
				s.Equal(0, got.Count())
				s.Equal("retry(0)", got.String())
			},
		},
		{
			name:   "retry has n retries",
			policy: engine.Retry(5),
			validateFunc: func(got engine.RetryPolicy) {
				s.Equal(5, got.Count())
				s.Equal("retry(5)", got.String())
			},
		},
		{
			name:   "a negative count is no retry",
			policy: engine.Retry(-2),
			validateFunc: func(got engine.RetryPolicy) {
				s.Equal(0, got.Count())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(tt.policy)
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
		name         string
		option       engine.PlanOption
		validateFunc func(*engine.PlanConfig)
	}{
		{
			name:   "on error sets strategy",
			option: engine.OnError(engine.Continue),
			validateFunc: func(cfg *engine.PlanConfig) {
				s.Equal(engine.Continue.String(), cfg.OnErrorStrategy.String())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			cfg := &engine.PlanConfig{}
			tt.option(cfg)
			tt.validateFunc(cfg)
		})
	}
}

func (s *OptionsPublicTestSuite) TestBackoffDelay() {
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
		{
			name:    "a max interval of 0 doubles below the default cap",
			initial: 100 * time.Millisecond,
			max:     0,
			attempt: 3,
			validateFunc: func(got time.Duration) {
				s.Equal(800*time.Millisecond, got)
			},
		},
		{
			name:    "a long run of retries with a max of 0 stops at 30s",
			initial: time.Second,
			max:     0,
			attempt: 100,
			validateFunc: func(got time.Duration) {
				s.Equal(30*time.Second, got)
			},
		},
		{
			name:    "an initial interval above the default cap is its own cap",
			initial: time.Minute,
			max:     0,
			attempt: 4,
			validateFunc: func(got time.Duration) {
				s.Equal(time.Minute, got)
			},
		},
		{
			name:    "a long run of retries is clamped to the max",
			initial: time.Second,
			max:     time.Minute,
			attempt: 100,
			validateFunc: func(got time.Duration) {
				s.Equal(time.Minute, got)
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

func (s *OptionsPublicTestSuite) TestJitter() {
	tests := []struct {
		name  string
		delay time.Duration
		low   time.Duration
		high  time.Duration
	}{
		{
			name:  "a delay is spread over its upper half",
			delay: 100 * time.Millisecond,
			low:   50 * time.Millisecond,
			high:  100 * time.Millisecond,
		},
		{
			name:  "a delay of 1ns or less is kept",
			delay: 1,
			low:   1,
			high:  1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			for range 200 {
				got := engine.ExportJitter(tt.delay)
				s.GreaterOrEqual(got, tt.low)
				s.LessOrEqual(got, tt.high)
			}
		})
	}
}

func (s *OptionsPublicTestSuite) TestIntervals() {
	initial, maxInterval := engine.Retry(
		2,
		engine.WithRetryBackoff(time.Second, time.Minute),
	).Intervals()

	s.Equal(time.Second, initial)
	s.Equal(time.Minute, maxInterval)
}
