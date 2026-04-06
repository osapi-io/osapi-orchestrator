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

func (s *OpsPublicTestSuite) TestNetworkDNSDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NetworkDNSDelete("_any", "eth0")
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

func (s *OpsPublicTestSuite) TestInterfaceList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.InterfaceList("_any")
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

func (s *OpsPublicTestSuite) TestInterfaceGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.InterfaceGet("_any", "eth0")
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

func (s *OpsPublicTestSuite) TestInterfaceCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.InterfaceCreate("_any", "eth0", osapi.InterfaceConfigOpts{})
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

func (s *OpsPublicTestSuite) TestInterfaceUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.InterfaceUpdate("_any", "eth0", osapi.InterfaceConfigOpts{})
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

func (s *OpsPublicTestSuite) TestInterfaceDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.InterfaceDelete("_any", "eth0")
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

func (s *OpsPublicTestSuite) TestRouteList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.RouteList("_any")
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

func (s *OpsPublicTestSuite) TestRouteGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.RouteGet("_any", "eth0")
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

func (s *OpsPublicTestSuite) TestRouteCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.RouteCreate("_any", "eth0", osapi.RouteConfigOpts{})
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

func (s *OpsPublicTestSuite) TestRouteUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.RouteUpdate("_any", "eth0", osapi.RouteConfigOpts{})
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

func (s *OpsPublicTestSuite) TestRouteDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.RouteDelete("_any", "eth0")
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

func (s *OpsPublicTestSuite) TestSysctlList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.SysctlList("_any")
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

func (s *OpsPublicTestSuite) TestSysctlGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.SysctlGet("_any", "net.ipv4.ip_forward")
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

func (s *OpsPublicTestSuite) TestSysctlCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.SysctlCreate("_any", osapi.SysctlCreateOpts{})
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

func (s *OpsPublicTestSuite) TestSysctlUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.SysctlUpdate("_any", "net.ipv4.ip_forward", osapi.SysctlUpdateOpts{})
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

func (s *OpsPublicTestSuite) TestSysctlDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.SysctlDelete("_any", "net.ipv4.ip_forward")
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

func (s *OpsPublicTestSuite) TestNTPGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NTPGet("_any")
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

func (s *OpsPublicTestSuite) TestNTPCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NTPCreate("_any", osapi.NtpCreateOpts{})
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

func (s *OpsPublicTestSuite) TestNTPUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NTPUpdate("_any", osapi.NtpUpdateOpts{})
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

func (s *OpsPublicTestSuite) TestNTPDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.NTPDelete("_any")
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

func (s *OpsPublicTestSuite) TestTimezoneGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.TimezoneGet("_any")
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

func (s *OpsPublicTestSuite) TestTimezoneUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.TimezoneUpdate("_any", osapi.TimezoneUpdateOpts{})
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

func (s *OpsPublicTestSuite) TestServiceList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceList("_any")
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

func (s *OpsPublicTestSuite) TestServiceGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceGet("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestServiceCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceCreate("_any", osapi.ServiceCreateOpts{})
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

func (s *OpsPublicTestSuite) TestServiceUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceUpdate("_any", "nginx", osapi.ServiceUpdateOpts{})
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

func (s *OpsPublicTestSuite) TestServiceDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceDelete("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestServiceStart() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceStart("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestServiceStop() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceStop("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestServiceRestart() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceRestart("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestServiceEnable() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceEnable("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestServiceDisable() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ServiceDisable("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestPackageList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PackageList("_any")
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

func (s *OpsPublicTestSuite) TestPackageGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PackageGet("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestPackageInstall() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PackageInstall("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestPackageRemove() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PackageRemove("_any", "nginx")
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

func (s *OpsPublicTestSuite) TestPackageUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PackageUpdate("_any")
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

func (s *OpsPublicTestSuite) TestPackageListUpdates() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PackageListUpdates("_any")
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

func (s *OpsPublicTestSuite) TestUserList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserList("_any")
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

func (s *OpsPublicTestSuite) TestUserGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserGet("_any", "admin")
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

func (s *OpsPublicTestSuite) TestUserCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserCreate("_any", osapi.UserCreateOpts{})
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

func (s *OpsPublicTestSuite) TestUserUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserUpdate("_any", "admin", osapi.UserUpdateOpts{})
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

func (s *OpsPublicTestSuite) TestUserDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserDelete("_any", "admin")
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

func (s *OpsPublicTestSuite) TestUserListKeys() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserListKeys("_any", "admin")
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

func (s *OpsPublicTestSuite) TestUserAddKey() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserAddKey("_any", "admin", osapi.SSHKeyAddOpts{})
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

func (s *OpsPublicTestSuite) TestUserRemoveKey() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserRemoveKey("_any", "admin", "SHA256:abc123")
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

func (s *OpsPublicTestSuite) TestUserChangePassword() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.UserChangePassword("_any", "admin", "newpass")
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

func (s *OpsPublicTestSuite) TestGroupList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.GroupList("_any")
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

func (s *OpsPublicTestSuite) TestGroupGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.GroupGet("_any", "admins")
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

func (s *OpsPublicTestSuite) TestGroupCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.GroupCreate("_any", osapi.GroupCreateOpts{})
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

func (s *OpsPublicTestSuite) TestGroupUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.GroupUpdate("_any", "admins", osapi.GroupUpdateOpts{})
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

func (s *OpsPublicTestSuite) TestGroupDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.GroupDelete("_any", "admins")
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

func (s *OpsPublicTestSuite) TestCertificateList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CertificateList("_any")
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

func (s *OpsPublicTestSuite) TestCertificateCreate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CertificateCreate("_any", osapi.CertificateCreateOpts{})
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

func (s *OpsPublicTestSuite) TestCertificateUpdate() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CertificateUpdate("_any", "my-ca", osapi.CertificateUpdateOpts{})
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

func (s *OpsPublicTestSuite) TestCertificateDelete() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.CertificateDelete("_any", "my-ca")
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

func (s *OpsPublicTestSuite) TestProcessList() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ProcessList("_any")
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

func (s *OpsPublicTestSuite) TestProcessGet() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ProcessGet("_any", 1234)
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

func (s *OpsPublicTestSuite) TestProcessSignal() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.ProcessSignal("_any", 1234, osapi.ProcessSignalOpts{})
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

func (s *OpsPublicTestSuite) TestPowerReboot() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PowerReboot("_any", osapi.PowerOpts{})
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

func (s *OpsPublicTestSuite) TestPowerShutdown() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.PowerShutdown("_any", osapi.PowerOpts{})
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

func (s *OpsPublicTestSuite) TestLogQuery() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.LogQuery("_any", osapi.LogQueryOpts{})
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

func (s *OpsPublicTestSuite) TestLogSources() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.LogSources("_any")
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

func (s *OpsPublicTestSuite) TestLogQueryUnit() {
	tests := []struct {
		name  string
		newFn func() *orchestrator.Step
	}{
		{
			name: "Returns non-nil step",
			newFn: func() *orchestrator.Step {
				return s.orch.LogQueryUnit("_any", "nginx.service", osapi.LogQueryOpts{})
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
