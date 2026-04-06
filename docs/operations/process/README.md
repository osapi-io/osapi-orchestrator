# Process Management

List, inspect, and signal running processes on target nodes.

## Operations

| Method                                          | Description                | Idempotent     |
| ----------------------------------------------- | -------------------------- | -------------- |
| [`ProcessList(target)`](list.md)                | List running processes     | Read-only      |
| [`ProcessGet(target, pid)`](get.md)             | Get a specific process     | Read-only      |
| [`ProcessSignal(target, pid, opts)`](signal.md) | Send a signal to a process | Non-idempotent |

## Permissions

| Operation         | Permission        |
| ----------------- | ----------------- |
| Read operations   | `process:read`    |
| Signal operations | `process:execute` |

## Example

See [`examples/operations/process.go`](../../examples/operations/process.go) for
a complete workflow example covering all operations.
