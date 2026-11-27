# DockerDelete

Remove a container.

Runs docker rm -f, which also removes a container that runs: a delete leaves no
container of the name. A container that the host does not have reports changed:
false.

## Usage

```go
step := o.DockerDelete("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.DockerResult` row. A list answers with one row
for each item.

```go
var row osapi.DockerResult
err := report.Decode("delete-docker", &row)
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

The token must have the `docker:write` permission. The step sends
`DELETE /api/node/{hostname}/docker/{name}`.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go).
