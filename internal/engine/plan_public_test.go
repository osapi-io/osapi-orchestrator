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

package engine_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

type PlanPublicTestSuite struct {
	suite.Suite
}

func TestPlanPublicTestSuite(t *testing.T) {
	suite.Run(t, new(PlanPublicTestSuite))
}

// taskFunc creates a TaskFn with the given changed value and optional
// side effect.
func taskFunc(
	changed bool,
	sideEffect func(),
) engine.TaskFn {
	return func(
		_ context.Context,
		_ *osapi.Client,
	) (*engine.Result, error) {
		if sideEffect != nil {
			sideEffect()
		}

		return &engine.Result{Changed: changed}, nil
	}
}

// failFunc creates a TaskFn that always returns the given error.
func failFunc(
	msg string,
) engine.TaskFn {
	return func(
		_ context.Context,
		_ *osapi.Client,
	) (*engine.Result, error) {
		return nil, fmt.Errorf("%s", msg)
	}
}

// statusMap builds a name→status map from a report for easy assertions.
func statusMap(
	report *engine.Report,
) map[string]engine.Status {
	m := make(map[string]engine.Status, len(report.Tasks))
	for _, r := range report.Tasks {
		m[r.Name] = r.Status
	}

	return m
}

// filterPrefix returns only strings that start with prefix.
func filterPrefix(
	ss []string,
	prefix string,
) []string {
	var out []string
	for _, s := range ss {
		if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
			out = append(out, s)
		}
	}

	return out
}

func (s *PlanPublicTestSuite) TestRun() {
	tests := []struct {
		name         string
		validateFunc func()
	}{
		{
			name: "linear chain executes in order",
			validateFunc: func() {
				var order []string
				plan := engine.NewPlan(nil)

				mk := func(name string, changed bool) *engine.Task {
					n := name

					return plan.TaskFunc(n, taskFunc(changed, func() {
						order = append(order, n)
					}))
				}

				a := mk("a", true)
				b := mk("b", true)
				c := mk("c", false)
				b.DependsOn(a)
				c.DependsOn(b)

				report, err := plan.Run(context.Background())
				s.Require().NoError(err)
				s.Equal([]string{"a", "b", "c"}, order)
				s.Len(report.Tasks, 3)
				s.Contains(report.Summary(), "2 changed")
				s.Contains(report.Summary(), "1 unchanged")
			},
		},
		{
			name: "parallel tasks run concurrently",
			validateFunc: func() {
				var concurrentMax atomic.Int32
				var concurrent atomic.Int32

				plan := engine.NewPlan(nil)

				for _, name := range []string{"a", "b", "c"} {
					plan.TaskFunc(name, func(
						_ context.Context,
						_ *osapi.Client,
					) (*engine.Result, error) {
						cur := concurrent.Add(1)
						for {
							prev := concurrentMax.Load()
							if cur > prev {
								if concurrentMax.CompareAndSwap(prev, cur) {
									break
								}
							} else {
								break
							}
						}
						concurrent.Add(-1)

						return &engine.Result{Changed: false}, nil
					})
				}

				report, err := plan.Run(context.Background())
				s.Require().NoError(err)
				s.Len(report.Tasks, 3)
				s.GreaterOrEqual(int(concurrentMax.Load()), 1)
			},
		},
		{
			name: "cycle detection returns error",
			validateFunc: func() {
				plan := engine.NewPlan(nil)
				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(false, nil))
				a.DependsOn(b)
				b.DependsOn(a)

				_, err := plan.Run(context.Background())
				s.Error(err)
				s.Contains(err.Error(), "cycle")
			},
		},
		{
			name: "stop all on error",
			validateFunc: func() {
				plan := engine.NewPlan(nil)
				didRun := false

				fail := plan.TaskFunc("fail", failFunc("boom"))
				next := plan.TaskFunc("next", taskFunc(false, func() {
					didRun = true
				}))
				next.DependsOn(fail)

				_, err := plan.Run(context.Background())
				s.Error(err)
				s.False(didRun)
			},
		},
		{
			name: "continue on error runs independent tasks",
			validateFunc: func() {
				plan := engine.NewPlan(
					nil,
					engine.OnError(engine.Continue),
				)
				didRun := false

				plan.TaskFunc("fail", failFunc("boom"))
				plan.TaskFunc("independent", taskFunc(true, func() {
					didRun = true
				}))

				report, err := plan.Run(context.Background())
				s.Error(err, "a task that failed with Continue makes Run fail")
				s.True(didRun)
				s.Len(report.Tasks, 2)
			},
		},
		{
			name: "continue skips dependents of failed task",
			validateFunc: func() {
				plan := engine.NewPlan(
					nil,
					engine.OnError(engine.Continue),
				)

				a := plan.TaskFunc("a", failFunc("a failed"))
				plan.TaskFunc("b", taskFunc(true, nil)).DependsOn(a)
				plan.TaskFunc("c", taskFunc(true, nil))

				report, err := plan.Run(context.Background())
				s.Error(err, "a task that failed with Continue makes Run fail")
				s.Len(report.Tasks, 3)
				m := statusMap(report)
				s.Equal(engine.StatusFailed, m["a"])
				s.Equal(engine.StatusSkipped, m["b"])
				s.Equal(engine.StatusChanged, m["c"])
			},
		},
		{
			name: "continue transitive skip",
			validateFunc: func() {
				plan := engine.NewPlan(
					nil,
					engine.OnError(engine.Continue),
				)

				a := plan.TaskFunc("a", failFunc("a failed"))
				b := plan.TaskFunc("b", taskFunc(true, nil))
				b.DependsOn(a)
				c := plan.TaskFunc("c", taskFunc(true, nil))
				c.DependsOn(b)

				report, err := plan.Run(context.Background())
				s.Error(err, "a task that failed with Continue makes Run fail")
				m := statusMap(report)
				s.Equal(engine.StatusFailed, m["a"])
				s.Equal(engine.StatusSkipped, m["b"])
				s.Equal(engine.StatusSkipped, m["c"])
			},
		},
		{
			name: "per-task continue override",
			validateFunc: func() {
				plan := engine.NewPlan(nil) // default StopAll

				a := plan.TaskFunc("a", failFunc("a failed"))
				a.OnError(engine.Continue)
				plan.TaskFunc("b", taskFunc(true, nil))

				report, err := plan.Run(context.Background())
				s.Error(err, "a task that failed with Continue makes Run fail")
				m := statusMap(report)
				s.Equal(engine.StatusFailed, m["a"])
				s.Equal(engine.StatusChanged, m["b"])
			},
		},
		{
			name: "retry succeeds after transient failure",
			validateFunc: func() {
				attempts := 0
				plan := engine.NewPlan(
					nil,
					engine.WithRetry(engine.Retry(2)),
				)

				plan.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					attempts++
					if attempts < 3 {
						return nil, fmt.Errorf("attempt %d failed", attempts)
					}

					return &engine.Result{Changed: true}, nil
				})

				report, err := plan.Run(context.Background())
				s.NoError(err)
				s.Equal(3, attempts)
				s.Equal(engine.StatusChanged, report.Tasks[0].Status)
			},
		},
		{
			name: "retry exhausted returns error",
			validateFunc: func() {
				plan := engine.NewPlan(
					nil,
					engine.WithRetry(engine.Retry(1)),
				)

				plan.TaskFunc("always-fail", failFunc("permanent failure"))

				report, err := plan.Run(context.Background())
				s.Error(err)
				s.Equal(engine.StatusFailed, report.Tasks[0].Status)
			},
		},
		{
			name: "per-task retry override",
			validateFunc: func() {
				attempts := 0
				plan := engine.NewPlan(nil) // default StopAll

				plan.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					attempts++
					if attempts < 2 {
						return nil, fmt.Errorf("attempt %d failed", attempts)
					}

					return &engine.Result{Changed: true}, nil
				}).Retry(engine.Retry(1))

				report, err := plan.Run(context.Background())
				s.NoError(err)
				s.Equal(2, attempts)
				s.Equal(engine.StatusChanged, report.Tasks[0].Status)
			},
		},
		{
			name: "retry with backoff delays between attempts",
			validateFunc: func() {
				attempts := 0
				var timestamps []time.Time

				plan := engine.NewPlan(
					nil,
					engine.WithRetry(engine.Retry(
						2,
						engine.WithRetryBackoff(
							50*time.Millisecond,
							200*time.Millisecond,
						),
					)),
				)

				plan.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					timestamps = append(timestamps, time.Now())
					attempts++
					if attempts < 3 {
						return nil, fmt.Errorf("attempt %d failed", attempts)
					}

					return &engine.Result{Changed: true}, nil
				})

				report, err := plan.Run(context.Background())
				s.NoError(err)
				s.Equal(3, attempts)
				s.Equal(engine.StatusChanged, report.Tasks[0].Status)

				// Verify delays between attempts: the jitter keeps at least half
				// of the initial interval.
				s.Require().Len(timestamps, 3)
				gap1 := timestamps[1].Sub(timestamps[0])
				s.GreaterOrEqual(gap1, 20*time.Millisecond)
			},
		},
		{
			name: "retry without backoff retries immediately",
			validateFunc: func() {
				attempts := 0
				var timestamps []time.Time

				plan := engine.NewPlan(
					nil,
					engine.WithRetry(engine.Retry(1)),
				)

				plan.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					timestamps = append(timestamps, time.Now())
					attempts++
					if attempts < 2 {
						return nil, fmt.Errorf("attempt %d failed", attempts)
					}

					return &engine.Result{Changed: true}, nil
				})

				report, err := plan.Run(context.Background())
				s.NoError(err)
				s.Equal(2, attempts)
				s.Equal(engine.StatusChanged, report.Tasks[0].Status)

				// Without backoff, retries should be near-instant.
				s.Require().Len(timestamps, 2)
				gap := timestamps[1].Sub(timestamps[0])
				s.Less(gap, 20*time.Millisecond)
			},
		},
		{
			name: "retry backoff respects context cancellation",
			validateFunc: func() {
				attempts := 0

				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
				defer cancel()

				plan := engine.NewPlan(
					nil,
					engine.WithRetry(engine.Retry(
						10,
						engine.WithRetryBackoff(
							5*time.Second,
							30*time.Second,
						),
					)),
				)

				plan.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					attempts++

					return nil, fmt.Errorf("attempt %d failed", attempts)
				})

				_, err := plan.Run(ctx)
				s.Error(err)
				// Should only get 1 attempt because backoff (5s) is longer
				// than context timeout (30ms).
				s.Equal(1, attempts)
			},
		},
		{
			name: "all hooks called in order",
			validateFunc: func() {
				var events []string

				hooks := allHooks(&events)
				plan := engine.NewPlan(nil, engine.WithHooks(hooks))

				a := plan.TaskFunc("a", taskFunc(true, nil))
				plan.TaskFunc("b", taskFunc(false, nil)).DependsOn(a)

				report, err := plan.Run(context.Background())
				s.NoError(err)
				s.NotNil(report)
				s.Equal([]string{
					"before-plan",
					"before-level-0",
					"before-a",
					"after-a",
					"after-level-0",
					"before-level-1",
					"before-b",
					"after-b",
					"after-level-1",
					"after-plan",
				}, events)
			},
		},
		{
			name: "a task that panics still calls AfterTask",
			validateFunc: func() {
				var events []string

				hooks := allHooks(&events)
				plan := engine.NewPlan(
					nil,
					engine.WithHooks(hooks),
					engine.OnError(engine.Continue),
				)
				plan.TaskFunc("p", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					panic("boom")
				})

				_, err := plan.Run(context.Background())
				s.ErrorContains(err, `task "p" panicked`)
				s.Contains(events, "before-p")
				s.Contains(events, "after-p")
			},
		},
		{
			name: "retry hook called with correct args",
			validateFunc: func() {
				var events []string
				attempts := 0

				hooks := allHooks(&events)
				plan := engine.NewPlan(
					nil,
					engine.WithHooks(hooks),
					engine.WithRetry(engine.Retry(2)),
				)

				plan.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					attempts++
					if attempts < 3 {
						return nil, fmt.Errorf("fail-%d", attempts)
					}

					return &engine.Result{Changed: true}, nil
				})

				_, err := plan.Run(context.Background())
				s.NoError(err)

				retries := filterPrefix(events, "retry-")
				s.Len(retries, 2)
				s.Contains(retries[0], "retry-flaky-1-fail-1")
				s.Contains(retries[1], "retry-flaky-2-fail-2")
			},
		},
		{
			name: "skip hook for dependency failure",
			validateFunc: func() {
				var events []string

				hooks := allHooks(&events)
				plan := engine.NewPlan(
					nil,
					engine.WithHooks(hooks),
					engine.OnError(engine.Continue),
				)

				a := plan.TaskFunc("a", failFunc("a failed"))
				plan.TaskFunc("b", taskFunc(true, nil)).DependsOn(a)

				_, err := plan.Run(context.Background())
				s.ErrorContains(err, `step "a": a failed`)

				skips := filterPrefix(events, "skip-")
				s.Len(skips, 1)
				s.Contains(skips[0], "skip-b-dependency failed")
			},
		},
		{
			name: "skip hook for guard",
			validateFunc: func() {
				var events []string

				hooks := allHooks(&events)
				plan := engine.NewPlan(nil, engine.WithHooks(hooks))

				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(true, nil))
				b.DependsOn(a)
				b.When(func(_ engine.Results) bool { return false })

				_, err := plan.Run(context.Background())
				s.NoError(err)

				skips := filterPrefix(events, "skip-")
				s.Len(skips, 1)
				s.Contains(skips[0], "guard returned false")
			},
		},
		{
			name: "skip hook for guard with custom reason",
			validateFunc: func() {
				var events []string

				hooks := allHooks(&events)
				plan := engine.NewPlan(nil, engine.WithHooks(hooks))

				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(true, nil))
				b.DependsOn(a)
				b.WhenWithReason(
					func(_ engine.Results) bool { return false },
					"host is unreachable",
				)

				_, err := plan.Run(context.Background())
				s.NoError(err)

				skips := filterPrefix(events, "skip-")
				s.Len(skips, 1)
				s.Contains(skips[0], "host is unreachable")
			},
		},
		{
			name: "skip hook for only-if-changed",
			validateFunc: func() {
				var events []string

				hooks := allHooks(&events)
				plan := engine.NewPlan(nil, engine.WithHooks(hooks))

				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(true, nil))
				b.DependsOn(a)
				b.OnlyIfChanged()

				_, err := plan.Run(context.Background())
				s.NoError(err)

				skips := filterPrefix(events, "skip-")
				s.Len(skips, 1)
				s.Contains(skips[0], "no dependencies changed")
			},
		},
		{
			name: "retry without hooks configured",
			validateFunc: func() {
				attempts := 0
				plan := engine.NewPlan(
					nil,
					engine.WithRetry(engine.Retry(1)),
				)

				plan.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
				) (*engine.Result, error) {
					attempts++
					if attempts < 2 {
						return nil, fmt.Errorf("attempt %d", attempts)
					}

					return &engine.Result{Changed: true}, nil
				})

				report, err := plan.Run(context.Background())
				s.NoError(err)
				s.Equal(2, attempts)
				s.Equal(engine.StatusChanged, report.Tasks[0].Status)
			},
		},
		{
			name: "skip without hooks configured",
			validateFunc: func() {
				plan := engine.NewPlan(
					nil,
					engine.OnError(engine.Continue),
				)

				a := plan.TaskFunc("a", failFunc("a failed"))
				plan.TaskFunc("b", taskFunc(true, nil)).DependsOn(a)

				report, err := plan.Run(context.Background())
				s.ErrorContains(err, `step "a": a failed`)
				m := statusMap(report)
				s.Equal(engine.StatusSkipped, m["b"])
			},
		},
		{
			name: "skips an OnlyIfChanged task when no dependency changed",
			validateFunc: func() {
				plan := engine.NewPlan(nil)
				ran := false

				dep := plan.TaskFunc("dep", taskFunc(false, nil))
				plan.TaskFunc("conditional", taskFunc(false, func() {
					ran = true
				})).DependsOn(dep).OnlyIfChanged()

				report, err := plan.Run(context.Background())
				s.Require().NoError(err)
				s.False(ran)
				s.Contains(report.Summary(), "skipped")
			},
		},
		{
			name: "runs an OnlyIfChanged task when a dependency changed",
			validateFunc: func() {
				plan := engine.NewPlan(nil)
				ran := false

				dep := plan.TaskFunc("dep", taskFunc(true, nil))
				plan.TaskFunc("conditional", taskFunc(false, func() {
					ran = true
				})).DependsOn(dep).OnlyIfChanged()

				report, err := plan.Run(context.Background())
				s.Require().NoError(err)
				s.True(ran)
				s.Equal(engine.StatusUnchanged, report.Tasks[1].Status)
			},
		},
		{
			name: "skips a task when its guard returns false",
			validateFunc: func() {
				plan := engine.NewPlan(nil)
				ran := false

				a := plan.TaskFunc("a", taskFunc(false, nil))
				plan.TaskFunc("b", taskFunc(true, func() {
					ran = true
				})).DependsOn(a).When(func(_ engine.Results) bool { return false })

				report, err := plan.Run(context.Background())
				s.Require().NoError(err)
				s.False(ran)
				s.Equal(engine.StatusSkipped, report.Tasks[1].Status)
			},
		},
		{
			name: "a guard that runs on failure and inspects the failure runs its task",
			validateFunc: func() {
				ran, status, _ := s.runAfterFailure(func(r engine.Results) bool {
					res := r.Get("fail")

					return res != nil && res.Status == engine.StatusFailed
				}, true)

				s.True(ran)
				s.Equal(engine.StatusChanged, status)
			},
		},
		{
			name: "a guard that runs on failure and returns false skips its task",
			validateFunc: func() {
				ran, status, reason := s.runAfterFailure(func(r engine.Results) bool {
					res := r.Get("fail")

					return res != nil && res.Status == engine.StatusChanged
				}, true)

				s.False(ran)
				s.Equal(engine.StatusSkipped, status)
				s.Equal("guard returned false", reason)
			},
		},
		{
			name: "a guard that does not run on failure is skipped after a failed dependency",
			validateFunc: func() {
				ran, status, reason := s.runAfterFailure(func(r engine.Results) bool {
					res := r.Get("fail")

					return res != nil && res.Status == engine.StatusFailed
				}, false)

				s.False(ran)
				s.Equal(engine.StatusSkipped, status)
				s.Equal("dependency failed", reason)
			},
		},
		{
			name: "an empty plan runs no level",
			validateFunc: func() {
				levels := 0
				plan := engine.NewPlan(nil, engine.WithHooks(engine.Hooks{
					BeforeLevel: func(int, []*engine.Task, bool) { levels++ },
				}))

				report, err := plan.Run(context.Background())
				s.Require().NoError(err)
				s.Empty(report.Tasks)
				s.Zero(levels)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, tc.validateFunc)
	}
}

// runAfterFailure runs a plan where the task "alert" depends on the task
// "fail", which fails, and has guard. It returns whether alert ran, its
// status, and its reason.
func (s *PlanPublicTestSuite) runAfterFailure(
	guard func(engine.Results) bool,
	runsOnFailure bool,
) (bool, engine.Status, string) {
	plan := engine.NewPlan(
		nil,
		engine.OnError(engine.Continue),
	)
	ran := false

	fail := plan.TaskFunc("fail", failFunc("boom"))
	alert := plan.TaskFunc("alert", taskFunc(true, func() {
		ran = true
	}))
	alert.DependsOn(fail)
	alert.When(guard)

	if runsOnFailure {
		alert.RunsOnFailure()
	}

	report, err := plan.Run(context.Background())
	s.Require().ErrorContains(err, `step "fail": boom`)

	reason := ""
	for _, tr := range report.Tasks {
		if tr.Name == "alert" {
			reason = tr.Reason
		}
	}

	return ran, statusMap(report)["alert"], reason
}

func (s *PlanPublicTestSuite) TestClient() {
	client := osapi.New("http://localhost", "token")

	tests := []struct {
		name         string
		client       *osapi.Client
		validateFunc func(*osapi.Client)
	}{
		{
			name:   "returns the client the plan was made with",
			client: client,
			validateFunc: func(got *osapi.Client) {
				s.Equal(client, got)
			},
		},
		{
			name: "returns nil when the plan has no client",
			validateFunc: func(got *osapi.Client) {
				s.Nil(got)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(engine.NewPlan(tt.client).Client())
		})
	}
}

func (s *PlanPublicTestSuite) TestConfig() {
	tests := []struct {
		name         string
		opts         []engine.PlanOption
		validateFunc func(engine.PlanConfig)
	}{
		{
			name: "with no options the strategy is stop_all",
			validateFunc: func(cfg engine.PlanConfig) {
				s.Equal("stop_all", cfg.OnErrorStrategy.String())
			},
		},
		{
			name: "returns the error strategy the plan was made with",
			opts: []engine.PlanOption{engine.OnError(engine.Continue)},
			validateFunc: func(cfg engine.PlanConfig) {
				s.Equal("continue", cfg.OnErrorStrategy.String())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(engine.NewPlan(nil, tt.opts...).Config())
		})
	}
}

func (s *PlanPublicTestSuite) TestTasks() {
	tests := []struct {
		name         string
		tasks        []string
		validateFunc func([]*engine.Task)
	}{
		{
			name: "a new plan has no tasks",
			validateFunc: func(got []*engine.Task) {
				s.Empty(got)
			},
		},
		{
			name:  "returns each task that was added",
			tasks: []string{"a", "b"},
			validateFunc: func(got []*engine.Task) {
				s.Require().Len(got, 2)
				s.Equal("a", got[0].Name())
				s.Equal("b", got[1].Name())
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan := engine.NewPlan(nil)
			for _, name := range tt.tasks {
				plan.TaskFunc(name, taskFunc(false, nil))
			}

			tt.validateFunc(plan.Tasks())
		})
	}
}

func (s *PlanPublicTestSuite) TestValidate() {
	tests := []struct {
		name         string
		setup        func(plan *engine.Plan)
		validateFunc func(err error)
	}{
		{
			name: "duplicate task name returns error",
			setup: func(plan *engine.Plan) {
				plan.TaskFunc("dup", taskFunc(false, nil))
				plan.TaskFunc("dup", taskFunc(false, nil))
			},
			validateFunc: func(err error) {
				s.Error(err)
				s.Contains(err.Error(), "duplicate task name")
			},
		},
		{
			name: "a dependency that is not in the plan returns error",
			setup: func(plan *engine.Plan) {
				other := engine.NewPlan(nil)
				outside := other.TaskFunc("outside", taskFunc(false, nil))
				plan.TaskFunc("a", taskFunc(false, nil)).DependsOn(outside)
			},
			validateFunc: func(err error) {
				s.EqualError(err, `task "a" depends on "outside", which is not in the plan`)
			},
		},
		{
			name: "valid plan returns nil",
			setup: func(plan *engine.Plan) {
				plan.TaskFunc("a", taskFunc(false, nil))
				plan.TaskFunc("b", taskFunc(false, nil))
			},
			validateFunc: func(err error) {
				s.NoError(err)
			},
		},
		{
			name: "a task with no function returns error",
			setup: func(plan *engine.Plan) {
				plan.TaskFunc("none", nil)
			},
			validateFunc: func(err error) {
				s.ErrorContains(err, `task "none" has no function`)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan := engine.NewPlan(nil)
			tt.setup(plan)
			tt.validateFunc(plan.Validate())
		})
	}
}

func (s *PlanPublicTestSuite) TestLevels() {
	tests := []struct {
		name         string
		setup        func(plan *engine.Plan)
		validateFunc func(levels [][]*engine.Task, err error)
	}{
		{
			name: "returns levels for valid plan",
			setup: func(plan *engine.Plan) {
				a := plan.TaskFunc("a", taskFunc(false, nil))
				plan.TaskFunc("b", taskFunc(false, nil)).DependsOn(a)
			},
			validateFunc: func(levels [][]*engine.Task, err error) {
				s.NoError(err)
				s.Len(levels, 2)
			},
		},
		{
			name:  "returns no levels for an empty plan",
			setup: func(*engine.Plan) {},
			validateFunc: func(levels [][]*engine.Task, err error) {
				s.NoError(err)
				s.Empty(levels)
			},
		},
		{
			name: "returns error for invalid plan",
			setup: func(plan *engine.Plan) {
				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(false, nil))
				a.DependsOn(b)
				b.DependsOn(a)
			},
			validateFunc: func(levels [][]*engine.Task, err error) {
				s.Error(err)
				s.Nil(levels)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan := engine.NewPlan(nil)
			tt.setup(plan)
			levels, err := plan.Levels()
			tt.validateFunc(levels, err)
		})
	}
}

func (s *PlanPublicTestSuite) TestExplain() {
	tests := []struct {
		name         string
		setup        func(plan *engine.Plan)
		validateFunc func(string)
	}{
		{
			name: "valid plan with dependencies and guards",
			setup: func(plan *engine.Plan) {
				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(false, nil))
				b.DependsOn(a)
				b.OnlyIfChanged()
				b.When(func(_ engine.Results) bool { return true })
				b.Retry(engine.Retry(2))
			},
			validateFunc: func(output string) {
				s.Contains(output, "Plan: 2 tasks, 2 levels")
				s.Contains(output, "Level 0:")
				s.Contains(output, "a [fn]")
				s.Contains(output, "Level 1:")
				s.Contains(output, "b [fn] <- a (only-if-changed, when, retry(2))")
			},
		},
		{
			name:  "an empty plan has no levels",
			setup: func(*engine.Plan) {},
			validateFunc: func(output string) {
				s.Contains(output, "Plan: 0 tasks, 0 levels")
			},
		},
		{
			name: "invalid plan returns error string",
			setup: func(plan *engine.Plan) {
				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(false, nil))
				a.DependsOn(b)
				b.DependsOn(a)
			},
			validateFunc: func(output string) {
				s.Contains(output, "invalid plan:")
				s.Contains(output, "cycle")
			},
		},
		{
			name: "parallel tasks shown as parallel",
			setup: func(plan *engine.Plan) {
				plan.TaskFunc("a", taskFunc(false, nil))
				plan.TaskFunc("b", taskFunc(false, nil))
			},
			validateFunc: func(output string) {
				s.Contains(output, "Plan: 2 tasks, 1 levels")
				s.Contains(output, "Level 0 (parallel):")
			},
		},
		{
			name: "guard shown in flags",
			setup: func(plan *engine.Plan) {
				a := plan.TaskFunc("a", taskFunc(false, nil))
				b := plan.TaskFunc("b", taskFunc(false, nil))
				b.DependsOn(a)
				b.When(func(_ engine.Results) bool { return true })
			},
			validateFunc: func(output string) {
				s.Contains(output, "when")
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			plan := engine.NewPlan(nil)
			tt.setup(plan)
			tt.validateFunc(plan.Explain())
		})
	}
}

// allHooks returns a Hooks struct that appends all events to the given
// slice, covering every hook type.
func allHooks(
	events *[]string,
) engine.Hooks {
	return engine.Hooks{
		BeforePlan: func(_ engine.PlanSummary) {
			*events = append(*events, "before-plan")
		},
		AfterPlan: func(_ *engine.Report) {
			*events = append(*events, "after-plan")
		},
		BeforeLevel: func(level int, _ []*engine.Task, _ bool) {
			*events = append(*events, fmt.Sprintf("before-level-%d", level))
		},
		AfterLevel: func(level int, _ []engine.TaskResult) {
			*events = append(*events, fmt.Sprintf("after-level-%d", level))
		},
		BeforeTask: func(task *engine.Task) {
			*events = append(*events, "before-"+task.Name())
		},
		AfterTask: func(_ *engine.Task, r engine.TaskResult) {
			*events = append(*events, "after-"+r.Name)
		},
		OnRetry: func(
			task *engine.Task,
			attempt int,
			err error,
		) {
			*events = append(
				*events,
				fmt.Sprintf("retry-%s-%d-%s", task.Name(), attempt, err),
			)
		},
		OnSkip: func(task *engine.Task, reason string) {
			*events = append(
				*events,
				fmt.Sprintf("skip-%s-%s", task.Name(), reason),
			)
		},
	}
}
