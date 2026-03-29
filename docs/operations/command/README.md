# Command Execution

Execute arbitrary commands on target nodes -- either directly or through a shell
interpreter.

## Operations

| Method                                   | Description                                 | Idempotent     |
| ---------------------------------------- | ------------------------------------------- | -------------- |
| [`CommandExec(target, opts)`](exec.md)   | Execute a command directly (no shell)       | Non-idempotent |
| [`CommandShell(target, opts)`](shell.md) | Execute via `/bin/sh -c` (pipes, redirects) | Non-idempotent |

## Permissions

| Operation      | Permission        |
| -------------- | ----------------- |
| All operations | `command:execute` |

## Example

See [`examples/operations/command.go`](../../examples/operations/command.go) for
a complete workflow example covering both operations.

The example demonstrates:

- Direct command execution with argument lists
- Shell command execution with pipes and redirects
- Decoding stdout, stderr, and exit code from results
