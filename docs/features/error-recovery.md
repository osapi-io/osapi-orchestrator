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

A broadcast step fails when one or more of its hosts fail. A skipped host is not
a failure. With the `Continue` strategy, the plan continues if some hosts fail.
The cleanup step runs only if one or more dependencies failed.

## What Run returns

`Run` returns an error when a step failed, also with `Continue`. `Continue` lets
the independent steps run, and it does not hide the failure. The error joins the
failure of each step that failed, with its name. Thus a program that exits with
an error code when `Run` returns an error does that also after a cleanup step
ran. To make a decision from the steps that failed, read the report.

## Error strategies

| Strategy   | Behavior                                                       |
| ---------- | -------------------------------------------------------------- |
| `StopAll`  | Stop at the first failure and cancel all other steps (default) |
| `Continue` | Skip the dependents and continue to run the independent tasks  |

With `StopAll`, the first failure cancels the context of the other steps in its
level. A step that waits on its context stops at once, and the plan runs no
later level. A step that the cancel stopped, and a step that did not start, is
skipped with the reason "plan stopped". It is not a failure, so the report
counts only the step that failed. A step that returns the cancel as an error
with `%v`, not `%w`, is also skipped with that reason. A step that fails for its
own reason after the cancel stays failed. The duration of a skipped step is the
time of the skip, not the time of the plan.

## Cancellation

When the context of `Run` is cancelled, the steps that run get the cancelled
context, and a step that ends with that error fails. A step that did not start
is skipped with the reason "plan stopped", and no later level starts. `Run`
returns the context error, so a cancelled run does not report a success.

## A step that panics or returns no result

A panic in a step, in its guard or in a hook of the step fails that step. The
error names the step and the panic value, and has the stack. The `AfterTask`
hook still runs for the step, and the error strategy then applies. A step that
returns no result and no error also fails.

A panic in a hook of the plan or of a level (`BeforePlan`, `AfterPlan`,
`BeforeLevel`, `AfterLevel`) is dropped, and the plan runs on. These hooks only
show the run, such as the output of the renderer. A panic in the skip hooks of a
step that the plan did not reach is also dropped.

`Run` refuses a plan with a step that has no function. The error names the step.

## A dependency that is not in the plan

`Run` refuses a plan in which a step depends on a step of a different plan. The
error names both steps. A dependency of that kind never runs, so the dependent
would run without it.

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
- `OnlyIfAnyHostFailed` examines each host. It triggers if one or more hosts
  failed.
- `OnlyIfAllHostsFailed` triggers only if all hosts failed.
- A host failed if its `Status` is `failed`, or if its `Error` is not empty and
  its `Status` is not `skipped`.

Refer to [Guards: Broadcast Guards](guards.md#broadcast-guards) for the full
reference.

## Example

Refer to
[`examples/features/error-recovery.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/error-recovery.go)
and
[`examples/features/broadcast-guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast-guards.go)
for complete examples.
