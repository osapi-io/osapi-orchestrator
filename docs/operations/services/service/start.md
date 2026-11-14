# ServiceStart

Starts a service on the target host.

## Usage

```go
step := o.ServiceStart("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | The name of the service to start.                         |

## Result type

```go
var result osapi.ServiceMutationResult
err := results.Decode("start-service-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The name of the started service.                             |
| `Changed` | `bool`   | `true` if the operation started the service.                 |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the service already runs, the operation returns
`Changed: false`.

## Permissions

The token must have the `service:write` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
