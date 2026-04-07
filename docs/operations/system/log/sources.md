# LogSources

Lists available log sources (systemd units) on the target node.

## Usage

```go
step := o.LogSources("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.LogSourceResult
err := results.Decode("list-log-sources-1", &result)
```

| Field      | Type       | Description                                      |
| ---------- | ---------- | ------------------------------------------------ |
| `Hostname` | `string`   | The node's hostname.                             |
| `Sources`  | `[]string` | List of available log source names.              |
| `Error`    | `string`   | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `log:read` permission.

## Example

See
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go)
for a complete working example.
