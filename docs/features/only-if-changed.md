# OnlyIfChanged

Skip a step unless at least one of its dependencies reported a change
(`Changed: true` in the result).

## Usage

```go
disk := o.NodeDiskGet("_any").After(health)

o.CommandExec("_any", "df", "-h").
    Named("run-df").
    After(disk).
    OnlyIfChanged()
```

The `run-df` step only executes if the disk query reported changes. Read-only
operations always return `Changed: false`, so this pattern is most useful with
write operations or `FileChanged` checks.

## OnlyIfAllChanged

A stricter variant that requires **all** dependencies to report changes:

```go
step.OnlyIfAllChanged()
```

## Example

See
[`examples/only-if-changed.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/only-if-changed.go)
for a complete working example.
