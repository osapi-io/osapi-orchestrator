# Operations

The orchestrator provides 101 typed constructors organized by domain. Each
method returns a `*Step` that can be chained with ordering, conditions, and
error handling.

## Domains

| Domain                        | Operations | Description                                       |
| ----------------------------- | ---------- | ------------------------------------------------- |
| [Node](node/)                 | 8          | Hostname, disk, memory, load, uptime, OS, status  |
| [Network](network/)           | 4          | DNS configuration, ping                           |
| [Interface](interface/)       | 5          | Network interface configuration via Netplan       |
| [Route](route/)               | 5          | Network route configuration via Netplan           |
| [Command](command/)           | 2          | Direct exec, shell commands                       |
| [Docker](docker/)             | 9          | Container lifecycle, exec, images                 |
| [Cron](cron/)                 | 5          | Cron drop-in file management                      |
| [Sysctl](sysctl/)             | 5          | Kernel parameter management                       |
| [NTP](ntp/)                   | 4          | NTP server configuration                          |
| [Timezone](timezone/)         | 2          | System timezone management                        |
| [Service](service/)           | 10         | Systemd service lifecycle and unit files          |
| [Package](package/)           | 6          | System package management                         |
| [User](user/)                 | 9          | User accounts, SSH keys, passwords                |
| [Group](group/)               | 5          | Local group management                            |
| [Certificate](certificate/)   | 4          | CA certificate trust store management             |
| [Process](process/)           | 3          | Process listing, inspection, signals              |
| [Power](power/)               | 2          | Reboot and shutdown                               |
| [Log](log/)                   | 3          | Systemd journal queries                           |
| [File](file/)                 | 5          | Upload, deploy, status, undeploy, drift detection |
| [Agent](agent/)               | 4          | Discovery, inspection, drain/undrain              |
| [Health](health/)             | 1          | Liveness check                                    |

## Idempotency

- **Read-only** operations never modify state and always return
  `Changed: false`.
- **Idempotent** write operations check current state before mutating and return
  `Changed: true` only if something actually changed.
- **Non-idempotent** operations (command exec/shell) always return
  `Changed: true`. Use guards (`When`, `OnlyIfChanged`) to control when they
  run.
