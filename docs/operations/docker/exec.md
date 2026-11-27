# DockerExec

Run a command in a container.

Runs docker exec, and answers what the command wrote to standard output and to
standard error, and its exit code. A command that exits with a code other than 0
is an answer, not a failure. It changes what the command changes, which osapi
cannot tell, so it reports changed: true. Its permission is docker:exec, which
only the admin role gets, apart from the docker:execute of start and stop. An
agent with agent.command.enabled false refuses it, as it refuses a command.

## Usage

```go
step := o.DockerExec("web-01", input)
```

## Parameters

| Parameter | Type                    | Description                                                   |
| --------- | ----------------------- | ------------------------------------------------------------- |
| `target`  | `string`                | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.DockerExecInput` | Refer to `DockerExecInput` below.                             |

### DockerExecInput

| Field        | Type                | Description                                                                                         |
| ------------ | ------------------- | --------------------------------------------------------------------------------------------------- |
| `Name`       | `string`            | The name of the path: the container's name or ID. It names the item in the path.                    |
| `Command`    | `[]string`          | The command and its arguments.                                                                      |
| `Env`        | `map[string]string` | Environment variables for the command. A secret: a reference, secret://<vault>/<key>, or the value. |
| `WorkingDir` | `string`            | The directory the command runs in.                                                                  |

## Result

Each host answers with one `osapi.DockerExecResult` row. A list answers with one
row for each item.

```go
var row osapi.DockerExecResult
err := report.Decode("exec-docker", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Always true: osapi cannot tell what the command changed.                                                                                                                                                                 |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `ExitCode` | `int`             | The command's exit code.                                                                                                                                                                                                 |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Stderr`   | `string`          | What the command wrote to standard error.                                                                                                                                                                                |
| `Stdout`   | `string`          | What the command wrote to standard output.                                                                                                                                                                               |

## Retries

A second run is not the same as one run. So `Retry` does not run this step again
unless the step adds `RetryUnsafe`.

## Permissions

The token must have the `docker:exec` permission. The step sends
`POST /api/node/{hostname}/docker/{name}/exec`.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go).
