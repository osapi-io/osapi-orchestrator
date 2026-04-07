# CertificateUpdate

Updates a CA certificate in the system trust store.

## Usage

```go
step := o.CertificateUpdate("web-01", "internal-ca", osapi.CertificateUpdateOpts{
    Object: "internal-ca-v2.pem",
})
```

## Parameters

| Parameter  | Type                    | Description                                               |
| ---------- | ----------------------- | --------------------------------------------------------- |
| `target`   | `string`                | Target host: `_any`, `_all`, hostname, or label selector. |
| `certName` | `string`                | Name of the certificate to update.                        |
| `opts`     | `CertificateUpdateOpts` | Update options (see below).                               |

### CertificateUpdateOpts

| Field    | Type     | Required | Description                                  |
| -------- | -------- | -------- | -------------------------------------------- |
| `Object` | `string` | Yes      | New Object Store reference for the PEM file. |

## Result Type

```go
var result osapi.CertificateCAMutationResult
err := results.Decode("update-certificate-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Certificate name.                                 |
| `Changed` | `bool`   | Whether the certificate was actually modified.    |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current PEM content against the new content.
Returns `Changed: true` only if the certificate was actually modified.

## Permissions

Requires `certificate:write` permission.

## Example

See
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete working example.
