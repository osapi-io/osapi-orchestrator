# UserList

Lists user accounts on the target host.

## Usage

```go
step := o.UserList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.UserInfoResult
err := results.Decode("list-user-1", &result)
```

| Field      | Type         | Description                                                  |
| ---------- | ------------ | ------------------------------------------------------------ |
| `Hostname` | `string`     | The hostname of the agent that ran the operation.            |
| `Users`    | `[]UserInfo` | List of user accounts.                                       |
| `Error`    | `string`     | The error message if the operation failed. Empty on success. |

### UserInfo

| Field    | Type       | Description                      |
| -------- | ---------- | -------------------------------- |
| `Name`   | `string`   | The username.                    |
| `UID`    | `int`      | The numeric user ID.             |
| `GID`    | `int`      | The primary group ID.            |
| `Home`   | `string`   | The home directory path.         |
| `Shell`  | `string`   | The login shell path.            |
| `Groups` | `[]string` | The supplementary group names.   |
| `Locked` | `bool`     | `true` if the account is locked. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `user:read` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
