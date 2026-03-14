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

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
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
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","hostname":"web-01","changed":true}`,
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
						`{"job_id":"550e8400-e29b-41d4-a716-446655440000","hostname":"web-01","path":"/etc/app/config.yaml","status":"in-sync","changed":false}`,
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

func (s *OpsTestSuite) TestHealthCheckNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "health-check",
			secondName: "health-check-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.HealthCheck(), orch.HealthCheck()

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNodeHostnameGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-hostname",
			secondName: "get-hostname-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NodeHostnameGet("_any"), orch.NodeHostnameGet("_any")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNodeStatusGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-status",
			secondName: "get-status-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NodeStatusGet("_any"), orch.NodeStatusGet("_any")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNodeUptimeGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-uptime",
			secondName: "get-uptime-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NodeUptimeGet("_any"), orch.NodeUptimeGet("_any")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNodeDiskGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-disk",
			secondName: "get-disk-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NodeDiskGet("_any"), orch.NodeDiskGet("_any")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNodeMemoryGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-memory",
			secondName: "get-memory-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NodeMemoryGet("_any"), orch.NodeMemoryGet("_any")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNodeLoadGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-load",
			secondName: "get-load-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NodeLoadGet("_any"), orch.NodeLoadGet("_any")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNetworkDNSGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-dns",
			secondName: "get-dns-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NetworkDNSGet("_any", "eth0"), orch.NetworkDNSGet("_any", "eth0")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNetworkDNSUpdateNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "update-dns",
			secondName: "update-dns-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NetworkDNSUpdate("_any", "eth0", []string{"8.8.8.8"}, nil),
				orch.NetworkDNSUpdate("_any", "eth0", []string{"8.8.8.8"}, nil)

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestNetworkPingDoNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "ping",
			secondName: "ping-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.NetworkPingDo("_any", "8.8.8.8"), orch.NetworkPingDo("_any", "8.8.8.8")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestCommandExecNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "run-uptime",
			secondName: "run-uptime-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.CommandExec("_any", "uptime"), orch.CommandExec("_any", "uptime")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestCommandShellNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "shell-echo hello",
			secondName: "shell-echo hello-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.CommandShell("_any", "echo hello"),
				orch.CommandShell("_any", "echo hello")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestFileDeployNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "deploy-file",
			secondName: "deploy-file-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			opts := osapi.FileDeployOpts{ObjectName: "f", Path: "/p", ContentType: "raw"}
			first, second := orch.FileDeploy("_any", opts), orch.FileDeploy("_any", opts)

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestFileStatusGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "file-status",
			secondName: "file-status-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.FileStatusGet("_any", "/p"), orch.FileStatusGet("_any", "/p")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestAgentListNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "list-agents",
			secondName: "list-agents-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.AgentList(), orch.AgentList()

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestAgentGetNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "get-agent",
			secondName: "get-agent-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.AgentGet("web-01"), orch.AgentGet("web-01")

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestFileUploadNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "upload-file",
			secondName: "upload-file-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.FileUpload("test.txt", "raw", []byte("content")),
				orch.FileUpload("test.txt", "raw", []byte("content"))

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (s *OpsTestSuite) TestFileChangedNameCounter() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name       string
		firstName  string
		secondName string
	}{
		{
			name:       "Duplicate name gets counter suffix",
			firstName:  "check-file",
			secondName: "check-file-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := orch.FileChanged("test.txt", []byte("content")),
				orch.FileChanged("test.txt", []byte("content"))

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
		})
	}
}

func (suite *OpsTestSuite) TestCommandError() {
	tests := []struct {
		name       string
		result     osapi.CommandResult
		validateFn func(string)
	}{
		{
			name:   "returns error string when set",
			result: osapi.CommandResult{Error: "connection refused"},
			validateFn: func(s string) {
				suite.Equal("connection refused", s)
			},
		},
		{
			name:   "returns exit code when non-zero",
			result: osapi.CommandResult{ExitCode: 127},
			validateFn: func(s string) {
				suite.Equal("exit code 127", s)
			},
		},
		{
			name:   "returns empty string on success",
			result: osapi.CommandResult{ExitCode: 0},
			validateFn: func(s string) {
				suite.Empty(s)
			},
		},
		{
			name: "error takes precedence over exit code",
			result: osapi.CommandResult{
				Error:    "timeout",
				ExitCode: 1,
			},
			validateFn: func(s string) {
				suite.Equal("timeout", s)
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.validateFn(commandError(tt.result))
		})
	}
}

func TestOpsTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OpsTestSuite))
}
