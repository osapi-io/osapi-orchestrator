# UserList

Lists user accounts on the target node.

## Usage

```go
step := o.UserList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.UserInfoResult
err := results.Decode("list-user-1", &result)
```

| Field      | Type         | Description                                      |
| ---------- | ------------ | ------------------------------------------------ |
| `Hostname` | `string`     | The node's hostname.                             |
| `Users`    | `[]UserInfo` | List of user accounts.                           |
| `Error`    | `string`     | Error message if query failed; empty on success. |

### UserInfo

| Field    | Type       | Description                    |
| -------- | ---------- | ------------------------------ |
| `Name`   | `string`   | Username.                      |
| `UID`    | `int`      | Numeric user ID.               |
| `GID`    | `int`      | Primary group ID.              |
| `Home`   | `string`   | Home directory path.           |
| `Shell`  | `string`   | Login shell path.              |
| `Groups` | `[]string` | Supplementary group names.     |
| `Locked` | `bool`     | Whether the account is locked. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `user:read` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
