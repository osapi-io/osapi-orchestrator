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

func (s *OpsTestSuite) TestOperationNameCounter() {
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
		callOp     func(orch *Orchestrator) (*Step, *Step)
		firstName  string
		secondName string
	}{
		{
			name: "HealthCheck",
			callOp: func(orch *Orchestrator) (*Step, *Step) {
				return orch.HealthCheck("_any"), orch.HealthCheck("_any")
			},
			firstName:  "health-check",
			secondName: "health-check-2",
		},
		{
			name: "AgentList",
			callOp: func(orch *Orchestrator) (*Step, *Step) {
				return orch.AgentList(), orch.AgentList()
			},
			firstName:  "list-agents",
			secondName: "list-agents-2",
		},
		{
			name: "AgentGet",
			callOp: func(orch *Orchestrator) (*Step, *Step) {
				return orch.AgentGet("web-01"), orch.AgentGet("web-01")
			},
			firstName:  "get-agent",
			secondName: "get-agent-2",
		},
		{
			name: "FileUpload",
			callOp: func(orch *Orchestrator) (*Step, *Step) {
				return orch.FileUpload("test.txt", "raw", []byte("content")),
					orch.FileUpload("test.txt", "raw", []byte("content"))
			},
			firstName:  "upload-file",
			secondName: "upload-file-2",
		},
		{
			name: "FileChanged",
			callOp: func(orch *Orchestrator) (*Step, *Step) {
				return orch.FileChanged("test.txt", []byte("content")),
					orch.FileChanged("test.txt", []byte("content"))
			},
			firstName:  "check-file",
			secondName: "check-file-2",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			first, second := tc.callOp(orch)

			s.Equal(tc.firstName, first.task.Name())
			s.Equal(tc.secondName, second.task.Name())
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
