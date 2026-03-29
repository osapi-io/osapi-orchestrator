# Cron Management

Manage cron drop-in files on target nodes -- create, update, list, get,
and delete scheduled tasks backed by the file provider.

## Operations

| Method | Description | Idempotent |
| ------ | ----------- | ---------- |
| [`CronList(target)`](list.md) | List all cron entries | Read-only |
| [`CronGet(target, name)`](get.md) | Get a specific cron entry | Read-only |
| [`CronCreate(target, opts)`](create.md) | Create a new cron entry | Non-idempotent |
| [`CronUpdate(target, opts)`](update.md) | Update an existing cron entry | Idempotent |
| [`CronDelete(target, name)`](delete.md) | Delete a cron entry | Idempotent |

## Permissions

| Operation | Permission |
| --------- | ---------- |
| Read operations | `cron:read` |
| Write operations | `cron:write` |

## Example

See
[`examples/operations/cron.go`](../../examples/operations/cron.go)
for a complete workflow example covering all operations.

The example demonstrates:
- Creating a cron entry with schedule and interval variants
- Listing and inspecting entries
- Deleting entries for cleanup
