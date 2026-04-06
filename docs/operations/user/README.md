# User Management

Manage local user accounts and SSH authorized keys on target nodes -- list, get,
create, update, delete users, and manage their SSH keys and passwords.

## Operations

| Method                                                                    | Description                  | Idempotent     |
| ------------------------------------------------------------------------- | ---------------------------- | -------------- |
| [`UserList(target)`](list.md)                                            | List user accounts           | Read-only      |
| [`UserGet(target, username)`](get.md)                                    | Get a specific user          | Read-only      |
| [`UserCreate(target, opts)`](create.md)                                  | Create a user account        | Non-idempotent |
| [`UserUpdate(target, username, opts)`](update.md)                        | Update a user account        | Idempotent     |
| [`UserDelete(target, username)`](delete.md)                              | Delete a user account        | Idempotent     |
| [`UserListKeys(target, username)`](list-keys.md)                         | List SSH authorized keys     | Read-only      |
| [`UserAddKey(target, username, opts)`](add-key.md)                       | Add an SSH authorized key    | Non-idempotent |
| [`UserRemoveKey(target, username, fingerprint)`](remove-key.md)          | Remove an SSH authorized key | Idempotent     |
| [`UserChangePassword(target, username, password)`](change-password.md)   | Change a user's password     | Non-idempotent |

## Permissions

| Operation        | Permission   |
| ---------------- | ------------ |
| Read operations  | `user:read`  |
| Write operations | `user:write` |

## Example

See [`examples/operations/user.go`](../../examples/operations/user.go) for a
complete workflow example covering all operations.
