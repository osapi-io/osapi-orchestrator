# CommandShell

Run a command line through /bin/sh -c.

Runs the command line through /bin/sh -c, so pipes, redirects and variables
work. It answers as exec does. Its permission is command:shell, apart from
exec's command:execute: a token that may run a program may not open a shell.
Only the admin role gets it.

## Usage

```go
step := o.CommandShell("web-01", input)
```

## Parameters

| Parameter | Type                      | Description                                                   |
| --------- | ------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                  | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.CommandShellInput` | Refer to `CommandShellInput` below.                           |

### CommandShellInput

| Field     | Type     | Description                                                                |
| --------- | -------- | -------------------------------------------------------------------------- |
| `Command` | `string` | The command line, which /bin/sh -c reads, such as ls -la /tmp \| grep log. |
| `Cwd`     | `string` | The directory to run it in; the agent's own when empty.                    |
| `Timeout` | `int`    | Seconds to let it run before it is stopped, up to 300; 30 when 0.          |

## Result

Each host answers with one `osapi.CommandResult` row. A list answers with one
row for each item.

```go
var row osapi.CommandResult
err := report.Decode("shell-<first word>", &row)
```

| Field        | Type              | Description                                                                                                                                                                                                              |
| ------------ | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`    | `bool`            | Whether the host changed: true, since a command can change anything.                                                                                                                                                     |
| `DurationMS` | `int64`           | How long it ran, in milliseconds.                                                                                                                                                                                        |
| `Error`      | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `ExitCode`   | `int`             | Its exit code.                                                                                                                                                                                                           |
| `Hostname`   | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`     | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Stderr`     | `string`          | What it wrote to standard error.                                                                                                                                                                                         |
| `Stdout`     | `string`          | What it wrote to standard output.                                                                                                                                                                                        |

## Retries

A second run is not the same as one run. So `Retry` does not run this step again
unless the step adds `RetryUnsafe`.

A non-zero exit code fails the host, so guards such as `OnlyIfAnyHostFailed` see
it.

## Permissions

The token must have the `command:shell` permission. The step sends
`POST /api/node/{hostname}/command/shell`.

## Example

Refer to
[`examples/features/only-if-changed.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/only-if-changed.go).
