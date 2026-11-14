# File management

Upload, deploy, and manage files on target hosts. The operations find drift with
SHA-256, render templates, and record the file state idempotently.

## Operations

| Method                                             | Description                                    | Idempotent |
| -------------------------------------------------- | ---------------------------------------------- | ---------- |
| [`FileUpload(name, contentType, data)`](upload.md) | Upload content to the Object Store             | Idempotent |
| [`FileDeploy(target, opts)`](deploy.md)            | Deploy a file from the Object Store to a host  | Idempotent |
| [`FileStatusGet(target, path)`](status-get.md)     | Get the status and SHA of a deployed file      | Read-only  |
| [`FileUndeploy(target, path)`](undeploy.md)        | Remove a deployed file                         | Idempotent |
| [`FileChanged(name, data)`](changed.md)            | Check if the Object Store content is different | Read-only  |

## Permissions

| Operation        | Permission   |
| ---------------- | ------------ |
| Read operations  | `file:read`  |
| Write operations | `file:write` |

## Examples

Refer to
[`examples/operations/file-deploy.go`](../../../../examples/operations/file-deploy.go)
for a complete workflow that uploads, deploys, and verifies a file.

Refer to
[`examples/operations/file-changed.go`](../../../../examples/operations/file-changed.go)
for an upload that runs only if the file drifted, with `OnlyIfChanged`.
