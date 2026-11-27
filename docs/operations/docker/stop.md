# DockerStop

Stop a container.

Runs docker stop, which waits timeout seconds before docker kills the container,
or docker's default of 10 when it is 0. A container that does not run reports
changed: false.

## Usage

```go
step := o.DockerStop("web-01", input)
```

## Parameters

| Parameter | Type                    | Description                                                   |
| --------- | ----------------------- | ------------------------------------------------------------- |
| `target`  | `string`                | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.DockerStopInput` | Refer to `DockerStopInput` below.                             |

### DockerStopInput

| Field     | Type     | Description                                                                      |
| --------- | -------- | -------------------------------------------------------------------------------- |
| `Name`    | `string` | The name of the path: the container's name or ID. It names the item in the path. |
| `Timeout` | `int`    | Seconds to wait before docker kills it; docker's default when 0.                 |

## Result

Each host answers with one `osapi.DockerResult` row. A list answers with one row
for each item.

```go
var row osapi.DockerResult
err := report.Decode("stop-docker", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `ID`       | `string`          | The container's full ID, for create.                                                                                                                                                                                     |
| `Name`     | `string`          | The container.                                                                                                                                                                                                           |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `docker:execute` permission. The step sends
`POST /api/node/{hostname}/docker/{name}/stop`.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go).
