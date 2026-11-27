# ServiceList

List the services of the host.

Lists each service that systemd has loaded, with its state, whether it starts at
boot and its description, through gohai. A service that osapi deployed shows by
the name it was deployed with, not its unit's osapi- name, and is managed. On a
Mac, it lists each daemon that launchd has loaded in the system domain, by its
label. A daemon that osapi deployed reads by the name it was deployed with, app
for io.osapi.app.

## Usage

```go
step := o.ServiceList("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.ServiceEntry` row. A list answers with one row
for each item.

```go
var row osapi.ServiceEntry
err := report.Decode("list-service", &row)
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
`GET /api/node/{hostname}/service`.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go).
