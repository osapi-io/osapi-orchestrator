# UserList

List the users of the host.

Lists each user with its IDs, home, shell, groups and whether it is locked,
through gohai: from /etc/passwd, /etc/group and /etc/shadow on Linux, and from
the directory service with dscl on macOS. The groups are those that list the
user as a member, besides its primary group.

## Usage

```go
step := o.UserList("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.UserAccount` row. A list answers with one row
for each item.

```go
var row osapi.UserAccount
err := report.Decode("list-user", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `GID`      | `int`             | Its primary group ID.                                                                                                                                                                                                    |
| `Groups`   | `[]string`        | The groups it is a member of, besides its primary group.                                                                                                                                                                 |
| `Home`     | `string`          | Its home directory.                                                                                                                                                                                                      |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Locked`   | `bool`            | Whether a login with its password is refused.                                                                                                                                                                            |
| `Name`     | `string`          | The user.                                                                                                                                                                                                                |
| `Shell`    | `string`          | Its login shell.                                                                                                                                                                                                         |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `UID`      | `int`             | Its numeric user ID.                                                                                                                                                                                                     |

## Permissions

The token must have the `user:read` permission. The step sends
`GET /api/node/{hostname}/user`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
