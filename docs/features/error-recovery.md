# Error Recovery

Use `OnError(Continue)` to keep independent tasks running when a step fails, and
`OnlyIfFailed` to trigger cleanup only on failure.

## Usage

```go
deploy := o.CommandExec("_all", "echo", "deploying").
    Named("deploy").
    OnError(orchestrator.Continue)

o.CommandExec("_any", "echo", "running-cleanup").
    Named("cleanup").
    After(deploy).
    OnlyIfFailed()
```

The `Continue` strategy lets the plan proceed even if some hosts fail. The
cleanup step only runs when at least one dependency failed.

## Error Strategies

| Strategy   | Behavior                                        |
| ---------- | ----------------------------------------------- |
| `StopAll`  | Fail fast, cancel everything (default)          |
| `Continue` | Skip dependents, keep running independent tasks |

## Example

See
[`examples/error-recovery.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/error-recovery.go)
for a complete working example.
