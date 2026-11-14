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
| `opts`     | `CertificateUpdateOpts` | Update options (refer to the table below).                |

### CertificateUpdateOpts

| Field    | Type     | Required | Description                                  |
| -------- | -------- | -------- | -------------------------------------------- |
| `Object` | `string` | Yes      | New Object Store reference for the PEM file. |

## Result type

```go
var result osapi.CertificateCAMutationResult
err := results.Decode("update-certificate-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Certificate name.                                            |
| `Changed` | `bool`   | `true` if the operation changed the certificate.             |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current PEM content with the new
content. It returns `Changed: true` only if it changed the certificate.

## Permissions

The token must have the `certificate:write` permission.

## Example

Refer to
[`examples/operations/certificate.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/certificate.go)
for a complete example.
