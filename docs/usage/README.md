# Usage

The orchestrator provides typed constructors for every OSAPI operation. Each
method returns a `*Step` that can be chained with ordering, conditions, and
error handling.

## Quick Start

```go
o := orchestrator.New(url, token)             // or: New(url, token, orchestrator.WithVerbose())

health := o.HealthCheck("_any")
hostname := o.NodeHostnameGet("_any").After(health)
o.CommandExec("_any", "whoami").After(hostname)

report, err := o.Run()
// report.Decode("command.exec.execute-1", &cmd) to extract typed results
```

## Operations

| Method                                      | Operation               | Idempotent | Category |
| ------------------------------------------- | ----------------------- | ---------- | -------- |
| [`HealthCheck`](health-check.md)            | Liveness probe          | Read-only  | Health   |
| [`NodeHostnameGet`](node-hostname-get.md)   | `node.hostname.get`     | Read-only  | Node     |
| [`NodeStatusGet`](node-status-get.md)       | `node.status.get`       | Read-only  | Node     |
| [`NodeUptimeGet`](node-uptime-get.md)       | `node.uptime.get`       | Read-only  | Node     |
| [`NodeDiskGet`](node-disk-get.md)           | `node.disk.get`         | Read-only  | Node     |
| [`NodeMemoryGet`](node-memory-get.md)       | `node.memory.get`       | Read-only  | Node     |
| [`NodeLoadGet`](node-load-get.md)           | `node.load.get`         | Read-only  | Node     |
| [`NetworkDNSGet`](network-dns-get.md)       | `network.dns.get`       | Read-only  | Network  |
| [`NetworkDNSUpdate`](network-dns-update.md) | `network.dns.update`    | Yes        | Network  |
| [`NetworkPingDo`](network-ping-do.md)       | `network.ping.do`       | Read-only  | Network  |
| [`CommandExec`](command-exec.md)            | `command.exec.execute`  | No         | Command  |
| [`CommandShell`](command-shell.md)          | `command.shell.execute` | No         | Command  |
| [`FileDeploy`](file-deploy.md)              | `file.deploy.execute`   | Yes        | File     |
| [`FileStatusGet`](file-status-get.md)       | `file.status.get`       | Read-only  | File     |
| [`FileUpload`](file-upload.md)              | Upload to Object Store  | No         | File     |
| [`FileChanged`](file-changed.md)            | Check file drift        | Read-only  | File     |
| [`AgentList`](agent-list.md)                | List active agents      | Read-only  | Agent    |
| [`AgentGet`](agent-get.md)                  | Get agent details       | Read-only  | Agent    |

## Step Chaining

Every operation returns a `*Step`. Chain methods to declare ordering,
conditions, and error handling:

| Method             | What it does                                  |
| ------------------ | --------------------------------------------- |
| `After`            | Run after the given steps complete            |
| `Retry`            | Retry on failure up to N times                |
| `OnlyIfChanged`    | Skip unless a dependency reported changes     |
| `OnlyIfFailed`     | Skip unless at least one dependency failed    |
| `OnlyIfAllChanged` | Skip unless all dependencies reported changes |
| `When`             | Guard — only run if predicate returns true    |
| `OnError`          | Set error strategy (`StopAll` or `Continue`)  |

## Typed Results

Decode step results into typed structs:

```go
var h orchestrator.HostnameResult
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

## Idempotency

- **Read-only** operations never modify state and always return
  `Changed: false`.
- **Idempotent** write operations check current state before mutating and return
  `Changed: true` only if something actually changed.
- **Non-idempotent** operations (command exec/shell) always return
  `Changed: true`. Use guards (`When`, `OnlyIfChanged`) to control when they
  run.
