# Interface Management

Manage network interface configuration via Netplan on target nodes -- list, get,
create, update, and delete interface configurations.

## Operations

| Method                                             | Description                       | Idempotent     |
| -------------------------------------------------- | --------------------------------- | -------------- |
| [`InterfaceList(target)`](list.md)                 | List network interfaces           | Read-only      |
| [`InterfaceGet(target, name)`](get.md)             | Get a specific interface          | Read-only      |
| [`InterfaceCreate(target, name, opts)`](create.md) | Create an interface configuration | Non-idempotent |
| [`InterfaceUpdate(target, name, opts)`](update.md) | Update an interface configuration | Idempotent     |
| [`InterfaceDelete(target, name)`](delete.md)       | Delete an interface configuration | Idempotent     |

## Permissions

| Operation        | Permission      |
| ---------------- | --------------- |
| Read operations  | `network:read`  |
| Write operations | `network:write` |

## Example

See [`examples/operations/interface.go`](../../examples/operations/interface.go)
for a complete workflow example covering all operations.
