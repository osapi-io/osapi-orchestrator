# CommandShell

Executes a shell command string on the target node. The command is passed to the
system shell, so shell features like pipes, redirects, environment variable
expansion, and globbing are available. Use `CommandExec` instead when shell
interpretation is not needed.

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
| `command` | `string` | Shell command string to execute.                          |

## Result Type

```go
var result orchestrator.CommandResult
err := results.Decode("command.shell.execute-1", &result)
```

| Field        | Type     | Description                         |
| ------------ | -------- | ----------------------------------- |
| `Stdout`     | `string` | Standard output from the command.   |
| `Stderr`     | `string` | Standard error from the command.    |
| `ExitCode`   | `int`    | Process exit code (0 = success).    |
| `DurationMs` | `int64`  | Execution duration in milliseconds. |

## Idempotency

**Not idempotent.** The command executes every time the step runs. Always
returns `Changed: true`. Use `When` or `OnlyIfChanged` guards to control when
the step runs:

```go
o.CommandShell("web-01", "echo 'hello'").
    OnlyIfChanged()
```

## Permissions

Requires `command:execute` permission.
