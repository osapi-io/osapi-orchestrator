# CertificateCreate

Add a certificate authority to the trust store.

Writes the PEM of the object as a file with the osapi- prefix, such as
osapi-internal-ca.crt for internal-ca, in the family's local directory:
/usr/local/share/ca-certificates on the Debian family, /etc/pki/trust/anchors on
SUSE, /etc/pki/ca-trust/source/anchors on RHEL and Fedora. Then it rebuilds the
trust store: with update-ca-certificates on the Debian family and SUSE, and with
update-ca-trust on the others. Alpine reports skipped. A CA that osapi deployed
already with the same content does not change and reports changed: false. A file
of that name that osapi did not write is refused. macOS reports skipped.

## Usage

```go
step := o.CertificateCreate("web-01", input)
```

## Parameters

| Parameter | Type                           | Description                                                   |
| --------- | ------------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                       | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.CertificateCreateInput` | Refer to `CertificateCreateInput` below.                      |

### CertificateCreateInput

| Field    | Type     | Description                                                                                       |
| -------- | -------- | ------------------------------------------------------------------------------------------------- |
| `Name`   | `string` | The name of the CA, such as internal-ca, which osapi writes as osapi-internal-ca.crt.             |
| `Object` | `string` | The object of the Object Store that has the PEM of the CA, as osapi client file upload stored it. |

## Result

Each host answers with one `osapi.CertificateResult` row. A list answers with
one row for each item.

```go
var row osapi.CertificateResult
err := report.Decode("create-certificate", &row)
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
`POST /api/node/{hostname}/certificate`.

## Example

Refer to
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go).
