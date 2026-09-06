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

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"
)

type OpsTestSuite struct {
	suite.Suite
}

func (s *OpsTestSuite) TestHealthCheck() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		closeServer bool
		expectErr   bool
		errContains string
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
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			if tc.closeServer {
				server.Close()
			}

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.HealthCheck()
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

func (s *OpsTestSuite) TestNodeHostnameGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with hostname data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get hostname",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeHostnameGet("_any")
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
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
			s.Len(result.HostResults, 1)
			s.Equal("web-01", result.HostResults[0].Hostname)
		})
	}
}

func (s *OpsTestSuite) TestNodeStatusGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with status data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get status",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeStatusGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNodeUptimeGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with uptime data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","uptime":"5d","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get uptime",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeUptimeGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNodeDiskGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with disk data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false,"disks":[{"name":"/","total":100,"used":50,"free":50}]}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get disk",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeDiskGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNodeMemoryGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with memory data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false,"memory":{"total":8192,"used":4096,"free":4096}}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get memory",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeMemoryGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNodeLoadGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with load data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false,"load_average":{"1min":0.5,"5min":0.3,"15min":0.2}}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get load",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeLoadGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNetworkDNSGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with DNS data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","servers":["8.8.8.8"],"changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get dns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NetworkDNSGet("_any", "eth0")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNetworkDNSUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with update result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","status":"updated","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update dns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NetworkDNSUpdate(
				"_any",
				"eth0",
				[]string{"8.8.8.8"},
				[]string{"example.com"},
			)
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNetworkPingDo() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with ping data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false,"packets_sent":3,"packets_received":3}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "ping",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NetworkPingDo("_any", "8.8.8.8")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCommandExec() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with command result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","stdout":"output","exit_code":0,"changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "exec command",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CommandExec("_any", "uptime", "-s")
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCommandShell() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with shell result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","stdout":"hello","exit_code":0,"changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "shell command",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CommandShell("_any", "echo hello")
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestFileDeploy() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with deploy result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"web-01","changed":true}]}`,
					),
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
			errContains: "deploy file",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.FileDeploy("_any", osapi.FileDeployOpts{
				ObjectName:  "config.yaml",
				Path:        "/etc/app/config.yaml",
				ContentType: "raw",
				Mode:        "0644",
			})
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
		})
	}
}

func (s *OpsTestSuite) TestFileStatusGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with status result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"web-01","path":"/etc/app/config.yaml","status":"in-sync","changed":false}]}`,
					),
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
			errContains: "file status",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.FileStatusGet("_any", "/etc/app/config.yaml")
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
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
		})
	}
}

func (s *OpsTestSuite) TestAgentList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
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
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

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
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
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
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

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

func (s *OpsTestSuite) TestFileUpload() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		opts        []UploadOption
		closeServer bool
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success without force",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				s.Empty(r.URL.Query().Get("force"), "force param should not be set")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write(
					[]byte(
						`{"name":"test.txt","sha256":"abc123","size":7,"changed":true,"content_type":"raw"}`,
					),
				)
			}),
		},
		{
			name: "Returns success with force",
			opts: []UploadOption{WithForce()},
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				s.Equal("true", r.URL.Query().Get("force"), "force param should be set")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write(
					[]byte(
						`{"name":"test.txt","sha256":"abc123","size":7,"changed":true,"content_type":"raw"}`,
					),
				)
			}),
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
			errContains: "upload file",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			if tc.closeServer {
				server.Close()
			}

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.FileUpload("test.txt", "raw", []byte("content"), tc.opts...)
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
			s.True(result.Changed)
			s.NotNil(result.Data)
		})
	}
}

func (s *OpsTestSuite) TestFileChanged() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		closeServer bool
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with changed data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write(
						[]byte(
							`{"name":"test.txt","sha256":"different","size":7,"content_type":"raw"}`,
						),
					)

					return
				}
				w.WriteHeader(http.StatusOK)
			}),
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
			errContains: "check file",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			if tc.closeServer {
				server.Close()
			}

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.FileChanged("test.txt", []byte("content"))
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
			s.True(result.Changed)
			s.NotNil(result.Data)
		})
	}
}

// TestOperationStepNames pins the default step name of every operation. These
// names are the keys a caller passes to Report.Decode, so renaming one breaks
// consumers rather than only rearranging internals. The counter suffix applied
// to a repeated name is TestNextOpName's subject, not this table's.
func (s *OpsTestSuite) TestOperationStepNames() {
	tests := []struct {
		name     string
		newFunc  func(o *Orchestrator) *Step
		expected string
	}{
		{
			name:     "HealthCheck",
			newFunc:  func(o *Orchestrator) *Step { return o.HealthCheck() },
			expected: "health-check",
		},
		{
			name:     "NodeHostnameGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeHostnameGet("_any") },
			expected: "get-hostname",
		},
		{
			name:     "NodeHostnameUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeHostnameUpdate("_any", "web-01") },
			expected: "update-hostname",
		},
		{
			name:     "NodeStatusGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeStatusGet("_any") },
			expected: "get-status",
		},
		{
			name:     "NodeUptimeGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeUptimeGet("_any") },
			expected: "get-uptime",
		},
		{
			name:     "NodeDiskGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeDiskGet("_any") },
			expected: "get-disk",
		},
		{
			name:     "NodeMemoryGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeMemoryGet("_any") },
			expected: "get-memory",
		},
		{
			name:     "NodeLoadGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeLoadGet("_any") },
			expected: "get-load",
		},
		{
			name:     "NodeOSGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NodeOSGet("_any") },
			expected: "get-os",
		},
		{
			name:     "NetworkDNSGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NetworkDNSGet("_any", "eth0") },
			expected: "get-dns",
		},
		{
			name:     "NetworkDNSUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.NetworkDNSUpdate("_any", "eth0", nil, nil) },
			expected: "update-dns",
		},
		{
			name:     "NetworkDNSDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.NetworkDNSDelete("_any", "eth0") },
			expected: "delete-dns",
		},
		{
			name:     "NetworkPingDo",
			newFunc:  func(o *Orchestrator) *Step { return o.NetworkPingDo("_any", "1.1.1.1") },
			expected: "ping",
		},
		{
			name:     "InterfaceList",
			newFunc:  func(o *Orchestrator) *Step { return o.InterfaceList("_any") },
			expected: "list-interface",
		},
		{
			name:     "InterfaceGet",
			newFunc:  func(o *Orchestrator) *Step { return o.InterfaceGet("_any", "eth0") },
			expected: "get-interface",
		},
		{
			name:     "InterfaceCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.InterfaceCreate("_any", "eth0", osapi.InterfaceConfigOpts{}) },
			expected: "create-interface",
		},
		{
			name:     "InterfaceUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.InterfaceUpdate("_any", "eth0", osapi.InterfaceConfigOpts{}) },
			expected: "update-interface",
		},
		{
			name:     "InterfaceDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.InterfaceDelete("_any", "eth0") },
			expected: "delete-interface",
		},
		{
			name:     "RouteList",
			newFunc:  func(o *Orchestrator) *Step { return o.RouteList("_any") },
			expected: "list-route",
		},
		{
			name:     "RouteGet",
			newFunc:  func(o *Orchestrator) *Step { return o.RouteGet("_any", "eth0") },
			expected: "get-route",
		},
		{
			name:     "RouteCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.RouteCreate("_any", "eth0", osapi.RouteConfigOpts{}) },
			expected: "create-route",
		},
		{
			name:     "RouteUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.RouteUpdate("_any", "eth0", osapi.RouteConfigOpts{}) },
			expected: "update-route",
		},
		{
			name:     "RouteDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.RouteDelete("_any", "eth0") },
			expected: "delete-route",
		},
		{
			name:     "CommandExec",
			newFunc:  func(o *Orchestrator) *Step { return o.CommandExec("_any", "/usr/bin/uptime") },
			expected: "run-uptime",
		},
		{
			name:     "CommandShell",
			newFunc:  func(o *Orchestrator) *Step { return o.CommandShell("_any", "df -h /") },
			expected: "shell-df",
		},
		{
			name:     "FileDeploy",
			newFunc:  func(o *Orchestrator) *Step { return o.FileDeploy("_any", osapi.FileDeployOpts{}) },
			expected: "deploy-file",
		},
		{
			name:     "FileStatusGet",
			newFunc:  func(o *Orchestrator) *Step { return o.FileStatusGet("_any", "/etc/app.conf") },
			expected: "file-status",
		},
		{
			name:     "FileUndeploy",
			newFunc:  func(o *Orchestrator) *Step { return o.FileUndeploy("_any", "/etc/app.conf") },
			expected: "undeploy-file",
		},
		{
			name:     "FileUpload",
			newFunc:  func(o *Orchestrator) *Step { return o.FileUpload("app.conf", "text/plain", nil) },
			expected: "upload-file",
		},
		{
			name:     "FileChanged",
			newFunc:  func(o *Orchestrator) *Step { return o.FileChanged("app.conf", nil) },
			expected: "check-file",
		},
		{
			name:     "DockerPull",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerPull("_any", osapi.DockerPullOpts{}) },
			expected: "docker-pull",
		},
		{
			name:     "DockerCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerCreate("_any", osapi.DockerCreateOpts{}) },
			expected: "docker-create",
		},
		{
			name:     "DockerStart",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerStart("_any", "container-1") },
			expected: "docker-start",
		},
		{
			name:     "DockerStop",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerStop("_any", "container-1", osapi.DockerStopOpts{}) },
			expected: "docker-stop",
		},
		{
			name:     "DockerRemove",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerRemove("_any", "container-1", nil) },
			expected: "docker-remove",
		},
		{
			name:     "DockerExec",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerExec("_any", "container-1", osapi.DockerExecOpts{}) },
			expected: "docker-exec",
		},
		{
			name:     "DockerInspect",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerInspect("_any", "container-1") },
			expected: "docker-inspect",
		},
		{
			name:     "DockerList",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerList("_any", nil) },
			expected: "docker-list",
		},
		{
			name:     "DockerImageRemove",
			newFunc:  func(o *Orchestrator) *Step { return o.DockerImageRemove("_any", "nginx:latest", nil) },
			expected: "docker-image-remove",
		},
		{
			name:     "CronList",
			newFunc:  func(o *Orchestrator) *Step { return o.CronList("_any") },
			expected: "list-cron",
		},
		{
			name:     "CronGet",
			newFunc:  func(o *Orchestrator) *Step { return o.CronGet("_any", "nightly-backup") },
			expected: "get-cron",
		},
		{
			name:     "CronCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.CronCreate("_any", osapi.CronCreateOpts{}) },
			expected: "create-cron",
		},
		{
			name:     "CronUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.CronUpdate("_any", "nightly-backup", osapi.CronUpdateOpts{}) },
			expected: "update-cron",
		},
		{
			name:     "CronDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.CronDelete("_any", "nightly-backup") },
			expected: "delete-cron",
		},
		{
			name:     "AgentList",
			newFunc:  func(o *Orchestrator) *Step { return o.AgentList() },
			expected: "list-agents",
		},
		{
			name:     "AgentGet",
			newFunc:  func(o *Orchestrator) *Step { return o.AgentGet("web-01") },
			expected: "get-agent",
		},
		{
			name:     "AgentDrain",
			newFunc:  func(o *Orchestrator) *Step { return o.AgentDrain("web-01") },
			expected: "drain-agent",
		},
		{
			name:     "AgentUndrain",
			newFunc:  func(o *Orchestrator) *Step { return o.AgentUndrain("web-01") },
			expected: "undrain-agent",
		},
		{
			name:     "SysctlList",
			newFunc:  func(o *Orchestrator) *Step { return o.SysctlList("_any") },
			expected: "list-sysctl",
		},
		{
			name:     "SysctlGet",
			newFunc:  func(o *Orchestrator) *Step { return o.SysctlGet("_any", "vm.swappiness") },
			expected: "get-sysctl",
		},
		{
			name:     "SysctlCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.SysctlCreate("_any", osapi.SysctlCreateOpts{}) },
			expected: "create-sysctl",
		},
		{
			name:     "SysctlUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.SysctlUpdate("_any", "vm.swappiness", osapi.SysctlUpdateOpts{}) },
			expected: "update-sysctl",
		},
		{
			name:     "SysctlDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.SysctlDelete("_any", "vm.swappiness") },
			expected: "delete-sysctl",
		},
		{
			name:     "NTPGet",
			newFunc:  func(o *Orchestrator) *Step { return o.NTPGet("_any") },
			expected: "get-ntp",
		},
		{
			name:     "NTPCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.NTPCreate("_any", osapi.NtpCreateOpts{}) },
			expected: "create-ntp",
		},
		{
			name:     "NTPUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.NTPUpdate("_any", osapi.NtpUpdateOpts{}) },
			expected: "update-ntp",
		},
		{
			name:     "NTPDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.NTPDelete("_any") },
			expected: "delete-ntp",
		},
		{
			name:     "TimezoneGet",
			newFunc:  func(o *Orchestrator) *Step { return o.TimezoneGet("_any") },
			expected: "get-timezone",
		},
		{
			name:     "TimezoneUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.TimezoneUpdate("_any", osapi.TimezoneUpdateOpts{}) },
			expected: "update-timezone",
		},
		{
			name:     "ServiceList",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceList("_any") },
			expected: "list-service",
		},
		{
			name:     "ServiceGet",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceGet("_any", "nginx") },
			expected: "get-service",
		},
		{
			name:     "ServiceCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceCreate("_any", osapi.ServiceCreateOpts{}) },
			expected: "create-service",
		},
		{
			name:     "ServiceUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceUpdate("_any", "nginx", osapi.ServiceUpdateOpts{}) },
			expected: "update-service",
		},
		{
			name:     "ServiceDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceDelete("_any", "nginx") },
			expected: "delete-service",
		},
		{
			name:     "ServiceStart",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceStart("_any", "nginx") },
			expected: "start-service",
		},
		{
			name:     "ServiceStop",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceStop("_any", "nginx") },
			expected: "stop-service",
		},
		{
			name:     "ServiceRestart",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceRestart("_any", "nginx") },
			expected: "restart-service",
		},
		{
			name:     "ServiceEnable",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceEnable("_any", "nginx") },
			expected: "enable-service",
		},
		{
			name:     "ServiceDisable",
			newFunc:  func(o *Orchestrator) *Step { return o.ServiceDisable("_any", "nginx") },
			expected: "disable-service",
		},
		{
			name:     "PackageList",
			newFunc:  func(o *Orchestrator) *Step { return o.PackageList("_any") },
			expected: "list-package",
		},
		{
			name:     "PackageGet",
			newFunc:  func(o *Orchestrator) *Step { return o.PackageGet("_any", "curl") },
			expected: "get-package",
		},
		{
			name:     "PackageInstall",
			newFunc:  func(o *Orchestrator) *Step { return o.PackageInstall("_any", "curl") },
			expected: "install-package",
		},
		{
			name:     "PackageRemove",
			newFunc:  func(o *Orchestrator) *Step { return o.PackageRemove("_any", "curl") },
			expected: "remove-package",
		},
		{
			name:     "PackageUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.PackageUpdate("_any") },
			expected: "update-package",
		},
		{
			name:     "PackageListUpdates",
			newFunc:  func(o *Orchestrator) *Step { return o.PackageListUpdates("_any") },
			expected: "list-package-updates",
		},
		{
			name:     "UserList",
			newFunc:  func(o *Orchestrator) *Step { return o.UserList("_any") },
			expected: "list-user",
		},
		{
			name:     "UserGet",
			newFunc:  func(o *Orchestrator) *Step { return o.UserGet("_any", "alice") },
			expected: "get-user",
		},
		{
			name:     "UserCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.UserCreate("_any", osapi.UserCreateOpts{}) },
			expected: "create-user",
		},
		{
			name:     "UserUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.UserUpdate("_any", "alice", osapi.UserUpdateOpts{}) },
			expected: "update-user",
		},
		{
			name:     "UserDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.UserDelete("_any", "alice") },
			expected: "delete-user",
		},
		{
			name:     "UserListKeys",
			newFunc:  func(o *Orchestrator) *Step { return o.UserListKeys("_any", "alice") },
			expected: "list-ssh-key",
		},
		{
			name:     "UserAddKey",
			newFunc:  func(o *Orchestrator) *Step { return o.UserAddKey("_any", "alice", osapi.SSHKeyAddOpts{}) },
			expected: "add-ssh-key",
		},
		{
			name:     "UserRemoveKey",
			newFunc:  func(o *Orchestrator) *Step { return o.UserRemoveKey("_any", "alice", "SHA256:abcdef") },
			expected: "remove-ssh-key",
		},
		{
			name:     "UserChangePassword",
			newFunc:  func(o *Orchestrator) *Step { return o.UserChangePassword("_any", "alice", "s3cret") },
			expected: "change-password",
		},
		{
			name:     "GroupList",
			newFunc:  func(o *Orchestrator) *Step { return o.GroupList("_any") },
			expected: "list-group",
		},
		{
			name:     "GroupGet",
			newFunc:  func(o *Orchestrator) *Step { return o.GroupGet("_any", "admins") },
			expected: "get-group",
		},
		{
			name:     "GroupCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.GroupCreate("_any", osapi.GroupCreateOpts{}) },
			expected: "create-group",
		},
		{
			name:     "GroupUpdate",
			newFunc:  func(o *Orchestrator) *Step { return o.GroupUpdate("_any", "admins", osapi.GroupUpdateOpts{}) },
			expected: "update-group",
		},
		{
			name:     "GroupDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.GroupDelete("_any", "admins") },
			expected: "delete-group",
		},
		{
			name:     "CertificateList",
			newFunc:  func(o *Orchestrator) *Step { return o.CertificateList("_any") },
			expected: "list-certificate",
		},
		{
			name:     "CertificateCreate",
			newFunc:  func(o *Orchestrator) *Step { return o.CertificateCreate("_any", osapi.CertificateCreateOpts{}) },
			expected: "create-certificate",
		},
		{
			name: "CertificateUpdate",
			newFunc: func(o *Orchestrator) *Step {
				return o.CertificateUpdate("_any", "internal-ca", osapi.CertificateUpdateOpts{})
			},
			expected: "update-certificate",
		},
		{
			name:     "CertificateDelete",
			newFunc:  func(o *Orchestrator) *Step { return o.CertificateDelete("_any", "internal-ca") },
			expected: "delete-certificate",
		},
		{
			name:     "ProcessList",
			newFunc:  func(o *Orchestrator) *Step { return o.ProcessList("_any") },
			expected: "list-process",
		},
		{
			name:     "ProcessGet",
			newFunc:  func(o *Orchestrator) *Step { return o.ProcessGet("_any", 1234) },
			expected: "get-process",
		},
		{
			name:     "ProcessSignal",
			newFunc:  func(o *Orchestrator) *Step { return o.ProcessSignal("_any", 1234, osapi.ProcessSignalOpts{}) },
			expected: "signal-process",
		},
		{
			name:     "PowerReboot",
			newFunc:  func(o *Orchestrator) *Step { return o.PowerReboot("_any", osapi.PowerOpts{}) },
			expected: "reboot",
		},
		{
			name:     "PowerShutdown",
			newFunc:  func(o *Orchestrator) *Step { return o.PowerShutdown("_any", osapi.PowerOpts{}) },
			expected: "shutdown",
		},
		{
			name:     "LogQuery",
			newFunc:  func(o *Orchestrator) *Step { return o.LogQuery("_any", osapi.LogQueryOpts{}) },
			expected: "query-log",
		},
		{
			name:     "LogSources",
			newFunc:  func(o *Orchestrator) *Step { return o.LogSources("_any") },
			expected: "list-log-sources",
		},
		{
			name:     "LogQueryUnit",
			newFunc:  func(o *Orchestrator) *Step { return o.LogQueryUnit("_any", "nginx.service", osapi.LogQueryOpts{}) },
			expected: "query-log-unit",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New("http://orchestrator.invalid", "test-token")

			s.Equal(tc.expected, tc.newFunc(orch).task.Name())
		})
	}
}

func (s *OpsTestSuite) TestCommandError() {
	tests := []struct {
		name         string
		result       osapi.CommandResult
		validateFunc func(got string)
	}{
		{
			name:   "returns error string when set",
			result: osapi.CommandResult{Error: "connection refused"},
			validateFunc: func(got string) {
				s.Equal("connection refused", got)
			},
		},
		{
			name:   "returns exit code when non-zero",
			result: osapi.CommandResult{ExitCode: 127},
			validateFunc: func(got string) {
				s.Equal("exit code 127", got)
			},
		},
		{
			name:   "returns empty string on success",
			result: osapi.CommandResult{ExitCode: 0},
			validateFunc: func(got string) {
				s.Empty(got)
			},
		},
		{
			name: "error takes precedence over exit code",
			result: osapi.CommandResult{
				Error:    "timeout",
				ExitCode: 1,
			},
			validateFunc: func(got string) {
				s.Equal("timeout", got)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(commandError(tt.result))
		})
	}
}

func (s *OpsTestSuite) TestDockerPull() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with pull result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","image_id":"sha256:abc","tag":"latest","size":1024,"changed":true}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker pull",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerPull("_any", osapi.DockerPullOpts{Image: "nginx:latest"})
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
			s.Len(result.HostResults, 1)
			s.Equal("h1", result.HostResults[0].Hostname)
		})
	}
}

func (s *OpsTestSuite) TestDockerCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with create result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","id":"c1","name":"web","image":"nginx","state":"created","changed":true}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker create",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerCreate("_any", osapi.DockerCreateOpts{Image: "nginx"})
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
			s.Len(result.HostResults, 1)
			s.Equal("h1", result.HostResults[0].Hostname)
		})
	}
}

func (s *OpsTestSuite) TestDockerStart() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with start result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","id":"c1","message":"done","changed":true}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker start",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerStart("_any", "c1")
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestDockerStop() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with stop result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","id":"c1","message":"done","changed":true}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker stop",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerStop("_any", "c1", osapi.DockerStopOpts{Timeout: 10})
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestDockerRemove() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with remove result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","id":"c1","message":"done","changed":true}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker remove",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerRemove(
				"_any",
				"c1",
				&osapi.DockerRemoveParams{Force: true},
			)
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestDockerExec() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with exec result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","stdout":"output","stderr":"","exit_code":0,"changed":true}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker exec",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerExec(
				"_any",
				"c1",
				osapi.DockerExecOpts{Command: []string{"echo", "hello"}},
			)
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestDockerInspect() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with inspect result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","id":"c1","name":"web","image":"nginx","state":"running"}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker inspect",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerInspect("_any", "c1")
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
			s.Len(result.HostResults, 1)
			s.Equal("h1", result.HostResults[0].Hostname)
		})
	}
}

func (s *OpsTestSuite) TestDockerList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with list result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"h1","containers":[{"id":"c1","name":"web","image":"nginx","state":"running"}]}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker list",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerList("_any", &osapi.DockerListParams{State: "running"})
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
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestDockerImageRemove() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with image remove result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"00000000-0000-0000-0000-000000000001","results":[{"hostname":"h1","id":"nginx:latest","message":"Image removed successfully","changed":true}]}`,
					),
				)
			}),
		},
		{
			name: "Returns error on server error",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}),
			expectErr:   true,
			errContains: "docker image remove",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.DockerImageRemove(
				"_any",
				"nginx:latest",
				&osapi.DockerImageRemoveParams{Force: true},
			)
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestFileUndeploy() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with undeploy result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"web-01","changed":true}]}`,
					),
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
			errContains: "undeploy file",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.FileUndeploy("_any", "/etc/app/config.yaml")
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
		})
	}
}

func (s *OpsTestSuite) TestNodeHostnameUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with hostname update result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"web-01","changed":true}]}`,
					),
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
			errContains: "update hostname",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeHostnameUpdate("_any", "new-hostname")
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
		})
	}
}

func (s *OpsTestSuite) TestNodeOSGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with OS info data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false,"name":"Ubuntu","version":"22.04"}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get os",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NodeOSGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestAgentDrain() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with drain result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message":"agent drained"}`))
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
			errContains: "drain agent",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.AgentDrain("web-01")
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

func (s *OpsTestSuite) TestAgentUndrain() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with undrain result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message":"agent undrained"}`))
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
			errContains: "undrain agent",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.AgentUndrain("web-01")
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

func (s *OpsTestSuite) TestCronList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with cron list",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list cron",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CronList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCronGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with cron entry",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get cron",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CronGet("_any", "backup")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCronCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with cron create result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"web-01","changed":true}]}`,
					),
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
			errContains: "create cron",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CronCreate("_any", osapi.CronCreateOpts{
				Name:     "backup",
				Object:   "backup.sh",
				Schedule: "0 2 * * *",
			})
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
		})
	}
}

func (s *OpsTestSuite) TestCronUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with cron update result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"web-01","changed":true}]}`,
					),
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
			errContains: "update cron",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CronUpdate("_any", "backup", osapi.CronUpdateOpts{
				Schedule: "0 3 * * *",
			})
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
		})
	}
}

func (s *OpsTestSuite) TestCronDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with cron delete result",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","results":[{"hostname":"web-01","changed":true}]}`,
					),
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
			errContains: "delete cron",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CronDelete("_any", "backup")
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
			s.True(result.Changed)
			s.NotNil(result.Data)
			s.Equal("550e8400-e29b-41d4-a716-446655440000", result.JobID)
		})
	}
}

func (s *OpsTestSuite) TestNetworkDNSDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with dns delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete dns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NetworkDNSDelete("_any", "eth0")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestInterfaceList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with interface list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list interfaces",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.InterfaceList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestInterfaceGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with interface data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get interface",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.InterfaceGet("_any", "eth0")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestInterfaceCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with interface create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create interface",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.InterfaceCreate("_any", "eth0", osapi.InterfaceConfigOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestInterfaceUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with interface update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update interface",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.InterfaceUpdate("_any", "eth0", osapi.InterfaceConfigOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestInterfaceDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with interface delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete interface",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.InterfaceDelete("_any", "eth0")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestRouteList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with route list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list routes",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.RouteList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestRouteGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with route data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get route",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.RouteGet("_any", "eth0")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestRouteCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with route create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create route",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.RouteCreate("_any", "eth0", osapi.RouteConfigOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestRouteUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with route update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update route",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.RouteUpdate("_any", "eth0", osapi.RouteConfigOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestRouteDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with route delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete route",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.RouteDelete("_any", "eth0")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestSysctlList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with sysctl list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list sysctl",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.SysctlList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestSysctlGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with sysctl data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get sysctl",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.SysctlGet("_any", "net.ipv4.ip_forward")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestSysctlCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with sysctl create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create sysctl",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.SysctlCreate("_any", osapi.SysctlCreateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestSysctlUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with sysctl update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update sysctl",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.SysctlUpdate("_any", "net.ipv4.ip_forward", osapi.SysctlUpdateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestSysctlDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with sysctl delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete sysctl",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.SysctlDelete("_any", "net.ipv4.ip_forward")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNTPGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with ntp data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get ntp",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NTPGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNTPCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with ntp create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create ntp",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NTPCreate("_any", osapi.NtpCreateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNTPUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with ntp update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update ntp",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NTPUpdate("_any", osapi.NtpUpdateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestNTPDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with ntp delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete ntp",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.NTPDelete("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestTimezoneGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with timezone data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get timezone",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.TimezoneGet("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestTimezoneUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with timezone update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update timezone",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.TimezoneUpdate("_any", osapi.TimezoneUpdateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list services",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceGet("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceCreate("_any", osapi.ServiceCreateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceUpdate("_any", "nginx", osapi.ServiceUpdateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceDelete("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceStart() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service start data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "start service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceStart("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceStop() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service stop data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "stop service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceStop("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceRestart() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service restart data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "restart service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceRestart("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceEnable() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service enable data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "enable service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceEnable("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestServiceDisable() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with service disable data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "disable service",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ServiceDisable("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPackageList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with package list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list packages",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PackageList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPackageGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with package data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get package",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PackageGet("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPackageInstall() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with package install data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "install package",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PackageInstall("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPackageRemove() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with package remove data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "remove package",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PackageRemove("_any", "nginx")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPackageUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with package update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update packages",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PackageUpdate("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPackageListUpdates() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with package updates data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list package updates",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PackageListUpdates("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with user list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list users",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with user data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get user",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserGet("_any", "admin")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with user create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create user",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserCreate("_any", osapi.UserCreateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with user update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update user",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserUpdate("_any", "admin", osapi.UserUpdateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with user delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete user",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserDelete("_any", "admin")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserListKeys() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with ssh keys data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list ssh keys",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserListKeys("_any", "admin")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserAddKey() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with add key data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "add ssh key",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserAddKey("_any", "admin", osapi.SSHKeyAddOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserRemoveKey() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with remove key data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "remove ssh key",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserRemoveKey("_any", "admin", "SHA256:abc123")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestUserChangePassword() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with change password data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "change password",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.UserChangePassword("_any", "admin", "newpass")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestGroupList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with group list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list groups",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.GroupList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestGroupGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with group data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get group",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.GroupGet("_any", "admins")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestGroupCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with group create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create group",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.GroupCreate("_any", osapi.GroupCreateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestGroupUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with group update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update group",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.GroupUpdate("_any", "admins", osapi.GroupUpdateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestGroupDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with group delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete group",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.GroupDelete("_any", "admins")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCertificateList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with certificate list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list certificates",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CertificateList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCertificateCreate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with certificate create data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "create certificate",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CertificateCreate("_any", osapi.CertificateCreateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCertificateUpdate() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with certificate update data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "update certificate",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CertificateUpdate("_any", "my-ca", osapi.CertificateUpdateOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestCertificateDelete() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with certificate delete data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "delete certificate",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.CertificateDelete("_any", "my-ca")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestProcessList() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with process list data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list processes",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ProcessList("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestProcessGet() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with process data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "get process",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ProcessGet("_any", 1234)
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestProcessSignal() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with process signal data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "signal process",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.ProcessSignal("_any", 1234, osapi.ProcessSignalOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPowerReboot() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with reboot data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "reboot",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PowerReboot("_any", osapi.PowerOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestPowerShutdown() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with shutdown data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":true}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "shutdown",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.PowerShutdown("_any", osapi.PowerOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestLogQuery() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with log query data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "query log",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.LogQuery("_any", osapi.LogQueryOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestLogSources() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with log sources data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "list log sources",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.LogSources("_any")
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func (s *OpsTestSuite) TestLogQueryUnit() {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		errContains string
	}{
		{
			name: "Returns success with log query unit data",
			handler: http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(
					[]byte(
						`{"results":[{"hostname":"web-01","changed":false}],"job_id":"550e8400-e29b-41d4-a716-446655440000"}`,
					),
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
			errContains: "query log unit",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := osapi.New(server.URL, "test-token")

			orch := New(server.URL, "test-token")
			step := orch.LogQueryUnit("_any", "nginx.service", osapi.LogQueryOpts{})
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
			s.NotNil(result.Data)
			s.Len(result.HostResults, 1)
		})
	}
}

func TestOpsTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OpsTestSuite))
}
