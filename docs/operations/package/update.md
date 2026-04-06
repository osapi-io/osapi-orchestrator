# PackageUpdate

Updates all packages on the target node.

## Usage

```go
step := o.PackageUpdate("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.PackageMutationResult
err := results.Decode("update-package-1", &result)
```

| Field     | Type     | Description                                        |
| --------- | -------- | -------------------------------------------------- |
| `Changed` | `bool`   | Whether any packages were updated.                 |
| `Error`   | `string` | Error message if update failed; empty on success.  |

## Idempotency

**Non-idempotent.** Always runs the package update process.

## Permissions

Requires `package:write` permission.

## Example

See
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete working example.
