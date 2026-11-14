# Certificate management

Manage CA certificates in the system trust store on target hosts. You can list,
create, update, and delete certificates. The file provider manages the
certificate files.

## Operations

| Method                                               | Description             | Idempotent     |
| ---------------------------------------------------- | ----------------------- | -------------- |
| [`CertificateList(target)`](list.md)                 | List CA certificates    | Read-only      |
| [`CertificateCreate(target, opts)`](create.md)       | Create a CA certificate | Non-idempotent |
| [`CertificateUpdate(target, name, opts)`](update.md) | Update a CA certificate | Idempotent     |
| [`CertificateDelete(target, name)`](delete.md)       | Delete a CA certificate | Idempotent     |

## Permissions

| Operation        | Permission          |
| ---------------- | ------------------- |
| Read operations  | `certificate:read`  |
| Write operations | `certificate:write` |

## Example

Refer to
[`examples/operations/certificate.go`](../../../../examples/operations/certificate.go)
for a complete workflow example that uses all the operations.
