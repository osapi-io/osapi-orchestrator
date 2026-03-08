# Guards

Guards control whether a step runs based on runtime conditions.

## When

`When` takes a predicate function that receives accumulated results from prior
steps. The step only runs if the predicate returns true:

```go
o.CommandExec("_any", "whoami").
    After(hostname).
    When(func(r orchestrator.Results) bool {
        var h orchestrator.HostnameResult
        if err := r.Decode("get-hostname", &h); err != nil {
            return false
        }
        return h.Hostname != ""
    })
```

## WhenFact

`WhenFact` is a specialized guard that inspects agent facts from a prior
`AgentList` step:

```go
o.CommandShell(target, "apt-get update -qq").
    After(agents).
    WhenFact("list-agents", func(a orchestrator.AgentResult) bool {
        return a.OSInfo != nil &&
            a.OSInfo.Distribution == "Ubuntu"
    })
```

## OnlyIfChanged

See [OnlyIfChanged](only-if-changed.md).

## Example

See
[`examples/guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/guards.go)
and
[`examples/when-fact.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/when-fact.go)
for complete working examples.
