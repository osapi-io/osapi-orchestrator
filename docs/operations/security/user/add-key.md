# UserAddKey

Adds an SSH authorized key for a user on the target host.

## Usage

```go
step := o.UserAddKey("web-01", "deploy", osapi.SSHKeyAddOpts{
    Key: "ssh-ed25519 AAAA... user@host",
})
```

## Parameters

| Parameter  | Type            | Description                                               |
| ---------- | --------------- | --------------------------------------------------------- |
| `target`   | `string`        | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string`        | Username to add the key for.                              |
| `opts`     | `SSHKeyAddOpts` | Add key options (refer to the table below).               |

### SSHKeyAddOpts

| Field | Type     | Required | Description               |
| ----- | -------- | -------- | ------------------------- |
| `Key` | `string` | Yes      | Full SSH public key line. |

## Result type

```go
var result osapi.SSHKeyMutationResult
err := results.Decode("add-ssh-key-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Changed` | `bool`   | `true` if the operation added the key.                       |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the key already exists, the operation returns an error.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
