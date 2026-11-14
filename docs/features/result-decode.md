# Result decode

After `o.Run()` returns a `Report`, you can decode the typed results of each
step.

## Usage

```go
report, err := o.Run(ctx)
if err != nil {
    log.Fatal(err)
}

var cmd osapi.CommandResult
if err := report.Decode("run-uptime", &cmd); err == nil {
    fmt.Printf("stdout: %s\n", cmd.Stdout)
}
```

## Status inspection

Find if a step succeeded, failed, or was skipped:

```go
status := report.Status("step-name")
// status is one of: "success", "failed", "skipped"
```

## Example

Refer to
[`examples/operations/command.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/command.go)
and
[`examples/features/task-func.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/task-func.go)
for complete examples.
