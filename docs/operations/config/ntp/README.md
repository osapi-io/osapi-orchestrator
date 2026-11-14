# NTP management

Manage the NTP server configuration on target hosts. Get the status, and create,
update, and delete the NTP configuration.

## Operations

| Method                                 | Description              | Idempotent     |
| -------------------------------------- | ------------------------ | -------------- |
| [`NTPGet(target)`](get.md)             | Get the NTP status       | Read-only      |
| [`NTPCreate(target, opts)`](create.md) | Create NTP configuration | Non-idempotent |
| [`NTPUpdate(target, opts)`](update.md) | Update NTP configuration | Idempotent     |
| [`NTPDelete(target)`](delete.md)       | Delete NTP configuration | Idempotent     |

## Permissions

| Operation        | Permission  |
| ---------------- | ----------- |
| Read operations  | `ntp:read`  |
| Write operations | `ntp:write` |

## Example

Refer to [`examples/operations/ntp.go`](../../../../examples/operations/ntp.go)
for a complete workflow example that uses all the operations.
