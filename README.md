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

## Features

- [Operations](docs/operations/README.md) -- 18 typed constructors for every
  OSAPI operation
- [Features](docs/features/README.md) -- Step chaining, guards, retry,
  broadcast, discovery, file workflows, and result decoding
- [API Reference](docs/gen/orchestrator.md) -- Auto-generated Go documentation

## Examples

Each example is a standalone Go file. Run with:

    cd examples
    OSAPI_TOKEN="<jwt>" go run basic.go

| Example                                        | What it shows                                       |
| ---------------------------------------------- | --------------------------------------------------- |
| [basic.go](examples/basic.go)                  | Simple DAG with health check and hostname query     |
| [parallel.go](examples/parallel.go)            | Five parallel queries depending on health check     |
| [retry.go](examples/retry.go)                  | Retry on failure with configurable attempts         |
| [command.go](examples/command.go)              | Command exec and shell with result decoding         |
| [verbose.go](examples/verbose.go)              | Verbose output with stdout/stderr/response data     |
| [guards.go](examples/guards.go)                | When predicate for conditional execution            |
| [only-if-changed.go](examples/only-if-changed.go) | Skip step unless dependency reported changes     |
| [error-recovery.go](examples/error-recovery.go) | Continue strategy with OnlyIfFailed cleanup        |
| [broadcast.go](examples/broadcast.go)          | Per-host results from broadcast operations          |
| [task-func.go](examples/task-func.go)          | Custom steps with typed result decoding             |
| [dns-update.go](examples/dns-update.go)        | Read-then-write pattern with DNS operations         |
| [file-deploy.go](examples/file-deploy.go)      | Upload, deploy, and verify a file end-to-end        |
| [file-changed.go](examples/file-changed.go)    | Conditional upload with FileChanged + OnlyIfChanged |
| [agent-facts.go](examples/agent-facts.go)      | List agents with OS, load, memory, and interfaces   |
| [discover.go](examples/discover.go)            | Find agents by OS and architecture predicates       |
| [group-by-fact.go](examples/group-by-fact.go)  | Group agents by distro, run per-group commands      |
| [when-fact.go](examples/when-fact.go)          | Fact-based guard on a step                          |
| [fact-predicates.go](examples/fact-predicates.go) | Compose multiple predicates for discovery        |
| [label-filter.go](examples/label-filter.go)    | Filter by labels and arbitrary fact values          |
| [condition-filter.go](examples/condition-filter.go) | Filter by node conditions (e.g., DiskPressure) |

## Contributing

See the [Development](docs/development.md) guide for prerequisites, setup,
and conventions. See the [Contributing](docs/contributing.md) guide before
submitting a PR.

## License

The [MIT][] License.

[OSAPI]: https://github.com/osapi-io/osapi
[osapi-sdk]: https://github.com/osapi-io/osapi/tree/main/pkg/sdk
[MIT]: LICENSE
