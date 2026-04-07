# ServiceStop

Stops a service on the target node.

## Usage

```go
step := o.ServiceStop("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to stop.                              |

## Result Type

```go
var result osapi.ServiceMutationResult
err := results.Decode("stop-service-1", &result)
```

| Field     | Type     | Description                                     |
| --------- | -------- | ----------------------------------------------- |
| `Name`    | `string` | Name of the stopped service.                    |
| `Changed` | `bool`   | Whether the service was stopped.                |
| `Error`   | `string` | Error message if stop failed; empty on success. |

## Idempotency

**Idempotent.** If the service is already stopped, returns `Changed: false`.

## Permissions

Requires `service:write` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
