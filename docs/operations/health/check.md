# HealthCheck

Runs a liveness probe against the OSAPI API server. This is a lightweight
connectivity check that confirms the server is reachable and responding. It does
not go through the job system -- the probe calls the health endpoint directly.

## Usage

```go
step := o.HealthCheck()
```

## Parameters

None. The health check always targets the configured API server.

## Result Type

`HealthCheck` does not return typed result data. The step succeeds if the server
responds with HTTP 200, and fails otherwise.

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `health:read` permission.

## Example

See
[`examples/features/basic.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/basic.go)
for a complete working example.
