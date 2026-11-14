# Power management

Reboot or shut down target hosts.

## Operations

| Method                                       | Description        | Idempotent     |
| -------------------------------------------- | ------------------ | -------------- |
| [`PowerReboot(target, opts)`](reboot.md)     | Reboot the host    | Non-idempotent |
| [`PowerShutdown(target, opts)`](shutdown.md) | Shut down the host | Non-idempotent |

## Permissions

| Operation      | Permission      |
| -------------- | --------------- |
| All operations | `power:execute` |

## Example

Refer to
[`examples/operations/power.go`](../../../../examples/operations/power.go) for a
complete workflow example.
