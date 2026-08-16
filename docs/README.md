# Documentation

Reference documentation for the orchestrator. Runnable programs live in
[`examples/`](../examples); contributor setup and conventions are in
[CONTRIBUTING.md](../CONTRIBUTING.md).

| Document                                     | Covers                                                                   |
| -------------------------------------------- | ------------------------------------------------------------------------ |
| [features/README.md](features/README.md)     | The DSL — how the DAG works, guards, retry, discovery, and typed results |
| [operations/README.md](operations/README.md) | Every typed operation constructor, grouped by domain                     |

Operations are grouped by domain: `agent`, `cmd`, `config`, `containers`,
`files`, `hardware`, `health`, `networking`, `node`, `security`, `services`, and
`software`.
