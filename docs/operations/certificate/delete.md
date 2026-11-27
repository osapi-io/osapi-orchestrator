# CertificateDelete

Remove a certificate authority that osapi deployed.

Removes the file of the CA and rebuilds the trust store: with
update-ca-certificates --fresh on the Debian family and SUSE, which removes the
links of the CA from /etc/ssl/certs, and update-ca-trust on the RHEL family. A
CA that osapi did not deploy does not change and reports changed: false: a
system CA is not osapi's to remove. macOS reports skipped.

## Usage

```go
step := o.CertificateDelete("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.CertificateResult` row. A list answers with
one row for each item.

```go
var row osapi.CertificateResult
err := report.Decode("delete-certificate", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The CA.                                                                                                                                                                                                                  |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `certificate:write` permission. The step sends
`DELETE /api/node/{hostname}/certificate/{name}`.

## Example

Refer to
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go).
