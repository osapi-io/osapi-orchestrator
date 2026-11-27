# Object Store

Upload content to the controller's Object Store, and check whether content
differs from the stored object. To put an object on a host, use
[NodeFileUpdate](../../node-file/update.md).

## Operations

| Method                                             | Description                                          |
| -------------------------------------------------- | ---------------------------------------------------- |
| [`FileUpload(name, contentType, data)`](upload.md) | Upload content to the Object Store                   |
| [`FileChanged(name, data)`](changed.md)            | Check whether content differs from the stored object |

## Permissions

| Operation   | Permission   |
| ----------- | ------------ |
| FileChanged | `file:read`  |
| FileUpload  | `file:write` |

## Examples

Refer to
[`examples/operations/file-deploy.go`](../../../../examples/operations/file-deploy.go)
for a complete workflow that uploads, deploys, and verifies a file.

Refer to
[`examples/operations/file-changed.go`](../../../../examples/operations/file-changed.go)
for a workflow that skips the upload when the content is unchanged.
