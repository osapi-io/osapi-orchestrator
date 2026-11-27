# MemoryGet

Get the memory usage of the host.

Reads the total, the free and the used memory, in bytes. Used is the total less
what is available, which includes the cache, as free(1) reports it.

## Usage

```go
step := o.MemoryGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.MemoryEntry` row. A list answers with one row
for each item.

```go
var row osapi.MemoryEntry
err := report.Decode("get-memory", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Memory`   | `Usage`           | Memory is the property memory.                                                                                                                                                                                           |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `memory:read` permission. The step sends
`GET /api/node/{hostname}/memory`.

## Example

Refer to
[`examples/features/parallel.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/parallel.go).
