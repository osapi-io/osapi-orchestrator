# Basic DAG

This is the simplest orchestrator pattern. Make a plan with steps in an order,
then run the plan.

## Usage

```go
o := orchestrator.New(url, token)

health := o.HealthCheck()
o.NodeHostnameGet("_any").After(health)

report, err := o.Run(ctx)
```

`.After()` connects the steps into a directed acyclic graph. The orchestrator
finds the order of the dependencies and runs the steps one level at a time.

## Example

Refer to
[`examples/features/basic.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/basic.go)
for a complete example.
