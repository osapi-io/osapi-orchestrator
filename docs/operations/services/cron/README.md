# Cron management

Manage cron drop-in files on target hosts. You can create, update, list, get,
and delete scheduled tasks. The file provider manages the drop-in files.

## Operations

| Method                                  | Description               | Idempotent     |
| --------------------------------------- | ------------------------- | -------------- |
| [`CronList(target)`](list.md)           | List all cron entries     | Read-only      |
| [`CronGet(target, name)`](get.md)       | Get a specific cron entry | Read-only      |
| [`CronCreate(target, opts)`](create.md) | Create a new cron entry   | Non-idempotent |
| [`CronUpdate(target, opts)`](update.md) | Update a cron entry       | Idempotent     |
| [`CronDelete(target, name)`](delete.md) | Delete a cron entry       | Idempotent     |

## Permissions

| Operation        | Permission   |
| ---------------- | ------------ |
| Read operations  | `cron:read`  |
| Write operations | `cron:write` |

## Example

Refer to
[`examples/operations/cron.go`](../../../../examples/operations/cron.go) for a
complete workflow example that uses all the operations.

The example shows how to:

- Create a cron entry with a schedule and with an interval
- List and examine entries
- Delete the entries at the end
