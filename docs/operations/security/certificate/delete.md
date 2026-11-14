# CertificateDelete

Deletes a CA certificate from the system trust store.

## Usage

```go
step := o.CertificateDelete("web-01", "internal-ca")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `certName` | `string` | The name of the certificate to delete.                    |

## Result type

```go
var result osapi.CertificateCAMutationResult
err := results.Decode("delete-certificate-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The certificate name.                                        |
| `Changed` | `bool`   | `true` if the operation removed the certificate.             |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the certificate does not exist, the operation returns
`Changed: false`. If the certificate exists, the operation removes it and
returns `Changed: true`.

## Permissions

The token must have the `certificate:write` permission.

## Example

Refer to
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete example.
