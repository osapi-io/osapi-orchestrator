# CertificateCreate

Creates a CA certificate in the system trust store. Upload the PEM file to the
NATS Object Store first. Refer to [FileUpload](../../files/file/upload.md).

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
| `opts`    | `CertificateCreateOpts` | Create options (refer to the table below).                |

### CertificateCreateOpts

| Field    | Type     | Required | Description                              |
| -------- | -------- | -------- | ---------------------------------------- |
| `Name`   | `string` | Yes      | Certificate name.                        |
| `Object` | `string` | Yes      | Object Store reference for the PEM file. |

## Result type

```go
var result osapi.CertificateCAMutationResult
err := results.Decode("create-certificate-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Certificate name.                                            |
| `Changed` | `bool`   | `true` if the operation created the certificate.             |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the certificate already exists, the operation returns an
error. Use [CertificateUpdate](update.md) to replace a certificate that exists.

## Permissions

The token must have the `certificate:write` permission.

## Example

Refer to
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete example.
