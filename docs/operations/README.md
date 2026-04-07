# Operations

The orchestrator provides 101 typed constructors organized by domain. Each
method returns a `*Step` that can be chained with ordering, conditions, and
error handling.

## Groups

### [Services](services/)

| Domain                            | Operations | Description                              |
| --------------------------------- | ---------- | ---------------------------------------- |
| [Service](services/service/)      | 10         | Systemd service lifecycle and unit files |
| [Cron](services/cron/)            | 5          | Cron drop-in file management             |

### [Software](software/)

| Domain                            | Operations | Description                  |
| --------------------------------- | ---------- | ---------------------------- |
| [Package](software/package/)      | 6          | System package management    |

### [Config](config/)

| Domain                            | Operations | Description                     |
| --------------------------------- | ---------- | ------------------------------- |
| [Hostname](config/hostname/)      | 2          | Hostname query and update       |
| [Sysctl](config/sysctl/)          | 5          | Kernel parameter management     |
| [NTP](config/ntp/)                | 4          | NTP server configuration        |
| [Timezone](config/timezone/)      | 2          | System timezone management      |

### [Node](node/)

| Domain                            | Operations | Description                              |
| --------------------------------- | ---------- | ---------------------------------------- |
| [Power](node/power/)              | 2          | Reboot and shutdown                      |
| [Process](node/process/)          | 3          | Process listing, inspection, signals     |
| [Log](node/log/)                  | 3          | Systemd journal queries                  |
| Node                              | 5          | Status, load, uptime, OS                 |

### [Networking](networking/)

| Domain                            | Operations | Description                          |
| --------------------------------- | ---------- | ------------------------------------ |
| [Network](networking/network/)    | 4          | DNS configuration, ping              |
| [Interface](networking/interface/)| 5          | Network interface config via Netplan |
| [Route](networking/route/)        | 5          | Network route config via Netplan     |

### [Security](security/)

| Domain                            | Operations | Description                          |
| --------------------------------- | ---------- | ------------------------------------ |
| [User](security/user/)            | 9          | User accounts, SSH keys, passwords   |
| [Group](security/group/)          | 5          | Local group management               |
| [Certificate](security/certificate/) | 4      | CA certificate trust store           |

### [Containers](containers/)

| Domain                            | Operations | Description                     |
| --------------------------------- | ---------- | ------------------------------- |
| [Docker](containers/docker/)      | 9          | Container lifecycle, exec, images |

### [Files](files/)

| Domain                            | Operations | Description                              |
| --------------------------------- | ---------- | ---------------------------------------- |
| [File](files/file/)               | 5          | Upload, deploy, status, undeploy, drift  |

### [Command](cmd/)

| Domain                            | Operations | Description                     |
| --------------------------------- | ---------- | ------------------------------- |
| [Command](cmd/command/)           | 2          | Direct exec, shell commands     |

### [Hardware](hardware/)

| Domain                            | Operations | Description                     |
| --------------------------------- | ---------- | ------------------------------- |
| Hardware                          | 2          | Disk and memory queries         |

### [Agent](agent/)

| Domain                            | Operations | Description                     |
| --------------------------------- | ---------- | ------------------------------- |
| [Agent](agent/)                   | 4          | Discovery, inspection, drain    |

### [Health](health/)

| Domain                            | Operations | Description                     |
| --------------------------------- | ---------- | ------------------------------- |
| [Health](health/)                 | 1          | Liveness check                  |

## Idempotency

- **Read-only** operations never modify state and always return
  `Changed: false`.
- **Idempotent** write operations check current state before mutating and return
  `Changed: true` only if something actually changed.
- **Non-idempotent** operations (command exec/shell) always return
  `Changed: true`. Use guards (`When`, `OnlyIfChanged`) to control when they
  run.
