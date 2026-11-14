# Operations

The orchestrator has 101 typed constructors, in groups by domain. Each method
returns a `*Step`. You can chain methods on the step to set the order, the
conditions, and what to do on an error.

## Groups

### [Services](services/)

| Domain                       | Operations | Description                     |
| ---------------------------- | ---------- | ------------------------------- |
| [Service](services/service/) | 10         | Systemd services and unit files |
| [Cron](services/cron/)       | 5          | Cron drop-in files              |

### [Software](software/)

| Domain                       | Operations | Description     |
| ---------------------------- | ---------- | --------------- |
| [Package](software/package/) | 6          | System packages |

### [Config](config/)

| Domain                       | Operations | Description                 |
| ---------------------------- | ---------- | --------------------------- |
| [Hostname](config/hostname/) | 2          | Get and update the hostname |
| [Sysctl](config/sysctl/)     | 5          | Kernel parameters           |
| [NTP](config/ntp/)           | 4          | NTP server configuration    |
| [Timezone](config/timezone/) | 2          | System timezone             |

### [Node](node/)

| Domain                   | Operations | Description                              |
| ------------------------ | ---------- | ---------------------------------------- |
| [Power](node/power/)     | 2          | Reboot and shutdown                      |
| [Process](node/process/) | 3          | List and examine processes, send signals |
| [Log](node/log/)         | 3          | Query the systemd journal                |
| Node                     | 4          | Status, load, uptime, OS                 |

### [Networking](networking/)

| Domain                             | Operations | Description                                  |
| ---------------------------------- | ---------- | -------------------------------------------- |
| [Network](networking/network/)     | 4          | DNS configuration, ping                      |
| [Interface](networking/interface/) | 5          | Network interface configuration with Netplan |
| [Route](networking/route/)         | 5          | Network route configuration with Netplan     |

### [Security](security/)

| Domain                               | Operations | Description                        |
| ------------------------------------ | ---------- | ---------------------------------- |
| [User](security/user/)               | 9          | User accounts, SSH keys, passwords |
| [Group](security/group/)             | 5          | Local groups                       |
| [Certificate](security/certificate/) | 4          | CA certificate trust store         |

### [Containers](containers/)

| Domain                       | Operations | Description                       |
| ---------------------------- | ---------- | --------------------------------- |
| [Docker](containers/docker/) | 9          | Container lifecycle, exec, images |

### [Files](files/)

| Domain              | Operations | Description                             |
| ------------------- | ---------- | --------------------------------------- |
| [File](files/file/) | 5          | Upload, deploy, status, undeploy, drift |

### [Command](cmd/)

| Domain                  | Operations | Description                    |
| ----------------------- | ---------- | ------------------------------ |
| [Command](cmd/command/) | 2          | Direct exec and shell commands |

### [Hardware](hardware/)

| Domain   | Operations | Description           |
| -------- | ---------- | --------------------- |
| Hardware | 2          | Query disk and memory |

### [Agent](agent/)

| Domain          | Operations | Description                         |
| --------------- | ---------- | ----------------------------------- |
| [Agent](agent/) | 4          | Discover, examine, and drain agents |

### [Health](health/)

| Domain            | Operations | Description    |
| ----------------- | ---------- | -------------- |
| [Health](health/) | 1          | Liveness check |

## Idempotency

- **Read-only.** The operation does not change the host. It always returns
  `Changed: false`.
- **Idempotent.** A write operation compares the current state with the state
  you give. It returns `Changed: true` only if it changed something.
- **Non-idempotent.** The operation (for example, command exec and shell) always
  returns `Changed: true`. Use guards (`When`, `OnlyIfChanged`) to control when
  the step runs.
