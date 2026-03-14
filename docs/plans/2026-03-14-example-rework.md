# Example Rework Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development
> (if subagents available) or superpowers:executing-plans to implement this
> plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rewrite all examples so each one demonstrates its feature working
correctly — proving idempotency, showing changed/unchanged transitions, and
exercising guards/error paths for real.

**Architecture:** Each example is a self-contained `main()` that creates one or
more `orchestrator.New()` instances, runs plans, and prints narrated output
showing the feature in action. Safe mutating operations run twice to prove
idempotency. Dangerous operations use `OnError(Continue)`.

**Tech Stack:** Go, osapi-orchestrator DSL, osapi SDK types

**Spec:** `docs/superpowers/specs/2026-03-14-example-rework-design.md`

**Commits:** Per-task commits during development. Squash into 1-3 commits before
PR (one per chunk, or a single squash) per project conventions.

---

## Chunk 1: Operations Examples

### Task 1: Rewrite command.go

**Files:**

- Modify: `examples/operations/command.go`

- [ ] **Step 1: Rewrite command.go**

Replace the full file content. Two phases: successful commands showing decoded
output, then a failing command showing error handling.

```go
// Package main demonstrates command execution and error handling.
//
// Phase 1: exec and shell commands run in parallel — both always
// report changed=true (commands are non-idempotent by design).
//
// Phase 2: a command that fails shows how to inspect
// CommandResult.Error and ExitCode.
//
// DAG (phase 1):
//
//	health-check
//	    ├── run-uptime (exec)
//	    └── shell-uname -a (shell)
//
// DAG (phase 2):
//
//	run-ls /nonexistent (exec, error expected)
//
// Run with: OSAPI_TOKEN="<jwt>" go run command.go
```

Phase 1 plan:

- HealthCheck → CommandExec("\_any", "uptime") + CommandShell("\_any", "uname
  -a") in parallel
- Capture report, decode both as `osapi.CommandResult`, print stdout

Phase 2 plan:

- New orchestrator instance
- CommandExec("\_any", "ls", "/nonexistent") with `OnError(Continue)`
- Capture report, decode as `osapi.CommandResult`, print Error and ExitCode

Key imports: `context`, `fmt`, `log`, `os`,
`osapi "github.com/retr0h/osapi/pkg/sdk/client"`,
`"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"`

- [ ] **Step 2: Run command.go and verify output**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/operations/command.go`

Expected:

- Phase 1: both steps show `changed=true`, uptime stdout printed, uname stdout
  printed
- Phase 2: ls step shows `changed=true` with error, Error/ExitCode printed

- [ ] **Step 3: Commit**

```bash
git add examples/operations/command.go
git commit -m "docs(examples): rewrite command.go to show error handling"
```

### Task 2: Rewrite dns-update.go

**Files:**

- Modify: `examples/operations/dns-update.go`

- [ ] **Step 1: Rewrite dns-update.go**

Add `OnError(Continue)` to all steps so it doesn't crash on macOS. Add a
TaskFunc at the end to decode and print results if available, or print errors.

```go
// Package main demonstrates the read-then-write DNS pattern.
//
// Reads current DNS config, then updates with new servers. All steps
// use OnError(Continue) so the example runs on any platform — on
// macOS (no eth0) the DNS steps fail gracefully.
//
// DAG:
//
//	health-check
//	    └── get-dns (continue on error)
//	            └── update-dns (continue on error)
//
// Run with: OSAPI_TOKEN="<jwt>" go run dns-update.go
```

Plan:

- HealthCheck
- NetworkDNSGet("\_any", iface).OnError(Continue)
- NetworkDNSUpdate("\_any", iface, servers, domains).OnError(Continue)
- Capture report. Decode `osapi.DNSConfigResult` — if succeeds, print servers.
  If decode fails, print a note that DNS operations need a valid interface.

Key: Keep the existing `OSAPI_INTERFACE` env var support.

- [ ] **Step 2: Run dns-update.go and verify output**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/operations/dns-update.go`

Expected on macOS: steps show `[failed]` but program doesn't crash. Prints note
about needing valid interface.

- [ ] **Step 3: Commit**

```bash
git add examples/operations/dns-update.go
git commit -m "docs(examples): add OnError(Continue) to dns-update.go"
```

### Task 3: Rewrite file-deploy.go

**Files:**

- Modify: `examples/operations/file-deploy.go`

- [ ] **Step 1: Rewrite file-deploy.go**

Full self-contained lifecycle proving idempotency. Four phases using separate
orchestrator instances.

```go
// Package main demonstrates file deployment with idempotency proof.
//
// Phase 1: cleanup — remove any previously deployed file.
// Phase 2: first deploy — upload + deploy + verify (changed=true).
// Phase 3: idempotency — same upload + deploy + verify (changed=false).
// Phase 4: cleanup — remove the deployed file.
//
// Run with: OSAPI_TOKEN="<jwt>" go run file-deploy.go
```

Helper function `newOrchestrator(url, token string) *orchestrator.Orchestrator`
to reduce boilerplate.

Phase 1: `CommandShell("_any", "rm -f /tmp/app-config.yaml").OnError(Continue)`
— print "Phase 1: Cleanup" Phase 2: FileUpload → FileDeploy → FileStatusGet —
print "Phase 2: First deploy (expect changed=true)", decode and print status
Phase 3: Same as phase 2 — print "Phase 3: Idempotency check (expect
changed=false)", decode and print status Phase 4:
`CommandShell("_any", "rm -f /tmp/app-config.yaml").OnError(Continue)` — print
"Phase 4: Cleanup"

Key imports: `context`, `fmt`, `log`, `os`, `osapi`, `orchestrator`

- [ ] **Step 2: Run file-deploy.go and verify output**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/operations/file-deploy.go`

Expected:

- Phase 2: upload and deploy show `changed=true`, status shows `in-sync`
- Phase 3: upload and deploy show `changed=false`, status still `in-sync`

- [ ] **Step 3: Commit**

```bash
git add examples/operations/file-deploy.go
git commit -m "docs(examples): rewrite file-deploy.go with idempotency proof"
```

### Task 4: Rewrite file-changed.go

**Files:**

- Modify: `examples/operations/file-changed.go`

- [ ] **Step 1: Rewrite file-changed.go**

Demonstrates drift detection both ways.

```go
// Package main demonstrates drift detection with FileChanged.
//
// Phase 1: upload known content to establish baseline.
// Phase 2: FileChanged with same content — no drift (changed=false),
//          OnlyIfChanged steps are skipped.
// Phase 3: FileChanged with different content — drift detected
//          (changed=true), OnlyIfChanged steps run.
//
// Run with: OSAPI_TOKEN="<jwt>" go run file-changed.go
```

Phase 1: FileUpload("app-config.yaml", "raw", contentA) — print "Phase 1: Upload
baseline" Phase 2: New orchestrator. FileChanged("app-config.yaml", contentA) →
FileUpload.OnlyIfChanged → FileDeploy.OnlyIfChanged — print "Phase 2: Check with
same content (expect no drift)". Decode FileChanged result, print changed=false.
Phase 3: New orchestrator. FileChanged("app-config.yaml", contentB) →
FileUpload.OnlyIfChanged → FileDeploy.OnlyIfChanged — print "Phase 3: Check with
different content (expect drift)". Decode FileChanged result, print
changed=true.

Two content constants:

```go
contentA := []byte("server:\n  port: 8080\n  debug: false\n")
contentB := []byte("server:\n  port: 9090\n  debug: true\n")
```

- [ ] **Step 2: Run file-changed.go and verify output**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/operations/file-changed.go`

Expected:

- Phase 2: check-file shows `changed=false`, upload/deploy show `[skip]`
- Phase 3: check-file shows `changed=true`, upload/deploy show `changed=true`

- [ ] **Step 3: Commit**

```bash
git add examples/operations/file-changed.go
git commit -m "docs(examples): rewrite file-changed.go to show drift detection both ways"
```

---

## Chunk 2: Features Examples — Guards and Error Handling

### Task 5: Rewrite basic.go

**Files:**

- Modify: `examples/features/basic.go`

- [ ] **Step 1: Rewrite basic.go**

Add result decoding and print so it's not silent.

Keep existing structure but capture report and decode hostname:

```go
report, err := o.Run(context.Background())
if err != nil {
    log.Fatal(err)
}

var h osapi.HostnameResult
if err := report.Decode("get-hostname", &h); err == nil {
    fmt.Printf("Hostname: %s\n", h.Hostname)
}
```

Add imports: `fmt`, `osapi "github.com/retr0h/osapi/pkg/sdk/client"`

- [ ] **Step 2: Run and verify**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/features/basic.go`

Expected: prints hostname at the end.

- [ ] **Step 3: Commit**

```bash
git add examples/features/basic.go
git commit -m "docs(examples): add result decode to basic.go"
```

### Task 6: Rewrite parallel.go

**Files:**

- Modify: `examples/features/parallel.go`

- [ ] **Step 1: Rewrite parallel.go**

Remove the `CommandShell("_all", "cat /proc/version")` step that fails on macOS.
Keep the 5 parallel read-only queries. Add result decoding for hostname and
memory to prove parallel results are accessible.

Keep `WithVerbose()`. Remove the `OnError(Continue)` import since we're removing
the failing command.

After run, decode and print:

```go
var h osapi.HostnameResult
if err := report.Decode("get-hostname", &h); err == nil {
    fmt.Printf("\nHostname: %s\n", h.Hostname)
}

var m osapi.MemoryResult
if err := report.Decode("get-memory", &m); err == nil {
    fmt.Printf("Memory: %.1f GB total\n", float64(m.Total)/(1024*1024*1024))
}
```

- [ ] **Step 2: Run and verify**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/features/parallel.go`

Expected: 5 parallel queries all show `changed=false`, hostname and memory
printed at end.

- [ ] **Step 3: Commit**

```bash
git add examples/features/parallel.go
git commit -m "docs(examples): remove failing /proc/version from parallel.go"
```

### Task 7: Rewrite guards.go

**Files:**

- Modify: `examples/features/guards.go`

- [ ] **Step 1: Rewrite guards.go**

Show guard both passing and blocking using two separate orchestrator instances.

```go
// Package main demonstrates When guards — conditional execution
// based on prior step results.
//
// Plan 1: guard passes — whoami runs because hostname is non-empty.
// Plan 2: guard blocks — whoami is skipped because hostname doesn't
// match "impossible".
//
// Run with: OSAPI_TOKEN="<jwt>" go run guards.go
```

Plan 1:

```go
fmt.Println("=== Plan 1: Guard should PASS (hostname != \"\") ===")
o1 := orchestrator.New(url, token)
health1 := o1.HealthCheck()
hostname1 := o1.NodeHostnameGet("_any").After(health1)
o1.CommandExec("_any", "whoami").
    After(hostname1).
    When(func(r orchestrator.Results) bool {
        var h osapi.HostnameResult
        if err := r.Decode("get-hostname", &h); err != nil {
            return false
        }
        return h.Hostname != ""
    })
report1, err := o1.Run(context.Background())
// ... print summary
```

Plan 2:

```go
fmt.Println("\n=== Plan 2: Guard should BLOCK (hostname == \"impossible\") ===")
o2 := orchestrator.New(url, token)
health2 := o2.HealthCheck()
hostname2 := o2.NodeHostnameGet("_any").After(health2)
o2.CommandExec("_any", "whoami").
    After(hostname2).
    When(func(r orchestrator.Results) bool {
        var h osapi.HostnameResult
        if err := r.Decode("get-hostname", &h); err != nil {
            return false
        }
        return h.Hostname == "impossible"
    })
// ... run and print summary
```

- [ ] **Step 2: Run and verify**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/features/guards.go`

Expected:

- Plan 1: whoami runs (3 tasks, 1 changed, 2 unchanged)
- Plan 2: whoami is skipped (3 tasks, 2 unchanged, 1 skipped)

- [ ] **Step 3: Commit**

```bash
git add examples/features/guards.go
git commit -m "docs(examples): show guard passing and blocking in guards.go"
```

### Task 8: Rewrite only-if-changed.go

**Files:**

- Modify: `examples/features/only-if-changed.go`

- [ ] **Step 1: Rewrite only-if-changed.go**

Use file operations that actually produce `changed=true` to trigger the guard,
then show it blocking when nothing changed.

```go
// Package main demonstrates OnlyIfChanged with file operations.
//
// Phase 1: cleanup — remove deployed file.
// Phase 2: deploy file (changed=true) → OnlyIfChanged command runs.
// Phase 3: deploy same file (changed=false) → OnlyIfChanged command
//          is skipped.
// Phase 4: cleanup.
//
// Run with: OSAPI_TOKEN="<jwt>" go run only-if-changed.go
```

Phase 1 & 4:
`CommandShell("_any", "rm -f /tmp/only-if-changed.txt").OnError(Continue)`

Phase 2:

```go
o2 := orchestrator.New(url, token)
upload2 := o2.FileUpload("only-if-changed.txt", "raw", content)
deploy2 := o2.FileDeploy("_any", osapi.FileDeployOpts{
    ObjectName: "only-if-changed.txt", Path: "/tmp/only-if-changed.txt",
    ContentType: "raw", Mode: "0644",
}).After(upload2)
o2.CommandExec("_any", "echo", "post-deploy-hook").
    Named("post-deploy").After(deploy2).OnlyIfChanged()
```

Phase 3: Same plan — deploy shows `changed=false`, so post-deploy is skipped.

- [ ] **Step 2: Run and verify**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/features/only-if-changed.go`

Expected:

- Phase 2: deploy `changed=true`, post-deploy runs
- Phase 3: deploy `changed=false`, post-deploy `[skip]`

- [ ] **Step 3: Commit**

```bash
git add examples/features/only-if-changed.go
git commit -m "docs(examples): use file ops in only-if-changed.go to trigger guard"
```

### Task 9: Rewrite error-recovery.go

**Files:**

- Modify: `examples/features/error-recovery.go`

- [ ] **Step 1: Rewrite error-recovery.go**

Show both infrastructure-level and host-level failure recovery.

```go
// Package main demonstrates error recovery at two levels.
//
// Plan 1: Infrastructure failure — TaskFunc returns an error.
//         OnlyIfFailed cleanup runs.
//
// Plan 2: Host-level failure — command exits non-zero.
//         OnlyIfAnyHostFailed cleanup runs.
//
// Run with: OSAPI_TOKEN="<jwt>" go run error-recovery.go
```

Plan 1 (infrastructure failure):

```go
o1 := orchestrator.New(url, token)
deploy1 := o1.TaskFunc("deploy", func(_ context.Context, _ orchestrator.Results) (*sdk.Result, error) {
    return nil, fmt.Errorf("simulated deployment failure")
}).OnError(orchestrator.Continue)
o1.CommandExec("_any", "echo", "running-infra-cleanup").
    Named("cleanup").After(deploy1).OnlyIfFailed()
```

Plan 2 (host-level failure):

```go
o2 := orchestrator.New(url, token)
deploy2 := o2.CommandShell("_all", "cat /nonexistent-file").
    Named("deploy").OnError(orchestrator.Continue)
o2.CommandExec("_any", "echo", "running-host-cleanup").
    Named("cleanup").After(deploy2).OnlyIfAnyHostFailed()
```

Imports: add `fmt`, `sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"`

- [ ] **Step 2: Run and verify**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/features/error-recovery.go`

Expected:

- Plan 1: deploy shows `[failed]`, cleanup runs
- Plan 2: deploy shows `[changed]` with host error, cleanup runs

- [ ] **Step 3: Commit**

```bash
git add examples/features/error-recovery.go
git commit -m "docs(examples): trigger real failures in error-recovery.go"
```

### Task 10: Rewrite broadcast-guards.go

**Files:**

- Modify: `examples/features/broadcast-guards.go`

- [ ] **Step 1: Rewrite broadcast-guards.go**

Use a command that actually fails to trigger host-level guards.

```go
// Package main demonstrates broadcast-aware host-level guards.
//
// A deploy step broadcasts a failing command to all hosts with
// Continue error strategy. Downstream guards react to per-host
// outcomes:
// - OnlyIfAnyHostFailed: cleanup runs because the command errored.
// - OnlyIfAllHostsChanged: verify runs because commands always
//   report changed=true (non-idempotent).
//
// DAG:
//
//	deploy (_all, continue on error)
//	    ├── cleanup (only-if-any-host-failed)
//	    └── verify (only-if-all-hosts-changed)
//
// Run with: OSAPI_TOKEN="<jwt>" go run broadcast-guards.go
```

Change the deploy command from `echo deploying` to `cat /nonexistent-file` so it
actually fails:

```go
deploy := o.CommandShell("_all", "cat /nonexistent-file").
    Named("deploy").
    OnError(orchestrator.Continue)
```

Keep cleanup and verify the same.

- [ ] **Step 2: Run and verify**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/features/broadcast-guards.go`

Expected: deploy shows error, cleanup runs (host failed), verify runs
(changed=true).

- [ ] **Step 3: Commit**

```bash
git add examples/features/broadcast-guards.go
git commit -m "docs(examples): trigger real host failure in broadcast-guards.go"
```

---

## Chunk 3: Features Examples — TaskFunc and Minor Tweaks

### Task 11: Rewrite task-func.go

**Files:**

- Modify: `examples/features/task-func.go`

- [ ] **Step 1: Rewrite task-func.go**

Fix the hostname decode that prints empty string. The TaskFunc already reads
hostname correctly via `r.Decode("get-hostname", &h)` inside the func. The issue
is in the post-execution decode — `report.Decode("get-hostname", &h)` may not
populate the same fields because `Result.Data` contains the raw API collection
response.

Fix: decode the TaskFunc's own result (which has clean `Data` set by the
TaskFunc itself) instead of trying to re-decode the hostname step.

```go
// Post-execution: decode the TaskFunc result to prove data flow.
var summary map[string]any
if err := report.Decode("summarize", &summary); err == nil {
    fmt.Printf("Report summary host: %v\n", summary["host"])
}
```

- [ ] **Step 2: Run and verify**

Run:
`OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run examples/features/task-func.go`

Expected: TaskFunc prints hostname during execution, post-execution decode
prints the same hostname from summary data.

- [ ] **Step 3: Commit**

```bash
git add examples/features/task-func.go
git commit -m "docs(examples): fix result decode in task-func.go"
```

### Task 12: Minor tweaks to discovery examples

**Files:**

- Modify: `examples/features/discover.go`
- Modify: `examples/features/condition-filter.go`
- Modify: `examples/features/fact-predicates.go`
- Modify: `examples/features/label-filter.go`

These examples depend on fleet composition and may find 0 agents on a
single-host setup. The existing code already prints discovery results. Verify
each prints clear output about what was found or not found.

- [ ] **Step 1: Review and tweak discovery examples**

For each file, ensure it prints a clear message when 0 agents are found. The
existing code already does this for most — e.g., `discover.go` prints
`"Discovered %d matching agents"`. Check each file and add clearer output where
missing. If all four already have good output, no changes needed.

- [ ] **Step 2: Run and verify**

```bash
for f in discover.go condition-filter.go fact-predicates.go label-filter.go; do
    echo "=== $f ==="
    OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run "examples/features/$f" 2>&1
    echo
done
```

- [ ] **Step 3: Commit if changes were needed**

```bash
git add examples/features/discover.go examples/features/condition-filter.go \
    examples/features/fact-predicates.go examples/features/label-filter.go
git commit -m "docs(examples): improve discovery example output"
```

### Task 13: Run all examples and verify

- [ ] **Step 1: Run all operations examples**

```bash
for f in examples/operations/*.go; do
    echo "=== $f ==="
    OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run "$f" 2>&1
    echo
done
```

Verify each shows the expected changed/unchanged/error behavior.

- [ ] **Step 2: Run all features examples**

```bash
for f in examples/features/*.go; do
    echo "=== $f ==="
    OSAPI_VERBOSE=1 OSAPI_TOKEN="$OSAPI_TOKEN" go run "$f" 2>&1
    echo
done
```

Verify each shows the expected behavior.

- [ ] **Step 3: Run linter**

```bash
just go::vet
```

Fix any lint issues.

- [ ] **Step 4: Run tests**

```bash
just go::unit
```

Ensure no regressions.

- [ ] **Step 5: Final commit if any fixes were needed**

```bash
git add -A
git commit -m "docs(examples): fix lint and test issues from example rework"
```
