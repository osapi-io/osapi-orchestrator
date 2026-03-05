# Fact DSL Extensions Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add fact-aware discovery, filtering, guards, and examples to the orchestrator DSL.

**Architecture:** New `Predicate` type (`func(AgentResult) bool`) with composable helpers (`OS`, `Arch`, `MinMemory`, `MinCPU`, `HasLabel`, `FactEquals`). Two new `Orchestrator` methods (`Discover`, `GroupByFact`) that query agents synchronously at plan-build time. One new `Step` method (`WhenFact`) for execution-time fact guards. Store `url` and `token` on `Orchestrator` so `Discover`/`GroupByFact` can create temporary plans.

**Tech Stack:** Go, testify/suite, osapi-sdk

---

### Task 1: Store URL and Token on Orchestrator

**Files:**
- Modify: `pkg/orchestrator/types.go:27-31`
- Modify: `pkg/orchestrator/orchestrator.go:36-56`

**Step 1: Add fields to Orchestrator struct**

In `pkg/orchestrator/types.go`, add `url` and `token` fields:

```go
type Orchestrator struct {
	url       string
	token     string
	plan      *sdk.Plan
	nameCount map[string]int
	renderer  renderer
}
```

**Step 2: Store them in New()**

In `pkg/orchestrator/orchestrator.go`, update `New()` to store the values:

```go
return &Orchestrator{
	url:       url,
	token:     token,
	plan:      plan,
	nameCount: make(map[string]int),
	renderer:  r,
}
```

**Step 3: Run tests to verify no regressions**

Run: `go test ./pkg/orchestrator/... -v -count=1`
Expected: All existing tests pass.

**Step 4: Commit**

```
feat(orchestrator): store url and token on Orchestrator
```

---

### Task 2: Add Predicate Type and Helpers

**Files:**
- Create: `pkg/orchestrator/predicate.go`
- Create: `pkg/orchestrator/predicate_public_test.go`

**Step 1: Write the failing tests**

Create `pkg/orchestrator/predicate_public_test.go`:

```go
package orchestrator_test

import (
	"testing"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type PredicatePublicTestSuite struct {
	suite.Suite
}

func (s *PredicatePublicTestSuite) TestOS() {
	tests := []struct {
		name   string
		agent  orchestrator.AgentResult
		distro string
		want   bool
	}{
		{
			name: "Matches exact distribution",
			agent: orchestrator.AgentResult{
				OSInfo: &orchestrator.AgentOSInfo{Distribution: "Ubuntu"},
			},
			distro: "Ubuntu",
			want:   true,
		},
		{
			name: "Case insensitive match",
			agent: orchestrator.AgentResult{
				OSInfo: &orchestrator.AgentOSInfo{Distribution: "ubuntu"},
			},
			distro: "Ubuntu",
			want:   true,
		},
		{
			name:   "Returns false when OSInfo is nil",
			agent:  orchestrator.AgentResult{},
			distro: "Ubuntu",
			want:   false,
		},
		{
			name: "Returns false for non-matching distribution",
			agent: orchestrator.AgentResult{
				OSInfo: &orchestrator.AgentOSInfo{Distribution: "Debian"},
			},
			distro: "Ubuntu",
			want:   false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			pred := orchestrator.OS(tc.distro)
			s.Equal(tc.want, pred(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestArch() {
	tests := []struct {
		name  string
		agent orchestrator.AgentResult
		arch  string
		want  bool
	}{
		{
			name:  "Matches architecture",
			agent: orchestrator.AgentResult{Architecture: "amd64"},
			arch:  "amd64",
			want:  true,
		},
		{
			name:  "Case insensitive match",
			agent: orchestrator.AgentResult{Architecture: "AMD64"},
			arch:  "amd64",
			want:  true,
		},
		{
			name:  "Returns false for non-matching architecture",
			agent: orchestrator.AgentResult{Architecture: "arm64"},
			arch:  "amd64",
			want:  false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			pred := orchestrator.Arch(tc.arch)
			s.Equal(tc.want, pred(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestMinMemory() {
	tests := []struct {
		name  string
		agent orchestrator.AgentResult
		bytes int
		want  bool
	}{
		{
			name: "Passes when memory exceeds minimum",
			agent: orchestrator.AgentResult{
				Memory: &orchestrator.AgentMemory{Total: 8000},
			},
			bytes: 4000,
			want:  true,
		},
		{
			name: "Passes when memory equals minimum",
			agent: orchestrator.AgentResult{
				Memory: &orchestrator.AgentMemory{Total: 4000},
			},
			bytes: 4000,
			want:  true,
		},
		{
			name: "Fails when memory below minimum",
			agent: orchestrator.AgentResult{
				Memory: &orchestrator.AgentMemory{Total: 2000},
			},
			bytes: 4000,
			want:  false,
		},
		{
			name:  "Fails when Memory is nil",
			agent: orchestrator.AgentResult{},
			bytes: 4000,
			want:  false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			pred := orchestrator.MinMemory(tc.bytes)
			s.Equal(tc.want, pred(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestMinCPU() {
	tests := []struct {
		name  string
		agent orchestrator.AgentResult
		count int
		want  bool
	}{
		{
			name:  "Passes when CPU count exceeds minimum",
			agent: orchestrator.AgentResult{CPUCount: 8},
			count: 4,
			want:  true,
		},
		{
			name:  "Passes when CPU count equals minimum",
			agent: orchestrator.AgentResult{CPUCount: 4},
			count: 4,
			want:  true,
		},
		{
			name:  "Fails when CPU count below minimum",
			agent: orchestrator.AgentResult{CPUCount: 2},
			count: 4,
			want:  false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			pred := orchestrator.MinCPU(tc.count)
			s.Equal(tc.want, pred(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestHasLabel() {
	tests := []struct {
		name  string
		agent orchestrator.AgentResult
		key   string
		value string
		want  bool
	}{
		{
			name: "Matches label",
			agent: orchestrator.AgentResult{
				Labels: map[string]string{"env": "prod"},
			},
			key:   "env",
			value: "prod",
			want:  true,
		},
		{
			name: "Returns false for wrong value",
			agent: orchestrator.AgentResult{
				Labels: map[string]string{"env": "staging"},
			},
			key:   "env",
			value: "prod",
			want:  false,
		},
		{
			name:  "Returns false for nil labels",
			agent: orchestrator.AgentResult{},
			key:   "env",
			value: "prod",
			want:  false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			pred := orchestrator.HasLabel(tc.key, tc.value)
			s.Equal(tc.want, pred(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestFactEquals() {
	tests := []struct {
		name  string
		agent orchestrator.AgentResult
		key   string
		value any
		want  bool
	}{
		{
			name: "Matches string fact",
			agent: orchestrator.AgentResult{
				Facts: map[string]any{"cloud": "aws"},
			},
			key:   "cloud",
			value: "aws",
			want:  true,
		},
		{
			name: "Matches numeric fact",
			agent: orchestrator.AgentResult{
				Facts: map[string]any{"gpu_count": float64(2)},
			},
			key:   "gpu_count",
			value: float64(2),
			want:  true,
		},
		{
			name: "Returns false for wrong value",
			agent: orchestrator.AgentResult{
				Facts: map[string]any{"cloud": "gcp"},
			},
			key:   "cloud",
			value: "aws",
			want:  false,
		},
		{
			name:  "Returns false for nil facts",
			agent: orchestrator.AgentResult{},
			key:   "cloud",
			value: "aws",
			want:  false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			pred := orchestrator.FactEquals(tc.key, tc.value)
			s.Equal(tc.want, pred(tc.agent))
		})
	}
}

func (s *PredicatePublicTestSuite) TestMatchAll() {
	tests := []struct {
		name       string
		agent      orchestrator.AgentResult
		predicates []orchestrator.Predicate
		want       bool
	}{
		{
			name: "All predicates match",
			agent: orchestrator.AgentResult{
				Architecture: "amd64",
				OSInfo:       &orchestrator.AgentOSInfo{Distribution: "Ubuntu"},
				CPUCount:     8,
			},
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
				orchestrator.Arch("amd64"),
				orchestrator.MinCPU(4),
			},
			want: true,
		},
		{
			name: "One predicate fails",
			agent: orchestrator.AgentResult{
				Architecture: "arm64",
				OSInfo:       &orchestrator.AgentOSInfo{Distribution: "Ubuntu"},
			},
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
				orchestrator.Arch("amd64"),
			},
			want: false,
		},
		{
			name:       "Empty predicates match everything",
			agent:      orchestrator.AgentResult{},
			predicates: nil,
			want:       true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Equal(tc.want, orchestrator.MatchAll(tc.agent, tc.predicates...))
		})
	}
}

func TestPredicatePublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(PredicatePublicTestSuite))
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./pkg/orchestrator/... -run TestPredicatePublicTestSuite -v -count=1`
Expected: Compilation error — `OS`, `Arch`, etc. not defined.

**Step 3: Write the implementation**

Create `pkg/orchestrator/predicate.go`:

```go
package orchestrator

import "strings"

// Predicate filters agents by their facts and properties.
type Predicate func(AgentResult) bool

// OS returns a predicate that matches agents running the given
// distribution (case-insensitive).
func OS(
	distribution string,
) Predicate {
	return func(a AgentResult) bool {
		if a.OSInfo == nil {
			return false
		}

		return strings.EqualFold(a.OSInfo.Distribution, distribution)
	}
}

// Arch returns a predicate that matches agents with the given
// architecture (case-insensitive).
func Arch(
	architecture string,
) Predicate {
	return func(a AgentResult) bool {
		return strings.EqualFold(a.Architecture, architecture)
	}
}

// MinMemory returns a predicate that matches agents with at least
// the given total memory (in the same unit as AgentMemory.Total).
func MinMemory(
	total int,
) Predicate {
	return func(a AgentResult) bool {
		if a.Memory == nil {
			return false
		}

		return a.Memory.Total >= total
	}
}

// MinCPU returns a predicate that matches agents with at least
// the given number of CPUs.
func MinCPU(
	count int,
) Predicate {
	return func(a AgentResult) bool {
		return a.CPUCount >= count
	}
}

// HasLabel returns a predicate that matches agents with the given
// label key-value pair.
func HasLabel(
	key string,
	value string,
) Predicate {
	return func(a AgentResult) bool {
		return a.Labels[key] == value
	}
}

// FactEquals returns a predicate that matches agents where the
// given fact key equals the expected value.
func FactEquals(
	key string,
	value any,
) Predicate {
	return func(a AgentResult) bool {
		return a.Facts[key] == value
	}
}

// MatchAll returns true if the agent matches all given predicates.
// Returns true if no predicates are provided.
func MatchAll(
	agent AgentResult,
	predicates ...Predicate,
) bool {
	for _, p := range predicates {
		if !p(agent) {
			return false
		}
	}

	return true
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./pkg/orchestrator/... -run TestPredicatePublicTestSuite -v -count=1`
Expected: All tests pass.

**Step 5: Run full test suite**

Run: `go test ./pkg/orchestrator/... -v -count=1`
Expected: All tests pass.

**Step 6: Run linter**

Run: `just go::vet`
Expected: No lint errors.

**Step 7: Commit**

```
feat(orchestrator): add Predicate type and helpers
```

---

### Task 3: Add Discover Method

**Files:**
- Create: `pkg/orchestrator/discover.go`
- Create: `pkg/orchestrator/discover_public_test.go`

**Step 1: Write the failing tests**

Create `pkg/orchestrator/discover_public_test.go`:

```go
package orchestrator_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type DiscoverPublicTestSuite struct {
	suite.Suite
}

func (s *DiscoverPublicTestSuite) agentListHandler(
	agents []orchestrator.AgentResult,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"agents": agents,
			"total":  len(agents),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *DiscoverPublicTestSuite) TestDiscover() {
	agents := []orchestrator.AgentResult{
		{
			Hostname:     "web-01",
			Architecture: "amd64",
			CPUCount:     8,
			OSInfo:       &orchestrator.AgentOSInfo{Distribution: "Ubuntu", Version: "22.04"},
			Memory:       &orchestrator.AgentMemory{Total: 16000},
			Labels:       map[string]string{"env": "prod"},
		},
		{
			Hostname:     "web-02",
			Architecture: "arm64",
			CPUCount:     4,
			OSInfo:       &orchestrator.AgentOSInfo{Distribution: "Debian", Version: "12"},
			Memory:       &orchestrator.AgentMemory{Total: 8000},
			Labels:       map[string]string{"env": "staging"},
		},
		{
			Hostname:     "web-03",
			Architecture: "amd64",
			CPUCount:     2,
			OSInfo:       &orchestrator.AgentOSInfo{Distribution: "Ubuntu", Version: "24.04"},
			Memory:       &orchestrator.AgentMemory{Total: 4000},
			Labels:       map[string]string{"env": "prod"},
		},
	}

	tests := []struct {
		name       string
		predicates []orchestrator.Predicate
		wantCount  int
		wantHosts  []string
	}{
		{
			name:       "No predicates returns all agents",
			predicates: nil,
			wantCount:  3,
			wantHosts:  []string{"web-01", "web-02", "web-03"},
		},
		{
			name:       "Filter by OS",
			predicates: []orchestrator.Predicate{orchestrator.OS("Ubuntu")},
			wantCount:  2,
			wantHosts:  []string{"web-01", "web-03"},
		},
		{
			name: "Filter by OS and Arch",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
				orchestrator.Arch("amd64"),
			},
			wantCount: 2,
			wantHosts: []string{"web-01", "web-03"},
		},
		{
			name: "Filter by OS and MinCPU",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
				orchestrator.MinCPU(4),
			},
			wantCount: 1,
			wantHosts: []string{"web-01"},
		},
		{
			name: "No agents match",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("CentOS"),
			},
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(s.agentListHandler(agents))
			defer server.Close()

			o := orchestrator.New(server.URL, "test-token")

			results, err := o.Discover(
				context.Background(),
				tc.predicates...,
			)
			s.Require().NoError(err)
			s.Len(results, tc.wantCount)

			if tc.wantHosts != nil {
				hostnames := make([]string, len(results))
				for i, a := range results {
					hostnames[i] = a.Hostname
				}
				s.Equal(tc.wantHosts, hostnames)
			}
		})
	}
}

func TestDiscoverPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(DiscoverPublicTestSuite))
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./pkg/orchestrator/... -run TestDiscoverPublicTestSuite -v -count=1`
Expected: Compilation error — `Discover` not defined.

**Step 3: Write the implementation**

Create `pkg/orchestrator/discover.go`:

```go
package orchestrator

import (
	"context"
	"fmt"

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/osapi-io/osapi-sdk/pkg/osapi"
)

// Discover queries active agents and returns those matching all
// predicates. Runs synchronously at plan-build time. With no
// predicates, returns all agents.
func (o *Orchestrator) Discover(
	ctx context.Context,
	predicates ...Predicate,
) ([]AgentResult, error) {
	list, err := o.fetchAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("discover: %w", err)
	}

	if len(predicates) == 0 {
		return list.Agents, nil
	}

	matched := make([]AgentResult, 0, len(list.Agents))
	for _, a := range list.Agents {
		if MatchAll(a, predicates...) {
			matched = append(matched, a)
		}
	}

	return matched, nil
}

// fetchAgents creates a temporary plan, runs AgentList, and decodes
// the result.
func (o *Orchestrator) fetchAgents(
	ctx context.Context,
) (*AgentListResult, error) {
	client := osapi.New(o.url, o.token)
	plan := sdk.NewPlan(client)

	plan.TaskFunc(
		"list-agents",
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Agent.List(ctx)
			if err != nil {
				return nil, fmt.Errorf("list agents: %w", err)
			}

			var data map[string]any
			if err := json.Unmarshal(resp.RawJSON(), &data); err != nil {
				return nil, fmt.Errorf("unmarshal agents: %w", err)
			}

			return &sdk.Result{
				Changed: false,
				Data:    data,
			}, nil
		},
	)

	report, err := plan.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("run agent list: %w", err)
	}

	var list AgentListResult
	for _, t := range report.Tasks {
		if t.Name == "list-agents" {
			b, err := json.Marshal(t.Data)
			if err != nil {
				return nil, fmt.Errorf("marshal agent data: %w", err)
			}

			if err := json.Unmarshal(b, &list); err != nil {
				return nil, fmt.Errorf("decode agent data: %w", err)
			}

			return &list, nil
		}
	}

	return nil, fmt.Errorf("agent list result not found")
}
```

Note: The `fetchAgents` implementation depends on the SDK's `plan.Run(ctx)` signature. If the SDK's `Run` method does not accept a context, use `plan.Run(context.Background())` instead and ignore the ctx parameter. Add `"encoding/json"` to the import block.

**Step 4: Run tests to verify they pass**

Run: `go test ./pkg/orchestrator/... -run TestDiscoverPublicTestSuite -v -count=1`
Expected: Tests pass (or may need adjustment based on SDK's HTTP routing — the test server needs to handle the agent list API path).

**Step 5: Adjust test server if needed**

The test server may need to match the OSAPI SDK's expected endpoint path (e.g., `/api/v1alpha1/agents`). Check the SDK client's `Agent.List()` method for the exact path and update the `agentListHandler` to route accordingly. If needed, use a mux:

```go
mux := http.NewServeMux()
mux.HandleFunc("/api/v1alpha1/agents", s.agentListHandler(agents))
server := httptest.NewServer(mux)
```

**Step 6: Run full test suite and linter**

Run: `go test ./pkg/orchestrator/... -v -count=1 && just go::vet`
Expected: All tests pass, no lint errors.

**Step 7: Commit**

```
feat(orchestrator): add Discover method
```

---

### Task 4: Add GroupByFact Method

**Files:**
- Modify: `pkg/orchestrator/discover.go`
- Create: `pkg/orchestrator/discover_internal_test.go` (for `factValue`)
- Modify: `pkg/orchestrator/discover_public_test.go`

**Step 1: Write the failing tests**

Add to `pkg/orchestrator/discover_public_test.go`:

```go
func (s *DiscoverPublicTestSuite) TestGroupByFact() {
	agents := []orchestrator.AgentResult{
		{
			Hostname: "web-01",
			OSInfo:   &orchestrator.AgentOSInfo{Distribution: "Ubuntu"},
		},
		{
			Hostname: "web-02",
			OSInfo:   &orchestrator.AgentOSInfo{Distribution: "Debian"},
		},
		{
			Hostname: "web-03",
			OSInfo:   &orchestrator.AgentOSInfo{Distribution: "Ubuntu"},
		},
		{
			Hostname:     "web-04",
			Architecture: "arm64",
		},
	}

	tests := []struct {
		name       string
		key        string
		predicates []orchestrator.Predicate
		wantGroups map[string][]string
	}{
		{
			name: "Group by os.distribution",
			key:  "os.distribution",
			wantGroups: map[string][]string{
				"Ubuntu": {"web-01", "web-03"},
				"Debian": {"web-02"},
			},
		},
		{
			name: "Group by architecture",
			key:  "architecture",
			wantGroups: map[string][]string{
				"arm64": {"web-04"},
			},
		},
		{
			name: "Group with predicate filter",
			key:  "os.distribution",
			predicates: []orchestrator.Predicate{
				orchestrator.OS("Ubuntu"),
			},
			wantGroups: map[string][]string{
				"Ubuntu": {"web-01", "web-03"},
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(s.agentListHandler(agents))
			defer server.Close()

			o := orchestrator.New(server.URL, "test-token")

			groups, err := o.GroupByFact(
				context.Background(),
				tc.key,
				tc.predicates...,
			)
			s.Require().NoError(err)
			s.Len(groups, len(tc.wantGroups))

			for key, wantHosts := range tc.wantGroups {
				agents, ok := groups[key]
				s.Require().True(ok, "missing group %q", key)

				hostnames := make([]string, len(agents))
				for i, a := range agents {
					hostnames[i] = a.Hostname
				}
				s.Equal(wantHosts, hostnames)
			}
		})
	}
}
```

Create `pkg/orchestrator/discover_internal_test.go` for `factValue`:

```go
package orchestrator

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DiscoverTestSuite struct {
	suite.Suite
}

func (s *DiscoverTestSuite) TestFactValue() {
	tests := []struct {
		name  string
		agent AgentResult
		key   string
		want  string
	}{
		{
			name: "os.distribution",
			agent: AgentResult{
				OSInfo: &AgentOSInfo{Distribution: "Ubuntu"},
			},
			key:  "os.distribution",
			want: "Ubuntu",
		},
		{
			name: "os.version",
			agent: AgentResult{
				OSInfo: &AgentOSInfo{Version: "22.04"},
			},
			key:  "os.version",
			want: "22.04",
		},
		{
			name:  "architecture",
			agent: AgentResult{Architecture: "amd64"},
			key:   "architecture",
			want:  "amd64",
		},
		{
			name:  "service_manager",
			agent: AgentResult{ServiceMgr: "systemd"},
			key:   "service_manager",
			want:  "systemd",
		},
		{
			name:  "package_manager",
			agent: AgentResult{PackageMgr: "apt"},
			key:   "package_manager",
			want:  "apt",
		},
		{
			name:  "kernel_version",
			agent: AgentResult{KernelVersion: "6.5.0"},
			key:   "kernel_version",
			want:  "6.5.0",
		},
		{
			name: "falls back to Facts map",
			agent: AgentResult{
				Facts: map[string]any{"cloud": "aws"},
			},
			key:  "cloud",
			want: "aws",
		},
		{
			name: "returns empty for nil OSInfo on os key",
			agent: AgentResult{},
			key:   "os.distribution",
			want:  "",
		},
		{
			name:  "returns empty for missing fact",
			agent: AgentResult{Facts: map[string]any{}},
			key:   "nonexistent",
			want:  "",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Equal(tc.want, factValue(tc.agent, tc.key))
		})
	}
}

func TestDiscoverTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(DiscoverTestSuite))
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./pkg/orchestrator/... -run "TestDiscoverPublicTestSuite/TestGroupByFact|TestDiscoverTestSuite" -v -count=1`
Expected: Compilation errors — `GroupByFact`, `factValue` not defined.

**Step 3: Write the implementation**

Add to `pkg/orchestrator/discover.go`:

```go
// GroupByFact queries agents, optionally filters by predicates,
// and groups results by the string value at the given key.
//
// Known keys map to typed AgentResult fields:
//
//   - "os.distribution" → OSInfo.Distribution
//   - "os.version"      → OSInfo.Version
//   - "architecture"    → Architecture
//   - "service_manager" → ServiceMgr
//   - "package_manager" → PackageMgr
//   - "kernel_version"  → KernelVersion
//
// Any other key falls back to Facts[key].
// Agents with an empty value for the key are excluded.
func (o *Orchestrator) GroupByFact(
	ctx context.Context,
	key string,
	predicates ...Predicate,
) (map[string][]AgentResult, error) {
	agents, err := o.Discover(ctx, predicates...)
	if err != nil {
		return nil, fmt.Errorf("group by fact: %w", err)
	}

	groups := make(map[string][]AgentResult)
	for _, a := range agents {
		v := factValue(a, key)
		if v == "" {
			continue
		}

		groups[v] = append(groups[v], a)
	}

	return groups, nil
}

// factValue extracts a string value from an agent for grouping.
func factValue(
	a AgentResult,
	key string,
) string {
	switch key {
	case "os.distribution":
		if a.OSInfo == nil {
			return ""
		}

		return a.OSInfo.Distribution
	case "os.version":
		if a.OSInfo == nil {
			return ""
		}

		return a.OSInfo.Version
	case "architecture":
		return a.Architecture
	case "service_manager":
		return a.ServiceMgr
	case "package_manager":
		return a.PackageMgr
	case "kernel_version":
		return a.KernelVersion
	default:
		if a.Facts == nil {
			return ""
		}

		v, _ := a.Facts[key].(string)

		return v
	}
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./pkg/orchestrator/... -run "TestDiscoverPublicTestSuite|TestDiscoverTestSuite" -v -count=1`
Expected: All tests pass.

**Step 5: Run full test suite and linter**

Run: `go test ./pkg/orchestrator/... -v -count=1 && just go::vet`
Expected: All pass.

**Step 6: Commit**

```
feat(orchestrator): add GroupByFact method
```

---

### Task 5: Add WhenFact Step Method

**Files:**
- Modify: `pkg/orchestrator/step.go`
- Modify: `pkg/orchestrator/step_public_test.go`

**Step 1: Write the failing tests**

Add to `pkg/orchestrator/step_public_test.go` in the `TestChaining` table:

```go
{
	name: "WhenFact returns same step",
	chainFn: func() *orchestrator.Step {
		return s.orch.NodeHostnameGet("_any").
			WhenFact("list-agents", func(
				_ orchestrator.AgentResult,
			) bool {
				return true
			})
	},
},
```

Also add a dedicated test method for WhenFact behavior to the suite:

```go
func (s *StepPublicTestSuite) TestWhenFact() {
	tests := []struct {
		name    string
		chainFn func() *orchestrator.Step
	}{
		{
			name: "WhenFact chains with After",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck("_any")

				return s.orch.NodeHostnameGet("_any").
					After(health).
					WhenFact("list-agents", func(
						_ orchestrator.AgentResult,
					) bool {
						return true
					})
			},
		},
		{
			name: "WhenFact in full chain",
			chainFn: func() *orchestrator.Step {
				health := s.orch.HealthCheck("_any")

				return s.orch.NodeHostnameGet("_any").
					After(health).
					WhenFact("list-agents", func(
						a orchestrator.AgentResult,
					) bool {
						return a.OSInfo != nil
					}).
					Retry(2).
					OnError(orchestrator.Continue)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := tc.chainFn()
			s.NotNil(step)
		})
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./pkg/orchestrator/... -run TestStepPublicTestSuite -v -count=1`
Expected: Compilation error — `WhenFact` not defined on `*Step`.

**Step 3: Write the implementation**

Add to `pkg/orchestrator/step.go`:

```go
// WhenFact adds a fact-based guard. The step runs only if the
// predicate returns true for the target agent. Requires a prior
// AgentList step referenced by name.
//
// For broadcast targets (_all, labels), the guard passes if at
// least one agent matches the predicate.
func (s *Step) WhenFact(
	agentListStep string,
	fn Predicate,
) *Step {
	s.task.WhenWithReason(func(sdkResults sdk.Results) bool {
		r := Results{results: sdkResults}

		var list AgentListResult
		if err := r.Decode(agentListStep, &list); err != nil {
			return false
		}

		target := ""
		if op := s.task.Operation(); op != nil {
			target = op.Target
		}

		for _, a := range list.Agents {
			if target == "_all" || target == a.Hostname {
				if fn(a) {
					return true
				}
			}
		}

		return false
	}, "when-fact: no matching agent")

	return s
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./pkg/orchestrator/... -run TestStepPublicTestSuite -v -count=1`
Expected: All tests pass.

**Step 5: Run full test suite and linter**

Run: `go test ./pkg/orchestrator/... -v -count=1 && just go::vet`
Expected: All pass.

**Step 6: Commit**

```
feat(orchestrator): add WhenFact step method
```

---

### Task 6: Add Examples

**Files:**
- Create: `examples/discover/main.go`
- Create: `examples/group-by-fact/main.go`
- Create: `examples/when-fact/main.go`
- Create: `examples/fact-predicates/main.go`

**Step 1: Create discover example**

Create `examples/discover/main.go`:

```go
// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package main demonstrates agent discovery with fact predicates.
// Discovers agents running Ubuntu on amd64, then retrieves the
// hostname from each matching host.
//
// DAG (per discovered host):
//
//	health-check
//	    └── get-hostname (target=<discovered host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

func main() {
	token := os.Getenv("OSAPI_TOKEN")
	if token == "" {
		log.Fatal("OSAPI_TOKEN is required")
	}

	url := os.Getenv("OSAPI_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	o := orchestrator.New(url, token)

	// Discover Ubuntu agents at plan-build time.
	agents, err := o.Discover(
		context.Background(),
		orchestrator.OS("Ubuntu"),
		orchestrator.Arch("amd64"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Discovered %d matching agents\n", len(agents))

	health := o.HealthCheck("_any")

	// Create a hostname step for each discovered agent.
	for _, a := range agents {
		o.NodeHostnameGet(a.Hostname).After(health)
	}

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s in %s\n", report.Summary(), report.Duration)
}
```

**Step 2: Create group-by-fact example**

Create `examples/group-by-fact/main.go`:

```go
// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package main demonstrates grouping agents by a fact value.
// Groups the fleet by OS distribution and runs a distro-specific
// package update command on each group.
//
// DAG (per group, per host):
//
//	health-check
//	    └── shell-<update-cmd> (target=<host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

func installCmd(
	distro string,
) string {
	switch distro {
	case "Ubuntu", "Debian":
		return "apt-get update -qq"
	case "CentOS", "Rocky", "AlmaLinux":
		return "yum check-update -q"
	default:
		return "echo unsupported"
	}
}

func main() {
	token := os.Getenv("OSAPI_TOKEN")
	if token == "" {
		log.Fatal("OSAPI_TOKEN is required")
	}

	url := os.Getenv("OSAPI_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	o := orchestrator.New(url, token)

	groups, err := o.GroupByFact(
		context.Background(),
		"os.distribution",
	)
	if err != nil {
		log.Fatal(err)
	}

	health := o.HealthCheck("_any")

	for distro, agents := range groups {
		cmd := installCmd(distro)
		fmt.Printf("Group %s (%d hosts): %s\n", distro, len(agents), cmd)

		for _, a := range agents {
			o.CommandShell(a.Hostname, cmd).After(health)
		}
	}

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s in %s\n", report.Summary(), report.Duration)
}
```

**Step 3: Create when-fact example**

Create `examples/when-fact/main.go`:

```go
// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package main demonstrates the WhenFact execution-time guard.
// Lists agents, then conditionally runs a command only if the target
// agent is running Ubuntu.
//
// DAG:
//
//	health-check
//	    └── list-agents
//	            └── shell-apt-update (when-fact: os == Ubuntu)
//
// Run with: OSAPI_TOKEN="<jwt>" OSAPI_TARGET="<hostname>" go run main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

func main() {
	token := os.Getenv("OSAPI_TOKEN")
	if token == "" {
		log.Fatal("OSAPI_TOKEN is required")
	}

	url := os.Getenv("OSAPI_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	target := os.Getenv("OSAPI_TARGET")
	if target == "" {
		target = "_any"
	}

	o := orchestrator.New(url, token)

	health := o.HealthCheck("_any")
	agents := o.AgentList().After(health)

	// Guard: only run if the target agent is Ubuntu.
	o.CommandShell(target, "apt-get update -qq").
		After(agents).
		WhenFact("list-agents", func(a orchestrator.AgentResult) bool {
			return a.OSInfo != nil &&
				a.OSInfo.Distribution == "Ubuntu"
		})

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s in %s\n", report.Summary(), report.Duration)
}
```

**Step 4: Create fact-predicates example**

Create `examples/fact-predicates/main.go`:

```go
// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package main demonstrates composing multiple fact predicates.
// Discovers agents that are Ubuntu, amd64, with at least 4 CPUs
// and 8GB memory, then queries their load averages.
//
// DAG (per matching host):
//
//	health-check
//	    └── get-load (target=<host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

func main() {
	token := os.Getenv("OSAPI_TOKEN")
	if token == "" {
		log.Fatal("OSAPI_TOKEN is required")
	}

	url := os.Getenv("OSAPI_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	o := orchestrator.New(url, token)

	// Compose predicates: Ubuntu + amd64 + 4 CPUs + 8GB.
	agents, err := o.Discover(
		context.Background(),
		orchestrator.OS("Ubuntu"),
		orchestrator.Arch("amd64"),
		orchestrator.MinCPU(4),
		orchestrator.MinMemory(8000),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d agents matching all predicates\n", len(agents))

	for _, a := range agents {
		fmt.Printf("  %s (%s %s, %d CPUs)\n",
			a.Hostname,
			a.OSInfo.Distribution,
			a.Architecture,
			a.CPUCount,
		)
	}

	health := o.HealthCheck("_any")

	for _, a := range agents {
		o.NodeLoadGet(a.Hostname).After(health)
	}

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s in %s\n", report.Summary(), report.Duration)
}
```

**Step 5: Verify examples compile**

Run: `go build ./examples/discover/... && go build ./examples/group-by-fact/... && go build ./examples/when-fact/... && go build ./examples/fact-predicates/...`
Expected: No compilation errors.

**Step 6: Run linter**

Run: `just go::vet`
Expected: No lint errors.

**Step 7: Commit**

```
feat(examples): add discover, group-by-fact, when-fact, and fact-predicates examples
```

---

### Task 7: Update README

**Files:**
- Modify: `README.md`

**Step 1: Add Agent Discovery section**

After the "Broadcast Results" section (line 171) and before "Examples" (line 173), add:

```markdown
### Agent Discovery

Query agents at plan-build time and filter by typed predicates:

```go
agents, err := o.Discover(ctx,
    orchestrator.OS("Ubuntu"),
    orchestrator.Arch("amd64"),
    orchestrator.MinCPU(4),
)

for _, a := range agents {
    o.CommandShell(a.Hostname, "apt upgrade -y").After(health)
}
```

| Method        | What it does                                           |
| ------------- | ------------------------------------------------------ |
| `Discover`    | Query agents filtered by predicates                    |
| `GroupByFact` | Group agents by a fact key (e.g., `os.distribution`)   |

### Predicates

Composable filters passed to `Discover` and `GroupByFact`:

| Predicate    | What it matches                                       |
| ------------ | ----------------------------------------------------- |
| `OS`         | Agent OS distribution (case-insensitive)              |
| `Arch`       | Agent architecture (case-insensitive)                 |
| `MinMemory`  | Minimum total memory                                  |
| `MinCPU`     | Minimum CPU count                                     |
| `HasLabel`   | Label key-value pair                                  |
| `FactEquals` | Arbitrary fact key-value equality                     |

### Fact Guards

Use `WhenFact` for execution-time fact checks with a prior `AgentList`
step:

```go
agents := o.AgentList().After(health)

o.CommandShell("web-01", "apt upgrade -y").
    After(agents).
    WhenFact("list-agents", func(a orchestrator.AgentResult) bool {
        return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
    })
```
```

**Step 2: Restructure examples table**

Replace the current examples section (lines 173-185) with a categorized
table following the SDK pattern:

```markdown
## 📋 Examples

Each example is a standalone Go program you can read and run.

### Core

| Example                                           | What it shows                                          |
| ------------------------------------------------- | ------------------------------------------------------ |
| [basic](examples/basic/main.go)                   | Simple DAG with health check and hostname query        |
| [parallel](examples/parallel/main.go)             | Five parallel queries depending on health check        |
| [retry](examples/retry/main.go)                   | Retry on failure with configurable attempts            |
| [command](examples/command/main.go)               | Command exec and shell with result decoding            |
| [verbose](examples/verbose/main.go)               | Verbose output with stdout/stderr/response data        |

### Guards and Conditions

| Example                                           | What it shows                                          |
| ------------------------------------------------- | ------------------------------------------------------ |
| [guards](examples/guards/main.go)                 | When predicate for conditional execution               |
| [only-if-changed](examples/only-if-changed/main.go) | Skip step unless dependency reported changes        |
| [error-recovery](examples/error-recovery/main.go) | Continue strategy with OnlyIfFailed cleanup            |

### Results

| Example                                           | What it shows                                          |
| ------------------------------------------------- | ------------------------------------------------------ |
| [broadcast](examples/broadcast/main.go)           | Per-host results from broadcast operations             |
| [task-func](examples/task-func/main.go)           | Custom steps with typed result decoding                |
| [dns-update](examples/dns-update/main.go)         | Read-then-write pattern with DNS operations            |

### Agent Discovery

| Example                                           | What it shows                                          |
| ------------------------------------------------- | ------------------------------------------------------ |
| [agent-facts](examples/agent-facts/main.go)       | List agents with OS, load, memory, and interfaces      |
| [discover](examples/discover/main.go)             | Find agents by OS and architecture predicates          |
| [group-by-fact](examples/group-by-fact/main.go)   | Group agents by distro, run per-group commands         |
| [when-fact](examples/when-fact/main.go)            | Fact-based guard on a step                             |
| [fact-predicates](examples/fact-predicates/main.go) | Compose multiple predicates for discovery            |

```bash
cd examples/discover
OSAPI_TOKEN="<jwt>" go run main.go
```
```

**Step 3: Add WhenFact to step chaining table**

In the step chaining table (around line 89-97), add after the `When` row:

```markdown
| `WhenFact`         | Guard — only run if agent fact predicate is true |
```

**Step 4: Add Targeting section**

After the install section (line 27) and before Features (line 29), add:

```markdown
## 🎯 Targeting

Most operations accept a `target` parameter to control which agents receive
the request:

| Target      | Behavior                                    |
| ----------- | ------------------------------------------- |
| `_any`      | Send to any available agent (load balanced) |
| `_all`      | Broadcast to every agent                    |
| `hostname`  | Send to a specific host                     |
| `key:value` | Send to agents matching a label             |
```

**Step 5: Run linter**

Run: `just go::vet`
Expected: No lint errors.

**Step 6: Commit**

```
docs(readme): add discovery, predicates, and categorized examples
```

---

### Task 8: Final Verification

**Step 1: Run full test suite**

Run: `just test`
Expected: All tests pass, linting clean, coverage acceptable.

**Step 2: Verify all examples compile**

Run: `go build ./examples/...`
Expected: No errors.

**Step 3: Review git log**

Run: `git log --oneline -10`
Expected: Clean commit history with conventional commit messages.
