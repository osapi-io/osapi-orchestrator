# CertificateList

Lists CA certificates in the system trust store on the target host.

## Usage

```go
step := o.CertificateList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.CertificateCAResult
err := results.Decode("list-certificate-1", &result)
```

| Field          | Type              | Description                                                  |
| -------------- | ----------------- | ------------------------------------------------------------ |
| `Hostname`     | `string`          | The hostname of the agent that ran the operation.            |
| `Certificates` | `[]CertificateCA` | List of CA certificate entries.                              |
| `Error`        | `string`          | The error message if the operation failed. Empty on success. |

### CertificateCA

| Field    | Type     | Description                                    |
| -------- | -------- | ---------------------------------------------- |
| `Name`   | `string` | The certificate name.                          |
| `Source` | `string` | The path to the certificate on the filesystem. |
| `Object` | `string` | The Object Store reference.                    |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `certificate:read` permission.

## Example

Refer to
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete example.
