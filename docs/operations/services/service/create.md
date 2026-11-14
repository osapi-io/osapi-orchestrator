# ServiceCreate

Creates a service unit file on the target host. Upload the unit file to the NATS
Object Store first. Refer to [FileUpload](../../files/file/upload.md).

## Usage

```go
step := o.ServiceCreate("web-01", osapi.ServiceCreateOpts{
    Name:   "myapp",
    Object: "myapp.service",
})
```

## Parameters

| Parameter | Type                | Description                                               |
| --------- | ------------------- | --------------------------------------------------------- |
| `target`  | `string`            | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `ServiceCreateOpts` | Create options (refer to the table below).                |

### ServiceCreateOpts

| Field    | Type     | Required | Description                               |
| -------- | -------- | -------- | ----------------------------------------- |
| `Name`   | `string` | Yes      | Service unit name.                        |
| `Object` | `string` | Yes      | Object Store reference for the unit file. |

## Result type

```go
var result osapi.ServiceMutationResult
err := results.Decode("create-service-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Name of the created service.                                 |
| `Changed` | `bool`   | `true` if the operation created the unit file.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the service already exists, the operation returns an
error. Use [ServiceUpdate](update.md) to change a unit file that exists.

## Permissions

The token must have the `service:write` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
