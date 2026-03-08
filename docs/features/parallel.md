# Parallel Execution

Steps at the same DAG level run concurrently. Give multiple steps the same
dependency to fan out:

## Usage

```go
health := o.HealthCheck("_any")

o.NodeHostnameGet("_any").After(health)
o.NodeDiskGet("_any").After(health)
o.NodeMemoryGet("_any").After(health)
o.NodeLoadGet("_any").After(health)
o.NodeUptimeGet("_any").After(health)
```

All five queries share the same dependency, so the orchestrator schedules them
concurrently after the health check passes.

## Example

See
[`examples/features/parallel.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/parallel.go)
for a complete working example.
