# Package Management

Manage system packages on target nodes -- list installed packages, get details,
install, remove, update all, and check for available updates.

## Operations

| Method                                              | Description                    | Idempotent     |
| --------------------------------------------------- | ------------------------------ | -------------- |
| [`PackageList(target)`](list.md)                    | List installed packages        | Read-only      |
| [`PackageGet(target, name)`](get.md)                | Get a specific package         | Read-only      |
| [`PackageInstall(target, name)`](install.md)        | Install a package              | Idempotent     |
| [`PackageRemove(target, name)`](remove.md)          | Remove a package               | Idempotent     |
| [`PackageUpdate(target)`](update.md)                | Update all packages            | Non-idempotent |
| [`PackageListUpdates(target)`](list-updates.md)     | List available package updates | Read-only      |

## Permissions

| Operation        | Permission      |
| ---------------- | --------------- |
| Read operations  | `package:read`  |
| Write operations | `package:write` |

## Example

See
[`examples/operations/package.go`](../../examples/operations/package.go) for a
complete workflow example covering all operations.
