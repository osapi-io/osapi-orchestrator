# FileUpload

Uploads file content to the NATS Object Store through the OSAPI REST API. The
operation returns the object name. Later [FileDeploy](deploy.md) steps use this
name. The operation is a wrapper that uses `TaskFunc` to call the file upload
API directly.

## Usage

```go
step := o.FileUpload(
    "app-config.yaml",
    "raw",
    []byte("server:\n  port: 8080\n"),
)
```

To upload each time, also if the content did not change:

```go
step := o.FileUpload(
    "app-config.yaml",
    "raw",
    []byte("server:\n  port: 8080\n"),
    orchestrator.WithForce(),
)
```

## Parameters

| Parameter     | Type              | Description                                            |
| ------------- | ----------------- | ------------------------------------------------------ |
| `name`        | `string`          | The object name for the stored file.                   |
| `contentType` | `string`          | The content type: `"raw"` or `"template"`.             |
| `data`        | `[]byte`          | The raw file content to upload.                        |
| `opts`        | `...UploadOption` | Optional. Use `WithForce()` to skip the SHA-256 check. |

## Result type

```go
var result osapi.FileUpload
err := results.Decode("upload-file", &result)
```

| Field         | Type     | Description                                         |
| ------------- | -------- | --------------------------------------------------- |
| `Name`        | `string` | The name of the object in the Object Store.         |
| `SHA256`      | `string` | The SHA-256 hash of the uploaded content.           |
| `Size`        | `int`    | The size of the uploaded content (bytes).           |
| `Changed`     | `bool`   | `true` if the operation changed the stored content. |
| `ContentType` | `string` | The content type of the uploaded file.              |

## Idempotency

**Idempotent by default.** The SDK calculates the SHA-256 hash locally and
compares it with the stored hash. If the content did not change, the SDK does
not upload it. To skip this check and always upload, use `WithForce()`.

To control the check yourself, use [FileChanged](changed.md) with
`OnlyIfChanged`:

```go
check := o.FileChanged("config.yaml", data)
o.FileUpload("config.yaml", "raw", data).
    After(check).
    OnlyIfChanged()
```

## Permissions

The token must have the `file:write` permission.

## Example

Refer to
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
for a complete example.
