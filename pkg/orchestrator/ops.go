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

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
)

// nextOpName generates a human-readable task name from a prefix.
// Appends a counter suffix on collision (e.g. "get-hostname-2").
func (o *Orchestrator) nextOpName(
	prefix string,
) string {
	o.nameCount[prefix]++
	if o.nameCount[prefix] > 1 {
		return fmt.Sprintf("%s-%d", prefix, o.nameCount[prefix])
	}

	return prefix
}

// HealthCheck creates a health check step.
func (o *Orchestrator) HealthCheck() *Step {
	name := o.nextOpName("health-check")

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
	name := o.nextOpName("get-hostname")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Hostname(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get hostname: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.HostnameResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NodeStatusGet creates a step that retrieves node status.
func (o *Orchestrator) NodeStatusGet(
	target string,
) *Step {
	name := o.nextOpName("get-status")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Status(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get status: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.NodeStatus) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NodeUptimeGet creates a step that retrieves system uptime.
func (o *Orchestrator) NodeUptimeGet(
	target string,
) *Step {
	name := o.nextOpName("get-uptime")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Uptime(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get uptime: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.UptimeResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NodeDiskGet creates a step that retrieves disk usage.
func (o *Orchestrator) NodeDiskGet(
	target string,
) *Step {
	name := o.nextOpName("get-disk")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Disk(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get disk: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DiskResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NodeMemoryGet creates a step that retrieves memory stats.
func (o *Orchestrator) NodeMemoryGet(
	target string,
) *Step {
	name := o.nextOpName("get-memory")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Memory(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get memory: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.MemoryResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NodeLoadGet creates a step that retrieves load averages.
func (o *Orchestrator) NodeLoadGet(
	target string,
) *Step {
	name := o.nextOpName("get-load")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Load(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get load: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.LoadResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NetworkDNSGet creates a step that retrieves DNS configuration.
func (o *Orchestrator) NetworkDNSGet(
	target string,
	interfaceName string,
) *Step {
	name := o.nextOpName("get-dns")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.GetDNS(ctx, target, interfaceName)
			if err != nil {
				return nil, fmt.Errorf("get dns: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DNSConfig) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NetworkDNSUpdate creates a step that updates DNS configuration.
func (o *Orchestrator) NetworkDNSUpdate(
	target string,
	interfaceName string,
	servers []string,
	searchDomains []string,
) *Step {
	name := o.nextOpName("update-dns")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.UpdateDNS(ctx, target, interfaceName, servers, searchDomains)
			if err != nil {
				return nil, fmt.Errorf("update dns: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DNSUpdateResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NetworkPingDo creates a step that pings an address.
func (o *Orchestrator) NetworkPingDo(
	target string,
	address string,
) *Step {
	name := o.nextOpName("ping")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Ping(ctx, target, address)
			if err != nil {
				return nil, fmt.Errorf("ping: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.PingResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CommandExec creates a step that executes a command.
func (o *Orchestrator) CommandExec(
	target string,
	command string,
	args ...string,
) *Step {
	name := o.nextOpName(fmt.Sprintf("run-%s", command))

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Exec(ctx, osapi.ExecRequest{
				Command: command,
				Args:    args,
				Target:  target,
			})
			if err != nil {
				return nil, fmt.Errorf("exec command: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CommandResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CommandShell creates a step that executes a shell command string.
func (o *Orchestrator) CommandShell(
	target string,
	command string,
) *Step {
	name := o.nextOpName(fmt.Sprintf("shell-%s", command))

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.Shell(ctx, osapi.ShellRequest{
				Command: command,
				Target:  target,
			})
			if err != nil {
				return nil, fmt.Errorf("shell command: %w", err)
			}

			return sdk.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CommandResult) sdk.HostResult {
					return sdk.HostResult{
						Hostname: r.Hostname,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// FileDeploy creates a step that deploys a file from the Object Store
// to the target agent's filesystem. The objectName must reference a
// file previously uploaded to the Object Store. ContentType should be
// "raw" for literal content or "template" for Go-template rendering
// with vars and agent facts.
func (o *Orchestrator) FileDeploy(
	target string,
	opts osapi.FileDeployOpts,
) *Step {
	name := o.nextOpName("deploy-file")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			opts.Target = target

			resp, err := c.Node.FileDeploy(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("deploy file: %w", err)
			}

			return &sdk.Result{
				JobID:   resp.Data.JobID,
				Changed: resp.Data.Changed,
				Data:    sdk.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// FileStatusGet creates a step that checks the status of a deployed
// file on the target agent. Returns whether the file is in-sync,
// drifted, or missing compared to the expected state.
func (o *Orchestrator) FileStatusGet(
	target string,
	path string,
) *Step {
	name := o.nextOpName("file-status")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.Node.FileStatus(ctx, target, path)
			if err != nil {
				return nil, fmt.Errorf("file status: %w", err)
			}

			return &sdk.Result{
				JobID:   resp.Data.JobID,
				Changed: resp.Data.Changed,
				Data:    sdk.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// FileUpload creates a step that uploads file content to the Object
// Store via the OSAPI REST API. Returns the object name that can be
// used in subsequent FileDeploy steps. This is a convenience wrapper
// that uses TaskFunc to call the file upload API directly. By default
// the upload is idempotent — the SDK compares SHA-256 digests and
// skips the upload when content is unchanged. Pass WithForce to always
// upload regardless of content changes.
func (o *Orchestrator) FileUpload(
	name string,
	contentType string,
	data []byte,
	opts ...UploadOption,
) *Step {
	cfg := &uploadConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	taskName := o.nextOpName("upload-file")

	task := o.plan.TaskFunc(
		taskName,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			var uploadOpts []osapi.UploadOption
			if cfg.force {
				uploadOpts = append(uploadOpts, osapi.WithForce())
			}

			resp, err := c.File.Upload(
				ctx,
				name,
				contentType,
				bytes.NewReader(data),
				uploadOpts...,
			)
			if err != nil {
				return nil, fmt.Errorf("upload file %s: %w", name, err)
			}

			return &sdk.Result{
				Changed: resp.Data.Changed,
				Data:    sdk.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// FileChanged creates a step that checks whether local content differs
// from the version stored in the Object Store. Computes SHA-256 locally
// and compares against the stored hash. Pairs with OnlyIfChanged to
// skip uploads when content is unchanged.
func (o *Orchestrator) FileChanged(
	name string,
	data []byte,
) *Step {
	taskName := o.nextOpName("check-file")

	task := o.plan.TaskFunc(
		taskName,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*sdk.Result, error) {
			resp, err := c.File.Changed(
				ctx,
				name,
				bytes.NewReader(data),
			)
			if err != nil {
				return nil, fmt.Errorf("check file %s: %w", name, err)
			}

			return &sdk.Result{
				Changed: resp.Data.Changed,
				Data:    sdk.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// AgentList creates a step that lists all active agents with their facts.
func (o *Orchestrator) AgentList() *Step {
	name := o.nextOpName("list-agents")

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
				Data:    sdk.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// AgentGet creates a step that retrieves detailed info about a specific agent.
func (o *Orchestrator) AgentGet(
	hostname string,
) *Step {
	name := o.nextOpName("get-agent")

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
				Data:    sdk.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}
