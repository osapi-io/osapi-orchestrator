# FileUpload

Uploads file content to the NATS Object Store via the OSAPI REST API. Returns
the object name that can be referenced in subsequent
[FileDeploy](file-deploy.md) steps. This is a convenience wrapper that uses
`TaskFunc` to call the file upload API directly.

> **Note:** Requires `osapi-sdk` FileService. This operation will return an
> error until the SDK file endpoints are merged.

## Usage

```go
step := o.FileUpload(
    "app-config.yaml",
    []byte("server:\n  port: 8080\n"),
)
```

## Parameters

| Parameter | Type     | Description                          |
| --------- | -------- | ------------------------------------ |
| `name`    | `string` | Object name to store the file under. |
| `data`    | `[]byte` | Raw file content to upload.          |

## Result Type

```go
var result orchestrator.FileUploadResult
err := results.Decode("upload-file", &result)
```

| Field  | Type     | Description                           |
| ------ | -------- | ------------------------------------- |
| `Name` | `string` | Object name in the Object Store.      |
| `Size` | `int`    | Size of the uploaded content (bytes). |

## Idempotency

**Not idempotent.** Overwrites any existing object with the same name. Always
returns `Changed: true`. Use `When` or `OnlyIfChanged` guards to control when
the upload runs:

```go
o.FileUpload("config.yaml", data).
    OnlyIfChanged()
```

## Permissions

Requires `file:write` permission.
