# Sysctl management

Manage kernel parameters on target hosts through sysctl drop-in files. Create,
update, list, get, and delete parameters. The file provider keeps the
parameters.

## Operations

| Method                                         | Description                           | Idempotent     |
| ---------------------------------------------- | ------------------------------------- | -------------- |
| [`SysctlList(target)`](list.md)                | List the managed sysctl parameters    | Read-only      |
| [`SysctlGet(target, key)`](get.md)             | Get one sysctl parameter              | Read-only      |
| [`SysctlCreate(target, opts)`](create.md)      | Create a new sysctl parameter         | Non-idempotent |
| [`SysctlUpdate(target, key, opts)`](update.md) | Update a sysctl parameter that exists | Idempotent     |
| [`SysctlDelete(target, key)`](delete.md)       | Delete a sysctl parameter             | Idempotent     |

## Permissions

| Operation        | Permission     |
| ---------------- | -------------- |
| Read operations  | `sysctl:read`  |
| Write operations | `sysctl:write` |

## Example

Refer to
[`examples/operations/sysctl.go`](../../../../examples/operations/sysctl.go) for
a complete workflow example that uses all the operations.
