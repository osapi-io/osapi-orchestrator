# CommandShell

Runs a shell command string on the target host. The operation gives the command
to the system shell. You can then use shell features, for example pipes,
redirects, environment variable expansion, and globs. If you do not need a shell
to read the command, use `CommandExec`.

## Usage

```go
step := o.CommandShell(
    "web-01",
    "cat /etc/os-release | grep PRETTY_NAME",
)
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `command` | `string` | The shell command string to run.                          |

## Result type

```go
var result osapi.CommandResult
err := results.Decode("command.shell.execute-1", &result)
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
o.CommandShell("web-01", "echo 'hello'").
    OnlyIfChanged()
```

## Permissions

The token must have the `command:execute` permission.

## Example

Refer to
[`examples/operations/command.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/command.go)
for a complete example.
