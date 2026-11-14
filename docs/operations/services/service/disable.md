# ServiceDisable

Disables a service on the target host, so that it does not start at boot.

## Usage

```go
step := o.ServiceDisable("web-01", "nginx")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string` | Name of the service to disable.                           |

## Result type

```go
var result osapi.ServiceMutationResult
err := results.Decode("disable-service-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Name of the disabled service.                                |
| `Changed` | `bool`   | `true` if the operation disabled the service.                |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the service is already disabled, the operation returns
`Changed: false`.

## Permissions

The token must have the `service:write` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
