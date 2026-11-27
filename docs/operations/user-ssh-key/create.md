# UserSSHKeyCreate

Add an authorized SSH key to a user.

Adds the key line to the end of ~/.ssh/authorized_keys of the user, and makes
the directory and the file if they are missing, with modes 700 and 600, owned by
the user. A key that the file has already, by its fingerprint, does not change
it. A line with no public key is refused. A user that the host does not have is
not found. It runs on the Debian family only, and elsewhere reports skipped.

## Usage

```go
step := o.UserSSHKeyCreate("web-01", input)
```

## Parameters

| Parameter | Type                          | Description                                                   |
| --------- | ----------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                      | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.UserSSHKeyCreateInput` | Refer to `UserSSHKeyCreateInput` below.                       |

### UserSSHKeyCreateInput

| Field  | Type     | Description                                                                                   |
| ------ | -------- | --------------------------------------------------------------------------------------------- |
| `Name` | `string` | The name of the path: the user, such as alice. It names the item in the path.                 |
| `Key`  | `string` | The public key line, such as ssh-ed25519 AAAA... alice@laptop, with options before it if any. |

## Result

Each host answers with one `osapi.UserSshkeyResult` row. A list answers with one
row for each item.

```go
var row osapi.UserSshkeyResult
err := report.Decode("create-user-ssh-key", &row)
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
`POST /api/node/{hostname}/user/{name}/sshkey`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
