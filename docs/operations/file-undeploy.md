# FileUndeploy

Removes a previously deployed file from the target agent's filesystem. The file
must have been deployed via [FileDeploy](file-deploy.md).

## Usage

```go
step := o.FileUndeploy("web-01", "/etc/myapp/config.yaml")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `path`    | `string` | Path of the deployed file to remove.                      |

## Result Type

```go
var result osapi.FileUndeployResult
err := results.Decode("undeploy-file-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Changed` | `bool`   | Whether the file was removed or already absent.     |
| `Error`   | `string` | Error message if removal failed; empty on success.  |

## Idempotency

**Idempotent.** If the file does not exist, the operation returns
`Changed: false`. If the file exists and is removed, returns `Changed: true`.

## Permissions

Requires `file:write` permission.

## Example

```go
plan := o.Plan("cleanup")
o.FileUndeploy("web-01", "/etc/myapp/config.yaml")
report := plan.Execute(ctx)
```
