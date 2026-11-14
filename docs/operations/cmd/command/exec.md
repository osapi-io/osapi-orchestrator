# CommandExec

Runs a command on the target host with its arguments kept separate. The
operation sends the command and each argument as separate values, so no shell
reads them. Use this operation to run a command when you do not need shell
features (pipes, redirects, globs).

## Usage

```go
step := o.CommandExec("web-01", "ls", "-la", "/var/log")
```

## Parameters

| Parameter | Type        | Description                                               |
| --------- | ----------- | --------------------------------------------------------- |
| `target`  | `string`    | Target host: `_any`, `_all`, hostname, or label selector. |
| `command` | `string`    | The command to run.                                       |
| `args`    | `...string` | The arguments to the command, as variadic values.         |

## Result type

```go
var result osapi.CommandResult
err := results.Decode("command.exec.execute-1", &result)
```

| Field        | Type     | Description                                                  |
| ------------ | -------- | ------------------------------------------------------------ |
| `Stdout`     | `string` | The standard output of the command.                          |
| `Stderr`     | `string` | The standard error of the command.                           |
| `ExitCode`   | `int`    | The exit code of the process (0 = success).                  |
| `DurationMs` | `int64`  | The time that the command ran, in milliseconds.              |
| `Error`      | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Not idempotent.** The command runs each time the step runs. The operation
always returns `Changed: true`. Use the `When` or `OnlyIfChanged` guards to
control when the step runs:

```go
o.CommandExec("web-01", "systemctl", "restart", "nginx").
    OnlyIfChanged()
```

## Permissions

The token must have the `command:execute` permission.

## Example

Refer to
[`examples/operations/command.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/command.go)
for a complete example.
