# CertificateGet

Get one certificate authority.

Reads one CA from the same list, by the name that list gives it. A CA that the
host does not trust is not found. A name with a slash, such as a system CA of
the Debian family, is in list only.

## Usage

```go
step := o.CertificateGet("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.CertificateEntry` row. A list answers with one
row for each item.

```go
var row osapi.CertificateEntry
err := report.Decode("get-certificate", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Enabled`     | `bool`            | Whether the host trusts it: false for a system CA that the trust store turns off.                                                                                                                                        |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Fingerprint` | `string`          | The SHA-256 of the certificate, in hex.                                                                                                                                                                                  |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Managed`     | `bool`            | Whether osapi deployed it, and its writes change it.                                                                                                                                                                     |
| `Name`        | `string`          | The CA: the name it was deployed with for a CA that osapi deployed, else its file's name.                                                                                                                                |
| `NotAfter`    | `string`          | When it expires.                                                                                                                                                                                                         |
| `Object`      | `string`          | The object that osapi deployed it from.                                                                                                                                                                                  |
| `Path`        | `string`          | The file, or on macOS the keychain, that has it.                                                                                                                                                                         |
| `Source`      | `string`          | system for a CA of the distribution, local for one an administrator added.                                                                                                                                               |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Subject`     | `string`          | Its subject, such as CN=ISRG Root X1,O=Internet Security Research Group,C=US.                                                                                                                                            |

## Permissions

The token must have the `certificate:read` permission. The step sends
`GET /api/node/{hostname}/certificate/{name}`.
