# CertificateUpdate

Write a certificate authority that osapi deployed again.

Writes the CA from the object, or from the one that it was deployed from when
the object is empty, and rebuilds the trust store when the file changed. A CA
that osapi did not deploy is not found. macOS reports skipped.

## Usage

```go
step := o.CertificateUpdate("web-01", input)
```

## Parameters

| Parameter | Type                           | Description                                                   |
| --------- | ------------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                       | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.CertificateUpdateInput` | Refer to `CertificateUpdateInput` below.                      |

### CertificateUpdateInput

| Field    | Type     | Description                                                                                           |
| -------- | -------- | ----------------------------------------------------------------------------------------------------- |
| `Name`   | `string` | The name of the path: the CA that osapi deployed, such as internal-ca. It names the item in the path. |
| `Object` | `string` | The object that has the PEM; the one it was deployed from when empty.                                 |

## Result

Each host answers with one `osapi.CertificateResult` row. A list answers with
one row for each item.

```go
var row osapi.CertificateResult
err := report.Decode("update-certificate", &row)
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
`PUT /api/node/{hostname}/certificate/{name}`.

## Example

Refer to
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go).
