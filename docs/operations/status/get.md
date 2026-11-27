# StatusGet

Get the status summary of the host.

Runs the get of the hostname, uptime, os, load, memory and disk behaviors on the
host, and answers with each. A read that fails is in errors, as its name and its
error, and the other fields still answer.

## Usage

```go
step := o.StatusGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.StatusEntry` row. A list answers with one row
for each item.

```go
var row osapi.StatusEntry
err := report.Decode("get-status", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Disks`       | `[]Disk`          | Each local filesystem.                                                                                                                                                                                                   |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Errors`      | `[]string`        | Each read that failed, as name: error. The other fields still answer.                                                                                                                                                    |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `LoadAverage` | `Averages`        | LoadAverage is the property load_average.                                                                                                                                                                                |
| `Memory`      | `Usage`           | Memory is the property memory.                                                                                                                                                                                           |
| `Name`        | `string`          | The hostname of the host, as the hostname command prints it.                                                                                                                                                             |
| `OSInfo`      | `Release`         | OSInfo is the property os_info.                                                                                                                                                                                          |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Uptime`      | `Entry`           | Uptime is the property uptime.                                                                                                                                                                                           |

## Permissions

The token must have the `status:read` permission. The step sends
`GET /api/node/{hostname}/status`.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go).
