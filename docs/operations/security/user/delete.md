# UserDelete

Deletes a user account from the target host.

## Usage

```go
step := o.UserDelete("web-01", "deploy")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string` | The username to delete.                                   |

## Result type

```go
var result osapi.UserMutationResult
err := results.Decode("delete-user-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The username.                                                |
| `Changed` | `bool`   | `true` if the operation removed the user.                    |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the user does not exist, the operation returns
`Changed: false`. If the user exists, the operation removes it and returns
`Changed: true`.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
