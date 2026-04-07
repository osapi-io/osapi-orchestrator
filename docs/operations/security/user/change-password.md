# UserChangePassword

Changes a user's password on the target node.

## Usage

```go
step := o.UserChangePassword("web-01", "deploy", "newpassword123")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string` | Username whose password to change.                        |
| `password` | `string` | New password.                                             |

## Result Type

```go
var result osapi.UserMutationResult
err := results.Decode("change-password-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Username.                                         |
| `Changed` | `bool`   | Whether the password was changed.                 |
| `Error`   | `string` | Error message if change failed; empty on success. |

## Idempotency

**Non-idempotent.** Always sets the password regardless of current state.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
