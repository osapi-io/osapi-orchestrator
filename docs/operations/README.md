# Operations

The orchestrator has 100 typed constructors. Each returns a `*Step`. You can
chain methods on the step to set the order, the conditions, and what to do on an
error.

Each operation that targets hosts is one call of the osapi SDK, with the same
name and the same parameters: `o.ServiceStart("web-01", "nginx")` calls
`ServiceStart` of the SDK's `Service` service. The SDK's own docs have the
details of each call.

## The controller

| Operation                            | Description                                          |
| ------------------------------------ | ---------------------------------------------------- |
| [HealthCheck](health/check.md)       | Check that the controller is alive                   |
| [AgentList](agent/list.md)           | List the agents                                      |
| [AgentGet](agent/get.md)             | Get one agent                                        |
| [AgentDrain](agent/drain.md)         | Drain an agent, so that it takes no new jobs         |
| [AgentUndrain](agent/undrain.md)     | Undrain an agent                                     |
| [FileUpload](files/file/upload.md)   | Upload content to the Object Store                   |
| [FileChanged](files/file/changed.md) | Check whether content differs from the stored object |

## The hosts

### Hostname

| Operation                            | Description                  |
| ------------------------------------ | ---------------------------- |
| [HostnameGet](hostname/get.md)       | Get the hostname of the host |
| [HostnameUpdate](hostname/update.md) | Set the hostname of the host |

### Status

| Operation                  | Description                        |
| -------------------------- | ---------------------------------- |
| [StatusGet](status/get.md) | Get the status summary of the host |

### Uptime

| Operation                  | Description                |
| -------------------------- | -------------------------- |
| [UptimeGet](uptime/get.md) | Get the uptime of the host |

### Disk

| Operation              | Description                    |
| ---------------------- | ------------------------------ |
| [DiskGet](disk/get.md) | Get the disk usage of the host |

### Memory

| Operation                  | Description                      |
| -------------------------- | -------------------------------- |
| [MemoryGet](memory/get.md) | Get the memory usage of the host |

### Load

| Operation              | Description                       |
| ---------------------- | --------------------------------- |
| [LoadGet](load/get.md) | Get the load averages of the host |

### OS

| Operation          | Description                          |
| ------------------ | ------------------------------------ |
| [OSGet](os/get.md) | Get the operating system of the host |

### DNS

| Operation                  | Description                                      |
| -------------------------- | ------------------------------------------------ |
| [DNSGet](dns/get.md)       | Get the DNS resolvers of an interface            |
| [DNSUpdate](dns/update.md) | Set the DNS resolvers of an interface            |
| [DNSDelete](dns/delete.md) | Drop the DNS resolvers osapi set on an interface |

### Connectivity

| Operation                                | Description                   |
| ---------------------------------------- | ----------------------------- |
| [ConnectivityPing](connectivity/ping.md) | Ping an address from the host |

### Interface

| Operation                              | Description                                        |
| -------------------------------------- | -------------------------------------------------- |
| [InterfaceList](interface/list.md)     | List the network interfaces of the host            |
| [InterfaceGet](interface/get.md)       | Get a network interface of the host                |
| [InterfaceCreate](interface/create.md) | Configure a network interface                      |
| [InterfaceUpdate](interface/update.md) | Change the config osapi has of a network interface |
| [InterfaceDelete](interface/delete.md) | Drop the config osapi has of a network interface   |

### InterfaceRoute

| Operation                                         | Description                                 |
| ------------------------------------------------- | ------------------------------------------- |
| [InterfaceRouteList](interface-route/list.md)     | List the routes of an interface             |
| [InterfaceRouteCreate](interface-route/create.md) | Add a static route to an interface          |
| [InterfaceRouteUpdate](interface-route/update.md) | Change the gateway or the metric of a route |
| [InterfaceRouteDelete](interface-route/delete.md) | Drop a route that osapi wrote               |

### NodeFile

| Operation                             | Description                                         |
| ------------------------------------- | --------------------------------------------------- |
| [NodeFileGet](node-file/get.md)       | Show whether a deployed file holds what osapi wrote |
| [NodeFileUpdate](node-file/update.md) | Deploy an object of the object store to a path      |
| [NodeFileDelete](node-file/delete.md) | Remove a file that osapi deployed                   |

### Docker

| Operation                        | Description                     |
| -------------------------------- | ------------------------------- |
| [DockerList](docker/list.md)     | List the containers of the host |
| [DockerGet](docker/get.md)       | Get one container               |
| [DockerCreate](docker/create.md) | Create a container              |
| [DockerStart](docker/start.md)   | Start a container               |
| [DockerStop](docker/stop.md)     | Stop a container                |
| [DockerExec](docker/exec.md)     | Run a command in a container    |
| [DockerDelete](docker/delete.md) | Remove a container              |

### Image

| Operation                      | Description     |
| ------------------------------ | --------------- |
| [ImagePull](image/pull.md)     | Pull an image   |
| [ImageDelete](image/delete.md) | Remove an image |

### Schedule

| Operation                            | Description                                   |
| ------------------------------------ | --------------------------------------------- |
| [ScheduleList](schedule/list.md)     | List the cron entries that osapi deployed     |
| [ScheduleGet](schedule/get.md)       | Get one cron entry that osapi deployed        |
| [ScheduleCreate](schedule/create.md) | Deploy a cron entry                           |
| [ScheduleUpdate](schedule/update.md) | Deploy a cron entry that osapi deployed again |
| [ScheduleDelete](schedule/delete.md) | Remove a cron entry that osapi deployed       |

### Sysctl

| Operation                        | Description                                                          |
| -------------------------------- | -------------------------------------------------------------------- |
| [SysctlList](sysctl/list.md)     | List the kernel's parameters and which osapi manages                 |
| [SysctlGet](sysctl/get.md)       | Get a kernel parameter's current value, and whether osapi manages it |
| [SysctlCreate](sysctl/create.md) | Manage a kernel parameter and apply it                               |
| [SysctlUpdate](sysctl/update.md) | Set a kernel parameter, manage it, and apply it                      |
| [SysctlDelete](sysctl/delete.md) | Stop the management of a kernel parameter                            |

### NTP

| Operation                  | Description                                   |
| -------------------------- | --------------------------------------------- |
| [NTPGet](ntp/get.md)       | Get the sync status of chrony and its sources |
| [NTPCreate](ntp/create.md) | Set the NTP servers that chrony uses          |
| [NTPUpdate](ntp/update.md) | Replace the NTP servers that osapi set        |
| [NTPDelete](ntp/delete.md) | Remove the NTP servers that osapi set         |

### Timezone

| Operation                            | Description                  |
| ------------------------------------ | ---------------------------- |
| [TimezoneGet](timezone/get.md)       | Get the timezone of the host |
| [TimezoneUpdate](timezone/update.md) | Set the timezone of the host |

### Service

| Operation                            | Description                                          |
| ------------------------------------ | ---------------------------------------------------- |
| [ServiceList](service/list.md)       | List the services of the host                        |
| [ServiceGet](service/get.md)         | Get one service                                      |
| [ServiceCreate](service/create.md)   | Deploy a service                                     |
| [ServiceUpdate](service/update.md)   | Deploy a new unit file for a service                 |
| [ServiceDelete](service/delete.md)   | Remove a service that osapi deployed                 |
| [ServiceStart](service/start.md)     | Start a service that osapi deployed                  |
| [ServiceStop](service/stop.md)       | Stop a service that osapi deployed                   |
| [ServiceRestart](service/restart.md) | Restart a service that osapi deployed                |
| [ServiceEnable](service/enable.md)   | Make a service that osapi deployed start at boot     |
| [ServiceDisable](service/disable.md) | Make a service that osapi deployed not start at boot |

### Package

| Operation                          | Description                             |
| ---------------------------------- | --------------------------------------- |
| [PackageList](package/list.md)     | List the packages installed on the host |
| [PackageGet](package/get.md)       | Get one installed package               |
| [PackageCreate](package/create.md) | Install a package                       |
| [PackageDelete](package/delete.md) | Remove a package                        |

### PackageUpdates

| Operation                                         | Description                                       |
| ------------------------------------------------- | ------------------------------------------------- |
| [PackageUpdatesList](package-updates/list.md)     | List the packages the index has newer versions of |
| [PackageUpdatesUpdate](package-updates/update.md) | Refresh the package index                         |

### User

| Operation                    | Description                |
| ---------------------------- | -------------------------- |
| [UserList](user/list.md)     | List the users of the host |
| [UserGet](user/get.md)       | Get one user               |
| [UserCreate](user/create.md) | Add a user                 |
| [UserUpdate](user/update.md) | Change a user              |
| [UserDelete](user/delete.md) | Remove a user              |

### UserPassword

| Operation                                     | Description                |
| --------------------------------------------- | -------------------------- |
| [UserPasswordUpdate](user-password/update.md) | Set the password of a user |

### UserSSHKey

| Operation                                  | Description                            |
| ------------------------------------------ | -------------------------------------- |
| [UserSSHKeyList](user-ssh-key/list.md)     | List the authorized SSH keys of a user |
| [UserSSHKeyCreate](user-ssh-key/create.md) | Add an authorized SSH key to a user    |
| [UserSSHKeyDelete](user-ssh-key/delete.md) | Remove an authorized SSH key of a user |

### Group

| Operation                      | Description                 |
| ------------------------------ | --------------------------- |
| [GroupList](group/list.md)     | List the groups of the host |
| [GroupGet](group/get.md)       | Get one group               |
| [GroupCreate](group/create.md) | Add a group                 |
| [GroupUpdate](group/update.md) | Set the members of a group  |
| [GroupDelete](group/delete.md) | Remove a group              |

### Certificate

| Operation                                  | Description                                             |
| ------------------------------------------ | ------------------------------------------------------- |
| [CertificateList](certificate/list.md)     | List the certificate authorities that the host trusts   |
| [CertificateGet](certificate/get.md)       | Get one certificate authority                           |
| [CertificateCreate](certificate/create.md) | Add a certificate authority to the trust store          |
| [CertificateUpdate](certificate/update.md) | Write a certificate authority that osapi deployed again |
| [CertificateDelete](certificate/delete.md) | Remove a certificate authority that osapi deployed      |

### Process

| Operation                          | Description                    |
| ---------------------------------- | ------------------------------ |
| [ProcessList](process/list.md)     | List the processes of the host |
| [ProcessGet](process/get.md)       | Get one process                |
| [ProcessSignal](process/signal.md) | Send a signal to a process     |

### Power

| Operation                          | Description        |
| ---------------------------------- | ------------------ |
| [PowerReboot](power/reboot.md)     | Reboot the host    |
| [PowerShutdown](power/shutdown.md) | Shut the host down |

### Log

| Operation              | Description                     |
| ---------------------- | ------------------------------- |
| [LogList](log/list.md) | List the entries of the journal |

### LogSource

| Operation                           | Description                     |
| ----------------------------------- | ------------------------------- |
| [LogSourceList](log-source/list.md) | List the sources of the journal |

### Command

| Operation                        | Description                           |
| -------------------------------- | ------------------------------------- |
| [CommandExec](command/exec.md)   | Run a program with its arguments      |
| [CommandShell](command/shell.md) | Run a command line through /bin/sh -c |
