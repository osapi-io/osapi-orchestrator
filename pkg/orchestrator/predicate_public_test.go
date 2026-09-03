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

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
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
		agent        osapi.Agent
		expected     bool
	}{
		{
			name:         "Matches exact distribution",
			distribution: osFamilyUbuntu,
			agent: osapi.Agent{
				OSInfo: &osapi.OSInfo{
					Distribution: osFamilyUbuntu,
				},
			},
			expected: true,
		},
		{
			name:         "Matches case-insensitive distribution",
			distribution: "Ubuntu",
			agent: osapi.Agent{
				OSInfo: &osapi.OSInfo{
					Distribution: osFamilyUbuntu,
				},
			},
			expected: true,
		},
		{
			name:         "Returns false when OSInfo is nil",
			distribution: osFamilyUbuntu,
			agent:        osapi.Agent{},
			expected:     false,
		},
		{
			name:         "Returns false for non-matching distribution",
			distribution: "debian",
			agent: osapi.Agent{
				OSInfo: &osapi.OSInfo{
					Distribution: osFamilyUbuntu,
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
		agent        osapi.Agent
		expected     bool
	}{
		{
			name:         "Matches architecture",
			architecture: archX8664,
			agent: osapi.Agent{
				Architecture: archX8664,
			},
			expected: true,
		},
		{
			name:         "Matches case-insensitive architecture",
			architecture: "X86_64",
			agent: osapi.Agent{
				Architecture: archX8664,
			},
			expected: true,
		},
		{
			name:         "Returns false for non-matching architecture",
			architecture: "arm64",
			agent: osapi.Agent{
				Architecture: archX8664,
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
		agent    osapi.Agent
		expected bool
	}{
		{
			name:  "Matches when memory exceeds minimum",
			total: 4096,
			agent: osapi.Agent{
				Memory: &osapi.Memory{
					Total: 8192,
				},
			},
			expected: true,
		},
		{
			name:  "Matches when memory equals minimum",
			total: 4096,
			agent: osapi.Agent{
				Memory: &osapi.Memory{
					Total: 4096,
				},
			},
			expected: true,
		},
		{
			name:  "Returns false when memory below minimum",
			total: 8192,
			agent: osapi.Agent{
				Memory: &osapi.Memory{
					Total: 4096,
				},
			},
			expected: false,
		},
		{
			name:     "Returns false when Memory is nil",
			total:    4096,
			agent:    osapi.Agent{},
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
		agent    osapi.Agent
		expected bool
	}{
		{
			name:  "Matches when CPU count exceeds minimum",
			count: 2,
			agent: osapi.Agent{
				CPUCount: 4,
			},
			expected: true,
		},
		{
			name:  "Matches when CPU count equals minimum",
			count: 4,
			agent: osapi.Agent{
				CPUCount: 4,
			},
			expected: true,
		},
		{
			name:  "Returns false when CPU count below minimum",
			count: 8,
			agent: osapi.Agent{
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
		agent    osapi.Agent
		expected bool
	}{
		{
			name:  "Matches label key-value pair",
			key:   labelEnv,
			value: envProd,
			agent: osapi.Agent{
				Labels: map[string]string{
					labelEnv: envProd,
					"team":   "infra",
				},
			},
			expected: true,
		},
		{
			name:  "Returns false for wrong value",
			key:   labelEnv,
			value: envProd,
			agent: osapi.Agent{
				Labels: map[string]string{
					labelEnv: "staging",
				},
			},
			expected: false,
		},
		{
			name:     "Returns false when labels are nil",
			key:      labelEnv,
			value:    envProd,
			agent:    osapi.Agent{},
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
		agent    osapi.Agent
		expected bool
	}{
		{
			name:  "Matches string fact",
			key:   labelDatacenter,
			value: regionUSEast1,
			agent: osapi.Agent{
				Facts: map[string]any{
					labelDatacenter: regionUSEast1,
				},
			},
			expected: true,
		},
		{
			name:  "Matches numeric fact (float64)",
			key:   "version",
			value: float64(3),
			agent: osapi.Agent{
				Facts: map[string]any{
					"version": float64(3),
				},
			},
			expected: true,
		},
		{
			name:  "Returns false for wrong value",
			key:   labelDatacenter,
			value: "us-west-2",
			agent: osapi.Agent{
				Facts: map[string]any{
					labelDatacenter: regionUSEast1,
				},
			},
			expected: false,
		},
		{
			name:     "Returns false when facts are nil",
			key:      labelDatacenter,
			value:    regionUSEast1,
			agent:    osapi.Agent{},
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

func (s *PredicatePublicTestSuite) TestHasCondition() {
	tests := []struct {
		name          string
		conditionType string
		agent         osapi.Agent
		expected      bool
	}{
		{
			name:          "Matches active condition",
			conditionType: conditionDiskPressure,
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: true},
				},
			},
			expected: true,
		},
		{
			name:          "No match when condition is inactive",
			conditionType: conditionDiskPressure,
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: false},
				},
			},
			expected: false,
		},
		{
			name:          "No match for wrong type",
			conditionType: conditionMemoryPressure,
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: true},
				},
			},
			expected: false,
		},
		{
			name:          "No match when conditions are nil",
			conditionType: conditionDiskPressure,
			agent:         osapi.Agent{},
			expected:      false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.HasCondition(tc.conditionType)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestNoCondition() {
	tests := []struct {
		name          string
		conditionType string
		agent         osapi.Agent
		expected      bool
	}{
		{
			name:          "No match when condition is active",
			conditionType: conditionDiskPressure,
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: true},
				},
			},
			expected: false,
		},
		{
			name:          "Matches when condition is inactive",
			conditionType: conditionDiskPressure,
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: false},
				},
			},
			expected: true,
		},
		{
			name:          "Matches when type is missing",
			conditionType: conditionMemoryPressure,
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: true},
				},
			},
			expected: true,
		},
		{
			name:          "Matches when conditions are nil",
			conditionType: conditionDiskPressure,
			agent:         osapi.Agent{},
			expected:      true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.NoCondition(tc.conditionType)
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestHealthy() {
	tests := []struct {
		name     string
		agent    osapi.Agent
		expected bool
	}{
		{
			name:     "Matches when no conditions",
			agent:    osapi.Agent{},
			expected: true,
		},
		{
			name: "Matches when all conditions inactive",
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: false},
					{Type: conditionMemoryPressure, Status: false},
				},
			},
			expected: true,
		},
		{
			name: "No match when one condition active",
			agent: osapi.Agent{
				Conditions: []osapi.Condition{
					{Type: conditionDiskPressure, Status: false},
					{Type: conditionMemoryPressure, Status: true},
				},
			},
			expected: false,
		},
		{
			name: "Matches with empty conditions slice",
			agent: osapi.Agent{
				Conditions: []osapi.Condition{},
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			predicate := orchestrator.Healthy()
			s.Equal(tc.expected, predicate(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestMatchAll() {
	tests := []struct {
		name       string
		agent      osapi.Agent
		predicates []orchestrator.Predicate
		expected   bool
	}{
		{
			name: "Returns true when all predicates match",
			agent: osapi.Agent{
				Architecture: archX8664,
				CPUCount:     8,
				OSInfo: &osapi.OSInfo{
					Distribution: osFamilyUbuntu,
				},
			},
			predicates: []orchestrator.Predicate{
				orchestrator.OS(osFamilyUbuntu),
				orchestrator.Arch(archX8664),
				orchestrator.MinCPU(4),
			},
			expected: true,
		},
		{
			name: "Returns false when one predicate fails",
			agent: osapi.Agent{
				Architecture: archX8664,
				CPUCount:     2,
				OSInfo: &osapi.OSInfo{
					Distribution: osFamilyUbuntu,
				},
			},
			predicates: []orchestrator.Predicate{
				orchestrator.OS(osFamilyUbuntu),
				orchestrator.MinCPU(4),
			},
			expected: false,
		},
		{
			name:       "Returns true when no predicates are provided",
			agent:      osapi.Agent{},
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
