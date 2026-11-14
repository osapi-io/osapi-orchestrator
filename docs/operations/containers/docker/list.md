# DockerList

Lists the containers on the target host. You can filter the containers by state.

## Usage

```go
step := o.DockerList("_any", &osapi.DockerListParams{
    State: "running",
    Limit: 50,
})
```

## Parameters

| Parameter | Type                      | Description                       |
| --------- | ------------------------- | --------------------------------- |
| `target`  | `string`                  | The target host or routing value. |
| `params`  | `*osapi.DockerListParams` | The filter parameters (optional). |

### DockerListParams fields

| Field   | Type     | Description                                                |
| ------- | -------- | ---------------------------------------------------------- |
| `State` | `string` | The state to filter by: `"running"`, `"stopped"`, `"all"`. |
| `Limit` | `int`    | The maximum number of containers to return.                |

## Result type

```go
var result osapi.DockerListResult
err := results.Decode("docker-list", &result)
```

| Field        | Type                        | Description                                                  |
| ------------ | --------------------------- | ------------------------------------------------------------ |
| `Containers` | `[]osapi.DockerSummaryItem` | The containers that match.                                   |
| `Error`      | `string`                    | The error message if the operation failed. Empty on success. |

### DockerSummaryItem fields

| Field     | Type     | Description                              |
| --------- | -------- | ---------------------------------------- |
| `ID`      | `string` | The ID of the container.                 |
| `Name`    | `string` | The name of the container.               |
| `Image`   | `string` | The image that the container runs.       |
| `State`   | `string` | The current state of the container.      |
| `Created` | `string` | The time when the container was created. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `docker:read` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
