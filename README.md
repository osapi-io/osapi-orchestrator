<p align="center">
  <picture>
    <source srcset="asset/logo-dark.svg" media="(prefers-color-scheme: dark)">
    <source srcset="asset/logo-light.svg" media="(prefers-color-scheme: light)">
    <img src="asset/logo-dark.svg" alt="osapi-orchestrator" width="929">
  </picture>
</p>

<p align="center">Compose operations across OSAPI-managed hosts as one plan.</p>

<p align="center">
  <a href="https://github.com/osapi-io/osapi-orchestrator/releases/latest"><img alt="release" src="https://img.shields.io/github/release/osapi-io/osapi-orchestrator.svg?style=for-the-badge"></a>
  <a href="https://codecov.io/gh/osapi-io/osapi-orchestrator"><img alt="codecov" src="https://img.shields.io/codecov/c/github/osapi-io/osapi-orchestrator?style=for-the-badge"></a>
  <a href="LICENSE"><img alt="license" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
  <a href="https://github.com/osapi-io/osapi-orchestrator/actions/workflows/go.yml"><img alt="build" src="https://img.shields.io/github/actions/workflow/status/osapi-io/osapi-orchestrator/go.yml?style=for-the-badge"></a>
  <a href="https://github.com/goreleaser"><img alt="powered by" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge"></a>
  <a href="https://conventionalcommits.org"><img alt="conventional commits" src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge"></a>
  <a href="https://just.systems"><img alt="built with just" src="https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white"></a>
  <img alt="gitHub commit activity" src="https://img.shields.io/github/commit-activity/m/osapi-io/osapi-orchestrator?style=for-the-badge">
  <a href="https://pkg.go.dev/github.com/osapi-io/osapi-orchestrator"><img alt="go reference" src="https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white"></a>
  <a href="https://pkg.go.dev/github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"><img alt="go.dev reference" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge"></a>
</p>

<p align="center">
<b>Describe the plan, not the order.</b>
</p>

<p align="center">
Declare the work and the dependencies between the steps. The orchestrator
finds the order and runs independent steps in parallel. It gives one result
for each host, not one answer for the fleet.
</p>

## 📦 Install

```bash
go install github.com/osapi-io/osapi-orchestrator@latest
```

To use it as a library dependency:

```bash
go get github.com/osapi-io/osapi-orchestrator
```

## ✨ Features

The orchestrator has a declarative DSL. You use it to put operations into
DAG-based plans with typed results, guards, retry, and discovery.

| Feature                                           | Description                                    |
| ------------------------------------------------- | ---------------------------------------------- |
| [Step Chaining](docs/features/basic.md)           | Sequential and parallel DAG runs               |
| [Guards](docs/features/guards.md)                 | Conditional steps (When, OnlyIfChanged...)     |
| [Error Recovery](docs/features/error-recovery.md) | Continue strategy with OnlyIfFailed cleanup    |
| [Broadcast](docs/features/broadcast.md)           | A result for each host of `_all`/label targets |
| [Host Status](docs/features/guards.md)            | Skipped and failed status for each host        |
| [Retry](docs/features/retry.md)                   | Automatic retry with exponential backoff       |
| [Discovery](docs/features/discovery.md)           | Find agents by OS, arch, labels, conditions    |
| [File Workflow](docs/features/file-workflow.md)   | Upload, deploy, drift check, undeploy          |
| [Result Decode](docs/features/result-decode.md)   | Decode step results into typed structs         |
| [TaskFunc](docs/features/task-func.md)            | Custom logic that can read earlier results     |

Refer to the [DSL reference](docs/features/README.md) for guards, predicates,
error strategies, and typed result tables.

## 🧪 Examples

[examples/operations/](examples/operations/) has a runnable workflow for each
domain. [examples/features/](examples/features/) has runnable examples of the
DSL features. Run an example with:

```
OSAPI_TOKEN="<jwt>" go run examples/features/basic.go
```

## 📚 Documentation

- [Package documentation] on pkg.go.dev. The API reference.

## 🎯 Targeting

Most operations accept a `target` parameter. It controls which agents receive
the request:

| Target      | Behavior                                    |
| ----------- | ------------------------------------------- |
| `_any`      | Send to one available agent (load balanced) |
| `_all`      | Broadcast to every agent                    |
| `hostname`  | Send to a specific host                     |
| `key:value` | Send to the agents that have the label      |

## ⚙️ Operations

100 typed constructors, in groups by the service of the osapi SDK:

| Service        | Operations                                       | Example                                                      |
| -------------- | ------------------------------------------------ | ------------------------------------------------------------ |
| Agent          | [4 operations](docs/operations/agent/)           | [agent-drain.go](examples/operations/agent-drain.go)         |
| Certificate    | [5 operations](docs/operations/certificate/)     | [certificate.go](examples/operations/certificate.go)         |
| Command        | [2 operations](docs/operations/command/)         | [command.go](examples/operations/command.go)                 |
| Connectivity   | [1 operation](docs/operations/connectivity/)     | [ping.go](examples/operations/ping.go)                       |
| Disk           | [1 operation](docs/operations/disk/)             | [node-info.go](examples/operations/node-info.go)             |
| DNS            | [3 operations](docs/operations/dns/)             | [dns-update.go](examples/operations/dns-update.go)           |
| Docker         | [7 operations](docs/operations/docker/)          | [docker.go](examples/operations/docker.go)                   |
| File           | [2 operations](docs/operations/files/file/)      | [file-changed.go](examples/operations/file-changed.go)       |
| Group          | [5 operations](docs/operations/group/)           | [group.go](examples/operations/group.go)                     |
| Health         | [1 operation](docs/operations/health/)           | [basic.go](examples/features/basic.go)                       |
| Hostname       | [2 operations](docs/operations/hostname/)        | [hostname-update.go](examples/operations/hostname-update.go) |
| Image          | [2 operations](docs/operations/image/)           | [docker.go](examples/operations/docker.go)                   |
| Interface      | [5 operations](docs/operations/interface/)       | [interface.go](examples/operations/interface.go)             |
| InterfaceRoute | [4 operations](docs/operations/interface-route/) | [route.go](examples/operations/route.go)                     |
| Load           | [1 operation](docs/operations/load/)             | [node-info.go](examples/operations/node-info.go)             |
| Log            | [1 operation](docs/operations/log/)              | [log.go](examples/operations/log.go)                         |
| LogSource      | [1 operation](docs/operations/log-source/)       | [log.go](examples/operations/log.go)                         |
| Memory         | [1 operation](docs/operations/memory/)           | [node-info.go](examples/operations/node-info.go)             |
| NodeFile       | [3 operations](docs/operations/node-file/)       | [file-deploy.go](examples/operations/file-deploy.go)         |
| NTP            | [4 operations](docs/operations/ntp/)             | [ntp.go](examples/operations/ntp.go)                         |
| OS             | [1 operation](docs/operations/os/)               | [node-info.go](examples/operations/node-info.go)             |
| Package        | [4 operations](docs/operations/package/)         | [package.go](examples/operations/package.go)                 |
| PackageUpdates | [2 operations](docs/operations/package-updates/) | [package.go](examples/operations/package.go)                 |
| Power          | [2 operations](docs/operations/power/)           | [power.go](examples/operations/power.go)                     |
| Process        | [3 operations](docs/operations/process/)         | [process.go](examples/operations/process.go)                 |
| Schedule       | [5 operations](docs/operations/schedule/)        | [schedule.go](examples/operations/schedule.go)               |
| Service        | [10 operations](docs/operations/service/)        | [service.go](examples/operations/service.go)                 |
| Status         | [1 operation](docs/operations/status/)           | [node-info.go](examples/operations/node-info.go)             |
| Sysctl         | [5 operations](docs/operations/sysctl/)          | [sysctl.go](examples/operations/sysctl.go)                   |
| Timezone       | [2 operations](docs/operations/timezone/)        | [timezone.go](examples/operations/timezone.go)               |
| Uptime         | [1 operation](docs/operations/uptime/)           | [node-info.go](examples/operations/node-info.go)             |
| User           | [5 operations](docs/operations/user/)            | [user.go](examples/operations/user.go)                       |
| UserPassword   | [1 operation](docs/operations/user-password/)    | [user.go](examples/operations/user.go)                       |
| UserSSHKey     | [3 operations](docs/operations/user-ssh-key/)    | [user.go](examples/operations/user.go)                       |

## 🤝 Contributing

Refer to the [Contributing](CONTRIBUTING.md) guide for the prerequisites, the
setup, the conventions, and the PR workflow.

## 📄 License

The [MIT] License.

[mit]: LICENSE
[package documentation]: https://pkg.go.dev/github.com/osapi-io/osapi-orchestrator/pkg/orchestrator
