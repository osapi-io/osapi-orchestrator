# ServiceGet

Retrieves a specific service by name from the target node.

## Usage

```go
step := o.ServiceGet("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to retrieve.                          |

## Result Type

```go
var result osapi.ServiceGetResult
err := results.Decode("get-service-1", &result)
```

| Field      | Type           | Description                                      |
| ---------- | -------------- | ------------------------------------------------ |
| `Hostname` | `string`       | The node's hostname.                             |
| `Service`  | `*ServiceInfo` | Service details (see ServiceList).               |
| `Error`    | `string`       | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `service:read` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
