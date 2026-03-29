# Guards

Guards control whether a step runs based on runtime conditions.

## When

`When` takes a predicate function that receives accumulated results from prior
steps. The step only runs if the predicate returns true:

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

`WhenFact` is a specialized guard that inspects agent facts from a prior
`AgentList` step:

```go
o.CommandShell(target, "apt-get update -qq").
    After(agents).
    WhenFact("list-agents", func(a osapi.Agent) bool {
        return a.OSInfo != nil &&
            a.OSInfo.Distribution == "Ubuntu"
    })
```

## OnlyIfChanged

Skip unless at least one dependency reported changes. See
[OnlyIfChanged](only-if-changed.md).

## OnlyIfAllChanged

Skip unless **all** dependencies reported changes. Uses the `Changed` bool, so a
broadcast partial failure (Status=Failed, Changed=true) still counts as changed:

```go
step.OnlyIfAllChanged()
```

## OnlyIfFailed

Skip unless at least one dependency has `Status == Failed`:

```go
deploy := o.CommandExec("_all", "deploy.sh").
    Named("deploy").
    OnError(orchestrator.Continue)

o.CommandExec("_any", "cleanup.sh").
    Named("cleanup").
    After(deploy).
    OnlyIfFailed()
```

## Broadcast Guards

For broadcast operations (`_all`, label selectors), individual hosts can
independently succeed, fail, or report changes. The four host-level guards
inspect `HostResults` from dependencies:

| Method                  | Semantics                                               |
| ----------------------- | ------------------------------------------------------- |
| `OnlyIfAnyHostFailed`   | Any host in any dependency has `Status == "failed"`     |
| `OnlyIfAllHostsFailed`  | Every host in every dependency has `Status == "failed"` |
| `OnlyIfAnyHostSkipped`  | Any host in any dependency has `Status == "skipped"`    |
| `OnlyIfAnyHostChanged`  | Any host in any dependency has `Changed == true`        |
| `OnlyIfAllHostsChanged` | Every host in every dependency has `Changed == true`    |

### Skipped vs Failed

- **Skipped** (`Status == "skipped"`) means the operation is not supported on
  this host (e.g., `ErrUnsupported` on Darwin). This is NOT an error — it
  indicates the host cannot perform the operation.
- **Failed** (`Status == "failed"`) means the operation was attempted but
  encountered an error. This IS an error.
- `OnlyIfAnyHostFailed` checks `Status`, not `Error`. It does **not** trigger
  for skipped hosts, only for hosts where the operation actually failed.

Edge cases:

- No dependencies: all guards return false
- Dependency has no `HostResults` (unicast): `OnlyIfAny*` skips the dependency,
  `OnlyIfAll*` returns false

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

## Guard Reference

### Task-Level Guards

| Method             | Inspects        | Trigger condition              |
| ------------------ | --------------- | ------------------------------ |
| `OnlyIfChanged`    | `Changed` bool  | Any dependency changed         |
| `OnlyIfAllChanged` | `Changed` bool  | All dependencies changed       |
| `OnlyIfFailed`     | `Status` field  | Any dependency failed          |
| `When`             | Custom function | Predicate returns true         |
| `WhenFact`         | Agent facts     | Predicate matches target agent |

### Host-Level Guards (Broadcast)

| Method                  | Inspects             | Trigger condition    |
| ----------------------- | -------------------- | -------------------- |
| `OnlyIfAnyHostFailed`   | `HostResult.Status`  | Any host failed      |
| `OnlyIfAllHostsFailed`  | `HostResult.Status`  | All hosts failed     |
| `OnlyIfAnyHostSkipped`  | `HostResult.Status`  | Any host was skipped |
| `OnlyIfAnyHostChanged`  | `HostResult.Changed` | Any host changed     |
| `OnlyIfAllHostsChanged` | `HostResult.Changed` | All hosts changed    |

## Example

See
[`examples/features/guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/guards.go),
[`examples/features/when-fact.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/when-fact.go),
and
[`examples/features/broadcast-guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast-guards.go)
for complete working examples.
