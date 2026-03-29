# Host Status Awareness + Missing Operations

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add per-host status awareness (ok/skipped/failed) to the orchestrator DSL, fix guards to distinguish skipped from failed, and add all missing SDK operations.

**Architecture:** The SDK's `HostResult` now has a `Status string` field. The orchestrator wraps SDK types with its own `HostResult` for added behavior (Decode). We add `Status` to the local type, define constants, update guards to use status instead of error-string heuristics, fix the renderer, and add typed constructors for all missing SDK operations.

**Tech Stack:** Go 1.25, osapi SDK, testify/suite, lipgloss

---

## File Structure

**Modified files:**
- `go.mod` — bump SDK dependency
- `pkg/orchestrator/result.go` — add Status to HostResult, add host status constants, update HostResults()
- `pkg/orchestrator/result_public_test.go` — test Status field, host status constants
- `pkg/orchestrator/step.go` — fix OnlyIfAnyHostFailed/OnlyIfAllHostsFailed, add OnlyIfAnyHostSkipped
- `pkg/orchestrator/step_public_test.go` — test new/fixed guards
- `pkg/orchestrator/ops.go` — add Status to all mappers, add new operation constructors
- `pkg/orchestrator/ops_public_test.go` — tests for new operations
- `pkg/orchestrator/renderer_lipgloss.go` — use Status for display
- `pkg/orchestrator/renderer_lipgloss_test.go` — test skipped rendering
- `docs/operations/README.md` — add new operations to index
- `docs/features/guards.md` — document OnlyIfAnyHostSkipped, clarify skipped semantics
- `docs/features/broadcast.md` — add Status field to HostResult docs
- `docs/features/README.md` — update step chaining table
- `README.md` — update operation counts and feature lists

**New files:**
- `pkg/orchestrator/host_status.go` — host status constants
- `pkg/orchestrator/host_status_public_test.go` — tests
- `docs/operations/node-hostname-update.md`
- `docs/operations/node-os-get.md`
- `docs/operations/file-undeploy.md`
- `docs/operations/cron-list.md`
- `docs/operations/cron-get.md`
- `docs/operations/cron-create.md`
- `docs/operations/cron-update.md`
- `docs/operations/cron-delete.md`
- `docs/operations/agent-drain.md`
- `docs/operations/agent-undrain.md`
- `examples/operations/hostname-update.go`
- `examples/operations/cron.go`
- `examples/features/host-status.go`

---

### Task 1: Update SDK dependency

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Create feature branch**

```bash
cd /Users/john/git/osapi-io/osapi-orchestrator
git checkout -b feat/host-status-missing-ops main
```

- [ ] **Step 2: Update SDK dependency to latest main**

```bash
cd /Users/john/git/osapi-io/osapi-orchestrator
go get github.com/retr0h/osapi@main
go mod tidy
```

- [ ] **Step 3: Verify it compiles**

```bash
go build ./...
```

Expected: clean build (no errors)

- [ ] **Step 4: Verify tests still pass**

```bash
just go::unit
```

Expected: all existing tests pass

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: bump osapi SDK to pick up HostResult.Status"
```

---

### Task 2: Add host status constants and Status field to HostResult

**Files:**
- Create: `pkg/orchestrator/host_status.go`
- Create: `pkg/orchestrator/host_status_public_test.go`
- Modify: `pkg/orchestrator/result.go`
- Modify: `pkg/orchestrator/result_public_test.go`

- [ ] **Step 1: Write the test for host status constants**

Create `pkg/orchestrator/host_status_public_test.go`:

```go
package orchestrator_test

import (
	"testing"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type HostStatusPublicTestSuite struct {
	suite.Suite
}

func TestHostStatusPublicTestSuite(t *testing.T) {
	suite.Run(t, new(HostStatusPublicTestSuite))
}

func (s *HostStatusPublicTestSuite) TestHostStatusConstants() {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{
			name:     "ok constant",
			status:   orchestrator.HostStatusOk,
			expected: "ok",
		},
		{
			name:     "skipped constant",
			status:   orchestrator.HostStatusSkipped,
			expected: "skipped",
		},
		{
			name:     "failed constant",
			status:   orchestrator.HostStatusFailed,
			expected: "failed",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Equal(tc.expected, tc.status)
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -run TestHostStatusPublicTestSuite -v ./pkg/orchestrator/...
```

Expected: FAIL — `HostStatusOk` not defined

- [ ] **Step 3: Implement host status constants**

Create `pkg/orchestrator/host_status.go`:

```go
package orchestrator

// Host-level status constants returned by the API for per-host results.
// These represent agent-side outcomes, distinct from task-level statuses
// (TaskStatusChanged, TaskStatusSkipped, etc.) which are DAG-level.
const (
	// HostStatusOk indicates the operation completed successfully on the host.
	HostStatusOk = "ok"
	// HostStatusSkipped indicates the operation is not supported on the host
	// (e.g., a Darwin host in a Linux fleet). This is NOT an error.
	HostStatusSkipped = "skipped"
	// HostStatusFailed indicates the operation failed on the host.
	HostStatusFailed = "failed"
)
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -run TestHostStatusPublicTestSuite -v ./pkg/orchestrator/...
```

Expected: PASS

- [ ] **Step 5: Add Status field to local HostResult and update HostResults()**

In `pkg/orchestrator/result.go`, add `Status string` to the `HostResult` struct (between Hostname and Changed):

```go
type HostResult struct {
	Hostname string
	Status   string
	Changed  bool
	Error    string
	Data     map[string]any
}
```

Update the `HostResults()` method to copy Status:

```go
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
			Status:   hr.Status,
			Changed:  hr.Changed,
			Error:    hr.Error,
			Data:     hr.Data,
		}
	}

	return hrs
}
```

- [ ] **Step 6: Add test for Status in HostResults**

In `pkg/orchestrator/result_public_test.go`, add a new test case to `TestHostResults` that verifies Status is copied through. Add a "with status fields" test case:

```go
{
	name: "with status fields",
	sdkResults: sdk.Results{
		"deploy": &sdk.Result{
			HostResults: []sdk.HostResult{
				{Hostname: "web-01", Status: "ok", Changed: true},
				{Hostname: "web-02", Status: "skipped", Error: "unsupported"},
				{Hostname: "web-03", Status: "failed", Error: "permission denied"},
			},
		},
	},
	stepName: "deploy",
	validateFunc: func(hrs []orchestrator.HostResult) {
		s.Require().Len(hrs, 3)
		s.Equal("ok", hrs[0].Status)
		s.Equal("skipped", hrs[1].Status)
		s.Equal("unsupported", hrs[1].Error)
		s.Equal("failed", hrs[2].Status)
		s.Equal("permission denied", hrs[2].Error)
	},
},
```

- [ ] **Step 7: Run all tests**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add pkg/orchestrator/host_status.go pkg/orchestrator/host_status_public_test.go \
       pkg/orchestrator/result.go pkg/orchestrator/result_public_test.go
git commit -m "feat: add host status constants and Status field to HostResult"
```

---

### Task 3: Fix guards to use Status instead of Error

**Files:**
- Modify: `pkg/orchestrator/step.go`
- Modify: `pkg/orchestrator/step_public_test.go` (or the file containing guard tests)

- [ ] **Step 1: Find the guard test file and read it**

```bash
grep -rn "OnlyIfAnyHostFailed\|OnlyIfAllHostsFailed" /Users/john/git/osapi-io/osapi-orchestrator/pkg/orchestrator/*_test.go
```

Read the test file to understand the existing test patterns.

- [ ] **Step 2: Write tests for the fixed guards**

Add test cases to the existing guard test method that verify:
- A host with `Status: "skipped", Error: "unsupported"` does NOT trigger OnlyIfAnyHostFailed
- A host with `Status: "failed", Error: "permission denied"` DOES trigger OnlyIfAnyHostFailed
- OnlyIfAllHostsFailed returns false when some hosts are skipped (not all failed)

The test should set up SDK results with HostResults that have Status fields:

```go
{
	name: "skipped hosts do not count as failed",
	sdkResults: sdk.Results{
		"deploy": &sdk.Result{
			HostResults: []sdk.HostResult{
				{Hostname: "web-01", Status: "ok"},
				{Hostname: "web-02", Status: "skipped", Error: "unsupported"},
			},
		},
	},
	// OnlyIfAnyHostFailed should return false — skipped is not failed
},
{
	name: "failed hosts trigger guard",
	sdkResults: sdk.Results{
		"deploy": &sdk.Result{
			HostResults: []sdk.HostResult{
				{Hostname: "web-01", Status: "ok"},
				{Hostname: "web-02", Status: "failed", Error: "permission denied"},
			},
		},
	},
	// OnlyIfAnyHostFailed should return true
},
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
go test -run "TestOnlyIfAnyHostFailed" -v ./pkg/orchestrator/...
```

Expected: FAIL — the skipped test case fires the guard because it checks `hr.Error != ""`

- [ ] **Step 4: Fix OnlyIfAnyHostFailed in step.go**

Change from checking `hr.Error != ""` to checking `hr.Status`:

```go
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
```

- [ ] **Step 5: Fix OnlyIfAllHostsFailed in step.go**

Same pattern — check `hr.Status == "failed"` instead of `hr.Error == ""`:

```go
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
```

- [ ] **Step 6: Run tests to verify they pass**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add pkg/orchestrator/step.go pkg/orchestrator/step_public_test.go
git commit -m "fix: guards check Status instead of Error to exclude skipped hosts"
```

---

### Task 4: Add OnlyIfAnyHostSkipped guard

**Files:**
- Modify: `pkg/orchestrator/step.go`
- Modify: `pkg/orchestrator/step_public_test.go`

- [ ] **Step 1: Write tests for OnlyIfAnyHostSkipped**

Add test cases:
- Returns false when no dependencies
- Returns false when no host results
- Returns false when all hosts are ok
- Returns true when any host has `Status: "skipped"`
- Returns false when hosts are failed but not skipped

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test -run "TestOnlyIfAnyHostSkipped" -v ./pkg/orchestrator/...
```

Expected: FAIL — method doesn't exist

- [ ] **Step 3: Implement OnlyIfAnyHostSkipped**

Add to `pkg/orchestrator/step.go`:

```go
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
```

- [ ] **Step 4: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/orchestrator/step.go pkg/orchestrator/step_public_test.go
git commit -m "feat: add OnlyIfAnyHostSkipped guard for broadcast operations"
```

---

### Task 5: Fix renderer to use Status

**Files:**
- Modify: `pkg/orchestrator/renderer_lipgloss.go`
- Modify: `pkg/orchestrator/renderer_lipgloss_test.go`

- [ ] **Step 1: Write test for skipped host rendering**

Add a test case that verifies a host with `Status: "skipped"` renders with yellow "skipped:" prefix instead of red "error:". Read the existing test file first to match the pattern.

- [ ] **Step 2: Run test to verify it fails**

Expected: FAIL — skipped hosts render as "error: unsupported"

- [ ] **Step 3: Fix printHostResults**

In `pkg/orchestrator/renderer_lipgloss.go`, update the `printHostResults` method to check `hr.Status`:

```go
for _, hr := range hostResults {
	status := r.green.Render("ok")

	switch hr.Status {
	case "skipped":
		msg := "skipped"
		if hr.Error != "" {
			msg = "skipped: " + hr.Error
		}
		status = r.yellow.Render(msg)
	case "failed":
		msg := "failed"
		if hr.Error != "" {
			msg = "failed: " + hr.Error
		}
		status = r.red.Render(msg)
	default:
		if hr.Error != "" {
			status = r.red.Render("error: " + hr.Error)
		}
	}
```

Note: check if `r.yellow` exists on the lipglossRenderer. If not, add a `yellow` style field and initialize it in the constructor.

- [ ] **Step 4: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/orchestrator/renderer_lipgloss.go pkg/orchestrator/renderer_lipgloss_test.go
git commit -m "fix: renderer shows skipped status in yellow, distinct from failed"
```

---

### Task 6: Add Status to all existing mappers in ops.go

**Files:**
- Modify: `pkg/orchestrator/ops.go`
- Modify: `pkg/orchestrator/ops_public_test.go`

- [ ] **Step 1: Add `Status: r.Status` to every `sdk.HostResult{}` in ops.go**

There are 18 mappers that create `sdk.HostResult{}` in ops.go. Each needs `Status: r.Status` added between Hostname and Changed. The operations are:

1. NodeHostnameGet
2. NodeStatusGet
3. NodeUptimeGet
4. NodeDiskGet
5. NodeMemoryGet
6. NodeLoadGet
7. NetworkDNSGet
8. NetworkDNSUpdate
9. NetworkPingDo
10. CommandExec
11. CommandShell
12. DockerPull
13. DockerCreate
14. DockerStart
15. DockerStop
16. DockerRemove
17. DockerExec
18. DockerInspect
19. DockerList
20. DockerImageRemove

For each, change from:
```go
return sdk.HostResult{
	Hostname: r.Hostname,
	Changed:  r.Changed,
	Error:    r.Error,
}
```

To:
```go
return sdk.HostResult{
	Hostname: r.Hostname,
	Status:   r.Status,
	Changed:  r.Changed,
	Error:    r.Error,
}
```

For CommandExec and CommandShell, the Error field uses `commandError(r)` — keep that, but add Status:
```go
return sdk.HostResult{
	Hostname: r.Hostname,
	Status:   r.Status,
	Changed:  r.Changed,
	Error:    commandError(r),
}
```

- [ ] **Step 2: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add pkg/orchestrator/ops.go
git commit -m "feat: pass Status through all existing operation mappers"
```

---

### Task 7: Add missing operations — NodeHostnameUpdate and NodeOSGet

**Files:**
- Modify: `pkg/orchestrator/ops.go`
- Modify: `pkg/orchestrator/ops_public_test.go`

- [ ] **Step 1: Write tests for NodeHostnameUpdate and NodeOSGet**

Follow the existing test pattern in `ops_public_test.go`. Each test should mock the SDK client call and verify the CollectionResult mapping.

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test -run "TestNodeHostnameUpdate|TestNodeOSGet" -v ./pkg/orchestrator/...
```

Expected: FAIL — methods don't exist

- [ ] **Step 3: Implement NodeHostnameUpdate**

Add to `pkg/orchestrator/ops.go`:

```go
// NodeHostnameUpdate creates a step that sets the system hostname.
func (o *Orchestrator) NodeHostnameUpdate(
	target string,
	hostname string,
) *Step {
	name := o.nextOpName("update-hostname")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.UpdateHostname(ctx, target, hostname)
			if err != nil {
				return nil, fmt.Errorf("update hostname: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.HostnameUpdateResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}
```

- [ ] **Step 4: Implement NodeOSGet**

Add to `pkg/orchestrator/ops.go`:

```go
// NodeOSGet creates a step that retrieves OS information.
func (o *Orchestrator) NodeOSGet(
	target string,
) *Step {
	name := o.nextOpName("get-os")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.OS(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get os: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.OSResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}
```

Note: Check the exact SDK type name — it may be `osapi.OSInfoResult` or `osapi.OSResult`. Read `pkg/sdk/client/node_types.go` in the osapi repo to confirm.

- [ ] **Step 5: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add pkg/orchestrator/ops.go pkg/orchestrator/ops_public_test.go
git commit -m "feat: add NodeHostnameUpdate and NodeOSGet operations"
```

---

### Task 8: Add missing operations — FileUndeploy

**Files:**
- Modify: `pkg/orchestrator/ops.go`
- Modify: `pkg/orchestrator/ops_public_test.go`

- [ ] **Step 1: Write test for FileUndeploy**

- [ ] **Step 2: Implement FileUndeploy**

Check `pkg/sdk/client/node.go` for the `FileUndeploy` signature and `pkg/sdk/client/file_types.go` for the result type. Follow the FileDeploy pattern:

```go
// FileUndeploy creates a step that removes a previously deployed file
// from the target agent's filesystem.
func (o *Orchestrator) FileUndeploy(
	target string,
	path string,
) *Step {
	name := o.nextOpName("undeploy-file")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.FileUndeploy(ctx, target, path)
			if err != nil {
				return nil, fmt.Errorf("undeploy file: %w", err)
			}

			return &sdk.Result{
				JobID:   resp.Data.JobID,
				Changed: resp.Data.Changed,
				Data:    sdk.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}
```

Note: Verify the SDK method signature — it may take different params. Read `c.Node.FileUndeploy` in `pkg/sdk/client/node.go`.

- [ ] **Step 3: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

- [ ] **Step 4: Commit**

```bash
git add pkg/orchestrator/ops.go pkg/orchestrator/ops_public_test.go
git commit -m "feat: add FileUndeploy operation"
```

---

### Task 9: Add missing operations — all 5 Cron operations

**Files:**
- Modify: `pkg/orchestrator/ops.go`
- Modify: `pkg/orchestrator/ops_public_test.go`

- [ ] **Step 1: Read SDK cron method signatures**

Read `pkg/sdk/client/schedule.go` and `pkg/sdk/client/schedule_types.go` in the osapi repo to get exact method signatures and type names.

- [ ] **Step 2: Write tests for all 5 cron operations**

CronList, CronGet, CronCreate, CronUpdate, CronDelete.

- [ ] **Step 3: Implement CronList**

```go
// CronList creates a step that lists all cron entries on the target.
func (o *Orchestrator) CronList(
	target string,
) *Step {
	name := o.nextOpName("list-cron")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Schedule.CronList(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list cron: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CronEntryResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  false,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}
```

- [ ] **Step 4: Implement CronGet**

```go
func (o *Orchestrator) CronGet(
	target string,
	entryName string,
) *Step {
```

- [ ] **Step 5: Implement CronCreate**

Takes `target string` and `opts osapi.CronCreateOpts`. Returns Collection[CronMutationResult].

- [ ] **Step 6: Implement CronUpdate**

Takes `target string`, `entryName string`, and `opts osapi.CronUpdateOpts`.

- [ ] **Step 7: Implement CronDelete**

Takes `target string` and `entryName string`.

- [ ] **Step 8: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

- [ ] **Step 9: Commit**

```bash
git add pkg/orchestrator/ops.go pkg/orchestrator/ops_public_test.go
git commit -m "feat: add CronList, CronGet, CronCreate, CronUpdate, CronDelete operations"
```

---

### Task 10: Add missing operations — AgentDrain and AgentUndrain

**Files:**
- Modify: `pkg/orchestrator/ops.go`
- Modify: `pkg/orchestrator/ops_public_test.go`

- [ ] **Step 1: Read SDK agent drain/undrain signatures**

Read `pkg/sdk/client/agent.go` in the osapi repo.

- [ ] **Step 2: Write tests and implement AgentDrain and AgentUndrain**

These are non-collection operations (they target a single agent by hostname, not broadcast). Follow the AgentGet pattern — return `&sdk.Result{Changed: true/false, Data: ...}`.

- [ ] **Step 3: Run tests**

```bash
go test -v ./pkg/orchestrator/...
```

- [ ] **Step 4: Commit**

```bash
git add pkg/orchestrator/ops.go pkg/orchestrator/ops_public_test.go
git commit -m "feat: add AgentDrain and AgentUndrain operations"
```

---

### Task 11: Run full test suite and verify coverage

- [ ] **Step 1: Run all tests with coverage**

```bash
just go::unit
go test -coverprofile=/tmp/orch_cov.out ./pkg/orchestrator/...
go tool cover -func=/tmp/orch_cov.out | grep -v "100.0%" | grep -v "total:"
```

- [ ] **Step 2: Fix any coverage gaps**

- [ ] **Step 3: Run linter**

```bash
just go::vet
```

- [ ] **Step 4: Run formatter**

```bash
just go::fmt
```

---

### Task 12: Update documentation

**Files:**
- Create: `docs/operations/node-hostname-update.md`
- Create: `docs/operations/node-os-get.md`
- Create: `docs/operations/file-undeploy.md`
- Create: `docs/operations/cron-list.md`
- Create: `docs/operations/cron-get.md`
- Create: `docs/operations/cron-create.md`
- Create: `docs/operations/cron-update.md`
- Create: `docs/operations/cron-delete.md`
- Create: `docs/operations/agent-drain.md`
- Create: `docs/operations/agent-undrain.md`
- Modify: `docs/operations/README.md`
- Modify: `docs/features/guards.md`
- Modify: `docs/features/broadcast.md`
- Modify: `docs/features/README.md`

- [ ] **Step 1: Create operation docs**

Follow the exact pattern of `docs/operations/node-hostname-get.md` for each new operation. Include: description, usage, parameters, result type, idempotency, permissions, example reference.

- [ ] **Step 2: Update operations README index**

Add new operations to the table in `docs/operations/README.md`:

| Method | Operation | Idempotent | Category |
|--------|-----------|------------|----------|
| `NodeHostnameUpdate(target, hostname)` | Set system hostname | Idempotent | Node |
| `NodeOSGet(target)` | Get OS info | Read-only | Node |
| `FileUndeploy(target, path)` | Remove deployed file | Idempotent | File |
| `CronList(target)` | List cron entries | Read-only | Schedule |
| `CronGet(target, name)` | Get cron entry | Read-only | Schedule |
| `CronCreate(target, opts)` | Create cron entry | Non-idempotent | Schedule |
| `CronUpdate(target, name, opts)` | Update cron entry | Idempotent | Schedule |
| `CronDelete(target, name)` | Delete cron entry | Idempotent | Schedule |
| `AgentDrain(hostname)` | Drain agent | Idempotent | Agent |
| `AgentUndrain(hostname)` | Undrain agent | Idempotent | Agent |

- [ ] **Step 3: Update guards.md**

Add `OnlyIfAnyHostSkipped` to the Broadcast Guards table and add a section explaining skipped vs failed semantics.

- [ ] **Step 4: Update broadcast.md**

Add `Status` field to the Per-Host Results fields table. Add a section on status values (ok, skipped, failed).

- [ ] **Step 5: Update features README**

Add `OnlyIfAnyHostSkipped` to the Step Chaining table.

- [ ] **Step 6: Commit**

```bash
git add docs/
git commit -m "docs: add operation docs, host status semantics, and OnlyIfAnyHostSkipped"
```

---

### Task 13: Add examples

**Files:**
- Create: `examples/operations/hostname-update.go`
- Create: `examples/operations/cron.go`
- Create: `examples/features/host-status.go`

- [ ] **Step 1: Create hostname-update example**

Follow the `examples/operations/command.go` pattern. Show NodeHostnameUpdate targeting `_all` with a health check dependency.

- [ ] **Step 2: Create cron example**

Show CronCreate + CronList + CronDelete workflow with FileUpload dependency.

- [ ] **Step 3: Create host-status example**

Show a broadcast operation where some hosts are skipped, with OnlyIfAnyHostSkipped guard triggering a notification step and OnlyIfAnyHostFailed guard triggering a recovery step.

- [ ] **Step 4: Verify examples compile**

```bash
cd examples/operations && go build ./... && cd ../features && go build ./...
```

- [ ] **Step 5: Commit**

```bash
git add examples/
git commit -m "feat: add hostname-update, cron, and host-status examples"
```

---

### Task 14: Update README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update operation counts and tables**

Update the README to reflect the new operation count (from 27 to 37 operations). Add new operations to the relevant sections. Add OnlyIfAnyHostSkipped to features list.

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: update README with new operations and host status features"
```

---

### Task 15: Final verification

- [ ] **Step 1: Run full test suite**

```bash
just test
```

Expected: all pass

- [ ] **Step 2: Verify formatting and lint**

```bash
just go::fmt-check
just go::vet
```

- [ ] **Step 3: Build**

```bash
go build ./...
```

- [ ] **Step 4: Review git log**

```bash
git log --oneline main..HEAD
```

Verify commit messages follow conventional commits.
