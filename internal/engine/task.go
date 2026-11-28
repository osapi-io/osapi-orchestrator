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
	"context"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
)

// TaskFn is the signature for functional tasks. The client
// parameter provides access to the OSAPI SDK for making API calls.
type TaskFn func(
	ctx context.Context,
	client *osapi.Client,
) (*Result, error)

// TaskFnWithResults is like TaskFn but receives completed task results
// for inter-task data access.
type TaskFnWithResults func(
	ctx context.Context,
	client *osapi.Client,
	results Results,
) (*Result, error)

// GuardFn is a predicate that determines if a task should run.
type GuardFn func(results Results) bool

// CheckFn is a guard that also gives its skip reason. The reason is
// used only when run is false.
type CheckFn func(results Results) (run bool, reason string)

// Task is a unit of work in an orchestration plan.
type Task struct {
	name           string
	fn             TaskFn
	fnr            TaskFnWithResults
	deps           []*Task
	guards         []CheckFn
	runsOnFailure  bool
	requiresChange bool
	errorStrategy  *ErrorStrategy
	retry          *RetryPolicy
	// unsafe is a task whose second run does not do what one run does,
	// such as a command or a reboot. The runner retries it only when
	// retryUnsafe is set.
	unsafe      bool
	retryUnsafe bool
}

// NewTaskFunc creates a functional task with custom logic.
func NewTaskFunc(
	name string,
	fn TaskFn,
) *Task {
	return &Task{
		name: name,
		fn:   fn,
	}
}

// NewTaskFuncWithResults creates a functional task that receives
// completed results from prior tasks.
func NewTaskFuncWithResults(
	name string,
	fn TaskFnWithResults,
) *Task {
	return &Task{
		name: name,
		fnr:  fn,
	}
}

// Name returns the task name.
func (t *Task) Name() string {
	return t.name
}

// SetName changes the task name.
func (t *Task) SetName(
	name string,
) {
	t.name = name
}

// IsFunc returns true if this is a functional task.
func (t *Task) IsFunc() bool {
	return t.fn != nil || t.fnr != nil
}

// Fn returns the task function, or nil if not set.
func (t *Task) Fn() TaskFn {
	return t.fn
}

// DependsOn sets this task's dependencies. Returns the task for
// chaining.
func (t *Task) DependsOn(
	deps ...*Task,
) *Task {
	t.deps = append(t.deps, deps...)

	return t
}

// Dependencies returns the task's dependencies.
func (t *Task) Dependencies() []*Task {
	return t.deps
}

// OnlyIfChanged marks this task to only run if at least one
// dependency reported Changed=true.
func (t *Task) OnlyIfChanged() {
	t.requiresChange = true
}

// RequiresChange returns true if OnlyIfChanged was set.
func (t *Task) RequiresChange() bool {
	return t.requiresChange
}

// When adds a guard. The task runs only when each of its guards
// returns true.
func (t *Task) When(
	fn GuardFn,
) {
	t.WhenWithReason(fn, "guard returned false")
}

// WhenWithReason adds a guard with the skip reason shown when the
// guard returns false.
func (t *Task) WhenWithReason(
	fn GuardFn,
	reason string,
) {
	t.Check(func(results Results) (bool, string) {
		return fn(results), reason
	})
}

// Check adds a guard that gives its own skip reason.
func (t *Task) Check(
	fn CheckFn,
) {
	t.guards = append(t.guards, fn)
}

// HasGuard reports whether the task has one or more guards.
func (t *Task) HasGuard() bool {
	return len(t.guards) > 0
}

// RunsOnFailure marks a task whose guards inspect a failure, such as
// OnlyIfFailed. Such a task runs after a failed dependency when its
// guards hold. Each other task is skipped after a failed dependency.
func (t *Task) RunsOnFailure() {
	t.runsOnFailure = true
}

// Evaluate runs the guards of the task in order. It returns false and
// the skip reason of the first guard that does not hold, or true when
// each guard holds.
func (t *Task) Evaluate(
	results Results,
) (bool, string) {
	for _, guard := range t.guards {
		if run, reason := guard(results); !run {
			return false, reason
		}
	}

	return true, ""
}

// Retry sets how many times the task runs again when it fails. It
// overrides the retry policy of the plan.
func (t *Task) Retry(
	policy RetryPolicy,
) {
	t.retry = &policy
}

// SetUnsafe marks the task as not idempotent: a second run does not do
// what one run does, such as a command, a reboot or a signal. The
// runner does not retry it unless SetRetryUnsafe asks for it.
func (t *Task) SetUnsafe() {
	t.unsafe = true
}

// Unsafe reports whether the task is not idempotent.
func (t *Task) Unsafe() bool {
	return t.unsafe
}

// SetRetryUnsafe lets the runner retry the task although it is not
// idempotent.
func (t *Task) SetRetryUnsafe() {
	t.retryUnsafe = true
}

// RetryPolicy returns the per-task retry policy, or nil to use the
// plan default.
func (t *Task) RetryPolicy() *RetryPolicy {
	return t.retry
}

// OnError sets a per-task error strategy override.
func (t *Task) OnError(
	strategy ErrorStrategy,
) {
	t.errorStrategy = &strategy
}

// ErrorStrategy returns the per-task error strategy, or nil to
// use the plan default.
func (t *Task) ErrorStrategy() *ErrorStrategy {
	return t.errorStrategy
}
