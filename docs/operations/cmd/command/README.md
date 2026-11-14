# Command execution

Run any command on target hosts. The command runs directly or through a shell.

## Operations

| Method                                   | Description                                 | Idempotent     |
| ---------------------------------------- | ------------------------------------------- | -------------- |
| [`CommandExec(target, opts)`](exec.md)   | Run a command directly (no shell)           | Non-idempotent |
| [`CommandShell(target, opts)`](shell.md) | Run through `/bin/sh -c` (pipes, redirects) | Non-idempotent |

## Permissions

| Operation      | Permission        |
| -------------- | ----------------- |
| All operations | `command:execute` |

## Example

Refer to
[`examples/operations/command.go`](../../../../examples/operations/command.go)
for a complete workflow example that uses the two operations.

The example shows how to:

- Run a command directly with a list of arguments
- Run a shell command with pipes and redirects
- Decode stdout, stderr, and the exit code from the results
