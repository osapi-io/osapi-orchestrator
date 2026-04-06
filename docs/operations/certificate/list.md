# CertificateList

Lists CA certificates in the system trust store on the target node.

## Usage

```go
step := o.CertificateList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.CertificateCAResult
err := results.Decode("list-certificate-1", &result)
```

| Field          | Type              | Description                                      |
| -------------- | ----------------- | ------------------------------------------------ |
| `Hostname`     | `string`          | The node's hostname.                             |
| `Certificates` | `[]CertificateCA` | List of CA certificate entries.                  |
| `Error`        | `string`          | Error message if query failed; empty on success. |

### CertificateCA

| Field    | Type     | Description                             |
| -------- | -------- | --------------------------------------- |
| `Name`   | `string` | Certificate name.                       |
| `Source` | `string` | Path to the certificate on the filesystem. |
| `Object` | `string` | Object Store reference.                 |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `certificate:read` permission.

## Example

See
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete working example.
