# Error Recovery

Use `OnError(Continue)` to keep independent tasks running when a step fails, and
`OnlyIfFailed` to trigger cleanup only on failure.

## Usage

```go
deploy := o.CommandExec("_all", "echo", "deploying").
    Named("deploy").
    OnError(orchestrator.Continue)

o.CommandExec("_any", "echo", "running-cleanup").
    Named("cleanup").
    After(deploy).
    OnlyIfFailed()
```

The `Continue` strategy lets the plan proceed even if some hosts fail. The
cleanup step only runs when at least one dependency failed.

## Error Strategies

| Strategy   | Behavior                                        |
| ---------- | ----------------------------------------------- |
| `StopAll`  | Fail fast, cancel everything (default)          |
| `Continue` | Skip dependents, keep running independent tasks |

## Broadcast Error Recovery

For broadcast operations, `OnlyIfFailed` triggers when the task-level status is
`Failed` — this happens when any host fails. To make finer-grained decisions
based on per-host results, use the host-level guards:

```go
deploy := o.CommandExec("_all", "deploy.sh").
    Named("deploy").
    OnError(orchestrator.Continue)

// Run cleanup when at least one host failed (not all).
o.CommandExec("_any", "partial-cleanup.sh").
    Named("partial-cleanup").
    After(deploy).
    OnlyIfAnyHostFailed()

// Run full rollback only when every host failed.
o.CommandExec("_all", "rollback.sh").
    Named("full-rollback").
    After(deploy).
    OnlyIfAllHostsFailed()
```

The difference from task-level `OnlyIfFailed`:

- `OnlyIfFailed` checks `Status == Failed` — triggers on any failure
- `OnlyIfAnyHostFailed` checks `HostResult.Error` — inspects individual hosts
- `OnlyIfAllHostsFailed` requires every host to have an error

See [Guards — Broadcast Guards](guards.md#broadcast-guards) for the full
reference.

## Example

See
[`examples/features/error-recovery.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/error-recovery.go)
and
[`examples/features/broadcast-guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast-guards.go)
for complete working examples.
