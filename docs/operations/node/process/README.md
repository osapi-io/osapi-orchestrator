# Process management

List, examine, and send signals to processes on target hosts.

## Operations

| Method                                          | Description                | Idempotent     |
| ----------------------------------------------- | -------------------------- | -------------- |
| [`ProcessList(target)`](list.md)                | List processes             | Read-only      |
| [`ProcessGet(target, pid)`](get.md)             | Get a specific process     | Read-only      |
| [`ProcessSignal(target, pid, opts)`](signal.md) | Send a signal to a process | Non-idempotent |

## Permissions

| Operation         | Permission        |
| ----------------- | ----------------- |
| Read operations   | `process:read`    |
| Signal operations | `process:execute` |

## Example

Refer to
[`examples/operations/process.go`](../../../../examples/operations/process.go)
for a complete workflow example that uses all the operations.
