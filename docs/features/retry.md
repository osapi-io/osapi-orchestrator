# Retry

Automatically retry a step on failure up to N times.

## Usage

```go
o.NodeLoadGet("_any").
    After(health).
    Retry(3)
```

The step will be attempted up to 3 additional times if it fails. Retries happen
immediately without backoff.

## Example

See
[`examples/features/retry.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/retry.go)
for a complete working example.
