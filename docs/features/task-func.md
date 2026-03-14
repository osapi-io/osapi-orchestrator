# TaskFunc

Insert custom logic into a plan with `TaskFunc`. The function receives
accumulated results from prior steps and returns a typed `*sdk.Result`.

## Usage

```go
o.TaskFunc(
    "summarize",
    func(_ context.Context, r orchestrator.Results) (*sdk.Result, error) {
        var h osapi.HostnameResult
        if err := r.Decode("get-hostname", &h); err != nil {
            return &sdk.Result{Changed: false}, nil
        }
        return &sdk.Result{
            Changed: true,
            Data:    map[string]any{"host": h.Hostname},
        }, nil
    },
).After(hostname)
```

TaskFunc results are available in the report via `report.Decode()`.

## Example

See
[`examples/features/task-func.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/task-func.go)
for a complete working example.
