# UserListKeys

Lists SSH authorized keys for a user on the target node.

## Usage

```go
step := o.UserListKeys("web-01", "deploy")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string` | Username whose SSH keys to list.                          |

## Result Type

```go
var result osapi.SSHKeyInfoResult
err := results.Decode("list-ssh-key-1", &result)
```

| Field      | Type           | Description                                      |
| ---------- | -------------- | ------------------------------------------------ |
| `Hostname` | `string`       | The node's hostname.                             |
| `Keys`     | `[]SSHKeyInfo` | List of SSH authorized keys.                     |
| `Error`    | `string`       | Error message if query failed; empty on success. |

### SSHKeyInfo

| Field         | Type     | Description            |
| ------------- | -------- | ---------------------- |
| `Type`        | `string` | Key type (e.g., `ssh-ed25519`). |
| `Fingerprint` | `string` | Key fingerprint.       |
| `Comment`     | `string` | Key comment.           |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `user:read` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
