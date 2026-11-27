# UserCreate

Add a user.

Adds the user with useradd, with its home directory, and sets its password with
chpasswd when one is given. The controller hashes the password, so the job and
NATS have only the hash. A user that the host has already does not change and
reports changed: false. It runs on the Debian family only, and elsewhere reports
skipped.

## Usage

```go
step := o.UserCreate("web-01", input)
```

## Parameters

| Parameter | Type                    | Description                                                   |
| --------- | ----------------------- | ------------------------------------------------------------- |
| `target`  | `string`                | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.UserCreateInput` | Refer to `UserCreateInput` below.                             |

### UserCreateInput

| Field      | Type       | Description                                                                                                                       |
| ---------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `GID`      | `int`      | Its primary group ID; a group of its own when it is 0.                                                                            |
| `Groups`   | `[]string` | The groups it is a member of, besides its primary group.                                                                          |
| `Home`     | `string`   | Its home directory; /home/<name> when empty.                                                                                      |
| `Name`     | `string`   | The user to add, such as alice.                                                                                                   |
| `Password` | `string`   | Its password, which osapi hashes on the controller; none when empty. A secret: a reference, secret://<vault>/<key>, or the value. |
| `Shell`    | `string`   | Its login shell, such as /bin/bash; the system's default when empty.                                                              |
| `System`   | `bool`     | Add it as a system user, with an ID from the system range.                                                                        |
| `UID`      | `int`      | Its numeric user ID; the system assigns one when it is 0.                                                                         |

## Result

Each host answers with one `osapi.UserResult` row. A list answers with one row
for each item.

```go
var row osapi.UserResult
err := report.Decode("create-user", &row)
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
`POST /api/node/{hostname}/user`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
