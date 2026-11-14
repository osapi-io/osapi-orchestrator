# Route management

Manage the network route configuration on target hosts through Netplan. You can
list, get, create, update, and delete routes.

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

Refer to
[`examples/operations/route.go`](../../../../examples/operations/route.go) for a
complete workflow example that uses all the operations.
