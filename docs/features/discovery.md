# Discovery

Find agents with fact predicates when you build the plan. Then make a step that
targets each host that matches.

## Discover

```go
agents, err := o.Discover(
    context.Background(),
    orchestrator.OS("Ubuntu"),
    orchestrator.Arch("amd64"),
)

for _, a := range agents {
    o.HostnameGet(a.Hostname).After(health)
}
```

## GroupByFact

Put agents into groups by a fact value, and run a different command for each
group:

```go
groups, err := o.GroupByFact(context.Background(), "os.distribution")

for distro, agents := range groups {
    for _, a := range agents {
        o.CommandShell(a.Hostname, installCmd(distro)).After(health)
    }
}
```

## Predicates

| Predicate      | What it matches                                              |
| -------------- | ------------------------------------------------------------ |
| `OS`           | The OS distribution of the agent (not case-sensitive)        |
| `Arch`         | The architecture of the agent (not case-sensitive)           |
| `MinMemory`    | A minimum total memory                                       |
| `MinCPU`       | A minimum CPU count                                          |
| `HasLabel`     | A label key-value pair                                       |
| `FactEquals`   | A fact that is equal to the given value, as JSON compares it |
| `HasCondition` | An agent that has an active condition of the given type      |
| `NoCondition`  | An agent that does NOT have the active condition             |
| `Healthy`      | An agent that has no active conditions                       |

The facts of an agent come from JSON. Thus `FactEquals` compares the value in
its JSON form: `FactEquals("cpu_count", 4)` matches the fact `4`, and a slice or
a map matches a fact with the same elements.

## Example

Refer to
[`examples/features/discover.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/discover.go),
[`examples/features/group-by-fact.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/group-by-fact.go),
[`examples/features/fact-predicates.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/fact-predicates.go),
[`examples/features/label-filter.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/label-filter.go),
and
[`examples/features/condition-filter.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/condition-filter.go)
for complete examples.
