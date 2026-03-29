# Broadcast

Target `_all` or a label selector to send a job to every matching agent. The
renderer automatically shows per-host results.

## Usage

```go
o.NodeHostnameGet("_all")
```

For label-based targeting:

```go
o.CommandExec("group:web", "uptime")
```

## Per-Host Results

Broadcast operations return a `HostResult` for each responding agent:

| Field      | Type             | Description                        |
| ---------- | ---------------- | ---------------------------------- |
| `Hostname` | `string`         | Agent hostname                     |
| `Status`   | `string`         | Host status (see below)            |
| `Changed`  | `bool`           | Whether this host reported changes |
| `Error`    | `string`         | Error message (empty on success)   |
| `Data`     | `map[string]any` | Host-specific response data        |

### Host Status Values

| Value     | Meaning                                  |
| --------- | ---------------------------------------- |
| `ok`      | Operation completed successfully         |
| `skipped` | Operation not supported on this host     |
| `failed`  | Operation failed with an error           |

The `Status` field distinguishes between hosts that failed (encountered an
error) and hosts that were skipped (operation unsupported). Guards like
`OnlyIfAnyHostFailed` inspect `Status`, not `Error`, so skipped hosts do not
trigger failure guards.

Access per-host results in a `When` guard or after execution:

```go
step.When(func(r orchestrator.Results) bool {
    for _, hr := range r.HostResults("deploy") {
        if hr.Error != "" {
            return true // trigger cleanup
        }
    }
    return false
})
```

## Partial Failure

When some hosts succeed and others fail in a broadcast operation, the task
result has:

- `Status = Failed` — the overall task is marked as failed
- `Changed = true` — at least one host reported a change
- `HostResults` — contains per-host details with individual `Error` and
  `Changed` fields

This allows downstream guards to distinguish between "all hosts failed" and
"some hosts failed":

```go
// Task-level: triggers on any failure (Status == Failed)
step.OnlyIfFailed()

// Host-level: triggers only if specific hosts failed
step.OnlyIfAnyHostFailed()
step.OnlyIfAllHostsFailed()
```

## Broadcast Guards

Four host-level guard methods inspect `HostResults` from broadcast dependencies.
See [Guards — Broadcast Guards](guards.md#broadcast-guards) for details.

## Example

See
[`examples/features/broadcast.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast.go)
and
[`examples/features/broadcast-guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast-guards.go)
for complete working examples.
