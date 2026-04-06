# ServiceUpdate

Updates a service unit file on the target node.

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
| `opts`        | `ServiceUpdateOpts` | Update options (see below).                               |

### ServiceUpdateOpts

| Field    | Type     | Required | Description                                    |
| -------- | -------- | -------- | ---------------------------------------------- |
| `Object` | `string` | Yes      | New Object Store reference for the unit file.  |

## Result Type

```go
var result osapi.ServiceMutationResult
err := results.Decode("update-service-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Name of the updated service.                      |
| `Changed` | `bool`   | Whether the unit file was actually modified.      |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current unit file against the new content. Returns
`Changed: true` only if the file was actually modified.

## Permissions

Requires `service:write` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
