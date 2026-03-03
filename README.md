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
go get github.com/osapi-io/osapi-orchestrator
```

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

## 📋 Examples

Each example is a standalone Go program you can read and run.

| Example                     | What it shows                                                                      |
| --------------------------- | ---------------------------------------------------------------------------------- |
| [all](examples/all/main.go) | Fleet discovery, custom steps, guards, recovery, verbose mode, and result decoding |

```bash
cd examples/all
OSAPI_TOKEN="<jwt>" go run main.go                # normal output
OSAPI_TOKEN="<jwt>" OSAPI_VERBOSE=1 go run main.go  # verbose output
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
