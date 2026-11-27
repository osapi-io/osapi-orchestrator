# TaskFunc

Use `TaskFunc` to put custom logic into a plan. The function receives the OSAPI
client and the results of the steps that ran before. It returns a typed
`*orchestrator.Result`. Use `TaskFunc` for an operation that has no typed
constructor. The client gives full access to the SDK, so you can call each API
endpoint.

## Usage

```go
o.TaskFunc(
    "summarize",
    func(_ context.Context, c *osapi.Client, r orchestrator.Results) (*orchestrator.Result, error) {
        var h osapi.HostnameEntry
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

With the client parameter, you can call an SDK operation directly:

```go
o.TaskFunc(
    "custom-check",
    func(ctx context.Context, c *osapi.Client, _ orchestrator.Results) (*orchestrator.Result, error) {
        resp, err := c.Health.Status(ctx, osapi.HealthStatusParams{})
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

To read a TaskFunc result from the report, use `report.Decode()`.

## Example

Refer to
[`examples/features/task-func.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/task-func.go)
for a complete example.
