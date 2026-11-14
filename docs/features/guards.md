# Guards

A guard decides if a step runs. It uses conditions that the plan finds when it
runs.

## When

`When` takes a predicate function. The function receives the results of the
steps that ran before. The step runs only if the predicate returns true:

```go
o.CommandExec("_any", "whoami").
    After(hostname).
    When(func(r orchestrator.Results) bool {
        var h osapi.HostnameResult
        if err := r.Decode("get-hostname", &h); err != nil {
            return false
        }
        return h.Hostname != ""
    })
```

## WhenFact

`WhenFact` is a special guard. It examines the agent facts from an `AgentList`
step that ran before:

```go
o.CommandShell(target, "apt-get update -qq").
    After(agents).
    WhenFact("list-agents", func(a osapi.Agent) bool {
        return a.OSInfo != nil &&
            a.OSInfo.Distribution == "Ubuntu"
    })
```

## OnlyIfChanged

Skip the step if no dependency reported a change. Refer to
[OnlyIfChanged](only-if-changed.md).

## OnlyIfAllChanged

Skip the step if one or more dependencies reported no change. All dependencies
must report a change. The guard uses the `Changed` bool. Thus a broadcast with a
partial failure (Status=Failed, Changed=true) counts as changed:

```go
step.OnlyIfAllChanged()
```

## OnlyIfFailed

Skip the step if no dependency has `Status == Failed`:

```go
deploy := o.CommandExec("_all", "deploy.sh").
    Named("deploy").
    OnError(orchestrator.Continue)

o.CommandExec("_any", "cleanup.sh").
    Named("cleanup").
    After(deploy).
    OnlyIfFailed()
```

## Broadcast guards

In a broadcast operation (`_all`, label selectors), each host can succeed, fail,
or report a change on its own. The host-level guards examine the `HostResults`
of the dependencies:

| Method                  | Semantics                                               |
| ----------------------- | ------------------------------------------------------- |
| `OnlyIfAnyHostFailed`   | Any host in any dependency has `Status == "failed"`     |
| `OnlyIfAllHostsFailed`  | Every host in every dependency has `Status == "failed"` |
| `OnlyIfAnyHostSkipped`  | Any host in any dependency has `Status == "skipped"`    |
| `OnlyIfAnyHostChanged`  | Any host in any dependency has `Changed == true`        |
| `OnlyIfAllHostsChanged` | Every host in every dependency has `Changed == true`    |

### Skipped vs failed

- **Skipped** (`Status == "skipped"`) means that the host does not support the
  operation (for example, `ErrUnsupported` on Darwin). This is NOT an error. It
  shows that the host cannot do the operation.
- **Failed** (`Status == "failed"`) means that the agent tried the operation and
  got an error. This IS an error.
- `OnlyIfAnyHostFailed` checks `Status`, not `Error`. It does **not** trigger
  for a skipped host. It triggers only for a host where the operation failed.

Special cases:

- If the step has no dependencies, all guards return false.
- If a dependency has no `HostResults` (unicast), `OnlyIfAny*` skips that
  dependency and `OnlyIfAll*` returns false.

### Example

```go
deploy := o.CommandExec("_all", "deploy.sh").
    Named("deploy").
    OnError(orchestrator.Continue)

// Run cleanup only if at least one host reported an error.
o.CommandExec("_any", "cleanup.sh").
    Named("cleanup").
    After(deploy).
    OnlyIfAnyHostFailed()

// Verify only if every host reported a change.
o.CommandExec("_all", "verify.sh").
    Named("verify").
    After(deploy).
    OnlyIfAllHostsChanged()
```

## Guard reference

### Task-level guards

| Method             | Inspects        | Trigger condition                      |
| ------------------ | --------------- | -------------------------------------- |
| `OnlyIfChanged`    | `Changed` bool  | One or more dependencies changed       |
| `OnlyIfAllChanged` | `Changed` bool  | All dependencies changed               |
| `OnlyIfFailed`     | `Status` field  | One or more dependencies failed        |
| `When`             | Custom function | The predicate returns true             |
| `WhenFact`         | Agent facts     | The predicate matches the target agent |

### Host-level guards (broadcast)

| Method                  | Inspects             | Trigger condition              |
| ----------------------- | -------------------- | ------------------------------ |
| `OnlyIfAnyHostFailed`   | `HostResult.Status`  | One or more hosts failed       |
| `OnlyIfAllHostsFailed`  | `HostResult.Status`  | All hosts failed               |
| `OnlyIfAnyHostSkipped`  | `HostResult.Status`  | One or more hosts were skipped |
| `OnlyIfAnyHostChanged`  | `HostResult.Changed` | One or more hosts changed      |
| `OnlyIfAllHostsChanged` | `HostResult.Changed` | All hosts changed              |

## Example

Refer to
[`examples/features/guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/guards.go),
[`examples/features/when-fact.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/when-fact.go),
and
[`examples/features/broadcast-guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast-guards.go)
for complete examples.
