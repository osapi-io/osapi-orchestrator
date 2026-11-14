# Service management

Manage systemd services on target hosts. You can list and examine services, and
create and update unit files. You can also start, stop, restart, enable, and
disable services.

## Operations

| Method                                           | Description                | Idempotent     |
| ------------------------------------------------ | -------------------------- | -------------- |
| [`ServiceList(target)`](list.md)                 | List services              | Read-only      |
| [`ServiceGet(target, name)`](get.md)             | Get a specific service     | Read-only      |
| [`ServiceCreate(target, opts)`](create.md)       | Create a service unit file | Non-idempotent |
| [`ServiceUpdate(target, name, opts)`](update.md) | Update a service unit file | Idempotent     |
| [`ServiceDelete(target, name)`](delete.md)       | Delete a service unit file | Idempotent     |
| [`ServiceStart(target, name)`](start.md)         | Start a service            | Idempotent     |
| [`ServiceStop(target, name)`](stop.md)           | Stop a service             | Idempotent     |
| [`ServiceRestart(target, name)`](restart.md)     | Restart a service          | Non-idempotent |
| [`ServiceEnable(target, name)`](enable.md)       | Enable a service on boot   | Idempotent     |
| [`ServiceDisable(target, name)`](disable.md)     | Disable a service on boot  | Idempotent     |

## Permissions

| Operation        | Permission      |
| ---------------- | --------------- |
| Read operations  | `service:read`  |
| Write operations | `service:write` |

## Example

Refer to
[`examples/operations/service.go`](../../../../examples/operations/service.go)
for a complete workflow example that uses all the operations.
