# ServiceDelete

Deletes a service unit file from the target host.

## Usage

```go
step := o.ServiceDelete("web-01", "myapp")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to delete.                            |

## Result type

```go
var result osapi.ServiceMutationResult
err := results.Decode("delete-service-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Name of the deleted service.                                 |
| `Changed` | `bool`   | `true` if the operation removed the unit file.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the service does not exist, the operation returns
`Changed: false`. If the service exists, the operation removes it and returns
`Changed: true`.

## Permissions

The token must have the `service:write` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
