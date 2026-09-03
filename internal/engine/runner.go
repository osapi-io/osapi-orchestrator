package engine

import (
	"context"
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

		if err != nil {
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

// runLevel executes all tasks in a level concurrently.
func (r *runner) runLevel(
	ctx context.Context,
	tasks []*Task,
) ([]TaskResult, error) {
	results := make([]TaskResult, len(tasks))

	var wg sync.WaitGroup

	for i, t := range tasks {
		wg.Go(func() {
			results[i] = r.runTask(ctx, t)
		})
	}

	wg.Wait()

	if err := r.fatalFailure(tasks, results); err != nil {
		return results, err
	}

	return results, nil
}

// fatalFailure returns the error of the first failed task whose strategy
// says to stop. A task set to continue does not end the run.
func (r *runner) fatalFailure(
	tasks []*Task,
	results []TaskResult,
) error {
	for i, tr := range results {
		if tr.Status != StatusFailed {
			continue
		}

		if r.effectiveStrategy(tasks[i]).kind != kindContinue {
			return tr.Error
		}
	}

	return nil
}

// runTask executes a single task with guard checks.
func (r *runner) runTask(
	ctx context.Context,
	t *Task,
) TaskResult {
	start := time.Now()

	if reason, skip := r.skipReason(t); skip {
		return r.recordSkip(t, reason, start)
	}

	r.callBeforeTask(t)

	result, err := r.execute(ctx, t)

	elapsed := time.Since(start)
	if err != nil {
		return r.recordFailure(t, result, err, elapsed)
	}

	return r.recordSuccess(t, result, elapsed)
}

// skipReason reports why a task should not run, and whether it should
// not. The checks are ordered: a failed dependency stops a task
// outright, requiresChange asks whether anything changed, and a guard
// has the last word.
func (r *runner) skipReason(t *Task) (string, bool) {
	// A task with a guard may deliberately inspect failure — an
	// alert-on-failure step, say — so a failed dependency does not skip it.
	if t.guard == nil && r.markFailedByDep(t) {
		return "dependency failed", true
	}

	if t.requiresChange && !r.anyDepChanged(t) {
		return "no dependencies changed", true
	}

	if t.guard != nil && !r.guardPasses(t) {
		return guardReason(t), true
	}

	return "", false
}

// guardReason says why a guard rejected a task, in the guard's own words
// when it left any.
func guardReason(t *Task) string {
	if t.guardReason != "" {
		return t.guardReason
	}

	return "guard returned false"
}

// markFailedByDep reports whether any dependency failed, marking this
// task failed too when one did. The read and the write are one critical
// section so a concurrent sibling cannot see the task in between.
func (r *runner) markFailedByDep(t *Task) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, dep := range t.deps {
		if r.failed[dep.name] {
			r.failed[t.name] = true

			return true
		}
	}

	return false
}

// anyDepChanged reports whether any dependency reported a change.
func (r *runner) anyDepChanged(t *Task) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, dep := range t.deps {
		if res := r.results.Get(dep.name); res != nil && res.Changed {
			return true
		}
	}

	return false
}

// guardPasses evaluates the task's guard against the results so far.
func (r *runner) guardPasses(t *Task) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return t.guard(r.results)
}

// execute runs the task, retrying per its error strategy. It returns
// the last result and error seen, so a partial failure keeps its result.
func (r *runner) execute(
	ctx context.Context,
	t *Task,
) (*Result, error) {
	strategy := r.effectiveStrategy(t)
	maxAttempts := attemptsFor(strategy)

	var (
		result *Result
		err    error
	)

	for attempt := range maxAttempts {
		result, err = r.invoke(ctx, t)
		if err == nil {
			return result, nil
		}

		if attempt == maxAttempts-1 {
			break
		}

		r.callOnRetry(t, attempt+1, err)

		if waitErr := backoff(ctx, strategy, attempt); waitErr != nil {
			return result, waitErr
		}
	}

	return result, err
}

// attemptsFor reports how many times a task may run under a strategy.
// Only a retry strategy gets more than the one.
func attemptsFor(strategy ErrorStrategy) int {
	if strategy.kind != kindRetry {
		return 1
	}

	return strategy.retryCount + 1
}

// invoke calls the task function, handing it the results so far when it
// is the kind that asks for them.
func (r *runner) invoke(
	ctx context.Context,
	t *Task,
) (*Result, error) {
	if t.fnr == nil {
		return t.fn(ctx, r.plan.client)
	}

	r.mu.Lock()
	results := r.results
	r.mu.Unlock()

	return t.fnr(ctx, r.plan.client, results)
}

// backoff waits between attempts, returning the context's error if it
// is cancelled first. A strategy with no interval configured waits not
// at all.
func backoff(
	ctx context.Context,
	strategy ErrorStrategy,
	attempt int,
) error {
	if strategy.initialInterval <= 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(strategy.backoffDelay(attempt)):
		return nil
	}
}

// recordSkip marks a task skipped and tells the callbacks why.
func (r *runner) recordSkip(
	t *Task,
	reason string,
	start time.Time,
) TaskResult {
	r.mu.Lock()
	r.results[t.name] = &Result{Status: StatusSkipped}
	r.mu.Unlock()

	tr := TaskResult{
		Name:     t.name,
		Status:   StatusSkipped,
		Duration: time.Since(start),
	}

	r.callOnSkip(t, reason)
	r.callAfterTask(t, tr)

	return tr
}

// recordFailure marks a task failed, keeping whatever result it
// produced. A broadcast command that failed on one host of four still
// carries the other three, and guards such as OnlyIfChanged read them.
func (r *runner) recordFailure(
	t *Task,
	result *Result,
	err error,
	elapsed time.Duration,
) TaskResult {
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

// recordSuccess stores a completed result and reports it.
func (r *runner) recordSuccess(
	t *Task,
	result *Result,
	elapsed time.Duration,
) TaskResult {
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

func levelize(
	tasks []*Task,
) [][]*Task {
	level := make(map[string]int, len(tasks))

	maxLevel := 0

	for _, t := range tasks {
		if l := computeLevel(t, level); l > maxLevel {
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

// computeLevel returns how deep a task sits — one past its deepest
// dependency — memoising each answer in level as it goes.
func computeLevel(
	t *Task,
	level map[string]int,
) int {
	if l, ok := level[t.name]; ok {
		return l
	}

	maxDep := -1

	for _, dep := range t.deps {
		if depLevel := computeLevel(dep, level); depLevel > maxDep {
			maxDep = depLevel
		}
	}

	level[t.name] = maxDep + 1

	return maxDep + 1
}
