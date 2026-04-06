# UserAddKey

Adds an SSH authorized key for a user on the target node.

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
| `opts`     | `SSHKeyAddOpts` | Add key options (see below).                              |

### SSHKeyAddOpts

| Field | Type     | Required | Description               |
| ----- | -------- | -------- | ------------------------- |
| `Key` | `string` | Yes      | Full SSH public key line. |

## Result Type

```go
var result osapi.SSHKeyMutationResult
err := results.Decode("add-ssh-key-1", &result)
```

| Field     | Type     | Description                                    |
| --------- | -------- | ---------------------------------------------- |
| `Changed` | `bool`   | Whether the key was added.                     |
| `Error`   | `string` | Error message if add failed; empty on success. |

## Idempotency

**Non-idempotent.** Adding a key that already exists returns an error.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
