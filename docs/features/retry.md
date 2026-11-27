# Retry

The orchestrator can retry a failed step up to N times. It waits between the
attempts, with exponential backoff, so that a service that recovers does not get
too many requests.

## Usage

### Retry with the default backoff

Retry up to 3 times. The first wait is 1s, and each wait doubles, up to 30s:

```go
o.TaskFunc("flaky-step", myFunc).
    Retry(3)
```

Each wait is a random time in the upper half of its delay: 0.5s to 1s, then 1s
to 2s. Thus clients that failed at the same time do not retry at the same time.

### Retry with custom backoff

Use `WithBackoff(initial, max)` to set the intervals. The delay doubles after
each attempt, but it does not go above the maximum interval. A maximum of 0 caps
the delay at 30s, or at the initial interval when that is longer. An initial
interval of 0 uses the default of 1s:

```go
o.TaskFunc("flaky-step", myFunc).
    Retry(5, orchestrator.WithBackoff(500*time.Millisecond, 5*time.Second))
```

### Retry, then continue

`Retry` and `ContinueOnError` work together. The step runs again up to N times.
If the last attempt fails, the plan continues with the steps that do not depend
on it:

```go
o.TaskFunc("flaky-step", myFunc).
    Retry(3).
    ContinueOnError()
```

## What a retry runs again

A retry runs only when another attempt can give a different result:

| The error                                              | Retried |
| ------------------------------------------------------ | ------- |
| An API answer of 5xx, 408 or 429                       | Yes     |
| A lost connection, a timeout, or an error with no type | Yes     |
| Any other API answer, such as 400, 401, 403, 404, 409  | No      |
| A cancelled context                                    | No      |
| An error that `orchestrator.Permanent` wraps           | No      |
| A `TaskFunc` that returns no result and no error       | No      |

A `TaskFunc` returns `orchestrator.Permanent(err)` for a failure that is the
same on each attempt, such as a bad input. `orchestrator.Retryable(err)` gives
the same answer that the runner uses.

Some ops are not idempotent: a second run does not do what one run does. They
are `CommandExec`, `CommandShell`, `DockerCreate`, `DockerExec`,
`ServiceRestart`, `ProcessSignal`, `PowerReboot` and `PowerShutdown`. `Retry`
does not run them again, also for an error that a retry can fix. Add
`RetryUnsafe` when two runs are safe, such as a command that checks before it
changes anything:

```go
o.CommandExec("_all", "/usr/local/bin/converge").
    Retry(3).
    RetryUnsafe()
```

### A retry for each step of the plan

`WithRetry` on `New` gives each step a retry, with the same backoff. A step that
calls `Retry` uses its own. `WithOnError` sets the error strategy of each step
in the same way:

```go
o := orchestrator.New(url, token,
    orchestrator.WithRetry(3),
    orchestrator.WithOnError(orchestrator.Continue),
)
```

## A retry of a broadcast

A broadcast step fails when one or more of its hosts fail. Its retry sends the
call to each host again, also to the hosts that succeeded. osapi's writes are
idempotent: a create of a thing that exists, or a delete of a thing that is
absent, answers unchanged. Thus a host that already made the change reports no
change on the retry. A retry does not send the call only to the hosts that
failed, because the target, such as `_all` or a label, can name different hosts
on the next attempt.

## Cancellation

A retry stops when the context is cancelled. The step fails with the error of
its last attempt, and no attempt starts after the cancel. A negative count is
the same as `Retry(0)`: the step runs one time.

## Transient poll errors

When the SDK polls for a job result, it retries transient HTTP errors (404, 500)
with exponential backoff. Sometimes the SDK polls before the agent writes its
result, and this retry handles that race. Other errors (401, 403, network
failures) fail immediately.

## Example

Refer to
[`examples/features/retry.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/retry.go)
for a complete example. The example simulates transient failures to show the
default backoff and a custom one.
