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
			"facts": {"env": "prod", "uptime_days": 12}
		},
		{
			"hostname": "web-02",
			"status": "Ready",
			"architecture": "arm64",
			"cpu_count": 4,
			"os_info": {"distribution": "Debian", "version": "12"},
			"memory": {"total": 8000, "used": 4000, "free": 4000},
			"facts": {"env": "staging", "uptime_days": 4}
		},
		{
			"hostname": "web-03",
			"status": "Ready",
			"architecture": "amd64",
			"cpu_count": 2,
			"os_info": {"distribution": "Ubuntu", "version": "20.04"},
			"memory": {"total": 4000, "used": 2000, "free": 2000},
			"facts": {"env": "prod", "uptime_days": 12}
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
		name        string
		predicates  []orchestrator.Predicate
		expected    []string
		expectErr   bool
		errContains string
		setupServer func() *httptest.Server
	}{
		{
			name:       "No predicates returns all agents",
			predicates: nil,
			expected:   []string{hostWeb01, hostWeb02, hostWeb03},
		},
		{
			name: "Filter by OS returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS(distributionUbuntu),
			},
			expected: []string{hostWeb01, hostWeb03},
		},
		{
			name: "Filter by OS and Arch returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Debian"),
				orchestrator.Arch("arm64"),
			},
			expected: []string{hostWeb02},
		},
		{
			name: "Filter by OS and MinCPU returns matching agents",
			predicates: []orchestrator.Predicate{
				orchestrator.OS(distributionUbuntu),
				orchestrator.MinCPU(4),
			},
			expected: []string{hostWeb01},
		},
		{
			name: "No agents match returns empty slice",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Fedora"),
			},
			expected: []string{},
		},
		{
			name:        "Returns error when server returns unauthorized",
			expectErr:   true,
			errContains: "discover",
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

			agents, err := o.Discover(s.ctx, tc.predicates...)

			if tc.expectErr {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)
				s.Nil(agents)

				return
			}

			s.Require().NoError(err)

			hostnames := make([]string, 0, len(agents))
			for _, a := range agents {
				hostnames = append(hostnames, a.Hostname)
			}

			s.Equal(tc.expected, hostnames)
		})
	}
}

func (s *DiscoverPublicTestSuite) TestGroupByFact() {
	tests := []struct {
		name        string
		key         string
		predicates  []orchestrator.Predicate
		expected    map[string][]string
		expectErr   bool
		errContains string
		setupServer func() *httptest.Server
	}{
		{
			name: "Group by os.distribution",
			key:  factOSDistribution,
			expected: map[string][]string{
				distributionUbuntu: {hostWeb01, hostWeb03},
				"Debian":           {hostWeb02},
			},
		},
		{
			// Facts are free-form JSON, so a value that is not a
			// string cannot be grouped on and reads as absent.
			name:     "Group by a fact whose value is not a string",
			key:      "uptime_days",
			expected: map[string][]string{},
		},
		{
			name: "Group by architecture",
			key:  "architecture",
			expected: map[string][]string{
				"amd64": {hostWeb01, hostWeb03},
				"arm64": {hostWeb02},
			},
		},
		{
			name: "Skips agents with empty fact value",
			key:  factOSDistribution,
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
			expected: map[string][]string{
				distributionUbuntu: {hostWeb01},
			},
		},
		{
			name: "Group with predicate filter",
			key:  factOSDistribution,
			predicates: []orchestrator.Predicate{
				orchestrator.OS(distributionUbuntu),
			},
			expected: map[string][]string{
				distributionUbuntu: {hostWeb01, hostWeb03},
			},
		},
		{
			name:        "Returns error when discover fails",
			key:         factOSDistribution,
			expectErr:   true,
			errContains: "group by fact",
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

			groups, err := o.GroupByFact(s.ctx, tc.key, tc.predicates...)

			if tc.expectErr {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)
				s.Nil(groups)

				return
			}

			s.Require().NoError(err)
			s.Len(groups, len(tc.expected))

			for key, expectedHostnames := range tc.expected {
				agents, ok := groups[key]
				s.True(ok, "expected group %q", key)

				hostnames := make([]string, 0, len(agents))
				for _, a := range agents {
					hostnames = append(hostnames, a.Hostname)
				}

				s.Equal(expectedHostnames, hostnames)
			}
		})
	}
}

func TestDiscoverPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(DiscoverPublicTestSuite))
}
