# PackageRemove

Removes a package from the target host.

## Usage

```go
step := o.PackageRemove("web-01", "nginx")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `pkgName` | `string` | Name of the package to remove.                            |

## Result type

```go
var result osapi.PackageMutationResult
err := results.Decode("remove-package-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Name of the package.                                         |
| `Changed` | `bool`   | `true` if the operation removed the package.                 |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the package is not installed, the operation returns
`Changed: false`. If the package is installed, the operation removes it and
returns `Changed: true`.

## Permissions

The token must have the `package:write` permission.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete example.
