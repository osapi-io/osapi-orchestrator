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

func (s *OpsTestSuite) TestHealthCheckFunc() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
		expectAlive bool
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
			errContains: "unhealthy: status 503",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

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

func (s *OpsTestSuite) TestHealthCheckFuncConnectionError() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	server.Close()

	client := osapi.New(server.URL, "test-token")

	orch := New(server.URL, "test-token")
	step := orch.HealthCheck("_any")
	fn := step.task.Fn()

	result, err := fn(context.Background(), client)

	s.Require().Error(err)
	s.Contains(err.Error(), "health check")
	s.Nil(result)
}

func (s *OpsTestSuite) TestHealthCheckAutoNaming() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	orch := New(server.URL, "test-token")

	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "First health check has no suffix",
			expected: "health-check",
		},
		{
			name:     "Second health check gets counter suffix",
			expected: "health-check-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := orch.HealthCheck("_any")
			s.Equal(tc.expected, step.task.Name())
		})
	}
}

func TestOpsTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OpsTestSuite))
}
