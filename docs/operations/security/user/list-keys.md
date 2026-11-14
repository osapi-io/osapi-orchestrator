# UserListKeys

Lists SSH authorized keys for a user on the target host.

## Usage

```go
step := o.UserListKeys("web-01", "deploy")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string` | The username whose SSH keys to list.                      |

## Result type

```go
var result osapi.SSHKeyInfoResult
err := results.Decode("list-ssh-key-1", &result)
```

| Field      | Type           | Description                                                  |
| ---------- | -------------- | ------------------------------------------------------------ |
| `Hostname` | `string`       | The hostname of the agent that ran the operation.            |
| `Keys`     | `[]SSHKeyInfo` | List of SSH authorized keys.                                 |
| `Error`    | `string`       | The error message if the operation failed. Empty on success. |

### SSHKeyInfo

| Field         | Type     | Description                                |
| ------------- | -------- | ------------------------------------------ |
| `Type`        | `string` | The key type (for example, `ssh-ed25519`). |
| `Fingerprint` | `string` | The key fingerprint.                       |
| `Comment`     | `string` | The key comment.                           |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `user:read` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
