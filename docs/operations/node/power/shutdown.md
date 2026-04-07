# PowerShutdown

Initiates a shutdown on the target node.

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
| `opts`    | `PowerOpts` | Power options (see below).                                |

### PowerOpts

| Field     | Type     | Required | Description                               |
| --------- | -------- | -------- | ----------------------------------------- |
| `Delay`   | `int`    | No       | Seconds to wait before shutting down.     |
| `Message` | `string` | No       | Message to broadcast before the shutdown. |

## Result Type

```go
var result osapi.PowerResult
err := results.Decode("shutdown-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Action`  | `string` | Action taken (e.g., `shutdown`).                    |
| `Delay`   | `int`    | Delay in seconds before execution.                  |
| `Changed` | `bool`   | Whether the shutdown was initiated.                 |
| `Error`   | `string` | Error message if shutdown failed; empty on success. |

## Idempotency

**Non-idempotent.** Always initiates a shutdown.

## Permissions

Requires `power:execute` permission.

## Example

See
[`examples/operations/power.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/power.go)
for a complete working example.
