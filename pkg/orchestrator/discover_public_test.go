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
	mux.HandleFunc("/agent", func(
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
		name       string
		predicates []orchestrator.Predicate
		expected   []string
	}{
		{
			name:       "No predicates returns all agents",
			predicates: nil,
			expected:   []string{"web-01", "web-02", "web-03"},
		},
		{
			name: "Filter by OS returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
			},
			expected: []string{"web-01", "web-03"},
		},
		{
			name: "Filter by OS and Arch returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Debian"),
				orchestrator.Arch("arm64"),
			},
			expected: []string{"web-02"},
		},
		{
			name: "Filter by OS and MinCPU returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
				orchestrator.MinCPU(4),
			},
			expected: []string{"web-01"},
		},
		{
			name: "No agents match returns empty slice",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Fedora"),
			},
			expected: []string{},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			o := orchestrator.New(s.server.URL, "test-token")

			agents, err := o.Discover(s.ctx, tc.predicates...)

			s.Require().NoError(err)

			hostnames := make([]string, 0, len(agents))
			for _, a := range agents {
				hostnames = append(hostnames, a.Hostname)
			}

			s.Equal(tc.expected, hostnames)
		})
	}
}

func (s *DiscoverPublicTestSuite) TestDiscoverError() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
		}),
	)
	defer server.Close()

	o := orchestrator.New(server.URL, "bad-token")

	agents, err := o.Discover(s.ctx)

	s.Error(err)
	s.Nil(agents)
	s.Contains(err.Error(), "discover")
}

func TestDiscoverPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(DiscoverPublicTestSuite))
}
