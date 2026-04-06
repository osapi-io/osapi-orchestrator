# ServiceDelete

Deletes a service unit file from the target node.

## Usage

```go
step := o.ServiceDelete("web-01", "myapp")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to delete.                            |

## Result Type

```go
var result osapi.ServiceMutationResult
err := results.Decode("delete-service-1", &result)
```

| Field     | Type     | Description                                          |
| --------- | -------- | ---------------------------------------------------- |
| `Name`    | `string` | Name of the deleted service.                         |
| `Changed` | `bool`   | Whether the unit file was removed or already absent. |
| `Error`   | `string` | Error message if deletion failed; empty on success.  |

## Idempotency

**Idempotent.** If the service does not exist, the operation returns
`Changed: false`. If it exists and is removed, returns `Changed: true`.

## Permissions

Requires `service:write` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
