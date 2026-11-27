# Broadcast

Target `_all` or a label selector to send a job to each agent that matches. The
renderer shows the result of each host.

## Usage

```go
o.HostnameGet("_all")
```

To target by label:

```go
o.CommandExec("group:web", "uptime")
```

## Per-host results

A broadcast operation returns a `HostResult` for each agent that responds:

| Field      | Type             | Description                                                  |
| ---------- | ---------------- | ------------------------------------------------------------ |
| `Hostname` | `string`         | The hostname of the agent                                    |
| `Status`   | `string`         | The host status (refer to the table below)                   |
| `Changed`  | `bool`           | `true` if this host reported a change                        |
| `Error`    | `string`         | The error message if the operation failed. Empty on success. |
| `Data`     | `map[string]any` | The response data from this host                             |

`Data` is the host's row as the osapi SDK encodes it to JSON. The SDK leaves out
a field that is zero, so an `exit_code` of 0 or an empty `stdout` is not in the
map. To read each field, use `Decode` into the SDK's type:

```go
var cmd osapi.CommandResult
if err := hr.Decode(&cmd); err == nil {
    fmt.Println(cmd.ExitCode) // 0 is here
}
```

### Host status values

| Value     | Meaning                                 |
| --------- | --------------------------------------- |
| `ok`      | The operation completed with no error   |
| `skipped` | The host does not support the operation |
| `failed`  | The operation failed with an error      |
| `timeout` | The host did not answer in time         |

The `Status` field shows the difference between a host that failed and a host
that the agent skipped. A failed host had an error. A skipped host does not
support the operation. A host that timed out did not answer, so nobody knows
whether the operation ran there. The orchestrator counts it as failed, with the
error "the host did not answer".

A failure guard such as `OnlyIfAnyHostFailed` counts a host as failed if its
`Status` is `failed` or `timeout`. It also counts a host as failed if its
`Error` is not empty and its `Status` is not `skipped`. Thus a skipped host does
not trigger a failure guard, and a step whose hosts timed out fails.

Read the result of each host in a `When` guard or after the plan runs:

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

## Partial failure

If some hosts succeed and other hosts fail in a broadcast operation, the task
result has these values:

- `Status = Failed`. The orchestrator marks the full task as failed.
- `Changed = true`. One or more hosts reported a change.
- `HostResults`. This has the details of each host, with its own `Error` and
  `Changed` fields.

Thus a guard on a later step can tell "all hosts failed" from "some hosts
failed":

```go
// Task-level: triggers on any failure (Status == Failed)
step.OnlyIfFailed()

// Host-level: triggers only if specific hosts failed
step.OnlyIfAnyHostFailed()
step.OnlyIfAllHostsFailed()
```

## Broadcast guards

Four host-level guard methods examine the `HostResults` of broadcast
dependencies. Refer to [Guards: Broadcast Guards](guards.md#broadcast-guards)
for the details.

## Example

Refer to
[`examples/features/broadcast.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast.go)
and
[`examples/features/broadcast-guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast-guards.go)
for complete examples.
