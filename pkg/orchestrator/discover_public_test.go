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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
)

const agentListJSON = `{
	"agents": [
		{
			"hostname": "web-01",
			"status": "Ready",
			"architecture": "amd64",
			"cpu_count": 8,
			"os_info": {"distribution": "Ubuntu", "version": "22.04"},
			"memory": {"total": 16000, "used": 8000, "free": 8000},
			"facts": {"env": "prod"}
		},
		{
			"hostname": "web-02",
			"status": "Ready",
			"architecture": "arm64",
			"cpu_count": 4,
			"os_info": {"distribution": "Debian", "version": "12"},
			"memory": {"total": 8000, "used": 4000, "free": 4000},
			"facts": {"env": "staging"}
		},
		{
			"hostname": "web-03",
			"status": "Ready",
			"architecture": "amd64",
			"cpu_count": 2,
			"os_info": {"distribution": "Ubuntu", "version": "20.04"},
			"memory": {"total": 4000, "used": 2000, "free": 2000},
			"facts": {"env": "prod"}
		}
	],
	"total": 3
}`

type DiscoverPublicTestSuite struct {
	suite.Suite

	server *httptest.Server
	ctx    context.Context
}

func (s *DiscoverPublicTestSuite) SetupTest() {
	s.ctx = context.Background()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/agent", func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(agentListJSON))
	})

	s.server = httptest.NewServer(mux)
}

func (s *DiscoverPublicTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *DiscoverPublicTestSuite) TestDiscover() {
	tests := []struct {
		name         string
		predicates   []orchestrator.Predicate
		setupServer  func() *httptest.Server
		validateFunc func([]osapi.Agent, error)
	}{
		{
			name:       "No predicates returns all agents",
			predicates: nil,
			validateFunc: func(agents []osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal([]string{"web-01", "web-02", "web-03"}, s.hostnames(agents))
			},
		},
		{
			name: "Filter by OS returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
			},
			validateFunc: func(agents []osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal([]string{"web-01", "web-03"}, s.hostnames(agents))
			},
		},
		{
			name: "Filter by OS and Arch returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Debian"),
				orchestrator.Arch("arm64"),
			},
			validateFunc: func(agents []osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal([]string{"web-02"}, s.hostnames(agents))
			},
		},
		{
			name: "Filter by OS and MinCPU returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
				orchestrator.MinCPU(4),
			},
			validateFunc: func(agents []osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal([]string{"web-01"}, s.hostnames(agents))
			},
		},
		{
			name: "No agents match returns empty slice",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Fedora"),
			},
			validateFunc: func(agents []osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal([]string{}, s.hostnames(agents))
			},
		},
		{
			name: "Returns error when server returns unauthorized",
			validateFunc: func(agents []osapi.Agent, err error) {
				s.Require().Error(err)
				s.Contains(err.Error(), "discover")
				s.Nil(agents)
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(
						w http.ResponseWriter,
						_ *http.Request,
					) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusUnauthorized)
						_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
					}),
				)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := s.server
			if tc.setupServer != nil {
				server = tc.setupServer()
				defer server.Close()
			}

			o := orchestrator.New(server.URL, "test-token")

			tc.validateFunc(o.Discover(s.ctx, tc.predicates...))
		})
	}
}

func (s *DiscoverPublicTestSuite) TestGroupByFact() {
	tests := []struct {
		name         string
		key          string
		predicates   []orchestrator.Predicate
		setupServer  func() *httptest.Server
		validateFunc func(map[string][]osapi.Agent, error)
	}{
		{
			name: "Group by os.distribution",
			key:  "os.distribution",
			validateFunc: func(groups map[string][]osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal(map[string][]string{
					"Ubuntu": {"web-01", "web-03"},
					"Debian": {"web-02"},
				}, s.groupHostnames(groups))
			},
		},
		{
			name: "Group by architecture",
			key:  "architecture",
			validateFunc: func(groups map[string][]osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal(map[string][]string{
					"amd64": {"web-01", "web-03"},
					"arm64": {"web-02"},
				}, s.groupHostnames(groups))
			},
		},
		{
			name: "Skips agents with empty fact value",
			key:  "os.distribution",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(
						w http.ResponseWriter,
						_ *http.Request,
					) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(`{
							"agents": [
								{"hostname":"web-01","os_info":{"distribution":"Ubuntu"}},
								{"hostname":"web-02"}
							],
							"total": 2
						}`))
					}),
				)
			},
			validateFunc: func(groups map[string][]osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal(map[string][]string{
					"Ubuntu": {"web-01"},
				}, s.groupHostnames(groups))
			},
		},
		{
			name: "Group with predicate filter",
			key:  "os.distribution",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
			},
			validateFunc: func(groups map[string][]osapi.Agent, err error) {
				s.Require().NoError(err)
				s.Equal(map[string][]string{
					"Ubuntu": {"web-01", "web-03"},
				}, s.groupHostnames(groups))
			},
		},
		{
			name: "Returns error when discover fails",
			key:  "os.distribution",
			validateFunc: func(groups map[string][]osapi.Agent, err error) {
				s.Require().Error(err)
				s.Contains(err.Error(), "group by fact")
				s.Nil(groups)
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(
						w http.ResponseWriter,
						_ *http.Request,
					) {
						w.WriteHeader(http.StatusUnauthorized)
					}),
				)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := s.server
			if tc.setupServer != nil {
				server = tc.setupServer()
				defer server.Close()
			}

			o := orchestrator.New(server.URL, "test-token")

			tc.validateFunc(o.GroupByFact(s.ctx, tc.key, tc.predicates...))
		})
	}
}

// hostnames reduces agents to the field the discovery cases assert on.
func (s *DiscoverPublicTestSuite) hostnames(
	agents []osapi.Agent,
) []string {
	out := make([]string, 0, len(agents))
	for _, a := range agents {
		out = append(out, a.Hostname)
	}

	return out
}

// groupHostnames applies hostnames to every group, so a case states the
// grouping it expects as plain data.
func (s *DiscoverPublicTestSuite) groupHostnames(
	groups map[string][]osapi.Agent,
) map[string][]string {
	out := make(map[string][]string, len(groups))
	for key, agents := range groups {
		out[key] = s.hostnames(agents)
	}

	return out
}

func TestDiscoverPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(DiscoverPublicTestSuite))
}
