# ServiceRestart

Restarts a service on the target node.

## Usage

```go
step := o.ServiceRestart("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to restart.                           |

## Result Type

```go
var result osapi.ServiceMutationResult
err := results.Decode("restart-service-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Name of the restarted service.                      |
| `Changed` | `bool`   | Whether the service was restarted.                  |
| `Error`   | `string` | Error message if restart failed; empty on success.  |

## Idempotency

**Non-idempotent.** Always restarts the service regardless of current state.

## Permissions

Requires `service:write` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
