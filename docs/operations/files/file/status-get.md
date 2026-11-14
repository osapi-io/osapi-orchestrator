# FileStatusGet

Checks the status of a deployed file on the target host. The operation compares
the file on disk with the expected state that the agent records. It reports if
the file is in sync, has drifted, or is missing.

## Usage

```go
step := o.FileStatusGet("web-01", "/etc/myapp/config.yaml")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `path`    | `string` | The absolute path of the file to check.                   |

## Result type

```go
var result osapi.FileStatusResult
err := results.Decode("file.status.get-1", &result)
```

| Field    | Type     | Description                                        |
| -------- | -------- | -------------------------------------------------- |
| `Path`   | `string` | The path of the checked file.                      |
| `Status` | `string` | One of `"in-sync"`, `"drifted"`, or `"missing"`.   |
| `SHA256` | `string` | The SHA-256 checksum of the file on disk (empty if |
|          |          | the file is missing).                              |

### Status values

| Status    | Meaning                                                         |
| --------- | --------------------------------------------------------------- |
| `in-sync` | The file on disk is the same as the expected content.           |
| `drifted` | The file exists, but its checksum is not the expected checksum. |
| `missing` | The file does not exist at the path.                            |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `file:read` permission.

## Example

Refer to
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
for a complete example.
