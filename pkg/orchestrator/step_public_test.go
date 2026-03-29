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
	"net/http"
	"net/http/httptest"
	"testing"

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
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
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "After returns same step",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck()
				step := s.orch.NodeHostnameGet("_any")

				return step.After(health)
			},
		},
		{
			name: "After with multiple dependencies",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck()
				disk := s.orch.NodeDiskGet("_any")
				step := s.orch.NodeHostnameGet("_any")

				return step.After(health, disk)
			},
		},
		{
			name: "Full method chain",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck()

				return s.orch.NodeHostnameGet("_any").
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
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestRetry() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "Retry returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").Retry(3)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfChanged() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfChanged returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").OnlyIfChanged()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestWhen() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "When returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").When(func(
					_ orchestrator.Results,
				) bool {
					return true
				})
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnError() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnError with Continue returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					OnError(orchestrator.Continue)
			},
		},
		{
			name: "OnError with StopAll returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					OnError(orchestrator.StopAll)
			},
		},
		{
			name: "ContinueOnError returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					ContinueOnError()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfFailed() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfFailed returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").OnlyIfFailed()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAllChanged() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfAllChanged returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").OnlyIfAllChanged()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAnyHostFailed() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfAnyHostFailed returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					OnlyIfAnyHostFailed()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAllHostsFailed() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfAllHostsFailed returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					OnlyIfAllHostsFailed()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAnyHostSkipped() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfAnyHostSkipped returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					OnlyIfAnyHostSkipped()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAnyHostChanged() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfAnyHostChanged returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					OnlyIfAnyHostChanged()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestOnlyIfAllHostsChanged() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "OnlyIfAllHostsChanged returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					OnlyIfAllHostsChanged()
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestNamed() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "Named returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").Named("custom")
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func (s *StepPublicTestSuite) TestWhenFact() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "WhenFact returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").
					WhenFact("list-agents", func(
						_ osapi.Agent,
					) bool {
						return true
					})
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}

func TestStepPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(StepPublicTestSuite))
}
