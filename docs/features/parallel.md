# Parallel execution

Steps at the same DAG level run at the same time. To run steps in parallel, give
them the same dependency:

## Usage

```go
health := o.HealthCheck()

o.NodeHostnameGet("_any").After(health)
o.NodeDiskGet("_any").After(health)
o.NodeMemoryGet("_any").After(health)
o.NodeLoadGet("_any").After(health)
o.NodeUptimeGet("_any").After(health)
```

The five queries have the same dependency. Thus the orchestrator runs them at
the same time after the health check passes.

## Example

Refer to
[`examples/features/parallel.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/parallel.go)
for a complete example.
