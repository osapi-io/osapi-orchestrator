# NodeFileDelete

Remove a file that osapi deployed.

Removes the file and keeps its record, marked undeployed, so status reports
unmanaged. A path that osapi did not deploy, or deleted already, does not change
and reports changed: false. The object stays in the object store.

## Usage

```go
step := o.NodeFileDelete("web-01", input)
```

## Parameters

| Parameter | Type                        | Description                                                   |
| --------- | --------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                    | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.NodeFileDeleteInput` | Refer to `NodeFileDeleteInput` below.                         |

### NodeFileDeleteInput

| Field  | Type     | Description                            |
| ------ | -------- | -------------------------------------- |
| `Path` | `string` | Absolute path of the file on the host. |

## Result

Each host answers with one `osapi.FileDeleteResult` row. A list answers with one
row for each item.

```go
var row osapi.FileDeleteResult
err := report.Decode("delete-node-file", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Path`     | `string`          | Absolute path of the file.                                                                                                                                                                                               |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `file:write` permission. The step sends
`DELETE /api/node/{hostname}/file`.

## Example

Refer to
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go).
