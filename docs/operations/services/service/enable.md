# ServiceEnable

Enables a service to start on boot on the target node.

## Usage

```go
step := o.ServiceEnable("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to enable.                            |

## Result Type

```go
var result osapi.ServiceMutationResult
err := results.Decode("enable-service-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Name of the enabled service.                      |
| `Changed` | `bool`   | Whether the service was enabled.                  |
| `Error`   | `string` | Error message if enable failed; empty on success. |

## Idempotency

**Idempotent.** If the service is already enabled, returns `Changed: false`.

## Permissions

Requires `service:write` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
