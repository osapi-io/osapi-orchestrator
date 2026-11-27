# ServiceUpdate

Deploy a new unit file for a service.

Writes the object as the unit file of a service that osapi deployed, then
reloads systemd when the file changed. A service that osapi did not deploy is
not found. On a Mac, it writes the plist, and a daemon that runs keeps its old
plist until restart.

## Usage

```go
step := o.ServiceUpdate("web-01", input)
```

## Parameters

| Parameter | Type                       | Description                                                   |
| --------- | -------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                   | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ServiceUpdateInput` | Refer to `ServiceUpdateInput` below.                          |

### ServiceUpdateInput

| Field    | Type     | Description                                                                                          |
| -------- | -------- | ---------------------------------------------------------------------------------------------------- |
| `Name`   | `string` | The name of the path: the service, such as app for osapi-app.service. It names the item in the path. |
| `Object` | `string` | The object that holds the unit file, from osapi client file upload.                                  |

## Result

Each host answers with one `osapi.ServiceResult` row. A list answers with one
row for each item.

```go
var row osapi.ServiceResult
err := report.Decode("update-service", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The service.                                                                                                                                                                                                             |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `service:write` permission. The step sends
`PUT /api/node/{hostname}/service/{name}`.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go).
