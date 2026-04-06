# Power Management

Reboot or shut down target nodes.

## Operations

| Method                                       | Description             | Idempotent     |
| -------------------------------------------- | ----------------------- | -------------- |
| [`PowerReboot(target, opts)`](reboot.md)     | Reboot the node         | Non-idempotent |
| [`PowerShutdown(target, opts)`](shutdown.md)  | Shut down the node      | Non-idempotent |

## Permissions

| Operation        | Permission      |
| ---------------- | --------------- |
| All operations   | `power:execute` |

## Example

See [`examples/operations/power.go`](../../examples/operations/power.go) for a
complete workflow example.
