# UserPasswordUpdate

Set the password of a user.

Sets the password with chpasswd -e, from its hash on standard input. The
controller hashes the password, so the job and NATS have only the hash. Each
call makes a new hash, so it always reports changed: true. A user that the host
does not have is not found. It runs on the Debian family only, and elsewhere
reports skipped.

## Usage

```go
step := o.UserPasswordUpdate("web-01", input)
```

## Parameters

| Parameter | Type                            | Description                                                   |
| --------- | ------------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                        | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.UserPasswordUpdateInput` | Refer to `UserPasswordUpdateInput` below.                     |

### UserPasswordUpdateInput

| Field      | Type     | Description                                                                                                          |
| ---------- | -------- | -------------------------------------------------------------------------------------------------------------------- |
| `Name`     | `string` | The name of the path: the user, such as alice. It names the item in the path.                                        |
| `Password` | `string` | The new password, which osapi hashes on the controller. A secret: a reference, secret://<vault>/<key>, or the value. |

## Result

Each host answers with one `osapi.UserPasswordResult` row. A list answers with
one row for each item.

```go
var row osapi.UserPasswordResult
err := report.Decode("update-user-password", &row)
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
`PUT /api/node/{hostname}/user/{name}/password`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
