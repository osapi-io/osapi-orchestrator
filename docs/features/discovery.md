# Discovery

Discover agents by fact predicates at plan-build time, then generate steps
targeting each matching host.

## Discover

```go
agents, err := o.Discover(
    context.Background(),
    orchestrator.OS("Ubuntu"),
    orchestrator.Arch("amd64"),
)

for _, a := range agents {
    o.NodeHostnameGet(a.Hostname).After(health)
}
```

## GroupByFact

Group agents by a fact value and run different commands per group:

```go
groups, err := o.GroupByFact(context.Background(), "os.distribution")

for distro, agents := range groups {
    for _, a := range agents {
        o.CommandShell(a.Hostname, installCmd(distro)).After(health)
    }
}
```

## Predicates

| Predicate      | What it matches                          |
| -------------- | ---------------------------------------- |
| `OS`           | Agent OS distribution (case-insensitive) |
| `Arch`         | Agent architecture (case-insensitive)    |
| `MinMemory`    | Minimum total memory                     |
| `MinCPU`       | Minimum CPU count                        |
| `HasLabel`     | Label key-value pair                     |
| `FactEquals`   | Arbitrary fact key-value equality        |
| `HasCondition` | Agent has active condition of given type |
| `NoCondition`  | Agent does NOT have active condition     |
| `Healthy`      | Agent has no active conditions           |

## Example

See
[`examples/features/discover.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/discover.go),
[`examples/features/group-by-fact.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/group-by-fact.go),
[`examples/features/fact-predicates.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/fact-predicates.go),
[`examples/features/label-filter.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/label-filter.go),
and
[`examples/features/condition-filter.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/condition-filter.go)
for complete working examples.
