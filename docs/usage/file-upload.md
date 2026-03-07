# FileUpload

Uploads file content to the NATS Object Store via the OSAPI REST API. Returns
the object name that can be referenced in subsequent
[FileDeploy](file-deploy.md) steps. This is a convenience wrapper that uses
`TaskFunc` to call the file upload API directly. Uses `WithForce` to always
upload regardless of content changes.

## Usage

```go
step := o.FileUpload(
    "app-config.yaml",
    "raw",
    []byte("server:\n  port: 8080\n"),
)
```

## Parameters

| Parameter     | Type     | Description                                        |
| ------------- | -------- | -------------------------------------------------- |
| `name`        | `string` | Object name to store the file under.               |
| `contentType` | `string` | Content type: `"raw"` or `"template"`.             |
| `data`        | `[]byte` | Raw file content to upload.                        |

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

**Not idempotent.** Always uploads with force and returns the server response.
Use [FileChanged](file-changed.md) with `OnlyIfChanged` to skip uploads when
content is unchanged:

```go
check := o.FileChanged("config.yaml", data)
o.FileUpload("config.yaml", "raw", data).
    After(check).
    OnlyIfChanged()
```

## Permissions

Requires `file:write` permission.
