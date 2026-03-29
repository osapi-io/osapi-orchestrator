# FileStatusGet

Checks the status of a deployed file on the target agent. Compares the file on
disk against the expected state tracked by the agent and reports whether it is
in sync, has drifted, or is missing.

## Usage

```go
step := o.FileStatusGet("web-01", "/etc/myapp/config.yaml")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `path`    | `string` | Absolute path of the file to check.                       |

## Result Type

```go
var result osapi.FileStatusResult
err := results.Decode("file.status.get-1", &result)
```

| Field    | Type     | Description                                      |
| -------- | -------- | ------------------------------------------------ |
| `Path`   | `string` | Path of the checked file.                        |
| `Status` | `string` | One of `"in-sync"`, `"drifted"`, or `"missing"`. |
| `SHA256` | `string` | SHA-256 checksum of the file on disk (empty when |
|          |          | missing).                                        |

### Status Values

| Status    | Meaning                                             |
| --------- | --------------------------------------------------- |
| `in-sync` | File on disk matches the expected content.          |
| `drifted` | File exists but its checksum differs from expected. |
| `missing` | File does not exist at the specified path.          |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `file:read` permission.

## Example

See
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
for a complete working example.
