# ServiceCreate

Deploy a service.

Gets the object from the object store and writes it as the unit file of the
service, such as /etc/systemd/system/osapi-app.service for app, then reloads
systemd. A unit file at that path that osapi did not write, such as the agent's
own osapi-agent.service, is refused. A unit file that already has the content
does not change. The service does not start: start and enable are their own
actions. On a Mac, the object is the plist of the daemon,
/Library/LaunchDaemons/io.osapi.app.plist for app, with the label io.osapi.app,
and launchd does not load it until start.

## Usage

```go
step := o.ServiceCreate("web-01", input)
```

## Parameters

| Parameter | Type                       | Description                                                   |
| --------- | -------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                   | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ServiceCreateInput` | Refer to `ServiceCreateInput` below.                          |

### ServiceCreateInput

| Field    | Type     | Description                                                         |
| -------- | -------- | ------------------------------------------------------------------- |
| `Name`   | `string` | The service, such as app for osapi-app.service.                     |
| `Object` | `string` | The object that holds the unit file, from osapi client file upload. |

## Result

Each host answers with one `osapi.ServiceResult` row. A list answers with one
row for each item.

```go
var row osapi.ServiceResult
err := report.Decode("create-service", &row)
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
`POST /api/node/{hostname}/service`.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go).
