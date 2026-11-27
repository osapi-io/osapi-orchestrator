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

`Decode` returns an error when the step has no result, or when its result has no
data, such as a step that was skipped. It also returns an error when the step
has more than one host result, such as a broadcast to `_all`, a label, or a
list. One host's data does not stand for the others: read each host with
`Results.HostResults` or `Report.HostResults`, and `HostResult.Decode`.
`Report.Decode` and the `Results.Decode` that a `TaskFunc` gets act the same, so
a step never goes on with zero values from a step that did not run.

## Status inspection

The report has the status of each step: `changed`, `unchanged`, `failed` or
`skipped`:

```go
for _, t := range report.Tasks {
    fmt.Printf("%s: %s\n", t.Name, t.Status)
}
```

A `TaskFunc` and a `When` guard read the status of a step that ran before them
with `Results.Status`. It returns `TaskStatusChanged`, `TaskStatusUnchanged`,
`TaskStatusFailed`, `TaskStatusSkipped`, or `TaskStatusUnknown` for a step with
no result:

```go
o.TaskFunc("report", func(
    ctx context.Context,
    c *osapi.Client,
    r orchestrator.Results,
) (*orchestrator.Result, error) {
    if r.Status("deploy") == orchestrator.TaskStatusFailed {
        // ...
    }

    return &orchestrator.Result{}, nil
}).After(deploy)
```

## Example

Refer to
[`examples/operations/command.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/command.go)
and
[`examples/features/task-func.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/task-func.go)
for complete examples.
