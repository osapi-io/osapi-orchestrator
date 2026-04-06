# NTP Management

Manage NTP server configuration on target nodes -- get status, create, update,
and delete NTP configuration.

## Operations

| Method                                    | Description              | Idempotent     |
| ----------------------------------------- | ------------------------ | -------------- |
| [`NTPGet(target)`](get.md)                | Get NTP status           | Read-only      |
| [`NTPCreate(target, opts)`](create.md)    | Create NTP configuration | Non-idempotent |
| [`NTPUpdate(target, opts)`](update.md)    | Update NTP configuration | Idempotent     |
| [`NTPDelete(target)`](delete.md)          | Delete NTP configuration | Idempotent     |

## Permissions

| Operation        | Permission  |
| ---------------- | ----------- |
| Read operations  | `ntp:read`  |
| Write operations | `ntp:write` |

## Example

See [`examples/operations/ntp.go`](../../examples/operations/ntp.go) for a
complete workflow example covering all operations.
