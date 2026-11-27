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
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/osapi-io/osapi/pkg/sdk/types"
)

// hostCall is an SDK call that answers with one row for each host.
type hostCall[T any] func(
	ctx context.Context,
	c *osapi.Client,
) (*osapi.Response[osapi.Collection[T]], error)

// collectionStep makes the step of an SDK call that answers with one
// row for each host. The step is named from prefix. An unsafe step is
// not idempotent, so the runner retries it only when the step asks
// with RetryUnsafe.
func collectionStep[T any](
	o *Orchestrator,
	prefix string,
	unsafe bool,
	call hostCall[T],
) *Step {
	return rowStep(o, prefix, unsafe, call, hostResultOf[T])
}

// rowStep is collectionStep with the mapping of a row to a host result.
func rowStep[T any](
	o *Orchestrator,
	prefix string,
	unsafe bool,
	call hostCall[T],
	toHost func(T) engine.HostResult,
) *Step {
	name := o.nextOpName(prefix)

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := call(ctx, c)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", prefix, err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(), toHost)
		},
	)

	if unsafe {
		task.SetUnsafe()
	}

	return &Step{task: task}
}

// hostResultOf reads the hostname, status, changed and error fields
// that each row of the SDK has, through its JSON.
func hostResultOf[T any](
	r T,
) engine.HostResult {
	data := engine.StructToMap(r)
	hr := engine.HostResult{Data: data}
	hr.Hostname, _ = data["hostname"].(string)
	hr.Status, _ = data["status"].(string)
	hr.Changed, _ = data["changed"].(bool)
	hr.Error, _ = data["error"].(string)

	return hr
}

// commandHostResult is hostResultOf for a command, where a non-zero exit
// code is a failure, so that guards like OnlyIfAnyHostFailed work for
// command steps. An error from the server takes precedence.
func commandHostResult(
	r osapi.CommandResult,
) engine.HostResult {
	hr := hostResultOf(r)
	if hr.Error == "" && r.ExitCode != 0 {
		hr.Error = fmt.Sprintf("exit code %d", r.ExitCode)
	}

	return hr
}

// ---------------------------------------------------------------------------
// Health
// ---------------------------------------------------------------------------

// HealthCheck creates a step that checks the controller is alive.
func (o *Orchestrator) HealthCheck() *Step {
	name := o.nextOpName("health-check")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			if _, err := c.Health.Liveness(ctx); err != nil {
				return nil, fmt.Errorf("health check: %w", err)
			}

			return &engine.Result{Changed: false}, nil
		},
	)

	return &Step{task: task}
}

// ---------------------------------------------------------------------------
// Command
// ---------------------------------------------------------------------------

// CommandExec creates a step that runs Command.Exec on the target: run a
// command with its arguments. The step is named after the command. A
// non-zero exit code fails the host.
func (o *Orchestrator) CommandExec(
	target string,
	input osapi.CommandExecInput,
) *Step {
	return rowStep(
		o, "run-"+filepath.Base(input.Command), true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.CommandResult]], error) {
			return c.Command.Exec(ctx, target, input)
		},
		commandHostResult,
	)
}

// CommandShell creates a step that runs Command.Shell on the target: run a
// command line in a shell. The step is named after the first word of the
// line. A non-zero exit code fails the host.
func (o *Orchestrator) CommandShell(
	target string,
	input osapi.CommandShellInput,
) *Step {
	first := "shell"
	if words := strings.Fields(input.Command); len(words) > 0 {
		first = filepath.Base(words[0])
	}

	return rowStep(
		o, "shell-"+first, true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.CommandResult]], error) {
			return c.Command.Shell(ctx, target, input)
		},
		commandHostResult,
	)
}

// ---------------------------------------------------------------------------
// File (Object Store)
// ---------------------------------------------------------------------------

// FileUpload creates a step that uploads data to the Object Store. The
// upload is idempotent: the SDK compares SHA-256 digests and skips it when
// the content is unchanged. Pass WithForce to replace the object always.
func (o *Orchestrator) FileUpload(
	name string,
	contentType types.ContentType,
	data []byte,
	opts ...UploadOption,
) *Step {
	cfg := &uploadConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	task := o.plan.TaskFunc(
		o.nextOpName("upload-file"),
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			upload := c.File.Upload
			if cfg.force {
				upload = c.File.Update
			}

			resp, err := upload(ctx, name, contentType, bytes.NewReader(data))
			if err != nil {
				return nil, fmt.Errorf("upload file %s: %w", name, err)
			}

			return &engine.Result{
				Changed: resp.Data.Changed,
				Data:    engine.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// FileChanged creates a step that checks whether data differs from the
// object stored under name. It pairs with OnlyIfChanged to skip an
// upload when the content is unchanged.
func (o *Orchestrator) FileChanged(
	name string,
	data []byte,
) *Step {
	task := o.plan.TaskFunc(
		o.nextOpName("check-file"),
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.File.GetChanged(ctx, name, bytes.NewReader(data))
			if err != nil {
				return nil, fmt.Errorf("check file %s: %w", name, err)
			}

			return &engine.Result{
				Changed: resp.Data.Changed,
				Data:    engine.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// ---------------------------------------------------------------------------
// Agent
// ---------------------------------------------------------------------------

// AgentList creates a step that lists the agents.
func (o *Orchestrator) AgentList() *Step {
	task := o.plan.TaskFunc(
		o.nextOpName("list-agents"),
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Agent.List(ctx)
			if err != nil {
				return nil, fmt.Errorf("list agents: %w", err)
			}

			return &engine.Result{Data: engine.StructToMap(resp.Data)}, nil
		},
	)

	return &Step{task: task}
}

// AgentGet creates a step that gets one agent.
func (o *Orchestrator) AgentGet(
	hostname string,
) *Step {
	task := o.plan.TaskFunc(
		o.nextOpName("get-agent"),
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Agent.Get(ctx, hostname)
			if err != nil {
				return nil, fmt.Errorf("get agent %s: %w", hostname, err)
			}

			return &engine.Result{Data: engine.StructToMap(resp.Data)}, nil
		},
	)

	return &Step{task: task}
}

// AgentDrain creates a step that drains an agent, so that it takes no
// new jobs.
func (o *Orchestrator) AgentDrain(
	hostname string,
) *Step {
	return o.agentAction("drain-agent", hostname, func(
		ctx context.Context,
		c *osapi.Client,
	) error {
		return c.Agent.Drain(ctx, hostname)
	})
}

// AgentUndrain creates a step that undrains an agent, so that it takes
// jobs again.
func (o *Orchestrator) AgentUndrain(
	hostname string,
) *Step {
	return o.agentAction("undrain-agent", hostname, func(
		ctx context.Context,
		c *osapi.Client,
	) error {
		return c.Agent.Undrain(ctx, hostname)
	})
}

// agentAction makes the step of an agent call that answers with no data.
func (o *Orchestrator) agentAction(
	prefix string,
	hostname string,
	call func(ctx context.Context, c *osapi.Client) error,
) *Step {
	task := o.plan.TaskFunc(
		o.nextOpName(prefix),
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			if err := call(ctx, c); err != nil {
				return nil, fmt.Errorf("%s %s: %w", prefix, hostname, err)
			}

			return &engine.Result{Changed: false}, nil
		},
	)

	return &Step{task: task}
}
