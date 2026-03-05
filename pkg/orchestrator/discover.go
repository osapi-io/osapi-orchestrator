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
	"encoding/json"
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
