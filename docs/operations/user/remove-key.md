# UserRemoveKey

Removes an SSH authorized key from a user on the target node.

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

## Result Type

```go
var result osapi.SSHKeyMutationResult
err := results.Decode("remove-ssh-key-1", &result)
```

| Field     | Type     | Description                                          |
| --------- | -------- | ---------------------------------------------------- |
| `Changed` | `bool`   | Whether the key was removed or already absent.       |
| `Error`   | `string` | Error message if removal failed; empty on success.   |

## Idempotency

**Idempotent.** If the key does not exist, returns `Changed: false`.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
