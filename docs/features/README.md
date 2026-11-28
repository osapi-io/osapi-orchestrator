# DSL reference

The orchestrator has a declarative DSL. You use it to put OSAPI operations into
DAG plans with typed results, guards, retry, and discovery.

## How the DAG works

Each operation method (`HostnameGet`, `CommandExec`, and the others) returns a
`*Step`. `After()` connects steps into a directed acyclic graph (DAG). The
orchestrator divides the DAG into levels and runs each level:

- Steps that have no dependencies run first.
- Steps at the same level (with the same dependencies) run in parallel.
- A step with `After()` dependencies waits until those dependencies complete.

```go
health := o.HealthCheck()

// These three depend on health, and run in parallel.
hostname := o.HostnameGet("_any").After(health)
disk := o.DiskGet("_any").After(health)
load := o.LoadGet("_any").After(health)

// This runs after hostname completes.
o.CommandExec("_any", osapi.CommandExecInput{Command: "whoami"}).After(hostname)
```

## Step chaining

Chain methods on a `*Step` to declare the order, the conditions, and what to do
on an error:

| Method                  | What it does                                                 | Guide                               |
| ----------------------- | ------------------------------------------------------------ | ----------------------------------- |
| `After`                 | Run after the given steps complete                           | [Basic DAG](basic.md)               |
| `Retry`                 | Retry after a failure, with exponential backoff              | [Retry](retry.md)                   |
| `OnlyIfChanged`         | Skip if no dependency reported a change                      | [OnlyIfChanged](only-if-changed.md) |
| `OnlyIfFailed`          | Skip if no dependency failed                                 | [Error Recovery](error-recovery.md) |
| `OnlyIfAllChanged`      | Skip if one or more dependencies reported no change          |                                     |
| `OnlyIfAnyHostFailed`   | Skip if no host in a dependency failed                       | [Guards](guards.md)                 |
| `OnlyIfAllHostsFailed`  | Skip if one or more hosts in the dependencies did not fail   |                                     |
| `OnlyIfAnyHostSkipped`  | Skip if no dependency skipped a host                         | [Guards](guards.md)                 |
| `OnlyIfAnyHostChanged`  | Skip if no host in a dependency changed                      |                                     |
| `OnlyIfAllHostsChanged` | Skip if one or more hosts in the dependencies did not change | [Guards](guards.md)                 |
| `When`                  | Guard. Run only if the predicate returns true                | [Guards](guards.md)                 |
| `WhenFact`              | Guard. Run only if the agent fact matches                    | [Guards](guards.md)                 |
| `ContinueOnError`       | Continue to run independent steps after a failure            | [Error Recovery](error-recovery.md) |
| `OnError`               | Set the error strategy (`StopAll` or `Continue`)             | [Error Recovery](error-recovery.md) |
| `Named`                 | Set a custom step name for result decoding                   | [Result Decode](result-decode.md)   |

## Typed results

Decode the step results into typed structs:

```go
var h osapi.HostnameEntry
err := results.Decode("get-hostname", &h)
fmt.Println(h.Name)
```

A step that targets hosts has one row for each host, and a list has one row for
each item on each host. `Decode` takes the first row. Use `HostResults` for all
of them. Each row is the type of its SDK call: the reads answer with an entry
type, such as `osapi.HostnameEntry`, and the writes with a result type, such as
`osapi.HostnameResult`, which also has `Changed`. Each operation's page names
its type and fields.

| Struct                   | Fields                                                   |
| ------------------------ | -------------------------------------------------------- |
| `HostnameEntry`          | `Hostname`, `Name`, `Status`, `Error`                    |
| `DiskEntry`              | `Hostname`, `Disks`, `Status`, `Error`                   |
| `MemoryEntry`            | `Hostname`, `Memory` (`Total`, `Used`, `Free`), `Status` |
| `LoadEntry`              | `Hostname`, `LoadAverage` (`Min1`, `Min5`, `Min15`)      |
| `CommandResult`          | `Stdout`, `Stderr`, `ExitCode`, `DurationMS`, `Error`    |
| `ConnectivityPingResult` | `PacketsSent`, `PacketsReceived`, `PacketLoss`, `Error`  |
| `DnsEntry`               | `Servers`, `SearchDomains`                               |
| `FileUpdateResult`       | `Changed`, `SHA256`, `Path`                              |
| `FileStatus`             | `Path`, `State`, `SHA256`                                |
| `FileUploadResponse`     | `Name`, `SHA256`, `Size`, `Changed`, `ContentType`       |
| `FileChanged`            | `Name`, `Changed`, `SHA256`                              |
| `ListAgentsResponse`     | `Items` (slice of `AgentInfo`), `TotalItems`             |
| `AgentInfo`              | `Hostname`, `State`, `Architecture`, `OSInfo`, `Memory`  |

## Predicates

A predicate is a filter. You give predicates to `Discover` and `GroupByFact`,
and you can use more than one together:

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

## Error strategies

| Strategy            | Behavior                                                      |
| ------------------- | ------------------------------------------------------------- |
| `StopAll` (default) | Stop at the first failure and cancel all other steps          |
| `Continue`          | Skip the dependents and continue to run the independent steps |

`.Retry(n)` sets a retry for each step. A retry is separate from the error
strategy. Thus `.Retry(3).ContinueOnError()` retries the step, and continues the
plan if the last attempt fails.

When a step fails, `Run` returns the report with the error. The report has each
step that ran. The steps that did not run are in it as skipped, with the reason
`plan stopped`. The error names each step that failed.

## Feature guides

| Guide                               | What it covers                                                |
| ----------------------------------- | ------------------------------------------------------------- |
| [Basic DAG](basic.md)               | How to make a simple DAG plan                                 |
| [Parallel Execution](parallel.md)   | Steps that run at the same time at one DAG level              |
| [Guards](guards.md)                 | When predicates, OnlyIfChanged, WhenFact                      |
| [OnlyIfChanged](only-if-changed.md) | Skip if no dependency reported a change                       |
| [Retry](retry.md)                   | Retry after a failure, with a number of attempts that you set |
| [Error Recovery](error-recovery.md) | The Continue strategy with OnlyIfFailed cleanup               |
| [Broadcast](broadcast.md)           | Results for each host from `_all` and label targets           |
| [TaskFunc](task-func.md)            | Custom logic in a plan                                        |
| [Result Decode](result-decode.md)   | How to decode and examine results after the plan runs         |
| [File Workflow](file-workflow.md)   | Upload, deploy, and verify a file                             |
| [Discovery](discovery.md)           | Agent discovery, GroupByFact, predicates                      |
| [Verbose Output](verbose.md)        | WithVerbose() debug output                                    |
| [DNS Update](dns-update.md)         | The pattern that reads, then writes                           |
