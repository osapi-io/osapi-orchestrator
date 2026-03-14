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

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"
)

type DiscoverTestSuite struct {
	suite.Suite
}

func (s *DiscoverTestSuite) TestFactValue() {
	tests := []struct {
		name     string
		agent    osapi.Agent
		key      string
		expected string
	}{
		{
			name: "os.distribution with OSInfo returns Distribution",
			agent: osapi.Agent{
				OSInfo: &osapi.OSInfo{
					Distribution: "Ubuntu",
					Version:      "22.04",
				},
			},
			key:      "os.distribution",
			expected: "Ubuntu",
		},
		{
			name: "os.version with OSInfo returns Version",
			agent: osapi.Agent{
				OSInfo: &osapi.OSInfo{
					Distribution: "Ubuntu",
					Version:      "22.04",
				},
			},
			key:      "os.version",
			expected: "22.04",
		},
		{
			name: "architecture returns Architecture",
			agent: osapi.Agent{
				Architecture: "amd64",
			},
			key:      "architecture",
			expected: "amd64",
		},
		{
			name: "service_manager returns ServiceMgr",
			agent: osapi.Agent{
				ServiceMgr: "systemd",
			},
			key:      "service_manager",
			expected: "systemd",
		},
		{
			name: "package_manager returns PackageMgr",
			agent: osapi.Agent{
				PackageMgr: "apt",
			},
			key:      "package_manager",
			expected: "apt",
		},
		{
			name: "kernel_version returns KernelVersion",
			agent: osapi.Agent{
				KernelVersion: "5.15.0",
			},
			key:      "kernel_version",
			expected: "5.15.0",
		},
		{
			name: "Falls back to Facts map for unknown keys",
			agent: osapi.Agent{
				Facts: map[string]any{
					"datacenter": "us-east-1",
				},
			},
			key:      "datacenter",
			expected: "us-east-1",
		},
		{
			name:     "Returns empty for nil OSInfo on os.distribution",
			agent:    osapi.Agent{},
			key:      "os.distribution",
			expected: "",
		},
		{
			name:     "Returns empty for nil OSInfo on os.version",
			agent:    osapi.Agent{},
			key:      "os.version",
			expected: "",
		},
		{
			name:     "Returns empty for missing fact",
			agent:    osapi.Agent{},
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

func TestDiscoverTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(DiscoverTestSuite))
}
