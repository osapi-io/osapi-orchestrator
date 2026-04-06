# Sysctl Management

Manage kernel parameters via sysctl drop-in files on target nodes -- create,
update, list, get, and delete parameters backed by the file provider.

## Operations

| Method                                         | Description                         | Idempotent     |
| ---------------------------------------------- | ----------------------------------- | -------------- |
| [`SysctlList(target)`](list.md)                | List managed sysctl parameters      | Read-only      |
| [`SysctlGet(target, key)`](get.md)             | Get a specific sysctl parameter     | Read-only      |
| [`SysctlCreate(target, opts)`](create.md)      | Create a new sysctl parameter       | Non-idempotent |
| [`SysctlUpdate(target, key, opts)`](update.md) | Update an existing sysctl parameter | Idempotent     |
| [`SysctlDelete(target, key)`](delete.md)       | Delete a sysctl parameter           | Idempotent     |

## Permissions

| Operation        | Permission     |
| ---------------- | -------------- |
| Read operations  | `sysctl:read`  |
| Write operations | `sysctl:write` |

## Example

See [`examples/operations/sysctl.go`](../../examples/operations/sysctl.go) for a
complete workflow example covering all operations.
