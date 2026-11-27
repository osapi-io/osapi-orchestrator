# UserDelete

Remove a user.

Removes the user and its home directory with userdel -r. A user that the host
does not have does not change and reports changed: false. It runs on the Debian
family only, and elsewhere reports skipped.

## Usage

```go
step := o.UserDelete("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.UserResult` row. A list answers with one row
for each item.

```go
var row osapi.UserResult
err := report.Decode("delete-user", &row)
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
`DELETE /api/node/{hostname}/user/{name}`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
