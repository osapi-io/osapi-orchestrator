# PowerReboot

Initiates a reboot on the target node.

## Usage

```go
step := o.PowerReboot("web-01", osapi.PowerOpts{
    Delay:   30,
    Message: "Scheduled maintenance reboot",
})
```

## Parameters

| Parameter | Type        | Description                                               |
| --------- | ----------- | --------------------------------------------------------- |
| `target`  | `string`    | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `PowerOpts` | Power options (see below).                                |

### PowerOpts

| Field     | Type     | Required | Description                                       |
| --------- | -------- | -------- | ------------------------------------------------- |
| `Delay`   | `int`    | No       | Seconds to wait before rebooting.                 |
| `Message` | `string` | No       | Message to broadcast before the reboot.           |

## Result Type

```go
var result osapi.PowerResult
err := results.Decode("reboot-1", &result)
```

| Field     | Type     | Description                                        |
| --------- | -------- | -------------------------------------------------- |
| `Action`  | `string` | Action taken (e.g., `reboot`).                     |
| `Delay`   | `int`    | Delay in seconds before execution.                 |
| `Changed` | `bool`   | Whether the reboot was initiated.                  |
| `Error`   | `string` | Error message if reboot failed; empty on success.  |

## Idempotency

**Non-idempotent.** Always initiates a reboot.

## Permissions

Requires `power:execute` permission.

## Example

See
[`examples/operations/power.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/power.go)
for a complete working example.
