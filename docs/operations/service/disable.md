# ServiceDisable

Disables a service from starting on boot on the target node.

## Usage

```go
step := o.ServiceDisable("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to disable.                           |

## Result Type

```go
var result osapi.ServiceMutationResult
err := results.Decode("disable-service-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Name of the disabled service.                       |
| `Changed` | `bool`   | Whether the service was disabled.                   |
| `Error`   | `string` | Error message if disable failed; empty on success.  |

## Idempotency

**Idempotent.** If the service is already disabled, returns `Changed: false`.

## Permissions

Requires `service:write` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
