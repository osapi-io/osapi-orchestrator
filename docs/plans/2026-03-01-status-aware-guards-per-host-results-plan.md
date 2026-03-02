# Status-Aware Guards and Per-Host Results Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to
> implement this plan task-by-task.

**Goal:** Add status inspection in guards, per-host broadcast results, error
field exposure, non-zero exit handling, post-execution result access, and
renderer verbosity to the orchestrator DSL.

**Architecture:** Changes span two repos. The SDK (`osapi-sdk`) gains a Status
field on Result, always-store semantics for skipped/failed tasks, per-host
result extraction for broadcasts, non-zero exit handling, Data on TaskResult,
and a `TaskFnWithResults` type for inter-task data access. The orchestrator
(`osapi-orchestrator`) wraps these with `TaskStatus`, `Results.Status()`,
`Results.HostResults()`, convenience methods, `TaskFunc` with Results access,
`Report.Decode()`, error fields on result types, `WithVerbose()`, and renderer
updates.

**Tech Stack:** Go 1.25, osapi-sdk, osapi-orchestrator, testify/suite,
charmbracelet/lipgloss

**Repos:**

- `osapi-io/osapi-sdk` — Tasks 1-5
- `osapi-io/osapi-orchestrator` — Tasks 6-14

---

## Task 1: Add Status Field to SDK Result

**Repo:** `osapi-io/osapi-sdk`

**Files:**

- Modify: `pkg/orchestrator/result.go`
- Modify: `pkg/orchestrator/result_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/result_test.go`:

```go
func (s *ResultTestSuite) TestResultStatusField() {
	r := &Result{
		Changed: true,
		Data:    map[string]any{"hostname": "web-01"},
		Status:  StatusChanged,
	}

	s.Equal(StatusChanged, r.Status)
	s.True(r.Changed)
}
```

**Step 2: Run test to verify it fails**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -run TestResultTestSuite/TestResultStatusField -v ./pkg/orchestrator/...
```

Expected: FAIL — `Result` has no field `Status`.

**Step 3: Add Status field to Result**

In `pkg/orchestrator/result.go`, add `Status` to the `Result` struct:

```go
// Result is the outcome of a single task execution.
type Result struct {
	Changed bool
	Data    map[string]any
	Status  Status
}
```

**Step 4: Run test to verify it passes**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -run TestResultTestSuite/TestResultStatusField -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
cd ~/git/osapi-io/osapi-sdk
git add pkg/orchestrator/result.go pkg/orchestrator/result_test.go
git commit -m "feat(orchestrator): add Status field to Result"
```

---

## Task 2: Always Store Results in Runner

Currently the runner only stores a `Result` in `r.results` on success (line 353
of `runner.go`). Skipped and failed tasks have no entry. This means guards
cannot inspect their status. Fix by always storing a Result.

**Repo:** `osapi-io/osapi-sdk`

**Files:**

- Modify: `pkg/orchestrator/runner.go`
- Modify: `pkg/orchestrator/runner_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/runner_test.go`:

```go
func (s *RunnerTestSuite) TestSkippedTaskStoredInResults() {
	a := NewTaskFunc("a", func(
		_ context.Context,
		_ *osapi.Client,
	) (*Result, error) {
		return &Result{Changed: false}, nil
	})

	b := NewTaskFunc("b", func(
		_ context.Context,
		_ *osapi.Client,
	) (*Result, error) {
		return &Result{Changed: true}, nil
	})
	b.DependsOn(a)
	b.OnlyIfChanged()

	plan := NewPlan(nil, OnError(StopAll))
	plan.AddTask(a)
	plan.AddTask(b)

	report, err := plan.Run(context.Background())

	s.Require().NoError(err)
	s.Len(report.Tasks, 2)
	s.Equal(StatusSkipped, report.Tasks[1].Status)

	// The key assertion: skipped task's Result must be in the
	// runner's results map so guards can inspect it.
	// We verify via a third task that checks b's status.
	c := NewTaskFunc("c", func(
		_ context.Context,
		_ *osapi.Client,
	) (*Result, error) {
		return &Result{Changed: false}, nil
	})
	c.DependsOn(b)
	c.When(func(results Results) bool {
		r := results.Get("b")
		return r != nil && r.Status == StatusSkipped
	})

	plan2 := NewPlan(nil, OnError(StopAll))
	plan2.AddTask(a)
	plan2.AddTask(b)
	plan2.AddTask(c)

	report2, err := plan2.Run(context.Background())

	s.Require().NoError(err)
	// c should have run because b was skipped and the guard checked for it
	s.Equal(StatusUnchanged, report2.Tasks[2].Status)
}

func (s *RunnerTestSuite) TestFailedTaskStoredInResults() {
	a := NewTaskFunc("a", func(
		_ context.Context,
		_ *osapi.Client,
	) (*Result, error) {
		return nil, fmt.Errorf("boom")
	})
	a.OnError(Continue)

	b := NewTaskFunc("b", func(
		_ context.Context,
		_ *osapi.Client,
	) (*Result, error) {
		return &Result{Changed: false}, nil
	})
	b.When(func(results Results) bool {
		r := results.Get("a")
		return r != nil && r.Status == StatusFailed
	})

	plan := NewPlan(nil, OnError(Continue))
	plan.AddTask(a)
	plan.AddTask(b)

	report, err := plan.Run(context.Background())

	s.Require().NoError(err)
	s.Equal(StatusFailed, report.Tasks[0].Status)
	// b should have run because the guard saw a's failed status
	s.Equal(StatusUnchanged, report.Tasks[1].Status)
}
```

Note: This test will require adding `"context"`, `"fmt"`, and
`"github.com/osapi-io/osapi-sdk/pkg/osapi"` to imports in `runner_test.go`.

**Step 2: Run test to verify it fails**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -run TestRunnerTestSuite/TestSkippedTaskStoredInResults -v ./pkg/orchestrator/...
```

Expected: FAIL — skipped task has no entry in results.

**Step 3: Update runner to always store results**

In `pkg/orchestrator/runner.go`, modify the three skip/fail paths in `runTask()`
to store a Result before returning:

After the dependency-failed skip (around line 239), before the `r.mu.Unlock()`:

```go
r.results[t.name] = &Result{Status: StatusSkipped}
r.failed[t.name] = true
r.mu.Unlock()
```

After the OnlyIfChanged skip (around line 270), add:

```go
if !anyChanged {
	r.mu.Lock()
	r.results[t.name] = &Result{Status: StatusSkipped}
	r.mu.Unlock()

	tr := TaskResult{...}
	// ...
}
```

After the guard-returned-false skip (around line 289), add:

```go
if !shouldRun {
	r.mu.Lock()
	r.results[t.name] = &Result{Status: StatusSkipped}
	r.mu.Unlock()

	tr := TaskResult{...}
	// ...
}
```

After the failure path (around line 336), store the result:

```go
if err != nil {
	r.mu.Lock()
	r.failed[t.name] = true
	r.results[t.name] = &Result{Status: StatusFailed}
	r.mu.Unlock()

	tr := TaskResult{...}
	// ...
}
```

For the success path (around line 352-359), set the status:

```go
status := StatusUnchanged
if result.Changed {
	status = StatusChanged
}
result.Status = status

r.mu.Lock()
r.results[t.name] = result
r.mu.Unlock()
```

**Step 4: Run all runner tests**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -run TestRunnerTestSuite -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
cd ~/git/osapi-io/osapi-sdk
git add pkg/orchestrator/runner.go pkg/orchestrator/runner_test.go
git commit -m "feat(orchestrator): always store Result for skipped and failed tasks"
```

---

## Task 3: Add HostResult Type to SDK

**Repo:** `osapi-io/osapi-sdk`

**Files:**

- Modify: `pkg/orchestrator/result.go`
- Modify: `pkg/orchestrator/result_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/result_test.go`:

```go
func (s *ResultTestSuite) TestResultHostResults() {
	r := &Result{
		Changed: true,
		Data:    map[string]any{},
		Status:  StatusChanged,
		HostResults: []HostResult{
			{
				Hostname: "web-01",
				Changed:  true,
				Data:     map[string]any{"stdout": "ok"},
			},
			{
				Hostname: "web-02",
				Changed:  false,
				Error:    "connection timeout",
			},
		},
	}

	s.Len(r.HostResults, 2)
	s.Equal("web-01", r.HostResults[0].Hostname)
	s.True(r.HostResults[0].Changed)
	s.Equal("web-02", r.HostResults[1].Hostname)
	s.Equal("connection timeout", r.HostResults[1].Error)
}
```

**Step 2: Run test to verify it fails**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -run TestResultTestSuite/TestResultHostResults -v ./pkg/orchestrator/...
```

Expected: FAIL — `HostResult` type not defined.

**Step 3: Add HostResult type and field**

In `pkg/orchestrator/result.go`:

```go
// HostResult represents a single host's response within a broadcast
// operation.
type HostResult struct {
	Hostname string
	Changed  bool
	Error    string
	Data     map[string]any
}

// Result is the outcome of a single task execution.
type Result struct {
	Changed     bool
	Data        map[string]any
	Status      Status
	HostResults []HostResult
}
```

**Step 4: Run test to verify it passes**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -run TestResultTestSuite/TestResultHostResults -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
cd ~/git/osapi-io/osapi-sdk
git add pkg/orchestrator/result.go pkg/orchestrator/result_test.go
git commit -m "feat(orchestrator): add HostResult type for broadcast results"
```

---

## Task 4: Extract Per-Host Results and Handle Non-Zero Exit in Runner

**Repo:** `osapi-io/osapi-sdk`

**Files:**

- Modify: `pkg/orchestrator/runner.go`
- Modify: `pkg/orchestrator/task.go` (add `IsBroadcastTarget` helper)
- Create: `pkg/orchestrator/runner_broadcast_test.go`

**Step 1: Write the failing tests**

Create `pkg/orchestrator/runner_broadcast_test.go` with tests that verify
per-host extraction and non-zero exit handling. These will need a test HTTP
server that returns broadcast-style collection responses and command responses
with non-zero exit codes.

The exact test shape depends on how the SDK creates jobs and polls — since these
tests hit the full executeOp/pollJob path, they require an httptest server that
mimics the OSAPI API. Use the patterns from existing `plan_test.go` or
`runner_test.go` if present.

**Step 2: Add IsBroadcastTarget helper**

In `pkg/orchestrator/task.go`:

```go
// IsBroadcastTarget returns true if the target addresses multiple
// agents (broadcast or label selector).
func IsBroadcastTarget(
	target string,
) bool {
	return target == "_all" || strings.Contains(target, ":")
}
```

**Step 3: Update executeOp for per-host extraction**

In `pkg/orchestrator/runner.go`, modify `executeOp` to extract per-host results
after `pollJob` returns:

```go
func (r *runner) executeOp(
	ctx context.Context,
	op *Op,
) (*Result, error) {
	// ... existing job creation code ...

	result, err := r.pollJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Extract per-host results for broadcast targets.
	if IsBroadcastTarget(op.Target) {
		result.HostResults = extractHostResults(result.Data)
	}

	// Non-zero exit for command operations = failure.
	if isCommandOp(op.Operation) {
		if exitCode, ok := result.Data["exit_code"].(float64); ok && exitCode != 0 {
			result.Status = StatusFailed

			return result, fmt.Errorf(
				"command exited with code %d",
				int(exitCode),
			)
		}
	}

	return result, nil
}
```

Add helpers:

```go
// isCommandOp returns true for command execution operations.
func isCommandOp(
	operation string,
) bool {
	return operation == "command.exec.execute" ||
		operation == "command.shell.execute"
}

// extractHostResults parses per-agent results from a broadcast
// collection response.
func extractHostResults(
	data map[string]any,
) []HostResult {
	resultsRaw, ok := data["results"]
	if !ok {
		return nil
	}

	items, ok := resultsRaw.([]any)
	if !ok {
		return nil
	}

	hostResults := make([]HostResult, 0, len(items))

	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		hr := HostResult{
			Data: m,
		}

		if h, ok := m["hostname"].(string); ok {
			hr.Hostname = h
		}

		if c, ok := m["changed"].(bool); ok {
			hr.Changed = c
		}

		if e, ok := m["error"].(string); ok {
			hr.Error = e
		}

		hostResults = append(hostResults, hr)
	}

	return hostResults
}
```

**Step 4: Run tests**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
cd ~/git/osapi-io/osapi-sdk
git add pkg/orchestrator/runner.go pkg/orchestrator/task.go pkg/orchestrator/runner_broadcast_test.go
git commit -m "feat(orchestrator): per-host result extraction and non-zero exit handling"
```

---

## Task 5: Add TaskFnWithResults and Data on TaskResult to SDK

**Repo:** `osapi-io/osapi-sdk`

**Files:**

- Modify: `pkg/orchestrator/task.go`
- Modify: `pkg/orchestrator/result.go`
- Modify: `pkg/orchestrator/runner.go`
- Modify: `pkg/orchestrator/runner_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/runner_test.go`:

```go
func (s *RunnerTestSuite) TestTaskFuncWithResultsReceivesResults() {
	a := NewTaskFunc("a", func(
		_ context.Context,
		_ *osapi.Client,
	) (*Result, error) {
		return &Result{
			Changed: true,
			Data:    map[string]any{"hostname": "web-01"},
		}, nil
	})

	var captured string
	b := NewTaskFuncWithResults("b", func(
		_ context.Context,
		_ *osapi.Client,
		results Results,
	) (*Result, error) {
		r := results.Get("a")
		if r != nil {
			if h, ok := r.Data["hostname"].(string); ok {
				captured = h
			}
		}

		return &Result{Changed: false}, nil
	})
	b.DependsOn(a)

	plan := NewPlan(nil, OnError(StopAll))
	plan.AddTask(a)
	plan.AddTask(b)

	_, err := plan.Run(context.Background())

	s.Require().NoError(err)
	s.Equal("web-01", captured)
}

func (s *RunnerTestSuite) TestTaskResultCarriesData() {
	a := NewTaskFunc("a", func(
		_ context.Context,
		_ *osapi.Client,
	) (*Result, error) {
		return &Result{
			Changed: true,
			Data:    map[string]any{"stdout": "hello"},
		}, nil
	})

	plan := NewPlan(nil, OnError(StopAll))
	plan.AddTask(a)

	report, err := plan.Run(context.Background())

	s.Require().NoError(err)
	s.Len(report.Tasks, 1)
	s.Equal("hello", report.Tasks[0].Data["stdout"])
}
```

**Step 2: Run test to verify it fails**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -run "TestRunnerTestSuite/TestTaskFuncWithResults" -v ./pkg/orchestrator/...
```

Expected: FAIL — `NewTaskFuncWithResults` and `TaskResult.Data` not defined.

**Step 3: Add TaskFnWithResults type and Data field**

In `pkg/orchestrator/task.go`, add the new type and constructor:

```go
// TaskFnWithResults is like TaskFn but receives completed task results
// for inter-task data access.
type TaskFnWithResults func(
	ctx context.Context,
	client *osapi.Client,
	results Results,
) (*Result, error)
```

Add `fnr` field to `Task`:

```go
type Task struct {
	name           string
	op             *Op
	fn             TaskFn
	fnr            TaskFnWithResults
	deps           []*Task
	guard          GuardFn
	requiresChange bool
	errorStrategy  *ErrorStrategy
}
```

Add constructor:

```go
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
```

In `pkg/orchestrator/result.go`, add `Data` to `TaskResult`:

```go
type TaskResult struct {
	Name     string
	Status   Status
	Changed  bool
	Duration time.Duration
	Error    error
	Data     map[string]any
}
```

In `pkg/orchestrator/runner.go`, update `runTask` to handle `fnr` and populate
`Data` on TaskResult:

Around the execution block (line 318):

```go
if t.fnr != nil {
	result, err = t.fnr(ctx, client, r.results)
} else if t.fn != nil {
	result, err = t.fn(ctx, client)
} else {
	result, err = r.executeOp(ctx, t.op)
}
```

Around the success TaskResult creation (line 361):

```go
tr := TaskResult{
	Name:     t.name,
	Status:   status,
	Changed:  result.Changed,
	Duration: elapsed,
	Data:     result.Data,
}
```

**Step 4: Run tests**

```bash
cd ~/git/osapi-io/osapi-sdk && go test -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
cd ~/git/osapi-io/osapi-sdk
git add pkg/orchestrator/task.go pkg/orchestrator/result.go pkg/orchestrator/runner.go pkg/orchestrator/runner_test.go
git commit -m "feat(orchestrator): add TaskFnWithResults and Data on TaskResult"
```

---

## Task 6: Update osapi-sdk Dependency in Orchestrator

**Repo:** `osapi-io/osapi-orchestrator`

After SDK tasks 1-5 are committed and pushed, update the dependency:

```bash
cd ~/git/osapi-io/osapi-orchestrator
GOPROXY=direct go get github.com/osapi-io/osapi-sdk@latest
go mod tidy
```

Verify existing tests still pass:

```bash
go test -v ./pkg/orchestrator/...
```

Commit:

```bash
git add go.mod go.sum
git commit -m "chore(deps): update osapi-sdk with status-aware results"
```

---

## Task 7: Add TaskStatus Type and Results.Status()

**Repo:** `osapi-io/osapi-orchestrator`

**Files:**

- Modify: `pkg/orchestrator/result.go`
- Modify: `pkg/orchestrator/result_public_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/result_public_test.go`:

```go
func (s *ResultPublicTestSuite) TestStatus() {
	tests := []struct {
		name       string
		results    sdk.Results
		lookupName string
		wantStatus orchestrator.TaskStatus
	}{
		{
			name: "Returns TaskStatusChanged for changed result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Changed: true,
					Status:  sdk.StatusChanged,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusChanged,
		},
		{
			name: "Returns TaskStatusUnchanged for unchanged result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Changed: false,
					Status:  sdk.StatusUnchanged,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusUnchanged,
		},
		{
			name: "Returns TaskStatusSkipped for skipped result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Status: sdk.StatusSkipped,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusSkipped,
		},
		{
			name: "Returns TaskStatusFailed for failed result",
			results: sdk.Results{
				"step-a": &sdk.Result{
					Status: sdk.StatusFailed,
				},
			},
			lookupName: "step-a",
			wantStatus: orchestrator.TaskStatusFailed,
		},
		{
			name:       "Returns TaskStatusUnknown for missing result",
			results:    sdk.Results{},
			lookupName: "nonexistent",
			wantStatus: orchestrator.TaskStatusUnknown,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			r := orchestrator.NewResults(tc.results)
			s.Equal(tc.wantStatus, r.Status(tc.lookupName))
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test -run TestResultPublicTestSuite/TestStatus -v ./pkg/orchestrator/...
```

Expected: FAIL — `TaskStatus` type and `Status` method not defined.

**Step 3: Implement TaskStatus and Results.Status()**

In `pkg/orchestrator/result.go`, add after the `Results` struct:

```go
// TaskStatus represents the outcome of a step for guard inspection.
type TaskStatus int

const (
	// TaskStatusUnknown indicates the step was not found or has not run.
	TaskStatusUnknown TaskStatus = iota
	// TaskStatusChanged indicates the step ran and reported changes.
	TaskStatusChanged
	// TaskStatusUnchanged indicates the step ran with no changes.
	TaskStatusUnchanged
	// TaskStatusSkipped indicates the step was skipped.
	TaskStatusSkipped
	// TaskStatusFailed indicates the step failed.
	TaskStatusFailed
)

// Status returns the terminal status of a completed dependency step.
func (r Results) Status(
	name string,
) TaskStatus {
	result := r.results.Get(name)
	if result == nil {
		return TaskStatusUnknown
	}

	switch result.Status {
	case sdk.StatusChanged:
		return TaskStatusChanged
	case sdk.StatusUnchanged:
		return TaskStatusUnchanged
	case sdk.StatusSkipped:
		return TaskStatusSkipped
	case sdk.StatusFailed:
		return TaskStatusFailed
	default:
		return TaskStatusUnknown
	}
}
```

**Step 4: Run test to verify it passes**

```bash
go test -run TestResultPublicTestSuite/TestStatus -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
git add pkg/orchestrator/result.go pkg/orchestrator/result_public_test.go
git commit -m "feat(orchestrator): add TaskStatus type and Results.Status()"
```

---

## Task 8: Add OnlyIfFailed and OnlyIfAllChanged

**Repo:** `osapi-io/osapi-orchestrator`

**Files:**

- Modify: `pkg/orchestrator/step.go`
- Modify: `pkg/orchestrator/step_public_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/step_public_test.go` inside `TestChaining`:

```go
{
	name: "OnlyIfFailed returns non-nil step",
	chainFn: func() *orchestrator.Step {
		return s.orch.NodeHostnameGet("_any").OnlyIfFailed()
	},
},
{
	name: "OnlyIfAllChanged returns non-nil step",
	chainFn: func() *orchestrator.Step {
		return s.orch.NodeHostnameGet("_any").OnlyIfAllChanged()
	},
},
```

**Step 2: Run test to verify it fails**

```bash
go test -run TestStepPublicTestSuite/TestChaining -v ./pkg/orchestrator/...
```

Expected: FAIL — `OnlyIfFailed` and `OnlyIfAllChanged` not defined.

**Step 3: Implement convenience methods**

In `pkg/orchestrator/step.go`, add:

```go
// OnlyIfFailed skips this step unless at least one dependency failed.
func (s *Step) OnlyIfFailed() *Step {
	s.task.When(func(sdkResults sdk.Results) bool {
		for _, dep := range s.task.Dependencies() {
			if r := sdkResults.Get(dep.Name()); r != nil && r.Status == sdk.StatusFailed {
				return true
			}
		}

		return false
	})

	return s
}

// OnlyIfAllChanged skips this step unless all dependencies reported
// changes.
func (s *Step) OnlyIfAllChanged() *Step {
	s.task.When(func(sdkResults sdk.Results) bool {
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
	})

	return s
}
```

**Step 4: Run test to verify it passes**

```bash
go test -run TestStepPublicTestSuite/TestChaining -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
git add pkg/orchestrator/step.go pkg/orchestrator/step_public_test.go
git commit -m "feat(orchestrator): add OnlyIfFailed and OnlyIfAllChanged"
```

---

## Task 9: Add HostResult Type and Results.HostResults()

**Repo:** `osapi-io/osapi-orchestrator`

**Files:**

- Modify: `pkg/orchestrator/result.go`
- Modify: `pkg/orchestrator/result_public_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/result_public_test.go`:

```go
func (s *ResultPublicTestSuite) TestHostResults() {
	tests := []struct {
		name       string
		results    sdk.Results
		lookupName string
		wantNil    bool
		wantLen    int
		validateFn func(hrs []orchestrator.HostResult)
	}{
		{
			name: "Returns per-host results",
			results: sdk.Results{
				"deploy": &sdk.Result{
					Changed: true,
					Status:  sdk.StatusChanged,
					HostResults: []sdk.HostResult{
						{
							Hostname: "web-01",
							Changed:  true,
							Data: map[string]any{
								"stdout": "deployed",
							},
						},
						{
							Hostname: "web-02",
							Changed:  false,
							Error:    "timeout",
						},
					},
				},
			},
			lookupName: "deploy",
			wantLen:    2,
			validateFn: func(hrs []orchestrator.HostResult) {
				s.Equal("web-01", hrs[0].Hostname)
				s.True(hrs[0].Changed)
				s.Empty(hrs[0].Error)
				s.Equal("web-02", hrs[1].Hostname)
				s.Equal("timeout", hrs[1].Error)
			},
		},
		{
			name:       "Returns nil for missing step",
			results:    sdk.Results{},
			lookupName: "nonexistent",
			wantNil:    true,
		},
		{
			name: "Returns nil for unicast result",
			results: sdk.Results{
				"get-host": &sdk.Result{
					Changed: false,
					Status:  sdk.StatusUnchanged,
				},
			},
			lookupName: "get-host",
			wantNil:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			r := orchestrator.NewResults(tc.results)
			hrs := r.HostResults(tc.lookupName)

			if tc.wantNil {
				s.Nil(hrs)
				return
			}

			s.Len(hrs, tc.wantLen)

			if tc.validateFn != nil {
				tc.validateFn(hrs)
			}
		})
	}
}

func (s *ResultPublicTestSuite) TestHostResultDecode() {
	r := orchestrator.NewResults(sdk.Results{
		"run-cmd": &sdk.Result{
			Changed: true,
			Status:  sdk.StatusChanged,
			HostResults: []sdk.HostResult{
				{
					Hostname: "web-01",
					Changed:  true,
					Data: map[string]any{
						"stdout":    "hello",
						"stderr":    "",
						"exit_code": float64(0),
					},
				},
			},
		},
	})

	hrs := r.HostResults("run-cmd")
	s.Require().Len(hrs, 1)

	var cmd orchestrator.CommandResult
	err := hrs[0].Decode(&cmd)

	s.Require().NoError(err)
	s.Equal("hello", cmd.Stdout)
	s.Equal(0, cmd.ExitCode)
}
```

**Step 2: Run test to verify it fails**

```bash
go test -run "TestResultPublicTestSuite/TestHostResults" -v ./pkg/orchestrator/...
```

Expected: FAIL — `HostResult` type and `HostResults` method not defined.

**Step 3: Implement HostResult and Results.HostResults()**

In `pkg/orchestrator/result.go`:

```go
// HostResult represents a single host's response within a broadcast
// operation.
type HostResult struct {
	Hostname string
	Changed  bool
	Error    string
	Data     map[string]any
}

// Decode unmarshals host-specific data into a typed result struct.
func (h HostResult) Decode(
	v any,
) error {
	b, err := json.Marshal(h.Data)
	if err != nil {
		return fmt.Errorf("marshal host result data: %w", err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("decode host result data: %w", err)
	}

	return nil
}

// HostResults returns per-host results for a broadcast operation.
// Returns nil for unicast operations or unknown step names.
func (r Results) HostResults(
	name string,
) []HostResult {
	result := r.results.Get(name)
	if result == nil || len(result.HostResults) == 0 {
		return nil
	}

	hrs := make([]HostResult, len(result.HostResults))
	for i, hr := range result.HostResults {
		hrs[i] = HostResult{
			Hostname: hr.Hostname,
			Changed:  hr.Changed,
			Error:    hr.Error,
			Data:     hr.Data,
		}
	}

	return hrs
}
```

**Step 4: Run test to verify it passes**

```bash
go test -run "TestResultPublicTestSuite/TestHostResult" -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
git add pkg/orchestrator/result.go pkg/orchestrator/result_public_test.go
git commit -m "feat(orchestrator): add HostResult type and Results.HostResults()"
```

---

## Task 10: Add Error Field to Result Types

**Repo:** `osapi-io/osapi-orchestrator`

**Files:**

- Modify: `pkg/orchestrator/result_types.go`
- Modify: `pkg/orchestrator/result_public_test.go`

**Step 1: Write the failing test**

Add to `pkg/orchestrator/result_public_test.go` inside `TestDecode`:

```go
{
	name: "Decodes command result with error field",
	results: sdk.Results{
		"run-cmd": &sdk.Result{
			Changed: true,
			Data: map[string]any{
				"stdout":      "partial output",
				"stderr":      "command not found",
				"exit_code":   float64(127),
				"duration_ms": float64(50),
				"error":       "exec failed",
			},
		},
	},
	lookupName: "run-cmd",
	target:     &orchestrator.CommandResult{},
	validateFunc: func() {
		t := s.T()
		r := orchestrator.NewResults(sdk.Results{
			"run-cmd": &sdk.Result{
				Changed: true,
				Data: map[string]any{
					"stdout":      "partial output",
					"stderr":      "command not found",
					"exit_code":   float64(127),
					"duration_ms": float64(50),
					"error":       "exec failed",
				},
			},
		})
		var cmd orchestrator.CommandResult
		err := r.Decode("run-cmd", &cmd)
		if err != nil {
			t.Fatal(err)
		}
		s.Equal("exec failed", cmd.Error)
		s.Equal(127, cmd.ExitCode)
		s.Equal("command not found", cmd.Stderr)
	},
},
```

**Step 2: Run test to verify it fails**

```bash
go test -run "TestResultPublicTestSuite/TestDecode/Decodes_command_result_with_error_field" -v ./pkg/orchestrator/...
```

Expected: FAIL — `CommandResult` has no field `Error`.

**Step 3: Add Error fields to result types**

In `pkg/orchestrator/result_types.go`:

```go
type CommandResult struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
}

type PingResult struct {
	PacketsSent     int     `json:"packets_sent"`
	PacketsReceived int     `json:"packets_received"`
	PacketLoss      float64 `json:"packet_loss"`
	Error           string  `json:"error,omitempty"`
}

type DNSUpdateResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
```

**Step 4: Run test to verify it passes**

```bash
go test -run TestResultPublicTestSuite/TestDecode -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Commit**

```bash
git add pkg/orchestrator/result_types.go pkg/orchestrator/result_public_test.go
git commit -m "feat(orchestrator): add Error field to result types"
```

---

## Task 11: Add TaskFunc and Report.Decode to Orchestrator

**Repo:** `osapi-io/osapi-orchestrator`

**Files:**

- Modify: `pkg/orchestrator/orchestrator.go`
- Modify: `pkg/orchestrator/result.go`
- Modify: `pkg/orchestrator/result_public_test.go`

**Step 1: Write the failing tests**

Add to `pkg/orchestrator/result_public_test.go`:

```go
func (s *ResultPublicTestSuite) TestReportDecode() {
	tests := []struct {
		name        string
		tasks       []sdk.TaskResult
		lookupName  string
		expectErr   bool
		errContains string
		validateFn  func(cmd orchestrator.CommandResult)
	}{
		{
			name: "Decodes task result from report",
			tasks: []sdk.TaskResult{
				{
					Name:    "run-cmd",
					Status:  sdk.StatusChanged,
					Changed: true,
					Data: map[string]any{
						"stdout":    "hello",
						"exit_code": float64(0),
					},
				},
			},
			lookupName: "run-cmd",
			validateFn: func(cmd orchestrator.CommandResult) {
				s.Equal("hello", cmd.Stdout)
				s.Equal(0, cmd.ExitCode)
			},
		},
		{
			name:        "Returns error for missing task",
			tasks:       []sdk.TaskResult{},
			lookupName:  "nonexistent",
			expectErr:   true,
			errContains: "no result for",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			report := &orchestrator.Report{
				Tasks: tc.tasks,
			}

			var cmd orchestrator.CommandResult
			err := report.Decode(tc.lookupName, &cmd)

			if tc.expectErr {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.errContains)

				return
			}

			s.Require().NoError(err)

			if tc.validateFn != nil {
				tc.validateFn(cmd)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test -run TestResultPublicTestSuite/TestReportDecode -v ./pkg/orchestrator/...
```

Expected: FAIL — `Report.Decode` not defined.

**Step 3: Implement Report.Decode()**

In `pkg/orchestrator/result.go`:

```go
// Decode retrieves the result of a named task from the report
// and decodes it into the given typed struct.
func (r *Report) Decode(
	name string,
	v any,
) error {
	for _, t := range r.Tasks {
		if t.Name == name {
			if t.Data == nil {
				return fmt.Errorf("no result data for %q", name)
			}

			b, err := json.Marshal(t.Data)
			if err != nil {
				return fmt.Errorf("marshal result data: %w", err)
			}

			if err := json.Unmarshal(b, v); err != nil {
				return fmt.Errorf("decode result data: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("no result for %q", name)
}
```

**Step 4: Add TaskFunc to Orchestrator**

In `pkg/orchestrator/orchestrator.go`:

```go
// TaskFunc creates a custom step that receives completed results
// from prior steps.
func (o *Orchestrator) TaskFunc(
	name string,
	fn func(ctx context.Context, r Results) (*sdk.Result, error),
) *Step {
	task := o.plan.TaskFuncWithResults(
		name,
		func(
			ctx context.Context,
			_ *osapi.Client,
			results sdk.Results,
		) (*sdk.Result, error) {
			return fn(ctx, Results{results: results})
		},
	)

	return &Step{task: task}
}
```

**Step 5: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 6: Commit**

```bash
git add pkg/orchestrator/orchestrator.go pkg/orchestrator/result.go pkg/orchestrator/result_public_test.go
git commit -m "feat(orchestrator): add TaskFunc with Results and Report.Decode()"
```

---

## Task 12: Add WithVerbose Option

**Repo:** `osapi-io/osapi-orchestrator`

**Files:**

- Modify: `pkg/orchestrator/types.go`
- Modify: `pkg/orchestrator/orchestrator.go`
- Modify: `pkg/orchestrator/renderer_lipgloss.go`
- Create: `pkg/orchestrator/options_public_test.go`

**Step 1: Write the failing test**

Create `pkg/orchestrator/options_public_test.go`:

```go
package orchestrator_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type OptionsPublicTestSuite struct {
	suite.Suite
}

func TestOptionsPublicTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsPublicTestSuite))
}

func (s *OptionsPublicTestSuite) TestWithVerbose() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	o := orchestrator.New(
		server.URL,
		"test-token",
		orchestrator.WithVerbose(),
	)

	s.NotNil(o)
}
```

**Step 2: Run test to verify it fails**

```bash
go test -run TestOptionsPublicTestSuite/TestWithVerbose -v ./pkg/orchestrator/...
```

Expected: FAIL — `WithVerbose` not defined, `New` doesn't accept options.

**Step 3: Implement WithVerbose and Option type**

In `pkg/orchestrator/options.go`, add:

```go
// Option configures the Orchestrator.
type Option func(*config)

type config struct {
	verbose bool
}

// WithVerbose enables verbose output showing stdout, stderr, and
// full response data for all tasks.
func WithVerbose() Option {
	return func(c *config) {
		c.verbose = true
	}
}
```

In `pkg/orchestrator/orchestrator.go`, update `New` to accept options:

```go
func New(
	url string,
	token string,
	opts ...Option,
) *Orchestrator {
	cfg := &config{}
	for _, o := range opts {
		o(cfg)
	}

	client := osapi.New(url, token)
	r := newLipglossRenderer()
	r.verbose = cfg.verbose
	plan := sdk.NewPlan(client, sdk.WithHooks(rendererHooks(r)))

	return &Orchestrator{
		plan:     plan,
		renderer: r,
	}
}
```

In `pkg/orchestrator/renderer_lipgloss.go`, add the `verbose` field:

```go
type lipglossRenderer struct {
	w       io.Writer
	verbose bool
	// ... existing style fields ...
}
```

**Step 4: Run test to verify it passes**

```bash
go test -run TestOptionsPublicTestSuite -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Run all tests to verify nothing broke**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 6: Commit**

```bash
git add pkg/orchestrator/options.go pkg/orchestrator/orchestrator.go \
       pkg/orchestrator/renderer_lipgloss.go pkg/orchestrator/types.go \
       pkg/orchestrator/options_public_test.go
git commit -m "feat(orchestrator): add WithVerbose option"
```

---

## Task 13: Update Renderer for Failure Detail and Verbose Output

**Repo:** `osapi-io/osapi-orchestrator`

**Files:**

- Modify: `pkg/orchestrator/renderer.go` (update interface)
- Modify: `pkg/orchestrator/renderer_lipgloss.go`
- Modify: `pkg/orchestrator/renderer_lipgloss_test.go`
- Modify: `pkg/orchestrator/orchestrator.go` (pass result data to hooks)

**Step 1: Write the failing tests**

Add to `pkg/orchestrator/renderer_lipgloss_test.go`:

```go
func (s *RendererLipglossTestSuite) TestTaskDoneShowsErrorOnFailure() {
	r := newLipglossRendererWithWriter(&s.buf)

	r.TaskDone(sdk.TaskResult{
		Name:     "deploy",
		Status:   sdk.StatusFailed,
		Duration: 45 * time.Millisecond,
		Error:    fmt.Errorf("command exited with code 1"),
	})

	output := s.buf.String()
	s.Contains(output, "[failed]")
	s.Contains(output, "deploy")
	s.Contains(output, "command exited with code 1")
}

func (s *RendererLipglossTestSuite) TestTaskDoneVerboseShowsStdout() {
	r := newLipglossRendererWithWriter(&s.buf)
	r.verbose = true

	r.TaskDone(sdk.TaskResult{
		Name:     "run-cmd",
		Status:   sdk.StatusChanged,
		Changed:  true,
		Duration: 100 * time.Millisecond,
		Data: map[string]any{
			"stdout":    "hello world",
			"exit_code": float64(0),
		},
	})

	output := s.buf.String()
	s.Contains(output, "[changed]")
	s.Contains(output, "stdout: hello world")
}

func (s *RendererLipglossTestSuite) TestTaskDoneNormalModeHidesStdout() {
	r := newLipglossRendererWithWriter(&s.buf)
	r.verbose = false

	r.TaskDone(sdk.TaskResult{
		Name:     "run-cmd",
		Status:   sdk.StatusChanged,
		Changed:  true,
		Duration: 100 * time.Millisecond,
		Data: map[string]any{
			"stdout":    "hello world",
			"exit_code": float64(0),
		},
	})

	output := s.buf.String()
	s.Contains(output, "[changed]")
	s.NotContains(output, "stdout:")
}
```

Note: `TaskResult` in the SDK does not currently have a `Data` field. The
renderer needs access to result data to display it. The `AfterTask` hook
receives `(task *sdk.Task, result sdk.TaskResult)`. Two approaches:

1. Add a `Data map[string]any` field to `sdk.TaskResult`
2. Pass result data through the renderer interface

Option 1 is simpler and extends the existing contract. Add `Data map[string]any`
to `sdk.TaskResult` in the SDK's `result.go`, and populate it in the runner
where `result` is available (line 361-366 in `runner.go`). This requires one
more small SDK change.

**Step 2: SDK change — add Data to TaskResult**

In SDK `pkg/orchestrator/result.go`:

```go
type TaskResult struct {
	Name     string
	Status   Status
	Changed  bool
	Duration time.Duration
	Error    error
	Data     map[string]any
}
```

In SDK `pkg/orchestrator/runner.go`, around line 361-366:

```go
tr := TaskResult{
	Name:     t.name,
	Status:   status,
	Changed:  result.Changed,
	Duration: elapsed,
	Data:     result.Data,
}
```

**Step 3: Implement renderer changes**

In `pkg/orchestrator/renderer_lipgloss.go`, update `TaskDone`:

```go
func (r *lipglossRenderer) TaskDone(
	result sdk.TaskResult,
) {
	label := fmt.Sprintf("[%s]", result.Status)

	var tag string
	if result.Status == sdk.StatusFailed {
		tag = padTag(r.red.Render(label), len(label))
	} else {
		tag = padTag(r.green.Render(label), len(label))
	}

	changedStr := fmt.Sprintf("changed=%v", result.Changed)
	if result.Changed {
		changedStr = r.greenB.Render(changedStr)
	}

	r.printf(
		"  %s %-*s %s  %s\n",
		tag,
		nameWidth,
		result.Name,
		changedStr,
		r.dim.Render(formatDuration(result.Duration)),
	)

	// Always show error detail on failure.
	if result.Status == sdk.StatusFailed && result.Error != nil {
		r.printf(
			"  %s %-*s %s\n",
			strings.Repeat(" ", tagWidth),
			nameWidth,
			"",
			r.red.Render(result.Error.Error()),
		)
	}

	// Verbose mode: show response data on success.
	if r.verbose && result.Data != nil {
		r.printResultData(result.Data)
	}
}

// printResultData renders key result fields as indented lines.
func (r *lipglossRenderer) printResultData(
	data map[string]any,
) {
	indent := strings.Repeat(" ", tagWidth+nameWidth+2)

	for _, key := range []string{"stdout", "stderr", "hostname", "error"} {
		if v, ok := data[key]; ok {
			str := fmt.Sprintf("%v", v)
			if str != "" {
				r.printf(
					"%s%s\n",
					indent,
					r.dim.Render(fmt.Sprintf("%s: %s", key, str)),
				)
			}
		}
	}
}
```

**Step 4: Run tests**

```bash
go test -run TestRendererLipglossTestSuite -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 5: Run all tests**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

**Step 6: Commit**

```bash
git add pkg/orchestrator/renderer_lipgloss.go pkg/orchestrator/renderer_lipgloss_test.go
git commit -m "feat(orchestrator): show failure detail and verbose output in renderer"
```

---

## Task 14: Final Verification

**Repo:** `osapi-io/osapi-orchestrator`

**Step 1: Run full test suite**

```bash
just test
```

Expected: PASS — all linting, formatting, and tests pass.

**Step 2: Run full test suite on SDK**

```bash
cd ~/git/osapi-io/osapi-sdk && just test
```

Expected: PASS

**Step 3: Verify example still compiles**

```bash
cd ~/git/osapi-io/osapi-orchestrator && go build ./examples/...
```

Expected: builds without errors.
