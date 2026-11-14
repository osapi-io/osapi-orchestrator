# LogSources

Lists available log sources (systemd units) on the target host.

## Usage

```go
step := o.LogSources("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.LogSourceResult
err := results.Decode("list-log-sources-1", &result)
```

| Field      | Type       | Description                                                  |
| ---------- | ---------- | ------------------------------------------------------------ |
| `Hostname` | `string`   | The hostname of the agent that ran the operation.            |
| `Sources`  | `[]string` | List of available log source names.                          |
| `Error`    | `string`   | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `log:read` permission.

## Example

Refer to
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go)
for a complete example.
