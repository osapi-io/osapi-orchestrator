# UserGet

Gets a specific user account by username.

## Usage

```go
step := o.UserGet("web-01", "deploy")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string` | The username to get.                                      |

## Result type

```go
var result osapi.UserInfoResult
err := results.Decode("get-user-1", &result)
```

| Field      | Type         | Description                                                  |
| ---------- | ------------ | ------------------------------------------------------------ |
| `Hostname` | `string`     | The hostname of the agent that ran the operation.            |
| `Users`    | `[]UserInfo` | The user details (one entry).                                |
| `Error`    | `string`     | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `user:read` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
