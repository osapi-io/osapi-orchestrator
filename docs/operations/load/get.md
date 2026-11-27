# LoadGet

Get the load averages of the host.

Reads the averages of the run queue over 1, 5 and 15 minutes.

## Usage

```go
step := o.LoadGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.LoadEntry` row. A list answers with one row
for each item.

```go
var row osapi.LoadEntry
err := report.Decode("get-load", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `LoadAverage` | `Averages`        | LoadAverage is the property load_average.                                                                                                                                                                                |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `load:read` permission. The step sends
`GET /api/node/{hostname}/load`.

## Example

Refer to
[`examples/features/fact-predicates.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/fact-predicates.go).
