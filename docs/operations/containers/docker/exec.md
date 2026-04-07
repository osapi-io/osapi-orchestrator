# DockerExec

Executes a command inside a running container on the target host.

## Usage

```go
step := o.DockerExec("web-01", "c1a2b3d4e5f6", osapi.DockerExecOpts{
    Command: []string{"nginx", "-t"},
})
```

## Parameters

| Parameter | Type                   | Description                        |
| --------- | ---------------------- | ---------------------------------- |
| `target`  | `string`               | Target host or routing value.      |
| `id`      | `string`               | Container ID or name to exec into. |
| `opts`    | `osapi.DockerExecOpts` | Command execution options.         |

### DockerExecOpts Fields

| Field        | Type       | Description                                           |
| ------------ | ---------- | ----------------------------------------------------- |
| `Command`    | `[]string` | Command and arguments to execute (required).          |
| `Env`        | `[]string` | Additional environment variables in KEY=VALUE format. |
| `WorkingDir` | `string`   | Working directory inside the container.               |

## Result Type

```go
var result osapi.DockerExecResult
err := results.Decode("docker-exec", &result)
```

| Field      | Type     | Description                       |
| ---------- | -------- | --------------------------------- |
| `Stdout`   | `string` | Standard output from the command. |
| `Stderr`   | `string` | Standard error from the command.  |
| `ExitCode` | `int`    | Process exit code (0 = success).  |
| `Error`    | `string` | Error if execution failed.        |

## Idempotency

**No.** Executes the command every time the step runs. Always returns
`Changed: true`. Use guards to control when the step runs.

## Permissions

Requires `docker:execute` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
