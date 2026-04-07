# CertificateCreate

Creates a CA certificate in the system trust store. The PEM file must be
uploaded to the NATS Object Store first (see [FileUpload](../file/upload.md)).

## Usage

```go
step := o.CertificateCreate("web-01", osapi.CertificateCreateOpts{
    Name:   "internal-ca",
    Object: "internal-ca.pem",
})
```

## Parameters

| Parameter | Type                    | Description                                               |
| --------- | ----------------------- | --------------------------------------------------------- |
| `target`  | `string`                | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `CertificateCreateOpts` | Create options (see below).                               |

### CertificateCreateOpts

| Field    | Type     | Required | Description                              |
| -------- | -------- | -------- | ---------------------------------------- |
| `Name`   | `string` | Yes      | Certificate name.                        |
| `Object` | `string` | Yes      | Object Store reference for the PEM file. |

## Result Type

```go
var result osapi.CertificateCAMutationResult
err := results.Decode("create-certificate-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Certificate name.                                   |
| `Changed` | `bool`   | Whether the certificate was created.                |
| `Error`   | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating a certificate that already exists returns an error.
Use [CertificateUpdate](update.md) to replace existing certificates.

## Permissions

Requires `certificate:write` permission.

## Example

See
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete working example.
