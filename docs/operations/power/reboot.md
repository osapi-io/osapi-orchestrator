# PowerReboot

Reboot the host.

Runs shutdown -r -- +N, with N the delay in whole minutes and at least one, so
the agent can answer before the host goes down. A message goes to the logged-in
users. Each run changes the host.

## Usage

```go
step := o.PowerReboot("web-01", input)
```

## Parameters

| Parameter | Type                     | Description                                                   |
| --------- | ------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                 | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.PowerRebootInput` | Refer to `PowerRebootInput` below.                            |

### PowerRebootInput

| Field     | Type     | Description                                                                  |
| --------- | -------- | ---------------------------------------------------------------------------- |
| `Delay`   | `int`    | Seconds before the host goes down. It rounds to whole minutes, at least one. |
| `Message` | `string` | A message to the logged-in users first.                                      |

## Result

Each host answers with one `osapi.PowerResult` row. A list answers with one row
for each item.

```go
var row osapi.PowerResult
err := report.Decode("reboot-power", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Action`   | `string`          | reboot or shutdown.                                                                                                                                                                                                      |
| `Changed`  | `bool`            | Always true: each run changes the host.                                                                                                                                                                                  |
| `Delay`    | `int`             | Seconds before the host goes down.                                                                                                                                                                                       |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Retries

A second run is not the same as one run. So `Retry` does not run this step again
unless the step adds `RetryUnsafe`.

## Permissions

The token must have the `power:execute` permission. The step sends
`POST /api/node/{hostname}/power/reboot`.

## Example

Refer to
[`examples/operations/power.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/power.go).
