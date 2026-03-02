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

type OpsPublicTestSuite struct {
	suite.Suite

	server *httptest.Server
	orch   *orchestrator.Orchestrator
}

func (s *OpsPublicTestSuite) SetupTest() {
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

func (s *OpsPublicTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *OpsPublicTestSuite) TestOperations() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "HealthCheck",
			newFn: func() *orchestrator.Step {
				return s.orch.HealthCheck("_any")
			},
		},
		{
			name: "NodeHostnameGet",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any")
			},
		},
		{
			name: "NodeStatusGet",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeStatusGet("_any")
			},
		},
		{
			name: "NodeUptimeGet",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeUptimeGet("_any")
			},
		},
		{
			name: "NodeDiskGet",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeDiskGet("_any")
			},
		},
		{
			name: "NodeMemoryGet",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeMemoryGet("_any")
			},
		},
		{
			name: "NodeLoadGet",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeLoadGet("_any")
			},
		},
		{
			name: "NetworkDNSGet",
			newFn: func() *orchestrator.Step {
				return s.orch.NetworkDNSGet("_any", "eth0")
			},
		},
		{
			name: "NetworkDNSUpdate",
			newFn: func() *orchestrator.Step {
				return s.orch.NetworkDNSUpdate(
					"_any",
					"eth0",
					[]string{"8.8.8.8"},
					[]string{"example.com"},
				)
			},
		},
		{
			name: "NetworkPingDo",
			newFn: func() *orchestrator.Step {
				return s.orch.NetworkPingDo("_any", "8.8.8.8")
			},
		},
		{
			name: "CommandExec",
			newFn: func() *orchestrator.Step {
				return s.orch.CommandExec("_any", "uptime", "-s")
			},
		},
		{
			name: "CommandShell",
			newFn: func() *orchestrator.Step {
				return s.orch.CommandShell("_any", "echo hello")
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.newFn()
			s.NotNil(step)
		})
	}
}

func TestOpsPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OpsPublicTestSuite))
}
