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

package engine

import (
	"fmt"
	"math"
	"time"
)

// ErrorStrategy defines what the runner does when a task fails: stop
// the plan, or continue with the tasks that do not depend on it. A
// retry is separate (RetryPolicy), so a task can retry and then
// continue.
type ErrorStrategy struct {
	kind string
}

// StopAll cancels all remaining tasks on first failure.
var StopAll = ErrorStrategy{kind: "stop_all"}

// Continue skips dependents of the failed task but continues
// independent tasks.
var Continue = ErrorStrategy{kind: "continue"}

// String returns a human-readable representation of the strategy.
func (e ErrorStrategy) String() string {
	return e.kind
}

// RetryPolicy is how many times the runner runs a failed task again,
// and the wait between the attempts.
type RetryPolicy struct {
	count           int
	initialInterval time.Duration
	maxInterval     time.Duration
}

// RetryOption configures retry behavior.
type RetryOption func(*RetryPolicy)

// WithRetryBackoff enables exponential backoff between retry attempts
// with the given initial and maximum intervals.
func WithRetryBackoff(
	initial time.Duration,
	maxInterval time.Duration,
) RetryOption {
	return func(p *RetryPolicy) {
		p.initialInterval = initial
		p.maxInterval = maxInterval
	}
}

// Retry returns a policy that runs a failed task again n times. A
// negative n is 0. Options configure backoff behavior.
func Retry(
	n int,
	opts ...RetryOption,
) RetryPolicy {
	p := RetryPolicy{count: max(n, 0)}
	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// Count returns the number of retries of the policy.
func (p RetryPolicy) Count() int {
	return p.count
}

// String returns a human-readable representation of the policy.
func (p RetryPolicy) String() string {
	return fmt.Sprintf("retry(%d)", p.count)
}

// Intervals returns the initial and the maximum interval of the
// backoff between attempts. An initial interval of 0 means no wait.
func (p RetryPolicy) Intervals() (time.Duration, time.Duration) {
	return p.initialInterval, p.maxInterval
}

// defaultMaxInterval caps the backoff of a policy whose maximum
// interval is 0 or less, so that a long run of retries never waits
// longer than this between attempts.
const defaultMaxInterval = 30 * time.Second

// backoffDelay returns the wait before retry attempt+1: the initial
// interval doubled for each earlier attempt, and never above the cap.
// The cap is the maximum interval, or defaultMaxInterval when the
// maximum is 0 or less. An initial interval above defaultMaxInterval is
// its own cap.
func (p RetryPolicy) backoffDelay(
	attempt int,
) time.Duration {
	limit := p.maxInterval
	if limit <= 0 {
		limit = max(defaultMaxInterval, p.initialInterval)
	}

	delay := float64(p.initialInterval) * math.Pow(2, float64(attempt))
	if delay >= float64(limit) {
		return limit
	}

	return time.Duration(delay)
}

// Hooks provides consumer-controlled callbacks for plan execution
// events. All fields are optional — nil callbacks are skipped.
// The SDK performs no logging; hooks are the only output mechanism.
type Hooks struct {
	BeforePlan  func(summary PlanSummary)
	AfterPlan   func(report *Report)
	BeforeLevel func(level int, tasks []*Task, parallel bool)
	AfterLevel  func(level int, results []TaskResult)
	BeforeTask  func(task *Task)
	AfterTask   func(task *Task, result TaskResult)
	OnRetry     func(task *Task, attempt int, err error)
	OnSkip      func(task *Task, reason string)
}

// PlanConfig holds plan-level configuration.
type PlanConfig struct {
	OnErrorStrategy ErrorStrategy
	Retry           RetryPolicy
	Hooks           *Hooks
}

// PlanOption is a functional option for NewPlan.
type PlanOption func(*PlanConfig)

// OnError returns a PlanOption that sets the default error strategy.
func OnError(
	strategy ErrorStrategy,
) PlanOption {
	return func(cfg *PlanConfig) {
		cfg.OnErrorStrategy = strategy
	}
}

// WithRetry returns a PlanOption that sets the default retry policy.
func WithRetry(
	policy RetryPolicy,
) PlanOption {
	return func(cfg *PlanConfig) {
		cfg.Retry = policy
	}
}

// WithHooks attaches lifecycle callbacks to plan execution.
func WithHooks(
	hooks Hooks,
) PlanOption {
	return func(cfg *PlanConfig) {
		cfg.Hooks = &hooks
	}
}
