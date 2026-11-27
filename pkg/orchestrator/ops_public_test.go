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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

// hostRow is the answer of the fake controller for a host op: one host
// that changed, and a job.
const hostRow = `{"results":[{"hostname":"web-01","status":"ok","changed":true}],"job_id":"job-1"}`

type OpsPublicTestSuite struct {
	suite.Suite
}

// request is what the fake controller received.
type request struct {
	method string
	path   string
	body   string
}

// serve starts a fake controller that answers each request with status
// and body, and records the requests.
func (s *OpsPublicTestSuite) serve(
	status int,
	body string,
) (*httptest.Server, *[]request) {
	var got []request

	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		b, _ := io.ReadAll(r.Body)
		got = append(got, request{method: r.Method, path: r.URL.Path, body: string(b)})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	s.T().Cleanup(server.Close)

	return server, &got
}

// run builds one step on an orchestrator that talks to server, and runs
// the plan.
func (s *OpsPublicTestSuite) run(
	server *httptest.Server,
	build func(*orchestrator.Orchestrator) *orchestrator.Step,
) (*orchestrator.Report, error) {
	o := orchestrator.New(server.URL, "test-token")
	build(o)

	return o.Run(context.Background())
}

// TestControllerOps runs the ops that call the controller itself: health,
// the Object Store and the agents.
func (s *OpsPublicTestSuite) TestControllerOps() {
	tests := []struct {
		name      string
		status    int
		body      string
		build     func(*orchestrator.Orchestrator) *orchestrator.Step
		method    string
		path      string
		step      string
		changed   bool
		data      string
		errSubstr string
	}{
		{
			name:   "HealthCheck",
			status: http.StatusOK,
			body:   `{"status":"ok"}`,
			build:  func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.HealthCheck() },
			method: http.MethodGet,
			path:   "/api/health",
			step:   "health-check",
		},
		{
			name:      "HealthCheck refused",
			status:    http.StatusServiceUnavailable,
			body:      `{"error":"down"}`,
			build:     func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.HealthCheck() },
			errSubstr: "health check",
		},
		{
			name:   "FileUpload with force replaces the object",
			status: http.StatusOK,
			body:   `{"name":"app.conf","changed":true}`,
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.FileUpload("app.conf", "raw", []byte("x"), orchestrator.WithForce())
			},
			method:  http.MethodPut,
			path:    "/api/file/app.conf",
			step:    "upload-file",
			changed: true,
			data:    "app.conf",
		},
		{
			name:   "FileUpload refused",
			status: http.StatusForbidden,
			body:   `{"error":"forbidden"}`,
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.FileUpload("app.conf", "raw", []byte("x"), orchestrator.WithForce())
			},
			errSubstr: "upload file app.conf",
		},
		{
			name:   "FileChanged",
			status: http.StatusOK,
			body:   `{"name":"app.conf","changed":true}`,
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.FileChanged("app.conf", []byte("x"))
			},
			method:  http.MethodGet,
			path:    "/api/file/app.conf",
			step:    "check-file",
			changed: true,
			data:    "app.conf",
		},
		{
			name:   "FileChanged refused",
			status: http.StatusForbidden,
			body:   `{"error":"forbidden"}`,
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.FileChanged("app.conf", []byte("x"))
			},
			errSubstr: "check file app.conf",
		},
		{
			name:   "AgentList",
			status: http.StatusOK,
			body:   `{"items":[{"hostname":"web-01"}],"total_items":1}`,
			build:  func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.AgentList() },
			method: http.MethodGet,
			path:   "/api/agent",
			step:   "list-agents",
		},
		{
			name:      "AgentList refused",
			status:    http.StatusForbidden,
			body:      `{"error":"forbidden"}`,
			build:     func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.AgentList() },
			errSubstr: "list agents",
		},
		{
			name:   "AgentGet",
			status: http.StatusOK,
			body:   `{"hostname":"web-01"}`,
			build:  func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.AgentGet("web-01") },
			method: http.MethodGet,
			path:   "/api/agent/web-01",
			step:   "get-agent",
			data:   "web-01",
		},
		{
			name:      "AgentGet refused",
			status:    http.StatusNotFound,
			body:      `{"error":"no agent"}`,
			build:     func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.AgentGet("web-01") },
			errSubstr: "get agent web-01",
		},
		{
			name:   "AgentDrain",
			status: http.StatusNoContent,
			build:  func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.AgentDrain("web-01") },
			method: http.MethodPost,
			path:   "/api/agent/web-01/drain",
			step:   "drain-agent",
		},
		{
			name:   "AgentUndrain",
			status: http.StatusNoContent,
			build:  func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.AgentUndrain("web-01") },
			method: http.MethodPost,
			path:   "/api/agent/web-01/undrain",
			step:   "undrain-agent",
		},
		{
			name:      "AgentDrain refused",
			status:    http.StatusNotFound,
			body:      `{"error":"no agent"}`,
			build:     func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.AgentDrain("web-01") },
			errSubstr: "drain-agent web-01",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server, got := s.serve(tc.status, tc.body)

			report, err := s.run(server, tc.build)
			if tc.errSubstr != "" {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errSubstr)

				return
			}

			s.Require().NoError(err)
			last := (*got)[len(*got)-1]
			s.Equal(tc.method, last.method)
			s.Equal(tc.path, last.path)
			s.Require().Len(report.Tasks, 1)
			s.Equal(tc.step, report.Tasks[0].Name)
			s.Equal(tc.changed, report.Tasks[0].Changed)

			if tc.data != "" {
				raw, _ := json.Marshal(report.Tasks[0].Data)
				s.Contains(string(raw), tc.data)
			}
		})
	}
}

// TestFileUploadSkipsTheSameContent covers the upload with no force: the
// SDK reads the stored object first, and a matching digest uploads nothing.
func (s *OpsPublicTestSuite) TestFileUploadSkipsTheSameContent() {
	// The SHA-256 of "x".
	const sum = "2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881"

	server, got := s.serve(http.StatusOK, `{"name":"app.conf","sha256":"`+sum+`"}`)

	report, err := s.run(server, func(o *orchestrator.Orchestrator) *orchestrator.Step {
		return o.FileUpload("app.conf", "raw", []byte("x"))
	})
	s.Require().NoError(err)
	s.Require().Len(*got, 1, "only the read of the stored object")
	s.Equal(http.MethodGet, (*got)[0].method)
	s.False(report.Tasks[0].Changed)
}

func TestOpsPublicTestSuite(t *testing.T) {
	suite.Run(t, new(OpsPublicTestSuite))
}
