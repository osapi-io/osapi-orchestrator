# DockerExec

Runs a command in a container on the target host. The container must be in the
`running` state.

## Usage

```go
step := o.DockerExec("web-01", "c1a2b3d4e5f6", osapi.DockerExecOpts{
    Command: []string{"nginx", "-t"},
})
```

## Parameters

| Parameter | Type                   | Description                                               |
| --------- | ---------------------- | --------------------------------------------------------- |
| `target`  | `string`               | Target host: `_any`, `_all`, hostname, or label selector. |
| `id`      | `string`               | The ID or name of the container that runs the command.    |
| `opts`    | `osapi.DockerExecOpts` | The options for the command.                              |

### DockerExecOpts fields

| Field        | Type       | Description                                      |
| ------------ | ---------- | ------------------------------------------------ |
| `Command`    | `[]string` | The command and its arguments (required).        |
| `Env`        | `[]string` | More environment variables, in KEY=VALUE format. |
| `WorkingDir` | `string`   | The working directory in the container.          |

## Result type

```go
var result osapi.DockerExecResult
err := results.Decode("docker-exec", &result)
```

| Field      | Type     | Description                                                  |
| ---------- | -------- | ------------------------------------------------------------ |
| `Stdout`   | `string` | The standard output of the command.                          |
| `Stderr`   | `string` | The standard error of the command.                           |
| `ExitCode` | `int`    | The exit code of the process (0 = success).                  |
| `Error`    | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Not idempotent.** The command runs each time the step runs. The operation
always returns `Changed: true`. Use the `When` or `OnlyIfChanged` guards to
control when the step runs.

## Permissions

The token must have the `docker:execute` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
