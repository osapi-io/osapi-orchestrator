# UserUpdate

Change a user.

Changes what the input gives and the user does not have already, with usermod:
its shell, its home, which usermod moves the files to, its groups, which replace
the groups it has, and whether it is locked. The same values report changed:
false. A user that the host does not have is not found. It runs on the Debian
family only, and elsewhere reports skipped.

## Usage

```go
step := o.UserUpdate("web-01", input)
```

## Parameters

| Parameter | Type                    | Description                                                   |
| --------- | ----------------------- | ------------------------------------------------------------- |
| `target`  | `string`                | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.UserUpdateInput` | Refer to `UserUpdateInput` below.                             |

### UserUpdateInput

| Field    | Type         | Description                                                                                                      |
| -------- | ------------ | ---------------------------------------------------------------------------------------------------------------- |
| `Name`   | `string`     | The name of the path: the user to remove, with its home directory, such as alice. It names the item in the path. |
| `Groups` | `[]string`   | The groups it is a member of, which replace the groups it has.                                                   |
| `Home`   | `string`     | Its home directory, which usermod moves the files to.                                                            |
| `Lock`   | `types.Lock` | Lock is whether a user account refuses logins with its password.                                                 |
| `Shell`  | `string`     | Its login shell, such as /bin/bash.                                                                              |

## Result

Each host answers with one `osapi.UserResult` row. A list answers with one row
for each item.

```go
var row osapi.UserResult
err := report.Decode("update-user", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The user.                                                                                                                                                                                                                |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `user:write` permission. The step sends
`PUT /api/node/{hostname}/user/{name}`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
