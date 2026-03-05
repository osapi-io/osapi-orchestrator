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
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

type PredicatePublicTestSuite struct {
	suite.Suite
}

func (s *PredicatePublicTestSuite) TestOS() {
	tests := []struct {
		name         string
		distribution string
		agent        orchestrator.AgentResult
		expected     bool
	}{
		{
			name:         "Matches exact distribution",
			distribution: "ubuntu",
			agent: orchestrator.AgentResult{
				OSInfo: &orchestrator.AgentOSInfo{
					Distribution: "ubuntu",
				},
			},
			expected: true,
		},
		{
			name:         "Matches case-insensitive distribution",
			distribution: "Ubuntu",
			agent: orchestrator.AgentResult{
				OSInfo: &orchestrator.AgentOSInfo{
					Distribution: "ubuntu",
				},
			},
			expected: true,
		},
		{
			name:         "Returns false when OSInfo is nil",
			distribution: "ubuntu",
			agent:        orchestrator.AgentResult{},
			expected:     false,
		},
		{
			name:         "Returns false for non-matching distribution",
			distribution: "debian",
			agent: orchestrator.AgentResult{
				OSInfo: &orchestrator.AgentOSInfo{
					Distribution: "ubuntu",
				},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.OS(tc.distribution)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestArch() {
	tests := []struct {
		name         string
		architecture string
		agent        orchestrator.AgentResult
		expected     bool
	}{
		{
			name:         "Matches architecture",
			architecture: "x86_64",
			agent: orchestrator.AgentResult{
				Architecture: "x86_64",
			},
			expected: true,
		},
		{
			name:         "Matches case-insensitive architecture",
			architecture: "X86_64",
			agent: orchestrator.AgentResult{
				Architecture: "x86_64",
			},
			expected: true,
		},
		{
			name:         "Returns false for non-matching architecture",
			architecture: "arm64",
			agent: orchestrator.AgentResult{
				Architecture: "x86_64",
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.Arch(tc.architecture)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestMinMemory() {
	tests := []struct {
		name     string
		total    int
		agent    orchestrator.AgentResult
		expected bool
	}{
		{
			name:  "Matches when memory exceeds minimum",
			total: 4096,
			agent: orchestrator.AgentResult{
				Memory: &orchestrator.AgentMemory{
					Total: 8192,
				},
			},
			expected: true,
		},
		{
			name:  "Matches when memory equals minimum",
			total: 4096,
			agent: orchestrator.AgentResult{
				Memory: &orchestrator.AgentMemory{
					Total: 4096,
				},
			},
			expected: true,
		},
		{
			name:  "Returns false when memory below minimum",
			total: 8192,
			agent: orchestrator.AgentResult{
				Memory: &orchestrator.AgentMemory{
					Total: 4096,
				},
			},
			expected: false,
		},
		{
			name:     "Returns false when Memory is nil",
			total:    4096,
			agent:    orchestrator.AgentResult{},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.MinMemory(tc.total)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestMinCPU() {
	tests := []struct {
		name     string
		count    int
		agent    orchestrator.AgentResult
		expected bool
	}{
		{
			name:  "Matches when CPU count exceeds minimum",
			count: 2,
			agent: orchestrator.AgentResult{
				CPUCount: 4,
			},
			expected: true,
		},
		{
			name:  "Matches when CPU count equals minimum",
			count: 4,
			agent: orchestrator.AgentResult{
				CPUCount: 4,
			},
			expected: true,
		},
		{
			name:  "Returns false when CPU count below minimum",
			count: 8,
			agent: orchestrator.AgentResult{
				CPUCount: 4,
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.MinCPU(tc.count)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestHasLabel() {
	tests := []struct {
		name     string
		key      string
		value    string
		agent    orchestrator.AgentResult
		expected bool
	}{
		{
			name:  "Matches label key-value pair",
			key:   "env",
			value: "prod",
			agent: orchestrator.AgentResult{
				Labels: map[string]string{
					"env":  "prod",
					"team": "infra",
				},
			},
			expected: true,
		},
		{
			name:  "Returns false for wrong value",
			key:   "env",
			value: "prod",
			agent: orchestrator.AgentResult{
				Labels: map[string]string{
					"env": "staging",
				},
			},
			expected: false,
		},
		{
			name:     "Returns false when labels are nil",
			key:      "env",
			value:    "prod",
			agent:    orchestrator.AgentResult{},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.HasLabel(tc.key, tc.value)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestFactEquals() {
	tests := []struct {
		name     string
		key      string
		value    any
		agent    orchestrator.AgentResult
		expected bool
	}{
		{
			name:  "Matches string fact",
			key:   "datacenter",
			value: "us-east-1",
			agent: orchestrator.AgentResult{
				Facts: map[string]any{
					"datacenter": "us-east-1",
				},
			},
			expected: true,
		},
		{
			name:  "Matches numeric fact (float64)",
			key:   "version",
			value: float64(3),
			agent: orchestrator.AgentResult{
				Facts: map[string]any{
					"version": float64(3),
				},
			},
			expected: true,
		},
		{
			name:  "Returns false for wrong value",
			key:   "datacenter",
			value: "us-west-2",
			agent: orchestrator.AgentResult{
				Facts: map[string]any{
					"datacenter": "us-east-1",
				},
			},
			expected: false,
		},
		{
			name:     "Returns false when facts are nil",
			key:      "datacenter",
			value:    "us-east-1",
			agent:    orchestrator.AgentResult{},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.FactEquals(tc.key, tc.value)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestMatchAll() {
	tests := []struct {
		name       string
		agent      orchestrator.AgentResult
		predicates []orchestrator.Predicate
		expected   bool
	}{
		{
			name: "Returns true when all predicates match",
			agent: orchestrator.AgentResult{
				Architecture: "x86_64",
				CPUCount:     8,
				OSInfo: &orchestrator.AgentOSInfo{
					Distribution: "ubuntu",
				},
			},
			predicates: []orchestrator.Predicate{
				orchestrator.OS("ubuntu"),
				orchestrator.Arch("x86_64"),
				orchestrator.MinCPU(4),
			},
			expected: true,
		},
		{
			name: "Returns false when one predicate fails",
			agent: orchestrator.AgentResult{
				Architecture: "x86_64",
				CPUCount:     2,
				OSInfo: &orchestrator.AgentOSInfo{
					Distribution: "ubuntu",
				},
			},
			predicates: []orchestrator.Predicate{
				orchestrator.OS("ubuntu"),
				orchestrator.MinCPU(4),
			},
			expected: false,
		},
		{
			name:       "Returns true when no predicates are provided",
			agent:      orchestrator.AgentResult{},
			predicates: nil,
			expected:   true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			result := orchestrator.MatchAll(tc.agent, tc.predicates...)
			s.Equal(tc.expected, result)
		})
	}
}

func TestPredicatePublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(PredicatePublicTestSuite))
}
