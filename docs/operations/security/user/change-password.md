# UserChangePassword

Changes a user's password on the target host.

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

## Result type

```go
var result osapi.UserMutationResult
err := results.Decode("change-password-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Username.                                                    |
| `Changed` | `bool`   | `true` if the operation changed the password.                |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** The operation sets the password each time the step runs. It
does not examine the current password first.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
