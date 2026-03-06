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

import "strings"

// Predicate filters agents by their facts and properties.
type Predicate func(AgentResult) bool

// OS returns a predicate that matches agents running the given
// distribution (case-insensitive).
func OS(
	distribution string,
) Predicate {
	return func(a AgentResult) bool {
		if a.OSInfo == nil {
			return false
		}

		return strings.EqualFold(a.OSInfo.Distribution, distribution)
	}
}

// Arch returns a predicate that matches agents with the given
// architecture (case-insensitive).
func Arch(
	architecture string,
) Predicate {
	return func(a AgentResult) bool {
		return strings.EqualFold(a.Architecture, architecture)
	}
}

// MinMemory returns a predicate that matches agents with at least
// the given total memory (in the same unit as AgentMemory.Total).
func MinMemory(
	total int,
) Predicate {
	return func(a AgentResult) bool {
		if a.Memory == nil {
			return false
		}

		return a.Memory.Total >= total
	}
}

// MinCPU returns a predicate that matches agents with at least
// the given number of CPUs.
func MinCPU(
	count int,
) Predicate {
	return func(a AgentResult) bool {
		return a.CPUCount >= count
	}
}

// HasLabel returns a predicate that matches agents with the given
// label key-value pair.
func HasLabel(
	key string,
	value string,
) Predicate {
	return func(a AgentResult) bool {
		return a.Labels[key] == value
	}
}

// FactEquals returns a predicate that matches agents where the
// given fact key equals the expected value.
func FactEquals(
	key string,
	value any,
) Predicate {
	return func(a AgentResult) bool {
		return a.Facts[key] == value
	}
}

// HasCondition returns a predicate that matches agents with an active
// condition (Status=true) of the given type.
func HasCondition(
	conditionType string,
) Predicate {
	return func(a AgentResult) bool {
		for _, c := range a.Conditions {
			if strings.EqualFold(c.Type, conditionType) && c.Status {
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
	return func(a AgentResult) bool {
		for _, c := range a.Conditions {
			if strings.EqualFold(c.Type, conditionType) && c.Status {
				return false
			}
		}

		return true
	}
}

// Healthy returns a predicate that matches agents with no active
// conditions (all conditions are false or the list is empty).
func Healthy() Predicate {
	return func(a AgentResult) bool {
		for _, c := range a.Conditions {
			if c.Status {
				return false
			}
		}

		return true
	}
}

// MatchAll returns true if the agent matches all given predicates.
// Returns true if no predicates are provided.
func MatchAll(
	agent AgentResult,
	predicates ...Predicate,
) bool {
	for _, p := range predicates {
		if !p(agent) {
			return false
		}
	}

	return true
}
