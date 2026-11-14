# Interface management

Manage the network interface configuration on target hosts through Netplan. You
can list, get, create, update, and delete interface configurations.

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

Refer to
[`examples/operations/interface.go`](../../../../examples/operations/interface.go)
for a complete workflow example that uses all the operations.
