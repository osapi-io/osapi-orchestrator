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

package orchestrator

import (
	"context"
	"fmt"

	engine "github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapi "github.com/retr0h/osapi/pkg/sdk/client"
)

// fetchAgentsDecodeName is the task name used when decoding the
// agent list result (injectable for testing).
var fetchAgentsDecodeName = "list-agents"

// Discover queries active agents and returns those matching all
// predicates. Runs synchronously at plan-build time. With no
// predicates, returns all agents.
func (o *Orchestrator) Discover(
	ctx context.Context,
	predicates ...Predicate,
) ([]osapi.Agent, error) {
	agents, err := o.fetchAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("discover: %w", err)
	}

	if len(predicates) == 0 {
		return agents, nil
	}

	matched := make([]osapi.Agent, 0, len(agents))
	for _, a := range agents {
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
) ([]osapi.Agent, error) {
	client := osapi.New(o.url, o.token)
	plan := engine.NewPlan(client)

	plan.TaskFunc(
		"list-agents",
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Agent.List(ctx)
			if err != nil {
				return nil, fmt.Errorf("list agents: %w", err)
			}

			return &engine.Result{
				Changed: false,
				Data:    engine.StructToMap(resp.Data),
			}, nil
		},
	)

	sdkReport, err := plan.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("run agent list: %w", err)
	}

	report := &Report{Tasks: sdkReport.Tasks, Duration: sdkReport.Duration}

	var list osapi.AgentList
	if err := report.Decode(fetchAgentsDecodeName, &list); err != nil {
		return nil, fmt.Errorf("decode agent list: %w", err)
	}

	return list.Agents, nil
}

// GroupByFact queries agents, optionally filters by predicates,
// and groups results by the string value at the given key.
func (o *Orchestrator) GroupByFact(
	ctx context.Context,
	key string,
	predicates ...Predicate,
) (map[string][]osapi.Agent, error) {
	agents, err := o.Discover(ctx, predicates...)
	if err != nil {
		return nil, fmt.Errorf("group by fact: %w", err)
	}

	groups := make(map[string][]osapi.Agent)
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
	a osapi.Agent,
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
