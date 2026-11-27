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

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
)

// Step represents a single operation in the plan. Users chain methods
// to declare ordering, conditions, and error handling.
type Step struct {
	task *engine.Task
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
	tasks := make([]*engine.Task, len(deps))
	for i, d := range deps {
		tasks[i] = d.task
	}

	s.task.DependsOn(tasks...)

	return s
}

// Retry sets the number of retry attempts on failure. It waits between
// attempts: 1s, doubling up to 30s, unless WithBackoff sets other
// intervals. Retry is separate
// from OnError: a step can retry, and then continue the plan when its
// last attempt fails. A retry runs only for an error that a retry can
// fix (see Retryable), and not for an op that is not idempotent, such
// as CommandExec or PowerReboot, unless the step also calls
// RetryUnsafe.
func (s *Step) Retry(
	n int,
	opts ...RetryOption,
) *Step {
	s.task.Retry(enginePolicy(n, opts))

	return s
}

// RetryUnsafe lets Retry run the step again although its op is not
// idempotent, such as CommandExec, CommandShell, DockerCreate,
// DockerExec, ServiceRestart, ProcessSignal, PowerReboot and
// PowerShutdown. Use it only when running the op twice is safe, such as
// a command that checks before it changes anything.
func (s *Step) RetryUnsafe() *Step {
	s.task.SetRetryUnsafe()

	return s
}

// OnlyIfChanged skips this step unless a dependency reported changes.
func (s *Step) OnlyIfChanged() *Step {
	s.task.OnlyIfChanged()

	return s
}

// OnlyIfFailed skips this step unless at least one dependency failed.
// The step runs after a failed dependency, which a step with other
// guards does not.
func (s *Step) OnlyIfFailed() *Step {
	s.task.RunsOnFailure()
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		for _, dep := range s.task.Dependencies() {
			if r := sdkResults.Get(dep.Name()); r != nil && r.Status == engine.StatusFailed {
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
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
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
	s.task.RunsOnFailure()
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
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
				if engine.HostFailed(hr) {
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
// for broadcast operations. Skipped hosts are NOT errors, they
// indicate the operation is not available on that OS family.
func (s *Step) OnlyIfAnyHostSkipped() *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
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
	s.task.RunsOnFailure()
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
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
				if !engine.HostFailed(hr) {
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
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
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
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
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

// When adds a guard condition: the step only runs if the predicate
// returns true. Each guard of a step must hold, so When joins the
// OnlyIf guards with AND. A step with a When guard is skipped after a
// failed dependency.
func (s *Step) When(
	fn func(Results) bool,
) *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		return fn(Results{results: sdkResults})
	}, "when: guard returned false")

	return s
}

// WhenFact adds a fact-based guard. The step runs only if the
// predicate returns true for at least one agent of the given
// AgentList step. WhenFact makes this step run after that step. The
// skip reason distinguishes between a decode failure and no matching
// agents.
func (s *Step) WhenFact(
	agentList *Step,
	fn Predicate,
) *Step {
	s.task.DependsOn(agentList.task)
	s.task.Check(func(sdkResults engine.Results) (bool, string) {
		r := Results{results: sdkResults}
		name := agentList.task.Name()

		var list osapi.ListAgentsResponse
		if err := r.Decode(name, &list); err != nil {
			return false, fmt.Sprintf(
				"when-fact: step %q not found or decode failed: %v",
				name,
				err,
			)
		}

		if len(list.Items) == 0 {
			return false, "when-fact: no agents returned"
		}

		for _, a := range list.Items {
			if fn(a) {
				return true, ""
			}
		}

		return false, "when-fact: no matching agent"
	})

	return s
}

// OnError sets the error strategy for this step.
func (s *Step) OnError(
	strategy ErrorStrategy,
) *Step {
	s.task.OnError(toSDKStrategy(strategy))

	return s
}

// ContinueOnError is a convenience for OnError(Continue). Independent
// tasks keep running when this step fails; only direct dependents are
// skipped.
func (s *Step) ContinueOnError() *Step {
	return s.OnError(Continue)
}
