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
	"fmt"

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
)

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

// Retry sets the number of retry attempts on failure. Options
// configure exponential backoff between attempts.
func (s *Step) Retry(
	n int,
	opts ...RetryOption,
) *Step {
	cfg := &retryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	var sdkOpts []sdk.RetryOption
	if cfg.initialInterval > 0 {
		sdkOpts = append(
			sdkOpts,
			sdk.WithRetryBackoff(cfg.initialInterval, cfg.maxInterval),
		)
	}

	s.task.OnError(sdk.Retry(n, sdkOpts...))

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
			if r == nil || !r.Changed {
				return false
			}
		}

		return true
	}, "only-if-all-changed: not all dependencies changed")

	return s
}

// OnlyIfAnyHostFailed skips this step unless any host in any
// dependency has an error. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAnyHostFailed() *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		deps := s.task.Dependencies()
		if len(deps) == 0 {
			return false
		}

		for _, dep := range deps {
			r := sdkResults.Get(dep.Name())
			if r == nil || len(r.HostResults) == 0 {
				continue
			}

			for _, hr := range r.HostResults {
				if hr.Status == "failed" {
					return true
				}
			}
		}

		return false
	}, "only-if-any-host-failed: no host failed")

	return s
}

// OnlyIfAnyHostSkipped skips this step unless any host in any
// dependency was skipped (unsupported operation). Only meaningful
// for broadcast operations. Skipped hosts are NOT errors — they
// indicate the operation is not available on that OS family.
func (s *Step) OnlyIfAnyHostSkipped() *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		deps := s.task.Dependencies()
		if len(deps) == 0 {
			return false
		}

		for _, dep := range deps {
			r := sdkResults.Get(dep.Name())
			if r == nil || len(r.HostResults) == 0 {
				continue
			}

			for _, hr := range r.HostResults {
				if hr.Status == "skipped" {
					return true
				}
			}
		}

		return false
	}, "only-if-any-host-skipped: no host skipped")

	return s
}

// OnlyIfAllHostsFailed skips this step unless every host in every
// dependency has an error. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAllHostsFailed() *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		deps := s.task.Dependencies()
		if len(deps) == 0 {
			return false
		}

		for _, dep := range deps {
			r := sdkResults.Get(dep.Name())
			if r == nil || len(r.HostResults) == 0 {
				return false
			}

			for _, hr := range r.HostResults {
				if hr.Status != "failed" {
					return false
				}
			}
		}

		return true
	}, "only-if-all-hosts-failed: not all hosts failed")

	return s
}

// OnlyIfAnyHostChanged skips this step unless any host in any
// dependency reported changes. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAnyHostChanged() *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		deps := s.task.Dependencies()
		if len(deps) == 0 {
			return false
		}

		for _, dep := range deps {
			r := sdkResults.Get(dep.Name())
			if r == nil || len(r.HostResults) == 0 {
				continue
			}

			for _, hr := range r.HostResults {
				if hr.Changed {
					return true
				}
			}
		}

		return false
	}, "only-if-any-host-changed: no host changed")

	return s
}

// OnlyIfAllHostsChanged skips this step unless every host in every
// dependency reported changes. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAllHostsChanged() *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		deps := s.task.Dependencies()
		if len(deps) == 0 {
			return false
		}

		for _, dep := range deps {
			r := sdkResults.Get(dep.Name())
			if r == nil || len(r.HostResults) == 0 {
				return false
			}

			for _, hr := range r.HostResults {
				if !hr.Changed {
					return false
				}
			}
		}

		return true
	}, "only-if-all-hosts-changed: not all hosts changed")

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

// WhenFact adds a fact-based guard. The step runs only if the
// predicate returns true for at least one agent. Requires a prior
// AgentList step referenced by name. The skip reason distinguishes
// between a missing/mistyped step name and no matching agents.
func (s *Step) WhenFact(
	agentListStep string,
	fn Predicate,
) *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		r := Results{results: sdkResults}

		var list osapi.AgentList
		if err := r.Decode(agentListStep, &list); err != nil {
			s.task.SetGuardReason(
				fmt.Sprintf(
					"when-fact: step %q not found or decode failed: %v",
					agentListStep,
					err,
				),
			)

			return false
		}

		if len(list.Agents) == 0 {
			s.task.SetGuardReason("when-fact: no agents returned")

			return false
		}

		for _, a := range list.Agents {
			if fn(a) {
				return true
			}
		}

		return false
	}, "when-fact: no matching agent")

	return s
}

// OnError sets the error strategy for this step.
func (s *Step) OnError(
	strategy ErrorStrategy,
) *Step {
	s.task.OnError(toSDKStrategy(strategy))

	return s
}
