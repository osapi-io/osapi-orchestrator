# PowerReboot

Starts a reboot on the target host.

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
| `opts`    | `PowerOpts` | The power options (refer to the table below).             |

### PowerOpts

| Field     | Type     | Required | Description                                      |
| --------- | -------- | -------- | ------------------------------------------------ |
| `Delay`   | `int`    | No       | The number of seconds to wait before the reboot. |
| `Message` | `string` | No       | The message to broadcast before the reboot.      |

## Result type

```go
var result osapi.PowerResult
err := results.Decode("reboot-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Action`  | `string` | The action that the agent did (for example, `reboot`).       |
| `Delay`   | `int`    | The delay in seconds before the reboot starts.               |
| `Changed` | `bool`   | `true` if the operation started the reboot.                  |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** The reboot starts each time the step runs.

## Permissions

The token must have the `power:execute` permission.

## Example

Refer to
[`examples/operations/power.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/power.go)
for a complete example.
