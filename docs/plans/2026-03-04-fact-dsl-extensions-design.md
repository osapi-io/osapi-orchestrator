# Fact DSL Extensions Design

Add fact-aware discovery, filtering, and guards to the orchestrator DSL.
Build on existing `AgentResult` types and `When` guard patterns.

## Core Type

```go
// Predicate filters agents by their facts and properties.
type Predicate func(AgentResult) bool
```

Reuses `AgentResult` directly — it already has typed fields (`OSInfo`,
`Architecture`, `Memory`, `CPUCount`, `Labels`, `Facts map[string]any`).
No new wrapper type.

## Predicate Helpers

Composable functions that return `Predicate`:

```go
func OS(distribution string) Predicate          // OSInfo.Distribution match
func Arch(architecture string) Predicate        // Architecture match
func MinMemory(bytes uint64) Predicate          // Memory.Total >= bytes
func MinCPU(count int) Predicate                // CPUCount >= count
func HasLabel(key, value string) Predicate      // Labels[key] == value
func FactEquals(key string, value any) Predicate // Facts[key] == value
```

All predicates are case-insensitive string comparisons where applicable.
Multiple predicates compose with AND semantics (all must match).

## Orchestrator Methods

### Discover

```go
func (o *Orchestrator) Discover(
    ctx context.Context,
    predicates ...Predicate,
) ([]AgentResult, error)
```

Queries active agents and returns those matching ALL predicates. Runs
synchronously at plan-build time. Internally creates a temporary mini-plan
with `AgentList`, executes it, decodes the result, and applies filters.

### GroupByFact

```go
func (o *Orchestrator) GroupByFact(
    ctx context.Context,
    key string,
    predicates ...Predicate,
) (map[string][]AgentResult, error)
```

Queries agents, optionally filters by predicates, and groups results by the
string value at the named key. The key uses dot notation to access nested
fields on `AgentResult`:

- `"os.distribution"` → `OSInfo.Distribution`
- `"architecture"` → `Architecture`
- `"service_manager"` → `ServiceMgr`
- `"package_manager"` → `PackageMgr`
- Any other key falls back to `Facts[key]`

## Step Method

### WhenFact

```go
func (s *Step) WhenFact(
    agentListStep string,
    fn Predicate,
) *Step
```

Fact-based execution guard. Requires a prior `AgentList` step as a
dependency (referenced by name). The guard:

1. Decodes `AgentListResult` from the named step
2. Finds the agent matching this step's target hostname
3. Runs the predicate against that agent
4. Returns true (run) or false (skip)

For `_all` broadcast targets, the guard passes if at least one agent
matches. For precise per-agent filtering on broadcasts, use `Discover` +
loop instead.

## Implementation Detail

`Discover` and `GroupByFact` need URL and token to create their mini-plans.
Add `url` and `token` fields to the `Orchestrator` struct — they are
already passed to `New()`, just store them alongside the plan.

## Usage Patterns

### Build-time discovery

```go
ubuntuHosts, _ := o.Discover(ctx, orchestrator.OS("Ubuntu"))
for _, a := range ubuntuHosts {
    o.CommandShell(a.Hostname, "apt upgrade -y").After(health)
}
```

### Multi-distro playbooks

```go
groups, _ := o.GroupByFact(ctx, "os.distribution")
for distro, agents := range groups {
    for _, a := range agents {
        o.CommandShell(a.Hostname, installCmd(distro)).After(health)
    }
}
```

### Execution-time guard

```go
listing := o.AgentList().After(health)
o.CommandShell("web-01", "apt upgrade -y").
    After(listing).
    WhenFact("list-agents", func(a orchestrator.AgentResult) bool {
        return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
    })
```

### Composed predicates

```go
hosts, _ := o.Discover(ctx,
    orchestrator.OS("Ubuntu"),
    orchestrator.MinMemory(4 * 1024 * 1024 * 1024),
    orchestrator.Arch("amd64"),
)
```

## New Examples

Each example is a standalone `main.go`:

| Directory                      | What it shows                              |
| ------------------------------ | ------------------------------------------ |
| `examples/discover/`           | Find agents by OS/arch/memory predicates   |
| `examples/group-by-fact/`      | Group agents by distro, run per-group cmds |
| `examples/when-fact/`          | Fact-based guard on a step                 |
| `examples/fact-predicates/`    | Compose multiple predicates                |

## README Updates

1. Add a **Targeting** section (like osapi-sdk) documenting `_any`, `_all`,
   hostname, and label selectors
2. Add an **Agent Discovery** section documenting `Discover`,
   `GroupByFact`, predicates, and `WhenFact`
3. Restructure the examples table into categories matching the SDK pattern
   (separate tables for core orchestration examples and fact/discovery
   examples)
4. Add all existing examples to the table (currently only "all" is listed)

## What This Does NOT Change

- No changes to the SDK or NATS targeting
- No agent-side filtering — all filtering is orchestrator-side
- No new API endpoints
- Existing `When` guard unchanged
- Labels remain the primary routing mechanism
