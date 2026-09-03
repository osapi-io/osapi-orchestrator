# DSL Reference

The orchestrator provides a declarative DSL for composing OSAPI operations into
DAG-based plans with typed results, guards, retry, and discovery.

## How the DAG Works

Every operation method (`NodeHostnameGet`, `CommandExec`, etc.) returns a
`*Step`. Steps are connected into a directed acyclic graph (DAG) using
`After()`. The orchestrator resolves the DAG into levels and executes each
level:

- Steps with no dependencies run first.
- Steps at the same level (sharing the same dependencies) run in parallel.
- Steps with `After()` dependencies wait until those dependencies complete.

```go
health := o.HealthCheck()

// These three run in parallel — all depend on health.
hostname := o.NodeHostnameGet("_any").After(health)
disk := o.NodeDiskGet("_any").After(health)
load := o.NodeLoadGet("_any").After(health)

// This runs after hostname completes.
o.CommandExec("_any", "whoami").After(hostname)
```

## Step Chaining

Chain methods on any `*Step` to declare ordering, conditions, and error
handling:

| Method                  | What it does                                       | Guide                               |
| ----------------------- | -------------------------------------------------- | ----------------------------------- |
| `After`                 | Run after the given steps complete                 | [Basic DAG](basic.md)               |
| `Retry`                 | Retry on failure with optional exponential backoff | [Retry](retry.md)                   |
| `OnlyIfChanged`         | Skip unless a dependency reported changes          | [OnlyIfChanged](only-if-changed.md) |
| `OnlyIfFailed`          | Skip unless at least one dependency failed         | [Error Recovery](error-recovery.md) |
| `OnlyIfAllChanged`      | Skip unless all dependencies reported changes      |                                     |
| `OnlyIfAnyHostFailed`   | Skip unless any host in a dependency failed        | [Guards](guards.md)                 |
| `OnlyIfAllHostsFailed`  | Skip unless all hosts in dependencies failed       |                                     |
| `OnlyIfAnyHostSkipped`  | Skip unless any host in a dependency was skipped   | [Guards](guards.md)                 |
| `OnlyIfAnyHostChanged`  | Skip unless any host in a dependency changed       |                                     |
| `OnlyIfAllHostsChanged` | Skip unless all hosts in dependencies changed      | [Guards](guards.md)                 |
| `When`                  | Guard -- only run if predicate returns true        | [Guards](guards.md)                 |
| `WhenFact`              | Guard -- only run if agent fact matches            | [Guards](guards.md)                 |
| `ContinueOnError`       | Keep running independent tasks on failure          | [Error Recovery](error-recovery.md) |
| `OnError`               | Set error strategy (`StopAll` or `Continue`)       | [Error Recovery](error-recovery.md) |
| `Named`                 | Set a custom step name for result decoding         | [Result Decode](result-decode.md)   |

## Typed Results

Decode step results into typed structs:

```go
var h osapi.HostnameResult
err := results.Decode("node.hostname.get-1", &h)
fmt.Println(h.Hostname)
```

| Struct              | Fields                                                   |
| ------------------- | -------------------------------------------------------- |
| `HostnameResult`    | `Hostname`, `Labels`                                     |
| `DiskResult`        | `Disks` (slice of `DiskUsage`)                           |
| `MemoryResult`      | `Total`, `Free`, `Cached`                                |
| `LoadResult`        | `Load1`, `Load5`, `Load15`                               |
| `CommandResult`     | `Stdout`, `Stderr`, `ExitCode`, `DurationMs`, `Error`    |
| `PingResult`        | `PacketsSent`, `PacketsReceived`, `PacketLoss`, `Error`  |
| `DNSConfigResult`   | `DNSServers`, `SearchDomains`                            |
| `DNSUpdateResult`   | `Success`, `Message`, `Error`                            |
| `FileDeployResult`  | `Changed`, `SHA256`, `Path`                              |
| `FileStatusResult`  | `Path`, `Status`, `SHA256`                               |
| `FileUploadResult`  | `Name`, `SHA256`, `Size`, `Changed`, `ContentType`       |
| `FileChangedResult` | `Name`, `Changed`, `SHA256`                              |
| `AgentListResult`   | `Agents` (slice of `AgentResult`), `Total`               |
| `AgentResult`       | `Hostname`, `Status`, `Architecture`, `OSInfo`, `Memory` |

## Predicates

Composable filters passed to `Discover` and `GroupByFact`:

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

## Error Strategies

| Strategy            | Behavior                                        |
| ------------------- | ----------------------------------------------- |
| `StopAll` (default) | Fail fast, cancel everything                    |
| `Continue`          | Skip dependents, keep running independent tasks |

Per-step retry is available via `.Retry(n)`.

## Feature Guides

| Guide                               | What it covers                                |
| ----------------------------------- | --------------------------------------------- |
| [Basic DAG](basic.md)               | Creating a simple DAG plan                    |
| [Parallel Execution](parallel.md)   | Concurrent tasks at the same DAG level        |
| [Guards](guards.md)                 | When predicates, OnlyIfChanged, WhenFact      |
| [OnlyIfChanged](only-if-changed.md) | Skip unless dependency reported changes       |
| [Retry](retry.md)                   | Retry on failure with configurable attempts   |
| [Error Recovery](error-recovery.md) | Continue strategy with OnlyIfFailed cleanup   |
| [Broadcast](broadcast.md)           | Per-host results from `_all`/label targets    |
| [TaskFunc](task-func.md)            | Custom logic in a plan                        |
| [Result Decode](result-decode.md)   | Post-execution result decoding and inspection |
| [File Workflow](file-workflow.md)   | Upload, deploy, verify file lifecycle         |
| [Discovery](discovery.md)           | Agent discovery, GroupByFact, predicates      |
| [Verbose Output](verbose.md)        | WithVerbose() debug output                    |
| [DNS Update](dns-update.md)         | Read-then-write pattern                       |
