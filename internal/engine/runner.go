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
	"errors"
	"fmt"
	"sync"
	"time"
)

// runner executes a validated plan.
type runner struct {
	plan    *Plan
	results Results
	failed  map[string]bool
	mu      sync.Mutex
}

// newRunner creates a runner for the plan.
func newRunner(
	plan *Plan,
) *runner {
	return &runner{
		plan:    plan,
		results: make(Results),
		failed:  make(map[string]bool),
	}
}

// run executes the plan by levelizing the DAG and running each
// level in parallel.
func (r *runner) run(
	ctx context.Context,
) (*Report, error) {
	start := time.Now()
	levels := levelize(r.plan.tasks)

	r.callBeforePlan(buildPlanSummary(r.plan.tasks, levels))

	var taskResults []TaskResult

	for i, level := range levels {
		r.callBeforeLevel(i, level, len(level) > 1)

		results, err := r.runLevel(ctx, level)
		taskResults = append(taskResults, results...)

		r.callAfterLevel(i, results)

		// A cancelled context stops the plan after this level, and the
		// run reports the cancellation, not a success.
		if err == nil {
			err = ctx.Err()
		}

		if err != nil {
			for _, rest := range levels[i+1:] {
				for _, t := range rest {
					taskResults = append(taskResults, r.skip(t, start, "plan stopped"))
				}
			}

			report := &Report{
				Tasks:    taskResults,
				Duration: time.Since(start),
			}

			r.callAfterPlan(report)

			return report, err
		}
	}

	report := &Report{
		Tasks:    taskResults,
		Duration: time.Since(start),
	}

	r.callAfterPlan(report)

	return report, nil
}

// hook returns the plan's hooks or nil.
func (r *runner) hook() *Hooks {
	return r.plan.config.Hooks
}

// callBeforePlan invokes the BeforePlan hook if set.
func (r *runner) callBeforePlan(
	summary PlanSummary,
) {
	if h := r.hook(); h != nil && h.BeforePlan != nil {
		h.BeforePlan(summary)
	}
}

// buildPlanSummary creates a PlanSummary from tasks and levels.
func buildPlanSummary(
	tasks []*Task,
	levels [][]*Task,
) PlanSummary {
	steps := make([]StepSummary, len(levels))
	for i, level := range levels {
		names := make([]string, len(level))
		for j, t := range level {
			names[j] = t.name
		}

		steps[i] = StepSummary{
			Tasks:    names,
			Parallel: len(level) > 1,
		}
	}

	return PlanSummary{
		TotalTasks: len(tasks),
		Steps:      steps,
	}
}

// callAfterPlan invokes the AfterPlan hook if set.
func (r *runner) callAfterPlan(
	report *Report,
) {
	if h := r.hook(); h != nil && h.AfterPlan != nil {
		h.AfterPlan(report)
	}
}

// callBeforeLevel invokes the BeforeLevel hook if set.
func (r *runner) callBeforeLevel(
	level int,
	tasks []*Task,
	parallel bool,
) {
	if h := r.hook(); h != nil && h.BeforeLevel != nil {
		h.BeforeLevel(level, tasks, parallel)
	}
}

// callAfterLevel invokes the AfterLevel hook if set.
func (r *runner) callAfterLevel(
	level int,
	results []TaskResult,
) {
	if h := r.hook(); h != nil && h.AfterLevel != nil {
		h.AfterLevel(level, results)
	}
}

// callBeforeTask invokes the BeforeTask hook if set.
func (r *runner) callBeforeTask(
	task *Task,
) {
	if h := r.hook(); h != nil && h.BeforeTask != nil {
		h.BeforeTask(task)
	}
}

// callAfterTask invokes the AfterTask hook if set.
func (r *runner) callAfterTask(
	task *Task,
	result TaskResult,
) {
	if h := r.hook(); h != nil && h.AfterTask != nil {
		h.AfterTask(task, result)
	}
}

// callOnRetry invokes the OnRetry hook if set.
func (r *runner) callOnRetry(
	task *Task,
	attempt int,
	err error,
) {
	if h := r.hook(); h != nil && h.OnRetry != nil {
		h.OnRetry(task, attempt, err)
	}
}

// callOnSkip invokes the OnSkip hook if set.
func (r *runner) callOnSkip(
	task *Task,
	reason string,
) {
	if h := r.hook(); h != nil && h.OnSkip != nil {
		h.OnSkip(task, reason)
	}
}

// skip records a task as skipped with its reason, and calls the OnSkip
// and AfterTask hooks.
func (r *runner) skip(
	t *Task,
	start time.Time,
	reason string,
) TaskResult {
	r.mu.Lock()
	r.results[t.name] = &Result{Status: StatusSkipped}
	r.mu.Unlock()

	tr := TaskResult{
		Name:     t.name,
		Status:   StatusSkipped,
		Duration: time.Since(start),
		Reason:   reason,
	}

	r.callOnSkip(t, reason)
	r.callAfterTask(t, tr)

	return tr
}

// effectiveStrategy returns the error strategy for a task,
// checking the per-task override before falling back to the
// plan-level default.
func (r *runner) effectiveStrategy(
	t *Task,
) ErrorStrategy {
	if t.errorStrategy != nil {
		return *t.errorStrategy
	}

	return r.plan.config.OnErrorStrategy
}

// runLevel executes all tasks in a level concurrently. The first task
// that fails with a strategy other than Continue cancels the context of
// its siblings, so StopAll stops the tasks that still run. The error
// joins each such failure of the level, with the name of its task. A
// sibling that only saw that cancel is not in the error.
func (r *runner) runLevel(
	ctx context.Context,
	tasks []*Task,
) ([]TaskResult, error) {
	results := make([]TaskResult, len(tasks))

	levelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		wg       sync.WaitGroup
		stopMu   sync.Mutex
		stopped  bool
		stopErrs []error
	)

	for i, t := range tasks {
		wg.Add(1)

		go func() {
			defer wg.Done()

			tr := r.safeRunTask(levelCtx, t)
			results[i] = tr

			if tr.Status != StatusFailed || r.effectiveStrategy(t).kind == "continue" {
				return
			}

			stopMu.Lock()
			defer stopMu.Unlock()

			if stopped && ctx.Err() == nil && errors.Is(tr.Error, context.Canceled) {
				return
			}

			stopped = true
			stopErrs = append(stopErrs, fmt.Errorf("step %q: %w", t.name, tr.Error))

			cancel()
		}()
	}

	wg.Wait()

	return results, errors.Join(stopErrs...)
}

// safeRunTask runs a task and turns a panic in the task, its guard or a
// hook into a failed result, so that one task cannot stop the process.
func (r *runner) safeRunTask(
	ctx context.Context,
	t *Task,
) (tr TaskResult) {
	start := time.Now()

	defer func() {
		if p := recover(); p != nil {
			r.mu.Lock()
			r.failed[t.name] = true
			r.results[t.name] = &Result{Status: StatusFailed}
			r.mu.Unlock()

			tr = TaskResult{
				Name:     t.name,
				Status:   StatusFailed,
				Duration: time.Since(start),
				Error:    fmt.Errorf("task %q panicked: %v", t.name, p),
			}
		}
	}()

	return r.runTask(ctx, t)
}

// snapshot returns a copy of the results, made under the lock. A task
// and a guard read the copy, while the tasks beside it write the map.
func (r *runner) snapshot() Results {
	r.mu.Lock()
	defer r.mu.Unlock()

	results := make(Results, len(r.results))
	for name, res := range r.results {
		results[name] = res
	}

	return results
}

// runTask executes a single task with guard checks.
func (r *runner) runTask(
	ctx context.Context,
	t *Task,
) TaskResult {
	start := time.Now()

	// Skip if any dependency failed, unless a guard of the task
	// inspects the failure (RunsOnFailure), such as OnlyIfFailed.
	if !t.runsOnFailure {
		r.mu.Lock()

		depFailed := false

		for _, dep := range t.deps {
			if r.failed[dep.name] {
				depFailed = true
				r.failed[t.name] = true

				break
			}
		}

		r.mu.Unlock()

		if depFailed {
			return r.skip(t, start, "dependency failed")
		}
	}

	if t.requiresChange {
		anyChanged := false

		r.mu.Lock()

		for _, dep := range t.deps {
			if res := r.results.Get(dep.name); res != nil && res.Changed {
				anyChanged = true

				break
			}
		}

		r.mu.Unlock()

		if !anyChanged {
			return r.skip(t, start, "no dependencies changed")
		}
	}

	// Each guard must hold. The first one that does not gives the skip
	// reason.
	if len(t.guards) > 0 {
		if run, reason := t.Evaluate(r.snapshot()); !run {
			return r.skip(t, start, reason)
		}
	}

	r.callBeforeTask(t)

	policy := r.plan.config.Retry
	if t.retry != nil {
		policy = *t.retry
	}

	maxAttempts := policy.count + 1

	var result *Result
	var err error

	client := r.plan.client

retryLoop:
	for attempt := range maxAttempts {
		// A cancelled context ends the task; a retry does not run.
		if err = ctx.Err(); err != nil {
			result = nil

			break
		}

		if t.fnr != nil {
			result, err = t.fnr(ctx, client, r.snapshot())
		} else {
			result, err = t.fn(ctx, client)
		}

		if err == nil && result == nil {
			err = Permanent(fmt.Errorf("task %q returned no result", t.name))
		}

		if err == nil {
			break
		}

		if ctx.Err() != nil {
			break
		}

		// A retry runs only for an error that a retry can fix, and only
		// for a task that is safe to run again.
		if !Retryable(err) || (t.unsafe && !t.retryUnsafe) {
			break
		}

		if attempt < maxAttempts-1 {
			r.callOnRetry(t, attempt+1, err)

			// Backoff between retries if configured.
			if policy.initialInterval > 0 {
				delay := policy.backoffDelay(attempt)

				select {
				case <-ctx.Done():
					err = ctx.Err()

					break retryLoop
				case <-time.After(delay):
				}
			}
		}
	}

	elapsed := time.Since(start)

	if err != nil {
		// Preserve the full result for partial failures (e.g. broadcast
		// commands where some hosts succeeded and others failed).
		// Guards like OnlyIfChanged and OnlyIfFailed inspect r.results,
		// so they need Changed, HostResults, etc.
		failedResult := &Result{Status: StatusFailed}
		if result != nil {
			result.Status = StatusFailed
			failedResult = result
		}

		r.mu.Lock()
		r.failed[t.name] = true
		r.results[t.name] = failedResult
		r.mu.Unlock()

		tr := TaskResult{
			Name:     t.name,
			Status:   StatusFailed,
			Duration: elapsed,
			Error:    err,
		}

		if result != nil {
			tr.JobID = result.JobID
			tr.JobDuration = result.JobDuration
			tr.Data = result.Data
			tr.HostResults = result.HostResults
			tr.Changed = result.Changed
		}

		r.callAfterTask(t, tr)

		return tr
	}

	status := StatusUnchanged
	if result.Changed {
		status = StatusChanged
	}

	result.Status = status

	r.mu.Lock()
	r.results[t.name] = result
	r.mu.Unlock()

	tr := TaskResult{
		JobID:       result.JobID,
		Name:        t.name,
		Status:      status,
		Changed:     result.Changed,
		Duration:    elapsed,
		JobDuration: result.JobDuration,
		Data:        result.Data,
		HostResults: result.HostResults,
	}

	r.callAfterTask(t, tr)

	return tr
}

// levelize groups tasks into levels where all tasks in a level can
// run concurrently (all dependencies are in earlier levels).
func levelize(
	tasks []*Task,
) [][]*Task {
	level := make(map[string]int, len(tasks))

	var computeLevel func(t *Task) int
	computeLevel = func(t *Task) int {
		if l, ok := level[t.name]; ok {
			return l
		}

		maxDep := -1

		for _, dep := range t.deps {
			depLevel := computeLevel(dep)
			if depLevel > maxDep {
				maxDep = depLevel
			}
		}

		level[t.name] = maxDep + 1

		return maxDep + 1
	}

	maxLevel := 0

	for _, t := range tasks {
		l := computeLevel(t)
		if l > maxLevel {
			maxLevel = l
		}
	}

	levels := make([][]*Task, maxLevel+1)

	for _, t := range tasks {
		l := level[t.name]
		levels[l] = append(levels[l], t)
	}

	return levels
}
