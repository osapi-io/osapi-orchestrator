[![release](https://img.shields.io/github/release/osapi-io/osapi-orchestrator.svg?style=for-the-badge)](https://github.com/osapi-io/osapi-orchestrator/releases/latest)
[![codecov](https://img.shields.io/codecov/c/github/osapi-io/osapi-orchestrator?style=for-the-badge)](https://codecov.io/gh/osapi-io/osapi-orchestrator)
[![go report card](https://goreportcard.com/badge/github.com/osapi-io/osapi-orchestrator?style=for-the-badge)](https://goreportcard.com/report/github.com/osapi-io/osapi-orchestrator)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/osapi-io/osapi-orchestrator/go.yml?style=for-the-badge)](https://github.com/osapi-io/osapi-orchestrator/actions/workflows/go.yml)
[![powered by](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![built with just](https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white)](https://just.systems)
![gitHub commit activity](https://img.shields.io/github/commit-activity/m/osapi-io/osapi-orchestrator?style=for-the-badge)
[![go reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/osapi-io/osapi-orchestrator/pkg/orchestrator)

# OSAPI Orchestrator

A Go package for orchestrating operations across [OSAPI][]-managed hosts --
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

## ⚙️ Operations

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

## ✨ Features

The orchestrator provides a declarative DSL for composing operations into
DAG-based plans with typed results, guards, retry, and discovery.

| Feature                                            | Description                                    |
| -------------------------------------------------- | ---------------------------------------------- |
| [Step Chaining](docs/features/basic.md)            | Sequential and parallel DAG execution          |
| [Guards](docs/features/guards.md)                  | Conditional execution (When, OnlyIfChanged...) |
| [Error Recovery](docs/features/error-recovery.md)  | Continue strategy with OnlyIfFailed cleanup    |
| [Broadcast](docs/features/broadcast.md)            | Per-host results from `_all`/label targets     |
| [Host Status](docs/features/guards.md)             | Skipped and failed detection per host          |
| [Retry](docs/features/retry.md)                    | Automatic retry with exponential backoff       |
| [Discovery](docs/features/discovery.md)            | Find agents by OS, arch, labels, conditions    |
| [File Workflow](docs/features/file-workflow.md)    | Upload, deploy, drift detection, undeploy      |
| [Result Decode](docs/features/result-decode.md)    | Typed struct decoding from step results        |
| [TaskFunc](docs/features/task-func.md)             | Custom logic with access to prior results      |

See the [DSL reference](docs/features/README.md) for guards, predicates, error
strategies, and typed result tables.

## 📖 Documentation

See the [package documentation][] on pkg.go.dev for API details.

## 📋 Examples

Runnable examples in [examples/operations/](examples/operations/) (per-domain
workflows) and [examples/features/](examples/features/) (DSL features). Run
with:

    OSAPI_TOKEN="<jwt>" go run examples/features/basic.go

## 🤝 Contributing

See the [Development](docs/development.md) guide for prerequisites, setup,
and conventions. See the [Contributing](docs/contributing.md) guide before
submitting a PR.

## 📄 License

The [MIT][] License.

[OSAPI]: https://github.com/osapi-io/osapi
[osapi-sdk]: https://github.com/osapi-io/osapi/tree/main/pkg/sdk
[package documentation]: https://pkg.go.dev/github.com/osapi-io/osapi-orchestrator/pkg/orchestrator
[MIT]: LICENSE
