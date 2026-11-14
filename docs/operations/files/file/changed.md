# FileChanged

Checks if the local content is different from the version in the Object Store.
The operation calculates the SHA-256 hash locally and compares it with the
stored hash. Use it with `OnlyIfChanged` to skip an upload or a deploy when the
content did not change.

## Usage

```go
check := o.FileChanged("config.yaml", localContent)
o.FileUpload("config.yaml", "raw", localContent).
    After(check).
    OnlyIfChanged()
```

## Parameters

| Parameter | Type     | Description                                 |
| --------- | -------- | ------------------------------------------- |
| `name`    | `string` | The name of the object in the Object Store. |
| `data`    | `[]byte` | The local file content to compare.          |

## Result type

```go
var result osapi.FileChanged
err := results.Decode("check-file", &result)
```

| Field     | Type     | Description                                                       |
| --------- | -------- | ----------------------------------------------------------------- |
| `Name`    | `string` | The name of the object that the operation checked.                |
| `Changed` | `bool`   | `true` if the local content is different from the stored content. |
| `SHA256`  | `string` | The SHA-256 hash of the local content.                            |

## Idempotency

**Read-only.** The operation does not change the host or the Object Store. It
returns `Changed: true` if the file is not in the Object Store, or if the
SHA-256 hash is different.

## Permissions

The token must have the `file:read` permission.

## Example

Refer to
[`examples/operations/file-changed.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-changed.go)
for a complete example.
