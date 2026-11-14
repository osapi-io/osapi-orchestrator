# PowerShutdown

Starts a shutdown on the target host.

## Usage

```go
step := o.PowerShutdown("web-01", osapi.PowerOpts{
    Delay:   60,
    Message: "System decommission",
})
```

## Parameters

| Parameter | Type        | Description                                               |
| --------- | ----------- | --------------------------------------------------------- |
| `target`  | `string`    | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `PowerOpts` | Power options (refer to the table below).                 |

### PowerOpts

| Field     | Type     | Required | Description                               |
| --------- | -------- | -------- | ----------------------------------------- |
| `Delay`   | `int`    | No       | Seconds to wait before the shutdown.      |
| `Message` | `string` | No       | Message to broadcast before the shutdown. |

## Result type

```go
var result osapi.PowerResult
err := results.Decode("shutdown-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Action`  | `string` | Action taken (for example, `shutdown`).                      |
| `Delay`   | `int`    | Delay in seconds before the shutdown starts.                 |
| `Changed` | `bool`   | `true` if the operation started the shutdown.                |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** The shutdown starts each time the step runs.

## Permissions

The token must have the `power:execute` permission.

## Example

Refer to
[`examples/operations/power.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/power.go)
for a complete example.
