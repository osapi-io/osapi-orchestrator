# Retry

Automatically retry a step on failure up to N times. Retries can be immediate or
use exponential backoff to avoid overwhelming a recovering service.

## Usage

### Immediate Retry

Retry up to 3 times with no delay between attempts:

```go
o.TaskFunc("flaky-step", myFunc).
    Retry(3)
```

### Retry with Default Exponential Backoff

Use `WithExponentialBackoff()` for sensible defaults (1s initial, 30s max):

```go
o.TaskFunc("flaky-step", myFunc).
    Retry(3, orchestrator.WithExponentialBackoff())
```

### Retry with Custom Backoff

Use `WithBackoff(initial, max)` for custom intervals. The delay doubles on each
attempt, clamped to the max interval:

```go
o.TaskFunc("flaky-step", myFunc).
    Retry(5, orchestrator.WithBackoff(500*time.Millisecond, 5*time.Second))
```

## Transient Poll Errors

Job polling automatically retries transient HTTP errors (404, 500) with
exponential backoff. This handles the race where the agent hasn't written
results yet when the SDK first polls. Non-transient errors (401, 403, network
failures) fail immediately.

## Example

See
[`examples/features/retry.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/retry.go)
for a complete working example that simulates transient failures to demonstrate
all three retry strategies.
