# Operations

The orchestrator provides 37 typed constructors organized by domain.
Each method returns a `*Step` that can be chained with ordering,
conditions, and error handling.

## Domains

| Domain | Operations | Description |
| ------ | ---------- | ----------- |
| [Node](node/) | 8 | Hostname, disk, memory, load, uptime, OS, status |
| [Network](network/) | 3 | DNS configuration, ping |
| [Command](command/) | 2 | Direct exec, shell commands |
| [Docker](docker/) | 9 | Container lifecycle, exec, images |
| [Cron](cron/) | 5 | Cron drop-in file management |
| [File](file/) | 5 | Upload, deploy, status, undeploy, drift detection |
| [Agent](agent/) | 4 | Discovery, inspection, drain/undrain |
| [Health](health/) | 1 | Liveness check |

## Idempotency

- **Read-only** operations never modify state and always return
  `Changed: false`.
- **Idempotent** write operations check current state before mutating
  and return `Changed: true` only if something actually changed.
- **Non-idempotent** operations (command exec/shell) always return
  `Changed: true`. Use guards (`When`, `OnlyIfChanged`) to control
  when they run.
