# DockerCreate

Create a container.

Runs docker create, which pulls the image when the host does not have it. The
environment goes on docker's standard input, not in its arguments, which each
user of the host can read. auto_start starts the container once it is made. A
container of the name that the host has already is not made again, and reports
changed: false.

## Usage

```go
step := o.DockerCreate("web-01", input)
```

## Parameters

| Parameter | Type                      | Description                                                   |
| --------- | ------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                  | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.DockerCreateInput` | Refer to `DockerCreateInput` below.                           |

### DockerCreateInput

| Field               | Type                | Description                                                                                              |
| ------------------- | ------------------- | -------------------------------------------------------------------------------------------------------- |
| `AutoStart`         | `bool`              | Start the container once it is created.                                                                  |
| `Command`           | `[]string`          | The command, in place of the image's own.                                                                |
| `ContainerHostname` | `string`            | The hostname inside the container.                                                                       |
| `DNS`               | `[]string`          | The DNS servers of the container.                                                                        |
| `Env`               | `map[string]string` | The environment variables of the container. A secret: a reference, secret://<vault>/<key>, or the value. |
| `Image`             | `string`            | The image to run, such as alpine:3.20; docker pulls it when the host does not have it.                   |
| `Name`              | `string`            | The container's name.                                                                                    |
| `Ports`             | `[]string`          | The published ports, each host:container, such as 8080:80.                                               |
| `Volumes`           | `[]string`          | The bind mounts, each host:container, such as /srv/data:/data.                                           |

## Result

Each host answers with one `osapi.DockerResult` row. A list answers with one row
for each item.

```go
var row osapi.DockerResult
err := report.Decode("create-docker", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `ID`       | `string`          | The container's full ID, for create.                                                                                                                                                                                     |
| `Name`     | `string`          | The container.                                                                                                                                                                                                           |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Retries

A second run is not the same as one run. So `Retry` does not run this step again
unless the step adds `RetryUnsafe`.

## Permissions

The token must have the `docker:write` permission. The step sends
`POST /api/node/{hostname}/docker`.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go).
