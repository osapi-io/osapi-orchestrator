# DockerList

List the containers of the host.

Reads each container through gohai's docker collector, which runs docker ps and
docker inspect. state limits the list to the running or the stopped ones. A host
whose docker does not answer reports skipped. The environment of a container
shows only the names of its variables: osapi never reads their values.

## Usage

```go
step := o.DockerList("web-01", params)
```

## Parameters

| Parameter | Type                     | Description                                                   |
| --------- | ------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                 | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `params`  | `osapi.DockerListParams` | Refer to `DockerListParams` below.                            |

### DockerListParams

| Field   | Type                   | Description                                                                                                                             |
| ------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| `State` | `types.ContainerState` | State is the state of the query: only the containers in this state: running, stopped or all, the default. The zero value leaves it out. |

## Result

Each host answers with one `osapi.DockerEntry` row. A list answers with one row
for each item.

```go
var row osapi.DockerEntry
err := report.Decode("list-docker", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Created`  | `string`          | When it was created, RFC 3339.                                                                                                                                                                                           |
| `EnvNames` | `[]string`        | The names of its environment variables; osapi never reads their values.                                                                                                                                                  |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Health`   | `string`          | The state of its health check: healthy, unhealthy or starting; empty without one.                                                                                                                                        |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `ID`       | `string`          | The container's full ID.                                                                                                                                                                                                 |
| `Image`    | `string`          | The image it runs.                                                                                                                                                                                                       |
| `Mounts`   | `[]string`        | The mounts, source:destination, with :ro when read only.                                                                                                                                                                 |
| `Name`     | `string`          | The container's name.                                                                                                                                                                                                    |
| `Networks` | `[]string`        | The networks, with the container's address and gateway on each.                                                                                                                                                          |
| `Ports`    | `[]string`        | The published ports, host:port->container port.                                                                                                                                                                          |
| `State`    | `string`          | running, exited, created, paused or restarting.                                                                                                                                                                          |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Summary`  | `string`          | What docker says of its state, such as Up 2 minutes.                                                                                                                                                                     |

## Permissions

The token must have the `docker:read` permission. The step sends
`GET /api/node/{hostname}/docker`.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go).
