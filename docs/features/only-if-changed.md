# OnlyIfChanged

Skip a step if none of its dependencies reported a change (`Changed: true` in
the result).

## Usage

```go
disk := o.DiskGet("_any").After(health)

o.CommandExec("_any", "df", "-h").
    Named("run-df").
    After(disk).
    OnlyIfChanged()
```

The `run-df` step runs only if the disk query reported a change. A read-only
operation always returns `Changed: false`. Thus use this pattern with write
operations or `FileChanged` checks.

## OnlyIfAllChanged

This variant is stricter. **All** dependencies must report a change:

```go
step.OnlyIfAllChanged()
```

## Host-level variants

For a broadcast operation, two host-level variants examine the `Changed` value
of each host. They do not use the task-level `Changed` bool:

```go
// Run if any host in the dependency reported changes.
step.OnlyIfAnyHostChanged()

// Run only if every host in the dependency reported changes.
step.OnlyIfAllHostsChanged()
```

Refer to [Guards: Broadcast Guards](guards.md#broadcast-guards) for the details.

## Example

Refer to
[`examples/features/only-if-changed.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/only-if-changed.go)
for a complete example.
