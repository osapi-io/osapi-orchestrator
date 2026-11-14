# UserRemoveKey

Removes an SSH authorized key from a user on the target host.

## Usage

```go
step := o.UserRemoveKey("web-01", "deploy", "SHA256:abc123...")
```

## Parameters

| Parameter     | Type     | Description                                               |
| ------------- | -------- | --------------------------------------------------------- |
| `target`      | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `username`    | `string` | Username to remove the key from.                          |
| `fingerprint` | `string` | Fingerprint of the key to remove.                         |

## Result type

```go
var result osapi.SSHKeyMutationResult
err := results.Decode("remove-ssh-key-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Changed` | `bool`   | `true` if the operation removed the key.                     |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the key does not exist, the operation returns
`Changed: false`. If the key exists, the operation removes it and returns
`Changed: true`.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
