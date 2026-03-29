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

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
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

func (s *OpsPublicTestSuite) TestHealthCheck() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.HealthCheck()
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

func (s *OpsPublicTestSuite) TestNodeHostnameUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameUpdate("_any", "new-hostname")
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

func (s *OpsPublicTestSuite) TestNodeOSGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeOSGet("_any")
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

func (s *OpsPublicTestSuite) TestNodeHostnameGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeHostnameGet("_any")
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

func (s *OpsPublicTestSuite) TestNodeStatusGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeStatusGet("_any")
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

func (s *OpsPublicTestSuite) TestNodeUptimeGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeUptimeGet("_any")
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

func (s *OpsPublicTestSuite) TestNodeDiskGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeDiskGet("_any")
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

func (s *OpsPublicTestSuite) TestNodeMemoryGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeMemoryGet("_any")
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

func (s *OpsPublicTestSuite) TestNodeLoadGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NodeLoadGet("_any")
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

func (s *OpsPublicTestSuite) TestNetworkDNSGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NetworkDNSGet("_any", "eth0")
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

func (s *OpsPublicTestSuite) TestNetworkDNSUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NetworkDNSUpdate(
					"_any",
					"eth0",
					[]string{"8.8.8.8"},
					[]string{"example.com"},
				)
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

func (s *OpsPublicTestSuite) TestNetworkPingDo() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NetworkPingDo("_any", "8.8.8.8")
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

func (s *OpsPublicTestSuite) TestCommandExec() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CommandExec("_any", "uptime", "-s")
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

func (s *OpsPublicTestSuite) TestCommandShell() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
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

func (s *OpsPublicTestSuite) TestAgentList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.AgentList()
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

func (s *OpsPublicTestSuite) TestAgentGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.AgentGet("server1")
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

func (s *OpsPublicTestSuite) TestAgentDrain() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.AgentDrain("server1")
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

func (s *OpsPublicTestSuite) TestAgentUndrain() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.AgentUndrain("server1")
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

func (s *OpsPublicTestSuite) TestFileDeploy() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step with raw content type",
			newFn: func() *orchestrator.Step {
				return s.orch.FileDeploy("_any", osapi.FileDeployOpts{
					ObjectName:  "config.yaml",
					Path:        "/etc/app/config.yaml",
					ContentType: "raw",
					Mode:        "0644",
					Owner:       "root",
					Group:       "root",
				})
			},
		},
		{
			name: "Returns non-nil step with template content type and vars",
			newFn: func() *orchestrator.Step {
				return s.orch.FileDeploy("_any", osapi.FileDeployOpts{
					ObjectName:  "config.tmpl",
					Path:        "/etc/app/config.yaml",
					ContentType: "template",
					Vars:        map[string]any{"env": "prod"},
				})
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

func (s *OpsPublicTestSuite) TestFileStatusGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.FileStatusGet("_any", "/etc/app/config.yaml")
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

func (s *OpsPublicTestSuite) TestFileUndeploy() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.FileUndeploy("_any", "/etc/app/config.yaml")
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

func (s *OpsPublicTestSuite) TestFileUpload() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.FileUpload("test.txt", "raw", []byte("hello"))
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

func (s *OpsPublicTestSuite) TestFileChanged() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.FileChanged("test.txt", []byte("hello"))
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

func (s *OpsPublicTestSuite) TestDockerPull() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerPull("_any", osapi.DockerPullOpts{
					Image: "nginx:latest",
				})
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

func (s *OpsPublicTestSuite) TestDockerCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerCreate("_any", osapi.DockerCreateOpts{
					Image: "nginx",
					Name:  "web",
				})
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

func (s *OpsPublicTestSuite) TestDockerStart() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerStart("_any", "c1")
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

func (s *OpsPublicTestSuite) TestDockerStop() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerStop("_any", "c1", osapi.DockerStopOpts{
					Timeout: 10,
				})
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

func (s *OpsPublicTestSuite) TestDockerRemove() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerRemove(
					"_any",
					"c1",
					&osapi.DockerRemoveParams{Force: true},
				)
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

func (s *OpsPublicTestSuite) TestDockerExec() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerExec(
					"_any",
					"c1",
					osapi.DockerExecOpts{
						Command: []string{"echo", "hello"},
					},
				)
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

func (s *OpsPublicTestSuite) TestDockerInspect() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerInspect("_any", "c1")
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

func (s *OpsPublicTestSuite) TestDockerList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerList("_any", &osapi.DockerListParams{
					State: "running",
				})
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

func (s *OpsPublicTestSuite) TestDockerImageRemove() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.DockerImageRemove(
					"_any",
					"nginx:latest",
					&osapi.DockerImageRemoveParams{Force: true},
				)
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

func (s *OpsPublicTestSuite) TestCronList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CronList("_any")
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

func (s *OpsPublicTestSuite) TestCronGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CronGet("_any", "backup")
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

func (s *OpsPublicTestSuite) TestCronCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CronCreate("_any", osapi.CronCreateOpts{
					Name:     "backup",
					Object:   "backup.sh",
					Schedule: "0 2 * * *",
				})
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

func (s *OpsPublicTestSuite) TestCronUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CronUpdate("_any", "backup", osapi.CronUpdateOpts{
					Schedule: "0 3 * * *",
				})
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

func (s *OpsPublicTestSuite) TestCronDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CronDelete("_any", "backup")
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
