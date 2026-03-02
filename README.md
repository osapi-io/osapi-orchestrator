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

## Install

```bash
go get github.com/osapi-io/osapi-orchestrator
```

## Features

Typed constructors, typed results, and chainable step methods. See the
[usage docs](docs/usage/README.md) for full details, examples, and
per-operation reference.

### Typed Operations

Every OSAPI operation has a strongly typed constructor — no raw maps or
string constants.

| Method | Operation | Docs | Source |
| ------ | --------- | ---- | ------ |
| `HealthCheck` | Liveness probe | [docs](docs/usage/health-check.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeHostnameGet` | `node.hostname.get` | [docs](docs/usage/node-hostname-get.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeStatusGet` | `node.status.get` | [docs](docs/usage/node-status-get.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeUptimeGet` | `node.uptime.get` | [docs](docs/usage/node-uptime-get.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeDiskGet` | `node.disk.get` | [docs](docs/usage/node-disk-get.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeMemoryGet` | `node.memory.get` | [docs](docs/usage/node-memory-get.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NodeLoadGet` | `node.load.get` | [docs](docs/usage/node-load-get.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NetworkDNSGet` | `network.dns.get` | [docs](docs/usage/network-dns-get.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NetworkDNSUpdate` | `network.dns.update` | [docs](docs/usage/network-dns-update.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `NetworkPingDo` | `network.ping.do` | [docs](docs/usage/network-ping-do.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `CommandExec` | `command.exec.execute` | [docs](docs/usage/command-exec.md) | [`ops.go`](pkg/orchestrator/ops.go) |
| `CommandShell` | `command.shell.execute` | [docs](docs/usage/command-shell.md) | [`ops.go`](pkg/orchestrator/ops.go) |

### Typed Results

Decode step results into typed structs instead of digging through
`map[string]any`. See [`result_types.go`](pkg/orchestrator/result_types.go).

| Struct | Fields |
| ------ | ------ |
| `HostnameResult` | `Hostname`, `Labels` |
| `DiskResult` | `Disks` (slice of `DiskUsage`) |
| `MemoryResult` | `Total`, `Free`, `Cached` |
| `LoadResult` | `Load1`, `Load5`, `Load15` |
| `CommandResult` | `Stdout`, `Stderr`, `ExitCode`, `DurationMs` |
| `PingResult` | `PacketsSent`, `PacketsReceived`, `PacketLoss` |
| `DNSConfigResult` | `DNSServers`, `SearchDomains` |
| `DNSUpdateResult` | `Success`, `Message` |

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

| Method | What it does |
| ------ | ------------ |
| `After` | Run after the given steps complete |
| `Retry` | Retry on failure up to N times |
| `OnlyIfChanged` | Skip unless a dependency reported changes |
| `When` | Guard — only run if predicate returns true |
| `OnError` | Set error strategy (`StopAll` or `Continue`) |

## Examples

Each example is a standalone Go program you can read and run.

| Example | What it shows |
| ------- | ------------- |
| [all](examples/all/main.go) | Fleet discovery with typed operations, chaining, conditions, retries, and result decoding |

```bash
cd examples/all
OSAPI_TOKEN="<jwt>" go run main.go
```

## Contributing

See the [Development](docs/development.md) guide for prerequisites, setup,
and conventions. See the [Contributing](docs/contributing.md) guide before
submitting a PR.

## License

The [MIT][] License.

[OSAPI]: https://github.com/osapi-io/osapi
[osapi-sdk]: https://github.com/osapi-io/osapi-sdk
[MIT]: LICENSE
