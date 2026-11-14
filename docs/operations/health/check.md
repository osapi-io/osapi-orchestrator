# HealthCheck

Runs a liveness probe on the OSAPI API server. The probe is a small connection
check. It makes sure that the server is reachable and that it responds. The
probe does not use the job system. It calls the health endpoint directly.

## Usage

```go
step := o.HealthCheck()
```

## Parameters

None. The health check always sends the probe to the configured API server.

## Result type

`HealthCheck` does not return typed result data. If the server responds with
HTTP 200, the step succeeds. If not, the step fails.

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `health:read` permission.

## Example

Refer to
[`examples/features/basic.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/basic.go)
for a complete example.
