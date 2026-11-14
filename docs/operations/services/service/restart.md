# ServiceRestart

Restarts a service on the target host.

## Usage

```go
step := o.ServiceRestart("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | The name of the service to restart.                       |

## Result type

```go
var result osapi.ServiceMutationResult
err := results.Decode("restart-service-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The name of the restarted service.                           |
| `Changed` | `bool`   | `true` if the operation restarted the service.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** The operation restarts the service each time the step runs.
It does not examine the service state first.

## Permissions

The token must have the `service:write` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
