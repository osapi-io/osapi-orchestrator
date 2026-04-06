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
| `certName` | `string` | Name of the certificate to delete.                        |

## Result Type

```go
var result osapi.CertificateCAMutationResult
err := results.Decode("delete-certificate-1", &result)
```

| Field     | Type     | Description                                              |
| --------- | -------- | -------------------------------------------------------- |
| `Name`    | `string` | Certificate name.                                        |
| `Changed` | `bool`   | Whether the certificate was removed or already absent.   |
| `Error`   | `string` | Error message if deletion failed; empty on success.      |

## Idempotency

**Idempotent.** If the certificate does not exist, returns `Changed: false`.

## Permissions

Requires `certificate:write` permission.

## Example

See
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete working example.
