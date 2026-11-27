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
	"strings"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
)

// Predicate filters agents by their facts and properties.
type Predicate func(osapi.AgentInfo) bool

// OS returns a predicate that matches agents running the given
// distribution (case-insensitive).
func OS(
	distribution string,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		return a.OSInfo.Distribution != "" &&
			strings.EqualFold(a.OSInfo.Distribution, distribution)
	}
}

// Arch returns a predicate that matches agents with the given
// architecture (case-insensitive).
func Arch(
	architecture string,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		return strings.EqualFold(a.Architecture, architecture)
	}
}

// MinMemory returns a predicate that matches agents with at least
// the given total memory (in the same unit as Memory.Total).
func MinMemory(
	total int,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		return a.Memory.Total > 0 && a.Memory.Total >= total
	}
}

// MinCPU returns a predicate that matches agents with at least
// the given number of CPUs.
func MinCPU(
	count int,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		return a.CPUCount >= count
	}
}

// HasLabel returns a predicate that matches agents with the given
// label key-value pair.
func HasLabel(
	key string,
	value string,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		return a.Labels[key] == value
	}
}

// FactEquals returns a predicate that matches agents where the
// given fact key equals the expected value.
func FactEquals(
	key string,
	value any,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		return a.Facts[key] == value
	}
}

// HasCondition returns a predicate that matches agents with an active
// condition of the given type. An agent lists only its active conditions.
func HasCondition(
	conditionType string,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		for _, c := range a.Conditions {
			if strings.EqualFold(string(c.Type), conditionType) {
				return true
			}
		}

		return false
	}
}

// NoCondition returns a predicate that matches agents that do NOT have
// an active condition of the given type.
func NoCondition(
	conditionType string,
) Predicate {
	return func(a osapi.AgentInfo) bool {
		for _, c := range a.Conditions {
			if strings.EqualFold(string(c.Type), conditionType) {
				return false
			}
		}

		return true
	}
}

// Healthy returns a predicate that matches agents with no active
// conditions.
func Healthy() Predicate {
	return func(a osapi.AgentInfo) bool {
		return len(a.Conditions) == 0
	}
}

// MatchAll returns true if the agent matches all given predicates.
// Returns true if no predicates are provided.
func MatchAll(
	agent osapi.AgentInfo,
	predicates ...Predicate,
) bool {
	for _, p := range predicates {
		if !p(agent) {
			return false
		}
	}

	return true
}
