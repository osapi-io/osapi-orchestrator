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

A Go package for orchestrating operations across [OSAPI][]-managed hosts --
typed operations, chaining, conditions, and result decoding built on top
of the [osapi-sdk][] engine.

## Install

```bash
go install github.com/osapi-io/osapi-orchestrator@latest
```

As a library dependency:

```bash
go get github.com/osapi-io/osapi-orchestrator
```

## Targeting

Most operations accept a `target` parameter to control which agents receive
the request:

| Target      | Behavior                                    |
| ----------- | ------------------------------------------- |
| `_any`      | Send to any available agent (load balanced) |
| `_all`      | Broadcast to every agent                    |
| `hostname`  | Send to a specific host                     |
| `key:value` | Send to agents matching a label             |

## Operations

37 typed constructors across 8 domains:

| Domain | Docs | Example |
| ------ | ---- | ------- |
| Node | [8 operations](docs/operations/node/) | [node-info.go](examples/operations/node-info.go) |
| Network | [3 operations](docs/operations/network/) | [dns-update.go](examples/operations/dns-update.go) |
| Command | [2 operations](docs/operations/command/) | [command.go](examples/operations/command.go) |
| Docker | [9 operations](docs/operations/docker/) | [docker.go](examples/operations/docker.go) |
| Cron | [5 operations](docs/operations/cron/) | [cron.go](examples/operations/cron.go) |
| File | [5 operations](docs/operations/file/) | [file-deploy.go](examples/operations/file-deploy.go) |
| Agent | [4 operations](docs/operations/agent/) | [agent-drain.go](examples/operations/agent-drain.go) |
| Health | [1 operation](docs/operations/health/) | [basic.go](examples/features/basic.go) |

## Features

- [Feature Guides](docs/features/README.md) -- Step chaining, guards,
  retry, broadcast, discovery, file workflows, host status awareness,
  and result decoding
- [API Reference](docs/gen/orchestrator.md) -- Auto-generated Go
  documentation

## Examples

Each example is a standalone Go file. Run with:

    cd examples/features
    OSAPI_TOKEN="<jwt>" go run basic.go

### Feature Examples

| Example                                                          | What it shows                                       |
| ---------------------------------------------------------------- | --------------------------------------------------- |
| [basic.go](examples/features/basic.go)                          | Simple DAG with health check and hostname query     |
| [parallel.go](examples/features/parallel.go)                    | Five parallel queries depending on health check     |
| [retry.go](examples/features/retry.go)                          | Retry on failure with configurable attempts         |
| [verbose.go](examples/features/verbose.go)                      | Verbose output with stdout/stderr/response data     |
| [guards.go](examples/features/guards.go)                        | When predicate for conditional execution            |
| [only-if-changed.go](examples/features/only-if-changed.go)      | Skip step unless dependency reported changes        |
| [error-recovery.go](examples/features/error-recovery.go)        | Continue strategy with OnlyIfFailed cleanup         |
| [broadcast.go](examples/features/broadcast.go)                  | Per-host results from broadcast operations          |
| [task-func.go](examples/features/task-func.go)                  | Custom steps with typed result decoding             |
| [agent-facts.go](examples/features/agent-facts.go)              | List agents with OS, load, memory, and interfaces   |
| [discover.go](examples/features/discover.go)                    | Find agents by OS and architecture predicates       |
| [group-by-fact.go](examples/features/group-by-fact.go)          | Group agents by distro, run per-group commands      |
| [when-fact.go](examples/features/when-fact.go)                  | Fact-based guard on a step                          |
| [fact-predicates.go](examples/features/fact-predicates.go)      | Compose multiple predicates for discovery           |
| [label-filter.go](examples/features/label-filter.go)            | Filter by labels and arbitrary fact values           |
| [condition-filter.go](examples/features/condition-filter.go)    | Filter by node conditions (e.g., DiskPressure)      |
| [host-status.go](examples/features/host-status.go)              | Host status guards (skipped and failed detection)   |
| [broadcast-guards.go](examples/features/broadcast-guards.go)    | Broadcast guards with per-host error and changed    |

### Operation Examples

| Example                                                          | What it shows                                       |
| ---------------------------------------------------------------- | --------------------------------------------------- |
| [command.go](examples/operations/command.go)                    | Command exec and shell with result decoding         |
| [dns-update.go](examples/operations/dns-update.go)              | Read-then-write pattern with DNS operations         |
| [file-deploy.go](examples/operations/file-deploy.go)            | Upload, deploy, and verify a file end-to-end        |
| [file-changed.go](examples/operations/file-changed.go)          | Conditional upload with FileChanged + OnlyIfChanged |
| [hostname-update.go](examples/operations/hostname-update.go)    | Read-then-write pattern with hostname broadcast     |
| [docker.go](examples/operations/docker.go)                      | Full Docker lifecycle with pull, create, exec       |
| [cron.go](examples/operations/cron.go)                          | Cron create, list, and delete lifecycle             |

## 🤝 Contributing

See the [Development](docs/development.md) guide for prerequisites, setup,
and conventions. See the [Contributing](docs/contributing.md) guide before
submitting a PR.

## 📄 License

The [MIT][] License.

[OSAPI]: https://github.com/osapi-io/osapi
[osapi-sdk]: https://github.com/osapi-io/osapi/tree/main/pkg/sdk
[MIT]: LICENSE
