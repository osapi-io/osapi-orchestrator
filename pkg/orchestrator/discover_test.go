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
	"testing"

	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
	"github.com/stretchr/testify/suite"
)

type DiscoverTestSuite struct {
	suite.Suite
}

func (s *DiscoverTestSuite) TestFactValue() {
	tests := []struct {
		name     string
		agent    AgentResult
		key      string
		expected string
	}{
		{
			name: "os.distribution with OSInfo returns Distribution",
			agent: AgentResult{
				OSInfo: &AgentOSInfo{
					Distribution: "Ubuntu",
					Version:      "22.04",
				},
			},
			key:      "os.distribution",
			expected: "Ubuntu",
		},
		{
			name: "os.version with OSInfo returns Version",
			agent: AgentResult{
				OSInfo: &AgentOSInfo{
					Distribution: "Ubuntu",
					Version:      "22.04",
				},
			},
			key:      "os.version",
			expected: "22.04",
		},
		{
			name: "architecture returns Architecture",
			agent: AgentResult{
				Architecture: "amd64",
			},
			key:      "architecture",
			expected: "amd64",
		},
		{
			name: "service_manager returns ServiceMgr",
			agent: AgentResult{
				ServiceMgr: "systemd",
			},
			key:      "service_manager",
			expected: "systemd",
		},
		{
			name: "package_manager returns PackageMgr",
			agent: AgentResult{
				PackageMgr: "apt",
			},
			key:      "package_manager",
			expected: "apt",
		},
		{
			name: "kernel_version returns KernelVersion",
			agent: AgentResult{
				KernelVersion: "5.15.0",
			},
			key:      "kernel_version",
			expected: "5.15.0",
		},
		{
			name: "Falls back to Facts map for unknown keys",
			agent: AgentResult{
				Facts: map[string]any{
					"datacenter": "us-east-1",
				},
			},
			key:      "datacenter",
			expected: "us-east-1",
		},
		{
			name:     "Returns empty for nil OSInfo on os.distribution",
			agent:    AgentResult{},
			key:      "os.distribution",
			expected: "",
		},
		{
			name:     "Returns empty for nil OSInfo on os.version",
			agent:    AgentResult{},
			key:      "os.version",
			expected: "",
		},
		{
			name:     "Returns empty for missing fact",
			agent:    AgentResult{},
			key:      "nonexistent",
			expected: "",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := factValue(tc.agent, tc.key)
			s.Equal(tc.expected, got)
		})
	}
}

func (s *DiscoverTestSuite) TestMustDecode() {
	tests := []struct {
		name   string
		report *Report
		task   string
		panics bool
	}{
		{
			name: "Succeeds with valid data",
			report: &Report{
				Tasks: []sdk.TaskResult{
					{
						Name: "test",
						Data: map[string]any{
							"hostname": "web-01",
						},
					},
				},
			},
			task: "test",
		},
		{
			name: "Panics on missing task",
			report: &Report{
				Tasks: []sdk.TaskResult{},
			},
			task:   "test",
			panics: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			if tc.panics {
				var result map[string]any
				s.Panics(func() {
					mustDecode(tc.report, tc.task, &result)
				})

				return
			}

			var result HostnameResult
			s.NotPanics(func() {
				mustDecode(tc.report, tc.task, &result)
			})
			s.Equal("web-01", result.Hostname)
		})
	}
}

func TestDiscoverTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(DiscoverTestSuite))
}
