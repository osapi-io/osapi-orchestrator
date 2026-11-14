# Error recovery

Use `OnError(Continue)` to continue to run independent tasks when a step fails.
Use `OnlyIfFailed` to start a cleanup step only after a failure.

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

With the `Continue` strategy, the plan continues if some hosts fail. The cleanup
step runs only if one or more dependencies failed.

## Error strategies

| Strategy   | Behavior                                                       |
| ---------- | -------------------------------------------------------------- |
| `StopAll`  | Stop at the first failure and cancel all other steps (default) |
| `Continue` | Skip the dependents and continue to run the independent tasks  |

## Broadcast error recovery

For a broadcast operation, `OnlyIfFailed` triggers when the task-level status is
`Failed`. The status is `Failed` when one or more hosts fail. To make a decision
from the result of each host, use the host-level guards:

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

These guards are different from the task-level `OnlyIfFailed`:

- `OnlyIfFailed` checks `Status == Failed`. It triggers on each failure.
- `OnlyIfAnyHostFailed` checks `HostResult.Error`. It examines each host.
- `OnlyIfAllHostsFailed` triggers only if each host has an error.

Refer to [Guards: Broadcast Guards](guards.md#broadcast-guards) for the full
reference.

## Example

Refer to
[`examples/features/error-recovery.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/error-recovery.go)
and
[`examples/features/broadcast-guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast-guards.go)
for complete examples.
