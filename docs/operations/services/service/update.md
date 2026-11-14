# ServiceUpdate

Updates a service unit file on the target host.

## Usage

```go
step := o.ServiceUpdate("web-01", "myapp", osapi.ServiceUpdateOpts{
    Object: "myapp-v2.service",
})
```

## Parameters

| Parameter     | Type                | Description                                               |
| ------------- | ------------------- | --------------------------------------------------------- |
| `target`      | `string`            | Target host: `_any`, `_all`, hostname, or label selector. |
| `serviceName` | `string`            | Name of the service to update.                            |
| `opts`        | `ServiceUpdateOpts` | Update options (refer to the table below).                |

### ServiceUpdateOpts

| Field    | Type     | Required | Description                                   |
| -------- | -------- | -------- | --------------------------------------------- |
| `Object` | `string` | Yes      | New Object Store reference for the unit file. |

## Result type

```go
var result osapi.ServiceMutationResult
err := results.Decode("update-service-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Name of the updated service.                                 |
| `Changed` | `bool`   | `true` if the operation changed the unit file.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current unit file with the new
content. It returns `Changed: true` only if it changed the file.

## Permissions

The token must have the `service:write` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
