# Example Rework Design

Rewrite all examples in `examples/operations/` and `examples/features/` so each
one demonstrates its feature working correctly — proving idempotency, showing
changed/unchanged transitions, and exercising guards/error paths for real.

## Principles

1. **Self-contained**: Each example cleans up before and after itself.
2. **Show the contrast**: Safe mutating operations run twice — first time shows
   `changed=true`, second time shows `changed=false` (idempotency proof).
3. **Narrate with prints**: Use `fmt.Println` headers to explain what each phase
   is doing and what to expect.
4. **Dangerous ops still run**: Operations like DNS update that don't work on
   macOS still execute but use `OnError(Continue)` so they don't crash.
5. **Multiple plans per example are fine**: A single `main()` can create
   multiple `orchestrator.New()` instances and call `Run()` on each to show
   before/after contrast. Each orchestrator instance has its own independent
   plan.

## Operations Examples

### command.go

Commands are non-idempotent — always `changed=true`. This is correct behavior.

- Phase 1: Run `uptime` (exec) and `uname -a` (shell) in parallel.
- Phase 2: Run a command that fails (e.g., `ls /nonexistent`) to show error
  handling — decode `CommandResult.Error` and `ExitCode`.
- Print decoded stdout, stderr, exit code for each.

No idempotency proof needed — commands are explicitly non-idempotent by design.

### dns-update.go

Dangerous — doesn't work on macOS (no `eth0`). Still runs, ignores failure.

- HealthCheck → NetworkDNSGet → NetworkDNSUpdate, all with `OnError(Continue)`.
- Print a message explaining the read-then-write pattern.
- If DNS get succeeds, print current config. If update succeeds, print result.
- If either fails, print the error but don't crash.

### file-deploy.go

Safe, works on macOS. Full self-contained lifecycle.

- Plan 1 (cleanup): `CommandShell("rm -f /tmp/app-config.yaml")` with
  `OnError(Continue)`.
- Plan 2 (first deploy): FileUpload → FileDeploy → FileStatusGet. Should show
  upload `changed=true`, deploy `changed=true`, status `in-sync`.
- Plan 3 (idempotency): FileUpload → FileDeploy → FileStatusGet with same
  content. Should show upload `changed=false`, deploy `changed=false`, status
  still `in-sync`.
- Plan 4 (cleanup): Remove the deployed file.

### file-changed.go

Safe, works on macOS. Demonstrates drift detection.

- Plan 1 (setup): FileUpload with content A.
- Plan 2 (no drift): FileChanged with content A → `changed=false` →
  OnlyIfChanged upload/deploy are skipped.
- Plan 3 (drift detected): FileChanged with content B → `changed=true` →
  OnlyIfChanged upload/deploy run.
- Plan 4 (cleanup): Remove uploaded object-store file if possible, or note that
  object-store uploads are ephemeral.
- Print FileChanged result (name, changed, sha256) for each phase.

## Features Examples

### basic.go

Minimal example — keep simple but add result decode + print so it's not silent.

- HealthCheck → NodeHostnameGet → decode and print hostname.

### parallel.go

Show parallel execution of read-only queries.

- HealthCheck → 5 parallel queries (hostname, disk, memory, load, uptime).
- Remove `cat /proc/version` (fails on macOS, adds noise).
- Decode and print one result (e.g., hostname) to prove it worked.

### guards.go

Show guard both passing and blocking. Each plan uses a separate
`orchestrator.New()` instance.

- Plan 1 (guard passes): HealthCheck → NodeHostnameGet → `whoami` guarded by
  `When(hostname != "")` → runs.
- Plan 2 (guard blocks): HealthCheck → NodeHostnameGet → `whoami` guarded by
  `When` that decodes hostname and checks for `"impossible"` → skipped.
- Print report summaries showing which commands ran and which were skipped.

### only-if-changed.go

Use file operations that actually produce `changed=true`.

- Plan 1 (cleanup): Remove deployed file with `OnError(Continue)`.
- Plan 2 (triggers guard): FileUpload → FileDeploy (`changed=true`) →
  CommandExec guarded by `OnlyIfChanged` → runs.
- Plan 3 (guard blocks): FileDeploy same content (`changed=false`) → CommandExec
  guarded by `OnlyIfChanged` → skipped.
- Plan 4 (cleanup).

### error-recovery.go

Demonstrate both levels of error recovery.

**Key distinction**: The SDK has two levels of failure:

- **Infrastructure failure** (`StatusFailed`): API/network errors. Checked by
  `OnlyIfFailed`.
- **Host-level error** (`HostResult.Error`): Command ran but exited non-zero.
  Checked by `OnlyIfAnyHostFailed`.

Show both:

- Part 1: TaskFunc "deploy" that returns an error with `OnError(Continue)`.
  Cleanup guarded by `OnlyIfFailed` → runs because deploy failed at
  infrastructure level.
- Part 2: CommandShell that exits non-zero with `OnError(Continue)`. Cleanup
  guarded by `OnlyIfAnyHostFailed` → runs because the command produced a
  host-level error.

### broadcast.go

Keep as-is — demonstrates `_all` targeting. Works with single host on Mac.

### broadcast-guards.go

Demonstrate host-level guards with both success and failure paths.

**Important**: `OnlyIfAnyHostFailed` checks `HostResult.Error` on the host
results, not `StatusFailed` on the step. A command that exits non-zero on the
agent DOES populate `HostResult.Error`, so `CommandExec`/`CommandShell` with a
failing command will trigger host-level failure guards correctly.

- Deploy step: CommandShell `_all` with a command that fails (e.g.,
  `cat /nonexistent`) with `OnError(Continue)`.
- Cleanup guarded by `OnlyIfAnyHostFailed` → runs (host has error).
- Verify guarded by `OnlyIfAllHostsChanged` → runs (commands are non-idempotent,
  always `changed=true` even on error).

### retry.go

Already works well with TaskFunc simulating transient failures. Keep as-is.

### verbose.go

Keep as-is. Demonstrates `WithVerbose()` output.

### task-func.go

Fix the hostname decode (printed empty string). Show TaskFunc reading prior
results, computing something, and returning data. Decode the TaskFunc result
post-execution to prove data flow works.

### agent-facts.go

Already works well. Keep as-is.

### discover.go, condition-filter.go, fact-predicates.go, label-filter.go

These depend on fleet composition — find 0 agents on a single Mac. Keep the
patterns but ensure they print what they found (or didn't find) clearly. These
are inherently fleet-demo examples.

### group-by-fact.go

Already works (groups by darwin on Mac). Keep as-is.

### when-fact.go

Already works (skips correctly on macOS since not Ubuntu). Keep as-is.

## Examples That Need No Changes

- `retry.go` — already demonstrates retry mechanics well
- `verbose.go` — demonstrates WithVerbose() output
- `agent-facts.go` — comprehensive agent fact display
- `group-by-fact.go` — groups by OS correctly
- `when-fact.go` — WhenFact guard works correctly
- `broadcast.go` — demonstrates \_all targeting

## Examples That Need Rework

- `command.go` — add error handling demonstration
- `dns-update.go` — add OnError(Continue) so it doesn't crash
- `file-deploy.go` — full lifecycle with idempotency proof
- `file-changed.go` — show drift detection both ways
- `basic.go` — add result decode + print
- `parallel.go` — remove failing /proc/version, add result decode
- `guards.go` — show guard passing and blocking
- `only-if-changed.go` — use operations that actually change
- `error-recovery.go` — actually cause a failure
- `broadcast-guards.go` — actually trigger host failure
- `task-func.go` — fix hostname decode

## Examples That Need Minor Tweaks

- `discover.go` — ensure clear output about what was found
- `condition-filter.go` — same
- `fact-predicates.go` — same
- `label-filter.go` — same
