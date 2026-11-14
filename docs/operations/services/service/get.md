# ServiceGet

Gets a specific service by name from the target host.

## Usage

```go
step := o.ServiceGet("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | The name of the service to get.                           |

## Result type

```go
var result osapi.ServiceGetResult
err := results.Decode("get-service-1", &result)
```

| Field      | Type           | Description                                                  |
| ---------- | -------------- | ------------------------------------------------------------ |
| `Hostname` | `string`       | The hostname of the agent that ran the operation.            |
| `Service`  | `*ServiceInfo` | The service details (refer to ServiceList).                  |
| `Error`    | `string`       | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `service:read` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
