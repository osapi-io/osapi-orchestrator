# PackageInstall

Installs a package on the target host.

## Usage

```go
step := o.PackageInstall("web-01", "nginx")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `pkgName` | `string` | The name of the package to install.                       |

## Result type

```go
var result osapi.PackageMutationResult
err := results.Decode("install-package-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The name of the package.                                     |
| `Changed` | `bool`   | `true` if the operation installed the package.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the package is already installed, the operation returns
`Changed: false`.

## Permissions

The token must have the `package:write` permission.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete example.
