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

package orchestrator_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

type StepPublicTestSuite struct {
	suite.Suite

	server *httptest.Server
	orch   *orchestrator.Orchestrator
}

func (s *StepPublicTestSuite) SetupTest() {
	s.server = httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	s.orch = orchestrator.New(s.server.URL, "test-token")
}

func (s *StepPublicTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *StepPublicTestSuite) TestAfter() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "After returns same step",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck()
				step := s.orch.HostnameGet("_any")

				return step.After(health)
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
		{
			name: "After with multiple dependencies",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck()
				disk := s.orch.DiskGet("_any")
				step := s.orch.HostnameGet("_any")

				return step.After(health, disk)
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
		{
			name: "Full method chain",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck()

				return s.orch.HostnameGet("_any").
					After(health).
					Retry(2).
					OnlyIfChanged().
					When(func(
						_ orchestrator.Results,
					) bool {
						return true
					}).
					OnError(orchestrator.Continue)
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestRetry() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "Retry returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").Retry(3)
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfChanged() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfChanged returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").OnlyIfChanged()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestWhen() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "When returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").When(func(
					_ orchestrator.Results,
				) bool {
					return true
				})
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnError() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnError with Continue returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					OnError(orchestrator.Continue)
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
		{
			name: "OnError with StopAll returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					OnError(orchestrator.StopAll)
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
		{
			name: "ContinueOnError returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					ContinueOnError()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfFailed() {
	var plan *orchestrator.Orchestrator
	ran := false

	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfFailed returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").OnlyIfFailed()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
		{
			name: "a cleanup after a skipped step does not run, because only direct dependencies count",
			chainFn: func() *orchestrator.Step {
				plan = orchestrator.New(s.server.URL, "test-token")
				fails := plan.TaskFunc("fails", func(
					_ context.Context,
					_ *osapi.Client,
					_ orchestrator.Results,
				) (*orchestrator.Result, error) {
					return nil, errors.New("boom")
				}).ContinueOnError()
				middle := plan.TaskFunc("middle", record(new(bool))).After(fails)

				return plan.TaskFunc("cleanup", record(&ran)).After(middle).OnlyIfFailed()
			},
			validateFunc: func(_ *orchestrator.Step) {
				report, _ := plan.Run(context.Background())

				s.False(ran)
				s.Require().NotNil(report)
				for _, tr := range report.Tasks {
					if tr.Name == "middle" || tr.Name == "cleanup" {
						s.Equal("skipped", string(tr.Status), tr.Name)
					}
				}
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAllChanged() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfAllChanged returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").OnlyIfAllChanged()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAnyHostFailed() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfAnyHostFailed returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					OnlyIfAnyHostFailed()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAllHostsFailed() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfAllHostsFailed returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					OnlyIfAllHostsFailed()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAnyHostSkipped() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfAnyHostSkipped returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					OnlyIfAnyHostSkipped()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAnyHostChanged() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfAnyHostChanged returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					OnlyIfAnyHostChanged()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAllHostsChanged() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "OnlyIfAllHostsChanged returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					OnlyIfAllHostsChanged()
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestNamed() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "Named returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").Named("custom")
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestWhenFact() {
	tests := []struct {
		name         string
		chainFn      func() *orchestrator.Step
		validateFunc func(*orchestrator.Step)
	}{
		{
			name: "WhenFact returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.HostnameGet("_any").
					WhenFact(s.orch.AgentList(), func(
						_ osapi.AgentInfo,
					) bool {
						return true
					})
			},
			validateFunc: func(step *orchestrator.Step) {
				s.NotNil(step)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(tc.chainFn())
		})
	}
}

func (s *StepPublicTestSuite) TestGuardsAndFailures() {
	unchanged := func(
		_ context.Context,
		_ *osapi.Client,
		_ orchestrator.Results,
	) (*orchestrator.Result, error) {
		return &orchestrator.Result{Changed: false}, nil
	}
	fails := func(
		_ context.Context,
		_ *osapi.Client,
		_ orchestrator.Results,
	) (*orchestrator.Result, error) {
		return nil, errors.New("boom")
	}
	always := func(_ orchestrator.Results) bool { return true }

	tests := []struct {
		name       string
		setup      func(o *orchestrator.Orchestrator, ran *bool)
		wantRan    bool
		wantReason string
		// wantErr is the error that Run returns: a step that failed with
		// ContinueOnError lets the plan go on, and Run still reports it.
		wantErr string
	}{
		{
			name: "When does not replace OnlyIfAllChanged",
			setup: func(o *orchestrator.Orchestrator, ran *bool) {
				dep := o.TaskFunc("dep", unchanged)
				o.TaskFunc("step", record(ran)).
					After(dep).
					OnlyIfAllChanged().
					When(always)
			},
			wantReason: "only-if-all-changed: not all dependencies changed",
		},
		{
			name: "a step with a When guard is skipped after a failed dependency",
			setup: func(o *orchestrator.Orchestrator, ran *bool) {
				dep := o.TaskFunc("dep", fails).ContinueOnError()
				o.TaskFunc("step", record(ran)).After(dep).When(always)
			},
			wantReason: "dependency failed",
			wantErr:    `step "dep": boom`,
		},
		{
			name: "OnlyIfFailed runs after a failed dependency",
			setup: func(o *orchestrator.Orchestrator, ran *bool) {
				dep := o.TaskFunc("dep", fails).ContinueOnError()
				o.TaskFunc("step", record(ran)).After(dep).OnlyIfFailed().When(always)
			},
			wantRan: true,
			wantErr: `step "dep": boom`,
		},
		{
			name: "a step retries and then continues",
			setup: func(o *orchestrator.Orchestrator, ran *bool) {
				o.TaskFunc("dep", fails).Retry(1).ContinueOnError()
				o.TaskFunc("step", record(ran))
			},
			wantRan: true,
			wantErr: `step "dep": boom`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			o := orchestrator.New(s.server.URL, "test-token")
			ran := false
			tt.setup(o, &ran)

			report, err := o.Run(context.Background())
			if tt.wantErr != "" {
				s.Require().ErrorContains(err, tt.wantErr)
			} else {
				s.Require().NoError(err)
			}
			s.Equal(tt.wantRan, ran)

			for _, tr := range report.Tasks {
				if tr.Name == "step" {
					s.Equal(tt.wantReason, tr.Reason)
				}
			}
		})
	}
}

// record returns a step function that sets ran.
func record(
	ran *bool,
) func(context.Context, *osapi.Client, orchestrator.Results) (*orchestrator.Result, error) {
	return func(
		_ context.Context,
		_ *osapi.Client,
		_ orchestrator.Results,
	) (*orchestrator.Result, error) {
		*ran = true

		return &orchestrator.Result{Changed: true}, nil
	}
}

func TestStepPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(StepPublicTestSuite))
}

// TestOnlyIfFailedReadsItsDirectDependencies pins the documented rule: a
// dependency that the plan skipped because its own dependency failed is
// Skipped, not Failed, so an OnlyIfFailed step after it does not run.
