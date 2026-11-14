# FileUndeploy

Removes a deployed file from the filesystem of the target agent. Use this
operation only for a file that [FileDeploy](deploy.md) deployed.

## Usage

```go
step := o.FileUndeploy("web-01", "/etc/myapp/config.yaml")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `path`    | `string` | The path of the deployed file to remove.                  |

## Result type

```go
var result osapi.FileUndeployResult
err := results.Decode("undeploy-file-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Changed` | `bool`   | `true` if the operation removed the file.                    |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If no file exists, the operation returns `Changed: false`. If
the file exists, the operation removes it and returns `Changed: true`.

## Permissions

The token must have the `file:write` permission.

## Example

Refer to
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
for a complete example.
