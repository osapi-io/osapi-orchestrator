[![release](https://img.shields.io/github/release/osapi-io/osapi-orchestrator.svg?style=for-the-badge)](https://github.com/osapi-io/osapi-orchestrator/releases/latest)
[![codecov](https://img.shields.io/codecov/c/github/osapi-io/osapi-orchestrator?style=for-the-badge)](https://codecov.io/gh/osapi-io/osapi-orchestrator)
[![go report card](https://goreportcard.com/badge/github.com/osapi-io/osapi-orchestrator?style=for-the-badge)](https://goreportcard.com/report/github.com/osapi-io/osapi-orchestrator)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/osapi-io/osapi-orchestrator/go.yml?style=for-the-badge)](https://github.com/osapi-io/osapi-orchestrator/actions/workflows/go.yml)
[![powered by](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![built with just](https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white)](https://just.systems)
![gitHub commit activity](https://img.shields.io/github/commit-activity/m/osapi-io/osapi-orchestrator?style=for-the-badge)

# OSAPI Orchestrator

A Go package for orchestrating operations across [OSAPI][]-managed hosts —
typed operations, chaining, conditions, and result decoding built on top
of the [osapi-sdk][] engine.

## 📦 Install

```bash
go install github.com/osapi-io/osapi-orchestrator@latest
```

As a library dependency:

```bash
go get github.com/osapi-io/osapi-orchestrator
```

## 🎯 Targeting

Most operations accept a `target` parameter to control which agents receive
the request:

| Target      | Behavior                                    |
| ----------- | ------------------------------------------- |
| `_any`      | Send to any available agent (load balanced) |
| `_all`      | Broadcast to every agent                    |
| `hostname`  | Send to a specific host                     |
| `key:value` | Send to agents matching a label             |

Agents expose labels (used for targeting) and extended system facts via
`AgentList` and `AgentGet`. Facts come from agent-side providers and
include OS, hardware, and network details.

## ✨ Features

Typed constructors, typed results, and chainable step methods. See the
[usage docs](docs/usage/README.md) for full details, examples, and
per-operation reference.

### Typed Operations

Every OSAPI operation has a strongly typed constructor — no raw maps or
string constants.

| Method             | Operation               | Docs                                     | Source                              |
| ------------------ | ----------------------- | ---------------------------------------- | ----------------------------------- |
| `HealthCheck`      | Liveness probe          | [docs](docs/usage/health-check.md)       | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeHostnameGet`  | `node.hostname.get`     | [docs](docs/usage/node-hostname-get.md)  | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeStatusGet`    | `node.status.get`       | [docs](docs/usage/node-status-get.md)    | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeUptimeGet`    | `node.uptime.get`       | [docs](docs/usage/node-uptime-get.md)    | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeDiskGet`      | `node.disk.get`         | [docs](docs/usage/node-disk-get.md)      | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeMemoryGet`    | `node.memory.get`       | [docs](docs/usage/node-memory-get.md)    | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeLoadGet`      | `node.load.get`         | [docs](docs/usage/node-load-get.md)      | [`ops.go`](pkg/orchestrator/ops.go) |
| `NetworkDNSGet`    | `network.dns.get`       | [docs](docs/usage/network-dns-get.md)    | [`ops.go`](pkg/orchestrator/ops.go) |
| `NetworkDNSUpdate` | `network.dns.update`    | [docs](docs/usage/network-dns-update.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NetworkPingDo`    | `network.ping.do`       | [docs](docs/usage/network-ping-do.md)    | [`ops.go`](pkg/orchestrator/ops.go) |
| `CommandExec`      | `command.exec.execute`  | [docs](docs/usage/command-exec.md)       | [`ops.go`](pkg/orchestrator/ops.go) |
| `CommandShell`     | `command.shell.execute` | [docs](docs/usage/command-shell.md)      | [`ops.go`](pkg/orchestrator/ops.go) |
| `AgentList`        | List active agents      | [docs](docs/usage/agent-list.md)         | [`ops.go`](pkg/orchestrator/ops.go) |
| `AgentGet`         | Get agent details       | [docs](docs/usage/agent-get.md)          | [`ops.go`](pkg/orchestrator/ops.go) |

### Typed Results

Decode step results into typed structs instead of digging through
`map[string]any`. See [`result_types.go`](pkg/orchestrator/result_types.go).

| Struct            | Fields                                                  |
| ----------------- | ------------------------------------------------------- |
| `HostnameResult`  | `Hostname`, `Labels`                                    |
| `DiskResult`      | `Disks` (slice of `DiskUsage`)                          |
| `MemoryResult`    | `Total`, `Free`, `Cached`                               |
| `LoadResult`      | `Load1`, `Load5`, `Load15`                              |
| `CommandResult`   | `Stdout`, `Stderr`, `ExitCode`, `DurationMs`, `Error`   |
| `PingResult`      | `PacketsSent`, `PacketsReceived`, `PacketLoss`, `Error` |
| `DNSConfigResult` | `DNSServers`, `SearchDomains`                           |
| `DNSUpdateResult` | `Success`, `Message`, `Error`                           |
| `AgentListResult` | `Agents` (slice of `AgentResult`), `Total`              |
| `AgentResult`     | `Hostname`, `Status`, `Architecture`, `OSInfo`, `Memory` |

### Step Chaining

Declare ordering, conditions, and error handling with chainable methods.
See [`step.go`](pkg/orchestrator/step.go).

```go
o.CommandExec("_any", "whoami").
    After(health, hostname).
    Retry(2).
    OnlyIfChanged().
    When(func(r orchestrator.Results) bool {
        var h orchestrator.HostnameResult
        r.Decode("get-hostname", &h)
        return h.Hostname != ""
    }).
    OnError(orchestrator.Continue)
```

| Method             | What it does                                  |
| ------------------ | --------------------------------------------- |
| `After`            | Run after the given steps complete            |
| `Retry`            | Retry on failure up to N times                |
| `OnlyIfChanged`    | Skip unless a dependency reported changes     |
| `OnlyIfFailed`     | Skip unless at least one dependency failed    |
| `OnlyIfAllChanged` | Skip unless all dependencies reported changes |
| `When`             | Guard — only run if predicate returns true    |
| `WhenFact`         | Guard — only run if agent fact predicate is true |
| `OnError`          | Set error strategy (`StopAll` or `Continue`)  |

### Custom Steps

Use `TaskFunc` to create custom steps that receive completed results from
prior steps — useful for decision logic, aggregation, or conditional
branching:

```go
o.TaskFunc("summarize", func(ctx context.Context, r orchestrator.Results) (*sdk.Result, error) {
    var h orchestrator.HostnameResult
    r.Decode("get-hostname", &h)

    return &sdk.Result{
        Changed: true,
        Data:    map[string]any{"summary": h.Hostname},
    }, nil
}).After(hostname)
```

### Configuration

Pass options to `New` to configure behavior:

```go
o := orchestrator.New(url, token, orchestrator.WithVerbose())
```

| Option          | What it does                                         |
| --------------- | ---------------------------------------------------- |
| `WithVerbose()` | Show stdout, stderr, and response data for all tasks |

### Post-Execution Results

After `Run()` completes, decode individual task results from the report:

```go
report, err := o.Run()

var cmd orchestrator.CommandResult
err = report.Decode("run-uptime", &cmd)
fmt.Println(cmd.Stdout)
```

### Status Inspection

Inside `When` guards, inspect the status of completed dependencies:

```go
step.When(func(r orchestrator.Results) bool {
    return r.Status("health-check") == orchestrator.TaskStatusChanged
})
```

| Constant              | Meaning                       |
| --------------------- | ----------------------------- |
| `TaskStatusUnknown`   | Step not found or has not run |
| `TaskStatusChanged`   | Step ran and reported changes |
| `TaskStatusUnchanged` | Step ran with no changes      |
| `TaskStatusSkipped`   | Step was skipped              |
| `TaskStatusFailed`    | Step failed                   |

### Broadcast Results

When targeting `_all` or label selectors, access per-host results:

```go
hrs := results.HostResults("deploy")
for _, hr := range hrs {
    fmt.Printf("host=%s changed=%v error=%s\n", hr.Hostname, hr.Changed, hr.Error)

    var cmd orchestrator.CommandResult
    hr.Decode(&cmd)
}
```

### Agent Discovery

Query agents at plan-build time and filter by typed predicates:

```go
agents, err := o.Discover(ctx,
    orchestrator.OS("Ubuntu"),
    orchestrator.Arch("amd64"),
    orchestrator.MinCPU(4),
)

for _, a := range agents {
    o.CommandShell(a.Hostname, "apt upgrade -y").After(health)
}
```

| Method        | What it does                                         |
| ------------- | ---------------------------------------------------- |
| `Discover`    | Query agents filtered by predicates                  |
| `GroupByFact` | Group agents by a fact key (e.g. `os.distribution`)  |

### Predicates

Composable filters passed to `Discover` and `GroupByFact`:

| Predicate    | What it matches                          |
| ------------ | ---------------------------------------- |
| `OS`         | Agent OS distribution (case-insensitive) |
| `Arch`       | Agent architecture (case-insensitive)    |
| `MinMemory`  | Minimum total memory                     |
| `MinCPU`     | Minimum CPU count                        |
| `HasLabel`   | Label key-value pair                     |
| `FactEquals` | Arbitrary fact key-value equality        |

### Fact Guards

Use `WhenFact` for execution-time fact checks with a prior `AgentList`
step:

```go
agents := o.AgentList().After(health)

o.CommandShell("web-01", "apt upgrade -y").
    After(agents).
    WhenFact("list-agents", func(a orchestrator.AgentResult) bool {
        return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
    })
```

## 📋 Examples

Each example is a standalone Go program you can read and run.

### Core

| Example                                              | What it shows                                   |
| ---------------------------------------------------- | ----------------------------------------------- |
| [basic](examples/basic/main.go)                      | Simple DAG with health check and hostname query |
| [parallel](examples/parallel/main.go)                | Five parallel queries depending on health check |
| [retry](examples/retry/main.go)                      | Retry on failure with configurable attempts     |
| [command](examples/command/main.go)                   | Command exec and shell with result decoding     |
| [verbose](examples/verbose/main.go)                   | Verbose output with stdout/stderr/response data |

### Guards and Conditions

| Example                                              | What it shows                                      |
| ---------------------------------------------------- | -------------------------------------------------- |
| [guards](examples/guards/main.go)                    | When predicate for conditional execution           |
| [only-if-changed](examples/only-if-changed/main.go)  | Skip step unless dependency reported changes       |
| [error-recovery](examples/error-recovery/main.go)    | Continue strategy with OnlyIfFailed cleanup        |

### Results

| Example                                              | What it shows                                      |
| ---------------------------------------------------- | -------------------------------------------------- |
| [broadcast](examples/broadcast/main.go)              | Per-host results from broadcast operations         |
| [task-func](examples/task-func/main.go)              | Custom steps with typed result decoding            |
| [dns-update](examples/dns-update/main.go)            | Read-then-write pattern with DNS operations        |

### Agent Discovery

| Example                                              | What it shows                                      |
| ---------------------------------------------------- | -------------------------------------------------- |
| [agent-facts](examples/agent-facts/main.go)          | List agents with OS, load, memory, and interfaces  |
| [discover](examples/discover/main.go)                | Find agents by OS and architecture predicates      |
| [group-by-fact](examples/group-by-fact/main.go)      | Group agents by distro, run per-group commands     |
| [when-fact](examples/when-fact/main.go)               | Fact-based guard on a step                         |
| [fact-predicates](examples/fact-predicates/main.go)   | Compose multiple predicates for discovery          |
| [label-filter](examples/label-filter/main.go)         | Filter by labels and arbitrary fact values         |

```bash
cd examples/discover
OSAPI_TOKEN="<jwt>" go run main.go
```

## 🤝 Contributing

See the [Development](docs/development.md) guide for prerequisites, setup,
and conventions. See the [Contributing](docs/contributing.md) guide before
submitting a PR.

## 📄 License

The [MIT][] License.

[OSAPI]: https://github.com/osapi-io/osapi
[osapi-sdk]: https://github.com/osapi-io/osapi-sdk
[MIT]: LICENSE
