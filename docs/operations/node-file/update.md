# NodeFileUpdate

Deploy an object of the object store to a path.

Writes the object at the path with its mode and ownership, and keeps a record of
it. A template object, or content_type template, is rendered first with .Facts,
.Vars and .Hostname. The record keeps the SHA-256 of the object as it was read.
A file that holds the content, mode and owner already, from the same object,
does not change and reports changed: false. osapi client file stale list
compares the object now with that SHA-256, so it shows the deployments whose
object changed since.

## Usage

```go
step := o.NodeFileUpdate("web-01", input)
```

## Parameters

| Parameter | Type                        | Description                                                   |
| --------- | --------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                    | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.NodeFileUpdateInput` | Refer to `NodeFileUpdateInput` below.                         |

### NodeFileUpdateInput

| Field         | Type                | Description                                                                                                                                   |
| ------------- | ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `ContentType` | `types.ContentType` | ContentType is how a deploy treats a file: raw writes its bytes as they are; template renders it with Go text/template and the agent's facts. |
| `Group`       | `string`            | Group that owns the file; left as written when empty.                                                                                         |
| `Mode`        | `string`            | Permissions in octal, such as 0644; when empty, a file keeps its mode and a new file gets 0644.                                               |
| `ObjectName`  | `string`            | The object of the object store, from osapi client file list.                                                                                  |
| `Owner`       | `string`            | User that owns the file; left as written when empty.                                                                                          |
| `Path`        | `string`            | Absolute path of the file on the host.                                                                                                        |
| `Vars`        | `map[string]string` | Variables of a template, as .Vars, each as key=value.                                                                                         |

## Result

Each host answers with one `osapi.FileUpdateResult` row. A list answers with one
row for each item.

```go
var row osapi.FileUpdateResult
err := report.Decode("update-node-file", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Path`     | `string`          | Absolute path of the file.                                                                                                                                                                                               |
| `SHA256`   | `string`          | SHA-256 of the content that the file holds.                                                                                                                                                                              |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `file:write` permission. The step sends
`PUT /api/node/{hostname}/file`.

## Example

Refer to
[`examples/features/only-if-changed.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/only-if-changed.go).
