# UserUpdate

Updates a user account on the target node. Only the fields specified in the
options are changed.

## Usage

```go
step := o.UserUpdate("web-01", "deploy", osapi.UserUpdateOpts{
    Shell: "/bin/zsh",
})
```

## Parameters

| Parameter  | Type             | Description                                               |
| ---------- | ---------------- | --------------------------------------------------------- |
| `target`   | `string`         | Target host: `_any`, `_all`, hostname, or label selector. |
| `username` | `string`         | Username to update.                                       |
| `opts`     | `UserUpdateOpts` | Update options (see below).                               |

### UserUpdateOpts

| Field    | Type       | Required | Description                                       |
| -------- | ---------- | -------- | ------------------------------------------------- |
| `Shell`  | `string`   | No       | New login shell path.                             |
| `Home`   | `string`   | No       | New home directory path.                          |
| `Groups` | `[]string` | No       | Supplementary group names (replaces existing).    |
| `Lock`   | `*bool`    | No       | Lock or unlock the account.                       |

## Result Type

```go
var result osapi.UserMutationResult
err := results.Decode("update-user-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Username.                                         |
| `Changed` | `bool`   | Whether the user was actually modified.           |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current account against the desired state. Returns
`Changed: true` only if the account was actually modified.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
