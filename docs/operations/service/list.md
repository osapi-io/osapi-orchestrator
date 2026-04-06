# ServiceList

Lists all services on the target node.

## Usage

```go
step := o.ServiceList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.ServiceInfoResult
err := results.Decode("list-service-1", &result)
```

| Field      | Type            | Description                                      |
| ---------- | --------------- | ------------------------------------------------ |
| `Hostname` | `string`        | The node's hostname.                             |
| `Services` | `[]ServiceInfo` | List of service entries.                         |
| `Error`    | `string`        | Error message if query failed; empty on success. |

### ServiceInfo

| Field         | Type     | Description                         |
| ------------- | -------- | ----------------------------------- |
| `Name`        | `string` | Service unit name.                  |
| `Status`      | `string` | Active status of the service.       |
| `Enabled`     | `bool`   | Whether the service starts on boot. |
| `Description` | `string` | Service description from unit file. |
| `PID`         | `int`    | Main process ID of the service.     |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `service:read` permission.

## Example

See
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete working example.
