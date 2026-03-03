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

import sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"

// Step represents a single operation in the plan. Users chain methods
// to declare ordering, conditions, and error handling.
type Step struct {
	task *sdk.Task
}

// Named overrides the auto-generated step name.
func (s *Step) Named(
	name string,
) *Step {
	s.task.SetName(name)

	return s
}

// After declares that this step runs after the given steps complete.
func (s *Step) After(
	deps ...*Step,
) *Step {
	tasks := make([]*sdk.Task, len(deps))
	for i, d := range deps {
		tasks[i] = d.task
	}

	s.task.DependsOn(tasks...)

	return s
}

// Retry sets the number of retry attempts on failure.
func (s *Step) Retry(
	n int,
) *Step {
	s.task.OnError(sdk.Retry(n))

	return s
}

// OnlyIfChanged skips this step unless a dependency reported changes.
func (s *Step) OnlyIfChanged() *Step {
	s.task.OnlyIfChanged()

	return s
}

// OnlyIfFailed skips this step unless at least one dependency failed.
func (s *Step) OnlyIfFailed() *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		for _, dep := range s.task.Dependencies() {
			if r := sdkResults.Get(dep.Name()); r != nil && r.Status == sdk.StatusFailed {
				return true
			}
		}

		return false
	}, "only-if-failed: no dependency failed")

	return s
}

// OnlyIfAllChanged skips this step unless all dependencies reported
// changes.
func (s *Step) OnlyIfAllChanged() *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		deps := s.task.Dependencies()
		if len(deps) == 0 {
			return false
		}

		for _, dep := range deps {
			r := sdkResults.Get(dep.Name())
			if r == nil || r.Status != sdk.StatusChanged {
				return false
			}
		}

		return true
	}, "only-if-all-changed: not all dependencies changed")

	return s
}

// When adds a guard condition — the step only runs if the predicate
// returns true.
func (s *Step) When(
	fn func(Results) bool,
) *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		return fn(Results{results: sdkResults})
	}, "when: guard returned false")

	return s
}

// OnError sets the error strategy for this step.
func (s *Step) OnError(
	strategy ErrorStrategy,
) *Step {
	s.task.OnError(toSDKStrategy(strategy))

	return s
}
