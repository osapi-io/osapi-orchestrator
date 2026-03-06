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

// mustRawToMap unmarshals raw JSON bytes into a map for sdk.Result.Data.
// Panics on error because the SDK guarantees the response body is valid
// JSON — an unmarshal failure here indicates a programming error.
func mustRawToMap(
	raw []byte,
) map[string]any {
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(fmt.Sprintf("unmarshal SDK response: %v", err))
	}

	return data
}

// Operation constants matching the OSAPI agent's supported operations.
const (
	opNodeHostnameGet   = "node.hostname.get"
	opNodeStatusGet     = "node.status.get"
	opNodeUptimeGet     = "node.uptime.get"
	opNodeDiskGet       = "node.disk.get"
	opNodeMemoryGet     = "node.memory.get"
	opNodeLoadGet       = "node.load.get"
	opNetworkDNSGet     = "network.dns.get"
	opNetworkDNSUpdate  = "network.dns.update"
	opNetworkPingDo     = "network.ping.do"
	opCommandExec       = "command.exec.execute"
	opCommandShell      = "command.shell.execute"
	opFileDeployExecute = "file.deploy.execute"
	opFileStatusGet     = "file.status.get"
)

func (o *Orchestrator) newStep(
	op *sdk.Op,
) *Step {
	name := o.nextName(op.Operation, op.Params)
	task := o.plan.Task(name, op)

	return &Step{task: task}
}

// HealthCheck creates a health check step against the given target.
func (o *Orchestrator) HealthCheck(
	_ string,
) *Step {
	prefix := "health-check"
	o.nameCount[prefix]++

	name := prefix
	if o.nameCount[prefix] > 1 {
		name = fmt.Sprintf("%s-%d", prefix, o.nameCount[prefix])
	}

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			_, err := c.Health.Liveness(ctx)
			if err != nil {
				return nil, fmt.Errorf("health check: %w", err)
			}

			return &sdk.Result{Changed: false}, nil
		},
	)

	return &Step{task: task}
}

// NodeHostnameGet creates a step that retrieves the hostname.
func (o *Orchestrator) NodeHostnameGet(
	target string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNodeHostnameGet,
		Target:    target,
	})
}

// NodeStatusGet creates a step that retrieves node status.
func (o *Orchestrator) NodeStatusGet(
	target string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNodeStatusGet,
		Target:    target,
	})
}

// NodeUptimeGet creates a step that retrieves system uptime.
func (o *Orchestrator) NodeUptimeGet(
	target string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNodeUptimeGet,
		Target:    target,
	})
}

// NodeDiskGet creates a step that retrieves disk usage.
func (o *Orchestrator) NodeDiskGet(
	target string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNodeDiskGet,
		Target:    target,
	})
}

// NodeMemoryGet creates a step that retrieves memory stats.
func (o *Orchestrator) NodeMemoryGet(
	target string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNodeMemoryGet,
		Target:    target,
	})
}

// NodeLoadGet creates a step that retrieves load averages.
func (o *Orchestrator) NodeLoadGet(
	target string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNodeLoadGet,
		Target:    target,
	})
}

// NetworkDNSGet creates a step that retrieves DNS configuration.
func (o *Orchestrator) NetworkDNSGet(
	target string,
	interfaceName string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNetworkDNSGet,
		Target:    target,
		Params: map[string]any{
			"interface": interfaceName,
		},
	})
}

// NetworkDNSUpdate creates a step that updates DNS configuration.
func (o *Orchestrator) NetworkDNSUpdate(
	target string,
	interfaceName string,
	servers []string,
	searchDomains []string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNetworkDNSUpdate,
		Target:    target,
		Params: map[string]any{
			"interface":      interfaceName,
			"servers":        servers,
			"search_domains": searchDomains,
		},
	})
}

// NetworkPingDo creates a step that pings an address.
func (o *Orchestrator) NetworkPingDo(
	target string,
	address string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opNetworkPingDo,
		Target:    target,
		Params: map[string]any{
			"address": address,
		},
	})
}

// CommandExec creates a step that executes a command.
func (o *Orchestrator) CommandExec(
	target string,
	command string,
	args ...string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opCommandExec,
		Target:    target,
		Params: map[string]any{
			"command": command,
			"args":    args,
		},
	})
}

// CommandShell creates a step that executes a shell command string.
func (o *Orchestrator) CommandShell(
	target string,
	command string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opCommandShell,
		Target:    target,
		Params: map[string]any{
			"command": command,
		},
	})
}

// FileDeploy creates a step that deploys a file from the Object Store
// to the target agent's filesystem. The objectName must reference a
// file previously uploaded to the Object Store. ContentType should be
// "raw" for literal content or "template" for Go-template rendering
// with vars and agent facts.
func (o *Orchestrator) FileDeploy(
	target string,
	opts FileDeployOpts,
) *Step {
	params := map[string]any{
		"object_name":  opts.ObjectName,
		"path":         opts.Path,
		"content_type": opts.ContentType,
	}
	if opts.Mode != "" {
		params["mode"] = opts.Mode
	}
	if opts.Owner != "" {
		params["owner"] = opts.Owner
	}
	if opts.Group != "" {
		params["group"] = opts.Group
	}
	if len(opts.Vars) > 0 {
		params["vars"] = opts.Vars
	}

	return o.newStep(&sdk.Op{
		Operation: opFileDeployExecute,
		Target:    target,
		Params:    params,
	})
}

// FileStatusGet creates a step that checks the status of a deployed
// file on the target agent. Returns whether the file is in-sync,
// drifted, or missing compared to the expected state.
func (o *Orchestrator) FileStatusGet(
	target string,
	path string,
) *Step {
	return o.newStep(&sdk.Op{
		Operation: opFileStatusGet,
		Target:    target,
		Params: map[string]any{
			"path": path,
		},
	})
}

// FileUpload creates a step that uploads file content to the Object
// Store via the OSAPI REST API. Returns the object name that can be
// used in subsequent FileDeploy steps. This is a convenience wrapper
// that uses TaskFunc to call the file upload API directly.
//
// NOTE: Requires osapi-sdk FileService (not yet available). This
// operation will return an error until the SDK file endpoints are
// merged.
func (o *Orchestrator) FileUpload(
	name string,
	_ []byte,
) *Step {
	prefix := "upload-file"
	o.nameCount[prefix]++

	taskName := prefix
	if o.nameCount[prefix] > 1 {
		taskName = fmt.Sprintf("%s-%d", prefix, o.nameCount[prefix])
	}

	task := o.plan.TaskFunc(
		taskName,
		func(
			_ context.Context,
			_ *osapi.Client,
		) (*sdk.Result, error) {
			// TODO(sdk): Wire to c.File.Upload() when osapi-sdk
			// FileService is available.
			return nil, fmt.Errorf(
				"upload file %s: SDK FileService not yet available",
				name,
			)
		},
	)

	return &Step{task: task}
}

// AgentList creates a step that lists all active agents with their facts.
func (o *Orchestrator) AgentList() *Step {
	prefix := "list-agents"
	o.nameCount[prefix]++

	name := prefix
	if o.nameCount[prefix] > 1 {
		name = fmt.Sprintf("%s-%d", prefix, o.nameCount[prefix])
	}

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Agent.List(ctx)
			if err != nil {
				return nil, fmt.Errorf("list agents: %w", err)
			}

			return &sdk.Result{
				Changed: false,
				Data:    mustRawToMap(resp.RawJSON()),
			}, nil
		},
	)

	return &Step{task: task}
}

// AgentGet creates a step that retrieves detailed info about a specific agent.
func (o *Orchestrator) AgentGet(
	hostname string,
) *Step {
	prefix := "get-agent"
	o.nameCount[prefix]++

	name := prefix
	if o.nameCount[prefix] > 1 {
		name = fmt.Sprintf("%s-%d", prefix, o.nameCount[prefix])
	}

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Agent.Get(ctx, hostname)
			if err != nil {
				return nil, fmt.Errorf("get agent %s: %w", hostname, err)
			}

			return &sdk.Result{
				Changed: false,
				Data:    mustRawToMap(resp.RawJSON()),
			}, nil
		},
	)

	return &Step{task: task}
}
