# Features

The orchestrator provides a declarative DSL for composing OSAPI operations into
DAG-based plans with typed results, guards, retry, and discovery.

## Step Chaining

Every operation returns a `*Step`. Chain methods to declare ordering,
conditions, and error handling:

| Method                  | What it does                                       |
| ----------------------- | -------------------------------------------------- |
| `After`                 | Run after the given steps complete                 |
| `Retry`                 | Retry on failure with optional exponential backoff |
| `OnlyIfChanged`         | Skip unless a dependency reported changes          |
| `OnlyIfFailed`          | Skip unless at least one dependency failed         |
| `OnlyIfAllChanged`      | Skip unless all dependencies reported changes      |
| `OnlyIfAnyHostFailed`   | Skip unless any host in a dependency has an error  |
| `OnlyIfAllHostsFailed`  | Skip unless all hosts in dependencies have errors  |
| `OnlyIfAnyHostChanged`  | Skip unless any host in a dependency changed       |
| `OnlyIfAllHostsChanged` | Skip unless all hosts in dependencies changed      |
| `When`                  | Guard -- only run if predicate returns true        |
| `WhenFact`              | Guard -- only run if agent fact matches            |
| `OnError`               | Set error strategy (`StopAll` or `Continue`)       |

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
