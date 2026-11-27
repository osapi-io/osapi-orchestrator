# UserSSHKeyDelete

Remove an authorized SSH key of a user.

Removes each line of the key with the fingerprint from ~/.ssh/authorized_keys of
the user, and keeps the other lines as they are. A key that the file does not
have, a user with no file, or a user that the host does not have, changes
nothing. It runs on the Debian family only, and elsewhere reports skipped.

## Usage

```go
step := o.UserSSHKeyDelete("web-01", input)
```

## Parameters

| Parameter | Type                          | Description                                                   |
| --------- | ----------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                      | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.UserSSHKeyDeleteInput` | Refer to `UserSSHKeyDeleteInput` below.                       |

### UserSSHKeyDeleteInput

| Field         | Type     | Description                                                                   |
| ------------- | -------- | ----------------------------------------------------------------------------- |
| `Name`        | `string` | The name of the path: the user, such as alice. It names the item in the path. |
| `Fingerprint` | `string` | The SHA256 fingerprint of the key, as list shows it.                          |

## Result

Each host answers with one `osapi.UserSshkeyResult` row. A list answers with one
row for each item.

```go
var row osapi.UserSshkeyResult
err := report.Decode("delete-user-ssh-key", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`     | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Fingerprint` | `string`          | The SHA256 fingerprint of the key.                                                                                                                                                                                       |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `user:write` permission. The step sends
`DELETE /api/node/{hostname}/user/{name}/sshkey`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
