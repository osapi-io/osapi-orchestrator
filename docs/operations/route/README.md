# Route Management

Manage network route configuration via Netplan on target nodes -- list, get,
create, update, and delete routes.

## Operations

| Method                                          | Description                 | Idempotent     |
| ----------------------------------------------- | --------------------------- | -------------- |
| [`RouteList(target)`](list.md)                  | List network routes         | Read-only      |
| [`RouteGet(target, iface)`](get.md)             | Get routes for an interface | Read-only      |
| [`RouteCreate(target, iface, opts)`](create.md) | Create route configuration  | Non-idempotent |
| [`RouteUpdate(target, iface, opts)`](update.md) | Update route configuration  | Idempotent     |
| [`RouteDelete(target, iface)`](delete.md)       | Delete route configuration  | Idempotent     |

## Permissions

| Operation        | Permission      |
| ---------------- | --------------- |
| Read operations  | `network:read`  |
| Write operations | `network:write` |

## Example

See [`examples/operations/route.go`](../../examples/operations/route.go) for a
complete workflow example covering all operations.
