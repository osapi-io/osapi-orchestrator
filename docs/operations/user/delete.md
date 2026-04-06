# UserDelete

Deletes a user account from the target node.

## Usage

```go
step := o.UserDelete("web-01", "deploy")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string` | Username to delete.                                       |

## Result Type

```go
var result osapi.UserMutationResult
err := results.Decode("delete-user-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Username.                                           |
| `Changed` | `bool`   | Whether the user was removed or already absent.     |
| `Error`   | `string` | Error message if deletion failed; empty on success. |

## Idempotency

**Idempotent.** If the user does not exist, returns `Changed: false`.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
