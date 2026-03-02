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

// Package orchestrator provides a declarative DSL for OSAPI infrastructure
// operations. Users declare what operations to run, where, in what order,
// and under what conditions. The orchestrator handles DAG execution,
// parallelism, retries, and reporting.
package orchestrator

import (
	"context"
	"fmt"

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/osapi-io/osapi-sdk/pkg/osapi"
)

// New creates an Orchestrator connected to the given OSAPI server.
func New(
	url string,
	token string,
) *Orchestrator {
	client := osapi.New(url, token)
	r := newLipglossRenderer()
	plan := sdk.NewPlan(client, sdk.WithHooks(rendererHooks(r)))

	return &Orchestrator{
		plan:     plan,
		renderer: r,
	}
}

// Run executes the plan and returns a report.
func (o *Orchestrator) Run() (*Report, error) {
	sdkReport, err := o.plan.Run(context.Background())
	if err != nil {
		return nil, err
	}

	return &Report{
		Tasks:    sdkReport.Tasks,
		Duration: sdkReport.Duration,
	}, nil
}

// TaskFunc creates a custom step that receives completed results
// from prior steps.
func (o *Orchestrator) TaskFunc(
	name string,
	fn func(ctx context.Context, r Results) (*sdk.Result, error),
) *Step {
	task := o.plan.TaskFuncWithResults(
		name,
		func(
			ctx context.Context,
			_ *osapi.Client,
			results sdk.Results,
		) (*sdk.Result, error) {
			return fn(ctx, Results{results: results})
		},
	)

	return &Step{task: task}
}

// nextName generates an auto-incremented task name from the operation.
func (o *Orchestrator) nextName(
	operation string,
) string {
	o.seq++

	return fmt.Sprintf("%s-%d", operation, o.seq)
}

// rendererHooks translates SDK hook callbacks into renderer calls.
func rendererHooks(
	r renderer,
) sdk.Hooks {
	levelParallel := make(map[int]bool)

	return sdk.Hooks{
		BeforePlan: func(
			summary sdk.PlanSummary,
		) {
			r.PlanStart(summary)
		},
		AfterPlan: func(
			report *sdk.Report,
		) {
			r.PlanDone(&Report{
				Tasks:    report.Tasks,
				Duration: report.Duration,
			})
		},
		BeforeLevel: func(
			level int,
			tasks []*sdk.Task,
			parallel bool,
		) {
			levelParallel[level] = parallel

			names := make([]string, len(tasks))
			for i, t := range tasks {
				names[i] = t.Name()
			}

			r.LevelStart(level, names, parallel)
		},
		AfterLevel: func(
			level int,
			results []sdk.TaskResult,
		) {
			changed := 0
			for _, res := range results {
				if res.Changed {
					changed++
				}
			}

			r.LevelDone(level, changed, len(results), levelParallel[level])
		},
		BeforeTask: func(
			task *sdk.Task,
		) {
			var detail string
			if op := task.Operation(); op != nil {
				detail = fmt.Sprintf(
					"operation=%s target=%s",
					op.Operation,
					op.Target,
				)
			} else {
				detail = "(custom function)"
			}

			r.TaskStart(task.Name(), detail)
		},
		AfterTask: func(
			_ *sdk.Task,
			result sdk.TaskResult,
		) {
			r.TaskDone(result)
		},
		OnRetry: func(
			task *sdk.Task,
			attempt int,
			err error,
		) {
			r.TaskRetry(task.Name(), attempt, err)
		},
		OnSkip: func(
			task *sdk.Task,
			reason string,
		) {
			r.TaskSkip(task.Name(), reason)
		},
	}
}
