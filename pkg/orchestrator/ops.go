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

	engine "github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapi "github.com/retr0h/osapi/pkg/sdk/client"
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
		) (*engine.Result, error) {
			_, err := c.Health.Liveness(ctx)
			if err != nil {
				return nil, fmt.Errorf("health check: %w", err)
			}

			return &engine.Result{Changed: false}, nil
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
		) (*engine.Result, error) {
			resp, err := c.Node.Hostname(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get hostname: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.HostnameResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.Status(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get status: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.NodeStatus) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.Uptime(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get uptime: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.UptimeResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.Disk(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get disk: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DiskResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.Memory(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get memory: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.MemoryResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.Load(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get load: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.LoadResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.GetDNS(ctx, target, interfaceName)
			if err != nil {
				return nil, fmt.Errorf("get dns: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DNSConfig) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.UpdateDNS(ctx, target, interfaceName, servers, searchDomains)
			if err != nil {
				return nil, fmt.Errorf("update dns: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DNSUpdateResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
		) (*engine.Result, error) {
			resp, err := c.Node.Ping(ctx, target, address)
			if err != nil {
				return nil, fmt.Errorf("ping: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.PingResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
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
	name := o.nextOpName(fmt.Sprintf("run-%s", filepath.Base(command)))

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Node.Exec(ctx, osapi.ExecRequest{
				Command: command,
				Args:    args,
				Target:  target,
			})
			if err != nil {
				return nil, fmt.Errorf("exec command: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CommandResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    commandError(r),
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
	shellName := strings.Fields(command)
	if len(shellName) > 0 {
		shellName[0] = filepath.Base(shellName[0])
	}

	name := o.nextOpName(fmt.Sprintf("shell-%s", shellName[0]))

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Node.Shell(ctx, osapi.ShellRequest{
				Command: command,
				Target:  target,
			})
			if err != nil {
				return nil, fmt.Errorf("shell command: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CommandResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    commandError(r),
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// commandError returns an error string for a command result. If the
// server set an explicit error, that takes precedence. Otherwise a
// non-zero exit code is treated as a failure so that guards like
// OnlyIfAnyHostFailed work naturally for command steps.
func commandError(
	r osapi.CommandResult,
) string {
	if r.Error != "" {
		return r.Error
	}

	if r.ExitCode != 0 {
		return fmt.Sprintf("exit code %d", r.ExitCode)
	}

	return ""
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
		) (*engine.Result, error) {
			opts.Target = target

			resp, err := c.Node.FileDeploy(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("deploy file: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.FileDeployResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
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
		) (*engine.Result, error) {
			resp, err := c.Node.FileStatus(ctx, target, path)
			if err != nil {
				return nil, fmt.Errorf("file status: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.FileStatusResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// FileUndeploy creates a step that removes a previously deployed file
// from the target agent's filesystem.
func (o *Orchestrator) FileUndeploy(
	target string,
	path string,
) *Step {
	name := o.nextOpName("undeploy-file")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Node.FileUndeploy(ctx, osapi.FileUndeployOpts{
				Target: target,
				Path:   path,
			})
			if err != nil {
				return nil, fmt.Errorf("undeploy file: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.FileUndeployResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
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
		) (*engine.Result, error) {
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

			return &engine.Result{
				Changed: resp.Data.Changed,
				Data:    engine.StructToMap(resp.Data),
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
		) (*engine.Result, error) {
			resp, err := c.File.Changed(
				ctx,
				name,
				bytes.NewReader(data),
			)
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

// NodeHostnameUpdate creates a step that sets the system hostname.
func (o *Orchestrator) NodeHostnameUpdate(
	target string,
	hostname string,
) *Step {
	name := o.nextOpName("update-hostname")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Node.UpdateHostname(ctx, target, hostname)
			if err != nil {
				return nil, fmt.Errorf("update hostname: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.HostnameUpdateResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NodeOSGet creates a step that retrieves OS information.
func (o *Orchestrator) NodeOSGet(
	target string,
) *Step {
	name := o.nextOpName("get-os")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Node.OS(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get os: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.OSInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
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
		) (*engine.Result, error) {
			resp, err := c.Agent.Get(ctx, hostname)
			if err != nil {
				return nil, fmt.Errorf("get agent %s: %w", hostname, err)
			}

			return &engine.Result{
				Changed: false,
				Data:    engine.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// AgentDrain creates a step that drains an agent, preventing it from
// accepting new jobs.
func (o *Orchestrator) AgentDrain(
	hostname string,
) *Step {
	name := o.nextOpName("drain-agent")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Agent.Drain(ctx, hostname)
			if err != nil {
				return nil, fmt.Errorf("drain agent %s: %w", hostname, err)
			}

			return &engine.Result{
				Changed: false,
				Data:    engine.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// AgentUndrain creates a step that undrains an agent, allowing it to accept
// new jobs again.
func (o *Orchestrator) AgentUndrain(
	hostname string,
) *Step {
	name := o.nextOpName("undrain-agent")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Agent.Undrain(ctx, hostname)
			if err != nil {
				return nil, fmt.Errorf("undrain agent %s: %w", hostname, err)
			}

			return &engine.Result{
				Changed: false,
				Data:    engine.StructToMap(resp.Data),
			}, nil
		},
	)

	return &Step{task: task}
}

// DockerPull creates a step that pulls a Docker image on the target host.
func (o *Orchestrator) DockerPull(
	target string,
	opts osapi.DockerPullOpts,
) *Step {
	name := o.nextOpName("docker-pull")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.Pull(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("docker pull: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerPullResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerCreate creates a step that creates a new container on the target host.
func (o *Orchestrator) DockerCreate(
	target string,
	opts osapi.DockerCreateOpts,
) *Step {
	name := o.nextOpName("docker-create")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("docker create: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerStart creates a step that starts a stopped container on the target host.
func (o *Orchestrator) DockerStart(
	target string,
	id string,
) *Step {
	name := o.nextOpName("docker-start")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.Start(ctx, target, id)
			if err != nil {
				return nil, fmt.Errorf("docker start: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerActionResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerStop creates a step that stops a running container on the target host.
func (o *Orchestrator) DockerStop(
	target string,
	id string,
	opts osapi.DockerStopOpts,
) *Step {
	name := o.nextOpName("docker-stop")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.Stop(ctx, target, id, opts)
			if err != nil {
				return nil, fmt.Errorf("docker stop: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerActionResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerRemove creates a step that removes a container from the target host.
func (o *Orchestrator) DockerRemove(
	target string,
	id string,
	params *osapi.DockerRemoveParams,
) *Step {
	name := o.nextOpName("docker-remove")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.Remove(ctx, target, id, params)
			if err != nil {
				return nil, fmt.Errorf("docker remove: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerActionResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerExec creates a step that executes a command inside a running container.
func (o *Orchestrator) DockerExec(
	target string,
	id string,
	opts osapi.DockerExecOpts,
) *Step {
	name := o.nextOpName("docker-exec")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.Exec(ctx, target, id, opts)
			if err != nil {
				return nil, fmt.Errorf("docker exec: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerExecResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerInspect creates a step that retrieves detailed info about a container.
func (o *Orchestrator) DockerInspect(
	target string,
	id string,
) *Step {
	name := o.nextOpName("docker-inspect")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.Inspect(ctx, target, id)
			if err != nil {
				return nil, fmt.Errorf("docker inspect: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerDetailResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerList creates a step that lists containers on the target host.
func (o *Orchestrator) DockerList(
	target string,
	params *osapi.DockerListParams,
) *Step {
	name := o.nextOpName("docker-list")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.List(ctx, target, params)
			if err != nil {
				return nil, fmt.Errorf("docker list: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerListResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CronList creates a step that lists cron entries on the target host.
func (o *Orchestrator) CronList(
	target string,
) *Step {
	name := o.nextOpName("list-cron")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Schedule.CronList(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list cron: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CronEntryResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  false,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CronGet creates a step that retrieves a specific cron entry on the target host.
func (o *Orchestrator) CronGet(
	target string,
	entryName string,
) *Step {
	name := o.nextOpName("get-cron")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Schedule.CronGet(ctx, target, entryName)
			if err != nil {
				return nil, fmt.Errorf("get cron: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CronEntryResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  false,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CronCreate creates a step that creates a new cron entry on the target host.
func (o *Orchestrator) CronCreate(
	target string,
	opts osapi.CronCreateOpts,
) *Step {
	name := o.nextOpName("create-cron")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Schedule.CronCreate(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create cron: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CronMutationResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CronUpdate creates a step that updates an existing cron entry on the target host.
func (o *Orchestrator) CronUpdate(
	target string,
	entryName string,
	opts osapi.CronUpdateOpts,
) *Step {
	name := o.nextOpName("update-cron")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Schedule.CronUpdate(ctx, target, entryName, opts)
			if err != nil {
				return nil, fmt.Errorf("update cron: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CronMutationResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CronDelete creates a step that deletes a cron entry on the target host.
func (o *Orchestrator) CronDelete(
	target string,
	entryName string,
) *Step {
	name := o.nextOpName("delete-cron")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Schedule.CronDelete(ctx, target, entryName)
			if err != nil {
				return nil, fmt.Errorf("delete cron: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.CronMutationResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// DockerImageRemove creates a step that removes a container image
// from the target host.
func (o *Orchestrator) DockerImageRemove(
	target string,
	imageName string,
	params *osapi.DockerImageRemoveParams,
) *Step {
	name := o.nextOpName("docker-image-remove")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Docker.ImageRemove(
				ctx,
				target,
				imageName,
				params,
			)
			if err != nil {
				return nil, fmt.Errorf("docker image remove: %w", err)
			}

			return engine.CollectionResult(resp.Data, resp.RawJSON(),
				func(r osapi.DockerActionResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Changed:  r.Changed,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}
