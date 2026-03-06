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

package orchestrator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osapi-io/osapi-sdk/pkg/osapi"
	"github.com/stretchr/testify/suite"
)

type OpsTestSuite struct {
	suite.Suite
}

func (s *OpsTestSuite) TestHealthCheck() {
	tests := []struct {
		name        string
		newOrch     bool
		handler     http.HandlerFunc
		closeServer bool
		expectErr   bool
		errContains string
		expectAlive bool
		expectName  string
	}{
		{
			name: "Returns success on 200",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"ok"}`))
			}),
			expectAlive: true,
		},
		{
			name: "Returns error on non-200 status",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(`{"status":"unhealthy"}`))
			}),
			expectErr:   true,
			errContains: "health check",
		},
		{
			name: "Returns error on connection failure",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			}),
			closeServer: true,
			expectErr:   true,
			errContains: "health check",
		},
		{
			name:    "First name has no suffix",
			newOrch: true,
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			}),
			expectName: "health-check",
		},
		{
			name: "Duplicate name gets counter suffix",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			}),
			expectName: "health-check-2",
		},
	}

	var sharedOrch *Orchestrator

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			if tc.closeServer {
				server.Close()
			}

			if tc.newOrch || sharedOrch == nil {
				sharedOrch = New(server.URL, "test-token")
			}

			if tc.expectName != "" {
				step := sharedOrch.HealthCheck("_any")
				s.Equal(tc.expectName, step.task.Name())

				return
			}

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.HealthCheck("_any")
			fn := step.task.Fn()
			s.Require().NotNil(fn)

			result, fnErr := fn(context.Background(), client)

			if tc.expectErr {
				s.Require().Error(fnErr)
				s.Contains(fnErr.Error(), tc.errContains)
				s.Nil(result)

				return
			}

			s.Require().NoError(fnErr)
			s.False(result.Changed)
		})
	}
}

func (s *OpsTestSuite) TestAgentList() {
	tests := []struct {
		name        string
		newOrch     bool
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
		expectName  string
	}{
		{
			name: "Returns success with agent data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(`{"agents":[{"hostname":"web-01","status":"Ready"}],"total":1}`),
				)
			}),
		},
		{
			name: "Returns error on auth failure",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
			}),
			expectErr:   true,
			errContains: "list agents",
		},
		{
			name:    "First name has no suffix",
			newOrch: true,
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			}),
			expectName: "list-agents",
		},
		{
			name: "Duplicate name gets counter suffix",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			}),
			expectName: "list-agents-2",
		},
	}

	var sharedOrch *Orchestrator

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			if tc.newOrch || sharedOrch == nil {
				sharedOrch = New(server.URL, "test-token")
			}

			if tc.expectName != "" {
				step := sharedOrch.AgentList()
				s.Equal(tc.expectName, step.task.Name())

				return
			}

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.AgentList()
			fn := step.task.Fn()
			s.Require().NotNil(fn)

			result, fnErr := fn(context.Background(), client)

			if tc.expectErr {
				s.Require().Error(fnErr)
				s.Contains(fnErr.Error(), tc.errContains)
				s.Nil(result)

				return
			}

			s.Require().NoError(fnErr)
			s.False(result.Changed)
			s.NotNil(result.Data)
		})
	}
}

func (s *OpsTestSuite) TestAgentGet() {
	tests := []struct {
		name        string
		newOrch     bool
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
		expectName  string
	}{
		{
			name: "Returns success with agent details",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(`{"hostname":"web-01","status":"Ready","architecture":"amd64"}`),
				)
			}),
		},
		{
			name: "Returns error on not found",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":"agent not found"}`))
			}),
			expectErr:   true,
			errContains: "get agent",
		},
		{
			name:    "First name has no suffix",
			newOrch: true,
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			}),
			expectName: "get-agent",
		},
		{
			name: "Duplicate name gets counter suffix",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			}),
			expectName: "get-agent-2",
		},
	}

	var sharedOrch *Orchestrator

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			if tc.newOrch || sharedOrch == nil {
				sharedOrch = New(server.URL, "test-token")
			}

			if tc.expectName != "" {
				step := sharedOrch.AgentGet("web-01")
				s.Equal(tc.expectName, step.task.Name())

				return
			}

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.AgentGet("web-01")
			fn := step.task.Fn()
			s.Require().NotNil(fn)

			result, fnErr := fn(context.Background(), client)

			if tc.expectErr {
				s.Require().Error(fnErr)
				s.Contains(fnErr.Error(), tc.errContains)
				s.Nil(result)

				return
			}

			s.Require().NoError(fnErr)
			s.False(result.Changed)
			s.NotNil(result.Data)
		})
	}
}

func (s *OpsTestSuite) TestMustRawToMap() {
	tests := []struct {
		name   string
		input  []byte
		panics bool
	}{
		{
			name:  "Valid JSON returns map",
			input: []byte(`{"key":"value"}`),
		},
		{
			name:   "Invalid JSON panics",
			input:  []byte(`not json`),
			panics: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			if tc.panics {
				s.Panics(func() { mustRawToMap(tc.input) })

				return
			}

			data := mustRawToMap(tc.input)
			s.NotNil(data)
			s.Equal("value", data["key"])
		})
	}
}

func TestOpsTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OpsTestSuite))
}
