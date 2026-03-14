# DockerList

Lists containers on the target host, optionally filtered by state.

## Usage

```go
step := o.DockerList("_any", &osapi.DockerListParams{
    State: "running",
    Limit: 50,
})
```

## Parameters

| Parameter | Type                      | Description                   |
| --------- | ------------------------- | ----------------------------- |
| `target`  | `string`                  | Target host or routing value. |
| `params`  | `*osapi.DockerListParams` | Optional filter parameters.   |

### DockerListParams Fields

| Field   | Type     | Description                                         |
| ------- | -------- | --------------------------------------------------- |
| `State` | `string` | Filter by state: `"running"`, `"stopped"`, `"all"`. |
| `Limit` | `int`    | Maximum number of containers to return.             |

## Result Type

```go
var result osapi.DockerListResult
err := results.Decode("docker-list", &result)
```

| Field        | Type                        | Description                  |
| ------------ | --------------------------- | ---------------------------- |
| `Containers` | `[]osapi.DockerSummaryItem` | List of matching containers. |
| `Error`      | `string`                    | Error if listing failed.     |

### DockerSummaryItem Fields

| Field     | Type     | Description                   |
| --------- | -------- | ----------------------------- |
| `ID`      | `string` | Container ID.                 |
| `Name`    | `string` | Container name.               |
| `Image`   | `string` | Image the container runs.     |
| `State`   | `string` | Current container state.      |
| `Created` | `string` | Container creation timestamp. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `docker:read` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
