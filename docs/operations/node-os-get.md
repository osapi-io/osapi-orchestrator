# NodeOSGet

Retrieves OS information from the target node, including distribution and
version.

## Usage

```go
step := o.NodeOSGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.OSInfoResult
err := results.Decode("get-os-1", &result)
```

| Field      | Type      | Description                                    |
| ---------- | --------- | ---------------------------------------------- |
| `Hostname` | `string`  | The node's hostname.                           |
| `OSInfo`   | `*OSInfo` | OS distribution and version (see below).       |
| `Error`    | `string`  | Error message if query failed; empty on success. |

### OSInfo

| Field          | Type     | Description                           |
| -------------- | -------- | ------------------------------------- |
| `Distribution` | `string` | OS distribution (e.g., `Ubuntu`).     |
| `Version`      | `string` | OS version (e.g., `24.04`).           |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `node:read` permission.

## Example

```go
plan := o.Plan("check-os")
o.NodeOSGet("_all")
report := plan.Execute(ctx)
```
