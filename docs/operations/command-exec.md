# CommandExec

Executes a command on the target node with explicit argument separation. The
command and its arguments are passed separately, avoiding shell interpretation.
This is the preferred method for running commands when shell features (pipes,
redirects, globbing) are not needed.

## Usage

```go
step := o.CommandExec("web-01", "ls", "-la", "/var/log")
```

## Parameters

| Parameter | Type        | Description                                               |
| --------- | ----------- | --------------------------------------------------------- |
| `target`  | `string`    | Target host: `_any`, `_all`, hostname, or label selector. |
| `command` | `string`    | The command to execute.                                   |
| `args`    | `...string` | Variadic arguments passed to the command.                 |

## Result Type

```go
var result orchestrator.CommandResult
err := results.Decode("command.exec.execute-1", &result)
```

| Field        | Type     | Description                                          |
| ------------ | -------- | ---------------------------------------------------- |
| `Stdout`     | `string` | Standard output from the command.                    |
| `Stderr`     | `string` | Standard error from the command.                     |
| `ExitCode`   | `int`    | Process exit code (0 = success).                     |
| `DurationMs` | `int64`  | Execution duration in milliseconds.                  |
| `Error`      | `string` | Error message if execution failed; empty on success. |

## Idempotency

**Not idempotent.** The command executes every time the step runs. Always
returns `Changed: true`. Use `When` or `OnlyIfChanged` guards to control when
the step runs:

```go
o.CommandExec("web-01", "systemctl", "restart", "nginx").
    OnlyIfChanged()
```

## Permissions

Requires `command:execute` permission.

## Example

See
[`examples/command.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/command.go)
for a complete working example.
