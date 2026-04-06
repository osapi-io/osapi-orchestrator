# ServiceCreate

Creates a service unit file on the target node. The unit file must be uploaded
to the NATS Object Store first (see [FileUpload](../file/upload.md)).

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
| `opts`    | `ServiceCreateOpts` | Create options (see below).                               |

### ServiceCreateOpts

| Field    | Type     | Required | Description                                    |
| -------- | -------- | -------- | ---------------------------------------------- |
| `Name`   | `string` | Yes      | Service unit name.                             |
| `Object` | `string` | Yes      | Object Store reference for the unit file.      |

## Result Type

```go
var result osapi.ServiceMutationResult
err := results.Decode("create-service-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Name of the created service.                        |
| `Changed` | `bool`   | Whether the unit file was created.                  |
| `Error`   | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating a service that already exists returns an error.
Use [ServiceUpdate](update.md) to modify existing unit files.

## Permissions

Requires `service:write` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
