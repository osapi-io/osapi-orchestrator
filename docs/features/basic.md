# Basic DAG

The simplest orchestrator pattern: create a plan with ordered steps and run it.

## Usage

```go
o := orchestrator.New(url, token)

health := o.HealthCheck()
o.NodeHostnameGet("_any").After(health)

report, err := o.Run()
```

Steps are connected with `.After()` to form a directed acyclic graph. The
orchestrator resolves the dependency order and executes steps level by level.

## Example

See
[`examples/features/basic.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/basic.go)
for a complete working example.
