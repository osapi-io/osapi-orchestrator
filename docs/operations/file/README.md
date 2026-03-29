# File Management

Upload, deploy, and manage files on target nodes -- with SHA-256 drift
detection, template rendering, and idempotent state tracking.

## Operations

| Method | Description | Idempotent |
| ------ | ----------- | ---------- |
| [`FileUpload(name, contentType, data)`](upload.md) | Upload content to the Object Store | Idempotent |
| [`FileDeploy(target, opts)`](deploy.md) | Deploy a file from Object Store to host | Idempotent |
| [`FileStatusGet(target, path)`](status-get.md) | Get deployed file status and SHA | Read-only |
| [`FileUndeploy(target, path)`](undeploy.md) | Remove a deployed file | Idempotent |
| [`FileChanged(name, data)`](changed.md) | Check if Object Store content differs | Read-only |

## Permissions

| Operation | Permission |
| --------- | ---------- |
| Read operations | `file:read` |
| Write operations | `file:write` |

## Examples

See
[`examples/operations/file-deploy.go`](../../examples/operations/file-deploy.go)
for a complete upload-deploy-verify workflow.

See
[`examples/operations/file-changed.go`](../../examples/operations/file-changed.go)
for conditional upload with drift detection using `OnlyIfChanged`.
