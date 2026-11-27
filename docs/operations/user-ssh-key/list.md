# UserSSHKeyList

List the authorized SSH keys of a user.

Reads ~/.ssh/authorized_keys of the user through gohai, with the type, the
SHA256 fingerprint, the comment and the options of each key. A user with no file
has no keys. A user that the host does not have is not found.

## Usage

```go
step := o.UserSSHKeyList("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.UserSshkeyKey` row. A list answers with one
row for each item.

```go
var row osapi.UserSshkeyKey
err := report.Decode("list-user-ssh-key", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Comment`     | `string`          | The comment after the key, such as alice@laptop.                                                                                                                                                                         |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Fingerprint` | `string`          | The SHA256 fingerprint of the key.                                                                                                                                                                                       |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Options`     | `string`          | The options before the key, such as from="10.0.0.0/8".                                                                                                                                                                   |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Type`        | `string`          | The key type, such as ssh-ed25519.                                                                                                                                                                                       |

## Permissions

The token must have the `user:read` permission. The step sends
`GET /api/node/{hostname}/user/{name}/sshkey`.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go).
