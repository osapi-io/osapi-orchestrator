# FileUpload

Uploads file content to the NATS Object Store via the OSAPI REST API. Returns
the object name that can be referenced in subsequent
[FileDeploy](file-deploy.md) steps. This is a convenience wrapper that uses
`TaskFunc` to call the file upload API directly.

## Usage

```go
step := o.FileUpload(
    "app-config.yaml",
    "raw",
    []byte("server:\n  port: 8080\n"),
)
```

To always upload regardless of content changes:

```go
step := o.FileUpload(
    "app-config.yaml",
    "raw",
    []byte("server:\n  port: 8080\n"),
    orchestrator.WithForce(),
)
```

## Parameters

| Parameter     | Type              | Description                                                  |
| ------------- | ----------------- | ------------------------------------------------------------ |
| `name`        | `string`          | Object name to store the file under.                         |
| `contentType` | `string`          | Content type: `"raw"` or `"template"`.                       |
| `data`        | `[]byte`          | Raw file content to upload.                                  |
| `opts`        | `...UploadOption` | Optional. Use `WithForce()` to bypass the SHA-256 pre-check. |

## Result Type

```go
var result orchestrator.FileUploadResult
err := results.Decode("upload-file", &result)
```

| Field         | Type     | Description                           |
| ------------- | -------- | ------------------------------------- |
| `Name`        | `string` | Object name in the Object Store.      |
| `SHA256`      | `string` | SHA-256 hash of the uploaded content. |
| `Size`        | `int`    | Size of the uploaded content (bytes). |
| `Changed`     | `bool`   | Whether the upload modified state.    |
| `ContentType` | `string` | Content type of the uploaded file.    |

## Idempotency

**Idempotent by default.** The SDK computes SHA-256 locally and compares against
the stored hash, skipping the upload when content is unchanged. Use
`WithForce()` to bypass this pre-check and always upload.

You can also use [FileChanged](file-changed.md) with `OnlyIfChanged` for
explicit pre-check control:

```go
check := o.FileChanged("config.yaml", data)
o.FileUpload("config.yaml", "raw", data).
    After(check).
    OnlyIfChanged()
```

## Permissions

Requires `file:write` permission.

## Example

See
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
for a complete working example.
