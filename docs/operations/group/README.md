# Group Management

Manage local groups on target nodes -- list, get, create, update, and delete
groups.

## Operations

| Method                                                    | Description            | Idempotent     |
| --------------------------------------------------------- | ---------------------- | -------------- |
| [`GroupList(target)`](list.md)                            | List groups            | Read-only      |
| [`GroupGet(target, name)`](get.md)                        | Get a specific group   | Read-only      |
| [`GroupCreate(target, opts)`](create.md)                  | Create a group         | Non-idempotent |
| [`GroupUpdate(target, name, opts)`](update.md)            | Update a group         | Idempotent     |
| [`GroupDelete(target, name)`](delete.md)                  | Delete a group         | Idempotent     |

## Permissions

| Operation        | Permission   |
| ---------------- | ------------ |
| Read operations  | `user:read`  |
| Write operations | `user:write` |

## Example

See [`examples/operations/user.go`](../../examples/operations/user.go) for a
complete workflow example covering group operations.
