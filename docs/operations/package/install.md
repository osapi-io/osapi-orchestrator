# PackageInstall

Installs a package on the target node.

## Usage

```go
step := o.PackageInstall("web-01", "nginx")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `pkgName` | `string` | Name of the package to install.                           |

## Result Type

```go
var result osapi.PackageMutationResult
err := results.Decode("install-package-1", &result)
```

| Field     | Type     | Description                                        |
| --------- | -------- | -------------------------------------------------- |
| `Name`    | `string` | Name of the package.                               |
| `Changed` | `bool`   | Whether the package was installed.                 |
| `Error`   | `string` | Error message if install failed; empty on success. |

## Idempotency

**Idempotent.** If the package is already installed, returns `Changed: false`.

## Permissions

Requires `package:write` permission.

## Example

See
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete working example.
