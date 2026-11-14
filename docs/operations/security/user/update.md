# UserUpdate

Updates a user account on the target host. The operation changes only the fields
that you set in the options.

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
| `opts`     | `UserUpdateOpts` | Update options (refer to the table below).                |

### UserUpdateOpts

| Field    | Type       | Required | Description                                                       |
| -------- | ---------- | -------- | ----------------------------------------------------------------- |
| `Shell`  | `string`   | No       | New login shell path.                                             |
| `Home`   | `string`   | No       | New home directory path.                                          |
| `Groups` | `[]string` | No       | Supplementary group names. This list replaces the current groups. |
| `Lock`   | `*bool`    | No       | Lock or unlock the account.                                       |

## Result type

```go
var result osapi.UserMutationResult
err := results.Decode("update-user-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Username.                                                    |
| `Changed` | `bool`   | `true` if the operation changed the user.                    |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current account with the values you
give. It returns `Changed: true` only if it changed the account.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
