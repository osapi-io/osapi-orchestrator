# NodeFileGet

Show whether a deployed file holds what osapi wrote.

Compares the file on the host with the record that osapi keeps of what it
deployed: in-sync when they agree, drifted when something changed the file,
missing when osapi deployed the file and it is gone, and unmanaged when osapi
did not deploy the path or undeployed it. A path that is a directory is refused.

## Usage

```go
step := o.NodeFileGet("web-01", path)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `path`    | `string` | The path of the file on the host.                             |

## Result

Each host answers with one `osapi.FileStatus` row. A list answers with one row
for each item.

```go
var row osapi.FileStatus
err := report.Decode("get-node-file", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Path`     | `string`          | Absolute path of the file.                                                                                                                                                                                               |
| `SHA256`   | `string`          | SHA-256 of the file on the host, when it is there.                                                                                                                                                                       |
| `State`    | `string`          | in-sync, drifted, missing or unmanaged.                                                                                                                                                                                  |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `file:read` permission. The step sends
`GET /api/node/{hostname}/file`.

## Example

Refer to
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go).
