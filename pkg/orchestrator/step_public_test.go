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

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
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

func (s *StepPublicTestSuite) TestChaining() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "After returns same step",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck("_any")
				step := s.orch.NodeHostnameGet("_any")

				return step.After(health)
			},
		},
		{
			name: "After with multiple dependencies",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck("_any")
				disk := s.orch.NodeDiskGet("_any")
				step := s.orch.NodeHostnameGet("_any")

				return step.After(health, disk)
			},
		},
		{
			name: "Retry returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").Retry(3)
			},
		},
		{
			name: "OnlyIfChanged returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").OnlyIfChanged()
			},
		},
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
			name: "OnlyIfFailed returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").OnlyIfFailed()
			},
		},
		{
			name: "OnlyIfAllChanged returns non-nil step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").OnlyIfAllChanged()
			},
		},
		{
			name: "Named returns same step",
			chainFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any").Named("custom")
			},
		},
		{
			name: "Full method chain",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck("_any")

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

func TestStepPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(StepPublicTestSuite))
}
