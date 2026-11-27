# ServiceGet

Get one service.

Reads one service from the same list, by its name with .service or without it. A
service that osapi deployed reads by the name it was deployed with, such as app,
as its writes take it, or by its unit's name, osapi-app. A service that systemd
has not loaded is not found. On a Mac, a daemon that osapi deployed and launchd
has not loaded is dead.

## Usage

```go
step := o.ServiceGet("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.ServiceEntry` row. A list answers with one row
for each item.

```go
var row osapi.ServiceEntry
err := report.Decode("get-service", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Description` | `string`          | The Description= of its unit.                                                                                                                                                                                            |
| `Enabled`     | `bool`            | Whether it starts at boot: the state of its unit file.                                                                                                                                                                   |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Managed`     | `bool`            | Whether osapi deployed it, so its writes and actions act on it.                                                                                                                                                          |
| `Name`        | `string`          | The service, without .service.                                                                                                                                                                                           |
| `State`       | `string`          | What systemd says it does now, such as running, exited or dead.                                                                                                                                                          |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `service:read` permission. The step sends
`GET /api/node/{hostname}/service/{name}`.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go).
