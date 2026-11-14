# ServiceList

Lists all services on the target host.

## Usage

```go
step := o.ServiceList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.ServiceInfoResult
err := results.Decode("list-service-1", &result)
```

| Field      | Type            | Description                                                  |
| ---------- | --------------- | ------------------------------------------------------------ |
| `Hostname` | `string`        | The hostname of the agent that ran the operation.            |
| `Services` | `[]ServiceInfo` | List of service entries.                                     |
| `Error`    | `string`        | The error message if the operation failed. Empty on success. |

### ServiceInfo

| Field         | Type     | Description                             |
| ------------- | -------- | --------------------------------------- |
| `Name`        | `string` | Service unit name.                      |
| `Status`      | `string` | Active status of the service.           |
| `Enabled`     | `bool`   | `true` if the service starts at boot.   |
| `Description` | `string` | Service description from the unit file. |
| `PID`         | `int`    | Main process ID of the service.         |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `service:read` permission.

## Example

Refer to
[`examples/operations/service.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/service.go)
for a complete example.
