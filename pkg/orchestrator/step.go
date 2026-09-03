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
	"slices"

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

	var sdkOpts []engine.RetryOption
	if cfg.initialInterval > 0 {
		sdkOpts = append(
			sdkOpts,
			engine.WithRetryBackoff(cfg.initialInterval, cfg.maxInterval),
		)
	}

	s.task.OnError(engine.Retry(n, sdkOpts...))

	return s
}

// OnlyIfChanged skips this step unless a dependency reported changes.
func (s *Step) OnlyIfChanged() *Step {
	s.task.OnlyIfChanged()

	return s
}

// OnlyIfFailed skips this step unless at least one dependency failed.
func (s *Step) OnlyIfFailed() *Step {
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
		return s.allDeps(sdkResults, func(r *engine.Result) bool {
			return r.Changed
		})
	}, "only-if-all-changed: not all dependencies changed")

	return s
}

// OnlyIfAnyHostFailed skips this step unless any host in any
// dependency has an error. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAnyHostFailed() *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		return s.anyHost(sdkResults, hostFailed)
	}, "only-if-any-host-failed: no host failed")

	return s
}

// OnlyIfAnyHostSkipped skips this step unless any host in any
// dependency was skipped (unsupported operation). Only meaningful
// for broadcast operations. Skipped hosts are NOT errors — they
// indicate the operation is not available on that OS family.
func (s *Step) OnlyIfAnyHostSkipped() *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		return s.anyHost(sdkResults, hostSkipped)
	}, "only-if-any-host-skipped: no host skipped")

	return s
}

// OnlyIfAllHostsFailed skips this step unless every host in every
// dependency has an error. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAllHostsFailed() *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		return s.allHosts(sdkResults, hostFailed)
	}, "only-if-all-hosts-failed: not all hosts failed")

	return s
}

// OnlyIfAnyHostChanged skips this step unless any host in any
// dependency reported changes. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAnyHostChanged() *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		return s.anyHost(sdkResults, hostChanged)
	}, "only-if-any-host-changed: no host changed")

	return s
}

// OnlyIfAllHostsChanged skips this step unless every host in every
// dependency reported changes. Only meaningful for broadcast operations.
func (s *Step) OnlyIfAllHostsChanged() *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		return s.allHosts(sdkResults, hostChanged)
	}, "only-if-all-hosts-changed: not all hosts changed")

	return s
}

// hostFailed reports whether a host errored. A skipped host is not a
// failure: it means the operation is unavailable on that OS family.
func hostFailed(hr engine.HostResult) bool {
	return hr.Status == string(engine.StatusFailed) ||
		(hr.Error != "" && hr.Status != string(engine.StatusSkipped))
}

// hostSkipped reports whether the operation was unavailable on a host.
func hostSkipped(hr engine.HostResult) bool {
	return hr.Status == string(engine.StatusSkipped)
}

// hostChanged reports whether a host reported a change.
func hostChanged(hr engine.HostResult) bool {
	return hr.Changed
}

// allDeps reports whether every dependency produced a result satisfying
// pred. A step with no dependencies satisfies nothing, so that a guard
// written for dependencies does not pass when there are none.
func (s *Step) allDeps(
	results engine.Results,
	pred func(*engine.Result) bool,
) bool {
	deps := s.task.Dependencies()
	if len(deps) == 0 {
		return false
	}

	for _, dep := range deps {
		r := results.Get(dep.Name())
		if r == nil || !pred(r) {
			return false
		}
	}

	return true
}

// hostsOf returns the hosts a dependency reported, or nil when the
// dependency produced no result at all.
func hostsOf(results engine.Results, dep *engine.Task) []engine.HostResult {
	r := results.Get(dep.Name())
	if r == nil {
		return nil
	}

	return r.HostResults
}

// everyHost reports whether every host in the slice satisfies pred. An
// empty slice satisfies it vacuously; callers decide what that means.
func everyHost(
	hosts []engine.HostResult,
	pred func(engine.HostResult) bool,
) bool {
	for _, hr := range hosts {
		if !pred(hr) {
			return false
		}
	}

	return true
}

// anyHost reports whether any host of any dependency satisfies pred. A
// dependency that reported no hosts is passed over rather than failing
// the test — there is nothing there to disagree with.
func (s *Step) anyHost(
	results engine.Results,
	pred func(engine.HostResult) bool,
) bool {
	deps := s.task.Dependencies()
	if len(deps) == 0 {
		return false
	}

	for _, dep := range deps {
		if slices.ContainsFunc(hostsOf(results, dep), pred) {
			return true
		}
	}

	return false
}

// allHosts reports whether every host of every dependency satisfies
// pred. Unlike anyHost, a dependency that reported no hosts fails the
// test: "all hosts changed" is not true of a dependency that ran on none.
func (s *Step) allHosts(
	results engine.Results,
	pred func(engine.HostResult) bool,
) bool {
	deps := s.task.Dependencies()
	if len(deps) == 0 {
		return false
	}

	for _, dep := range deps {
		hosts := hostsOf(results, dep)
		if len(hosts) == 0 || !everyHost(hosts, pred) {
			return false
		}
	}

	return true
}

// When adds a guard condition — the step only runs if the predicate
// returns true.
func (s *Step) When(
	fn func(Results) bool,
) *Step {
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
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
	s.task.WhenWithReason(func(sdkResults engine.Results) bool {
		agents, ok := s.agentsFrom(sdkResults, agentListStep)
		if !ok {
			return false
		}

		return slices.ContainsFunc(agents, fn)
	}, "when-fact: no matching agent")

	return s
}

// agentsFrom reads the agent list a named earlier step produced. It
// reports false, and records why, when the step is missing or listed no
// agents — the two are worth telling apart in the skip reason.
func (s *Step) agentsFrom(
	results engine.Results,
	agentListStep string,
) ([]osapi.Agent, bool) {
	var list osapi.AgentList

	r := Results{results: results}
	if err := r.Decode(agentListStep, &list); err != nil {
		s.task.SetGuardReason(
			fmt.Sprintf(
				"when-fact: step %q not found or decode failed: %v",
				agentListStep,
				err,
			),
		)

		return nil, false
	}

	if len(list.Agents) == 0 {
		s.task.SetGuardReason("when-fact: no agents returned")

		return nil, false
	}

	return list.Agents, true
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
