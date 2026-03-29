# TaskFunc

Insert custom logic into a plan with `TaskFunc`. The function receives the OSAPI
client and accumulated results from prior steps, and returns a typed
`*orchestrator.Result`. Use this for operations not covered by the typed constructors —
the client provides full access to the SDK for calling any API endpoint.

## Usage

```go
o.TaskFunc(
    "summarize",
    func(_ context.Context, c *osapi.Client, r orchestrator.Results) (*orchestrator.Result, error) {
        var h osapi.HostnameResult
        if err := r.Decode("get-hostname", &h); err != nil {
            return &orchestrator.Result{Changed: false}, nil
        }
        return &orchestrator.Result{
            Changed: true,
            Data:    map[string]any{"host": h.Hostname},
        }, nil
    },
).After(hostname)
```

The client parameter lets you call any SDK operation directly:

```go
o.TaskFunc(
    "custom-check",
    func(ctx context.Context, c *osapi.Client, _ orchestrator.Results) (*orchestrator.Result, error) {
        resp, err := c.Health.Status(ctx)
        if err != nil {
            return nil, err
        }
        return &orchestrator.Result{
            Changed: false,
            Data:    orchestrator.StructToMap(resp.Data),
        }, nil
    },
)
```

TaskFunc results are available in the report via `report.Decode()`.

## Example

See
[`examples/features/task-func.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/task-func.go)
for a complete working example.
