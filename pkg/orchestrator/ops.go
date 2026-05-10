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

// ---------------------------------------------------------------------------
// Health
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Node: Hostname
// ---------------------------------------------------------------------------

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
			resp, err := c.Hostname.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get hostname: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Hostname.Update(ctx, target, hostname)
			if err != nil {
				return nil, fmt.Errorf("update hostname: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// ---------------------------------------------------------------------------
// Node: Status, Uptime, Disk, Memory, Load, OS
// ---------------------------------------------------------------------------

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
			resp, err := c.Status.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get status: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Uptime.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get uptime: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Disk.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get disk: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Memory.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get memory: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Load.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get load: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.OS.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get os: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// ---------------------------------------------------------------------------
// Network: DNS
// ---------------------------------------------------------------------------

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
			resp, err := c.DNS.Get(ctx, target, interfaceName)
			if err != nil {
				return nil, fmt.Errorf("get dns: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.DNS.Update(
				ctx, target, interfaceName, servers, searchDomains, false,
			)
			if err != nil {
				return nil, fmt.Errorf("update dns: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// NetworkDNSDelete creates a step that deletes DNS configuration.
func (o *Orchestrator) NetworkDNSDelete(
	target string,
	interfaceName string,
) *Step {
	name := o.nextOpName("delete-dns")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.DNS.Delete(ctx, target, interfaceName)
			if err != nil {
				return nil, fmt.Errorf("delete dns: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.DNSDeleteResult) engine.HostResult {
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
			resp, err := c.Ping.Do(ctx, target, address)
			if err != nil {
				return nil, fmt.Errorf("ping: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// ---------------------------------------------------------------------------
// Network: Interface
// ---------------------------------------------------------------------------

// InterfaceList creates a step that lists network interfaces.
func (o *Orchestrator) InterfaceList(
	target string,
) *Step {
	name := o.nextOpName("list-interface")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Interface.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list interfaces: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.InterfaceListResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// InterfaceGet creates a step that retrieves a specific network interface.
func (o *Orchestrator) InterfaceGet(
	target string,
	ifaceName string,
) *Step {
	name := o.nextOpName("get-interface")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Interface.Get(ctx, target, ifaceName)
			if err != nil {
				return nil, fmt.Errorf("get interface: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.InterfaceGetResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// InterfaceCreate creates a step that creates a network interface
// configuration.
func (o *Orchestrator) InterfaceCreate(
	target string,
	ifaceName string,
	opts osapi.InterfaceConfigOpts,
) *Step {
	name := o.nextOpName("create-interface")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Interface.Create(ctx, target, ifaceName, opts)
			if err != nil {
				return nil, fmt.Errorf("create interface: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.InterfaceMutationResult) engine.HostResult {
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

// InterfaceUpdate creates a step that updates a network interface
// configuration.
func (o *Orchestrator) InterfaceUpdate(
	target string,
	ifaceName string,
	opts osapi.InterfaceConfigOpts,
) *Step {
	name := o.nextOpName("update-interface")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Interface.Update(ctx, target, ifaceName, opts)
			if err != nil {
				return nil, fmt.Errorf("update interface: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.InterfaceMutationResult) engine.HostResult {
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

// InterfaceDelete creates a step that deletes a network interface
// configuration.
func (o *Orchestrator) InterfaceDelete(
	target string,
	ifaceName string,
) *Step {
	name := o.nextOpName("delete-interface")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Interface.Delete(ctx, target, ifaceName)
			if err != nil {
				return nil, fmt.Errorf("delete interface: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.InterfaceMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Network: Route
// ---------------------------------------------------------------------------

// RouteList creates a step that lists network routes.
func (o *Orchestrator) RouteList(
	target string,
) *Step {
	name := o.nextOpName("list-route")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Route.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list routes: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.RouteListResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// RouteGet creates a step that retrieves routes for a specific interface.
func (o *Orchestrator) RouteGet(
	target string,
	interfaceName string,
) *Step {
	name := o.nextOpName("get-route")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Route.Get(ctx, target, interfaceName)
			if err != nil {
				return nil, fmt.Errorf("get route: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.RouteGetResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// RouteCreate creates a step that creates route configuration.
func (o *Orchestrator) RouteCreate(
	target string,
	interfaceName string,
	opts osapi.RouteConfigOpts,
) *Step {
	name := o.nextOpName("create-route")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Route.Create(ctx, target, interfaceName, opts)
			if err != nil {
				return nil, fmt.Errorf("create route: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.RouteMutationResult) engine.HostResult {
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

// RouteUpdate creates a step that updates route configuration.
func (o *Orchestrator) RouteUpdate(
	target string,
	interfaceName string,
	opts osapi.RouteConfigOpts,
) *Step {
	name := o.nextOpName("update-route")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Route.Update(ctx, target, interfaceName, opts)
			if err != nil {
				return nil, fmt.Errorf("update route: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.RouteMutationResult) engine.HostResult {
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

// RouteDelete creates a step that deletes route configuration.
func (o *Orchestrator) RouteDelete(
	target string,
	interfaceName string,
) *Step {
	name := o.nextOpName("delete-route")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Route.Delete(ctx, target, interfaceName)
			if err != nil {
				return nil, fmt.Errorf("delete route: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.RouteMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Command
// ---------------------------------------------------------------------------

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
			resp, err := c.Command.Exec(ctx, osapi.ExecRequest{
				Command: command,
				Args:    args,
				Target:  target,
			})
			if err != nil {
				return nil, fmt.Errorf("exec command: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Command.Shell(ctx, osapi.ShellRequest{
				Command: command,
				Target:  target,
			})
			if err != nil {
				return nil, fmt.Errorf("shell command: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// ---------------------------------------------------------------------------
// File
// ---------------------------------------------------------------------------

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

			resp, err := c.FileDeploy.Deploy(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("deploy file: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.FileDeploy.Status(ctx, target, path)
			if err != nil {
				return nil, fmt.Errorf("file status: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.FileDeploy.Undeploy(ctx, osapi.FileUndeployOpts{
				Target: target,
				Path:   path,
			})
			if err != nil {
				return nil, fmt.Errorf("undeploy file: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// ---------------------------------------------------------------------------
// Docker
// ---------------------------------------------------------------------------

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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// ---------------------------------------------------------------------------
// Cron
// ---------------------------------------------------------------------------

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
			resp, err := c.Cron.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list cron: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Cron.Get(ctx, target, entryName)
			if err != nil {
				return nil, fmt.Errorf("get cron: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Cron.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create cron: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Cron.Update(ctx, target, entryName, opts)
			if err != nil {
				return nil, fmt.Errorf("update cron: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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
			resp, err := c.Cron.Delete(ctx, target, entryName)
			if err != nil {
				return nil, fmt.Errorf("delete cron: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
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

// ---------------------------------------------------------------------------
// Agent
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Sysctl
// ---------------------------------------------------------------------------

// SysctlList creates a step that lists managed sysctl parameters.
func (o *Orchestrator) SysctlList(
	target string,
) *Step {
	name := o.nextOpName("list-sysctl")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Sysctl.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list sysctl: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SysctlEntryResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// SysctlGet creates a step that retrieves a sysctl parameter by key.
func (o *Orchestrator) SysctlGet(
	target string,
	key string,
) *Step {
	name := o.nextOpName("get-sysctl")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Sysctl.Get(ctx, target, key)
			if err != nil {
				return nil, fmt.Errorf("get sysctl: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SysctlEntryResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// SysctlCreate creates a step that creates a sysctl parameter.
func (o *Orchestrator) SysctlCreate(
	target string,
	opts osapi.SysctlCreateOpts,
) *Step {
	name := o.nextOpName("create-sysctl")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Sysctl.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create sysctl: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SysctlMutationResult) engine.HostResult {
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

// SysctlUpdate creates a step that updates a sysctl parameter.
func (o *Orchestrator) SysctlUpdate(
	target string,
	key string,
	opts osapi.SysctlUpdateOpts,
) *Step {
	name := o.nextOpName("update-sysctl")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Sysctl.Update(ctx, target, key, opts)
			if err != nil {
				return nil, fmt.Errorf("update sysctl: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SysctlMutationResult) engine.HostResult {
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

// SysctlDelete creates a step that deletes a sysctl parameter.
func (o *Orchestrator) SysctlDelete(
	target string,
	key string,
) *Step {
	name := o.nextOpName("delete-sysctl")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Sysctl.Delete(ctx, target, key)
			if err != nil {
				return nil, fmt.Errorf("delete sysctl: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SysctlMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// NTP
// ---------------------------------------------------------------------------

// NTPGet creates a step that retrieves NTP status.
func (o *Orchestrator) NTPGet(
	target string,
) *Step {
	name := o.nextOpName("get-ntp")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.NTP.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get ntp: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.NtpStatusResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// NTPCreate creates a step that creates NTP configuration.
func (o *Orchestrator) NTPCreate(
	target string,
	opts osapi.NtpCreateOpts,
) *Step {
	name := o.nextOpName("create-ntp")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.NTP.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create ntp: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.NtpMutationResult) engine.HostResult {
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

// NTPUpdate creates a step that updates NTP configuration.
func (o *Orchestrator) NTPUpdate(
	target string,
	opts osapi.NtpUpdateOpts,
) *Step {
	name := o.nextOpName("update-ntp")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.NTP.Update(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("update ntp: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.NtpMutationResult) engine.HostResult {
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

// NTPDelete creates a step that deletes NTP configuration.
func (o *Orchestrator) NTPDelete(
	target string,
) *Step {
	name := o.nextOpName("delete-ntp")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.NTP.Delete(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("delete ntp: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.NtpMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Timezone
// ---------------------------------------------------------------------------

// TimezoneGet creates a step that retrieves the system timezone.
func (o *Orchestrator) TimezoneGet(
	target string,
) *Step {
	name := o.nextOpName("get-timezone")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Timezone.Get(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("get timezone: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.TimezoneResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// TimezoneUpdate creates a step that sets the system timezone.
func (o *Orchestrator) TimezoneUpdate(
	target string,
	opts osapi.TimezoneUpdateOpts,
) *Step {
	name := o.nextOpName("update-timezone")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Timezone.Update(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("update timezone: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.TimezoneMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Service
// ---------------------------------------------------------------------------

// ServiceList creates a step that lists services on the target host.
func (o *Orchestrator) ServiceList(
	target string,
) *Step {
	name := o.nextOpName("list-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list services: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// ServiceGet creates a step that retrieves a specific service.
func (o *Orchestrator) ServiceGet(
	target string,
	serviceName string,
) *Step {
	name := o.nextOpName("get-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Get(ctx, target, serviceName)
			if err != nil {
				return nil, fmt.Errorf("get service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceGetResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// ServiceCreate creates a step that creates a service unit file.
func (o *Orchestrator) ServiceCreate(
	target string,
	opts osapi.ServiceCreateOpts,
) *Step {
	name := o.nextOpName("create-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ServiceUpdate creates a step that updates a service unit file.
func (o *Orchestrator) ServiceUpdate(
	target string,
	serviceName string,
	opts osapi.ServiceUpdateOpts,
) *Step {
	name := o.nextOpName("update-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Update(ctx, target, serviceName, opts)
			if err != nil {
				return nil, fmt.Errorf("update service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ServiceDelete creates a step that deletes a service unit file.
func (o *Orchestrator) ServiceDelete(
	target string,
	serviceName string,
) *Step {
	name := o.nextOpName("delete-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Delete(ctx, target, serviceName)
			if err != nil {
				return nil, fmt.Errorf("delete service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ServiceStart creates a step that starts a service.
func (o *Orchestrator) ServiceStart(
	target string,
	serviceName string,
) *Step {
	name := o.nextOpName("start-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Start(ctx, target, serviceName)
			if err != nil {
				return nil, fmt.Errorf("start service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ServiceStop creates a step that stops a service.
func (o *Orchestrator) ServiceStop(
	target string,
	serviceName string,
) *Step {
	name := o.nextOpName("stop-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Stop(ctx, target, serviceName)
			if err != nil {
				return nil, fmt.Errorf("stop service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ServiceRestart creates a step that restarts a service.
func (o *Orchestrator) ServiceRestart(
	target string,
	serviceName string,
) *Step {
	name := o.nextOpName("restart-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Restart(ctx, target, serviceName)
			if err != nil {
				return nil, fmt.Errorf("restart service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ServiceEnable creates a step that enables a service to start on boot.
func (o *Orchestrator) ServiceEnable(
	target string,
	serviceName string,
) *Step {
	name := o.nextOpName("enable-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Enable(ctx, target, serviceName)
			if err != nil {
				return nil, fmt.Errorf("enable service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ServiceDisable creates a step that disables a service from starting on boot.
func (o *Orchestrator) ServiceDisable(
	target string,
	serviceName string,
) *Step {
	name := o.nextOpName("disable-service")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Service.Disable(ctx, target, serviceName)
			if err != nil {
				return nil, fmt.Errorf("disable service: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ServiceMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Package
// ---------------------------------------------------------------------------

// PackageList creates a step that lists installed packages.
func (o *Orchestrator) PackageList(
	target string,
) *Step {
	name := o.nextOpName("list-package")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Package.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list packages: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PackageInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// PackageGet creates a step that retrieves a specific package.
func (o *Orchestrator) PackageGet(
	target string,
	pkgName string,
) *Step {
	name := o.nextOpName("get-package")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Package.Get(ctx, target, pkgName)
			if err != nil {
				return nil, fmt.Errorf("get package: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PackageInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// PackageInstall creates a step that installs a package.
func (o *Orchestrator) PackageInstall(
	target string,
	pkgName string,
) *Step {
	name := o.nextOpName("install-package")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Package.Install(ctx, target, pkgName)
			if err != nil {
				return nil, fmt.Errorf("install package: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PackageMutationResult) engine.HostResult {
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

// PackageRemove creates a step that removes a package.
func (o *Orchestrator) PackageRemove(
	target string,
	pkgName string,
) *Step {
	name := o.nextOpName("remove-package")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Package.Remove(ctx, target, pkgName)
			if err != nil {
				return nil, fmt.Errorf("remove package: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PackageMutationResult) engine.HostResult {
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

// PackageUpdate creates a step that updates all packages.
func (o *Orchestrator) PackageUpdate(
	target string,
) *Step {
	name := o.nextOpName("update-package")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Package.Update(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("update packages: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PackageMutationResult) engine.HostResult {
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

// PackageListUpdates creates a step that lists available package updates.
func (o *Orchestrator) PackageListUpdates(
	target string,
) *Step {
	name := o.nextOpName("list-package-updates")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Package.ListUpdates(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list package updates: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PackageUpdateResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// ---------------------------------------------------------------------------
// User
// ---------------------------------------------------------------------------

// UserList creates a step that lists user accounts.
func (o *Orchestrator) UserList(
	target string,
) *Step {
	name := o.nextOpName("list-user")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list users: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.UserInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// UserGet creates a step that retrieves a specific user account.
func (o *Orchestrator) UserGet(
	target string,
	username string,
) *Step {
	name := o.nextOpName("get-user")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.Get(ctx, target, username)
			if err != nil {
				return nil, fmt.Errorf("get user: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.UserInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// UserCreate creates a step that creates a user account.
func (o *Orchestrator) UserCreate(
	target string,
	opts osapi.UserCreateOpts,
) *Step {
	name := o.nextOpName("create-user")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.UserMutationResult) engine.HostResult {
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

// UserUpdate creates a step that updates a user account.
func (o *Orchestrator) UserUpdate(
	target string,
	username string,
	opts osapi.UserUpdateOpts,
) *Step {
	name := o.nextOpName("update-user")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.Update(ctx, target, username, opts)
			if err != nil {
				return nil, fmt.Errorf("update user: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.UserMutationResult) engine.HostResult {
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

// UserDelete creates a step that deletes a user account.
func (o *Orchestrator) UserDelete(
	target string,
	username string,
) *Step {
	name := o.nextOpName("delete-user")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.Delete(ctx, target, username)
			if err != nil {
				return nil, fmt.Errorf("delete user: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.UserMutationResult) engine.HostResult {
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

// UserListKeys creates a step that lists SSH authorized keys for a user.
func (o *Orchestrator) UserListKeys(
	target string,
	username string,
) *Step {
	name := o.nextOpName("list-ssh-key")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.ListKeys(ctx, target, username)
			if err != nil {
				return nil, fmt.Errorf("list ssh keys: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SSHKeyInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// UserAddKey creates a step that adds an SSH authorized key for a user.
func (o *Orchestrator) UserAddKey(
	target string,
	username string,
	opts osapi.SSHKeyAddOpts,
) *Step {
	name := o.nextOpName("add-ssh-key")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.AddKey(ctx, target, username, opts)
			if err != nil {
				return nil, fmt.Errorf("add ssh key: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SSHKeyMutationResult) engine.HostResult {
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

// UserRemoveKey creates a step that removes an SSH authorized key for a user.
func (o *Orchestrator) UserRemoveKey(
	target string,
	username string,
	fingerprint string,
) *Step {
	name := o.nextOpName("remove-ssh-key")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.RemoveKey(ctx, target, username, fingerprint)
			if err != nil {
				return nil, fmt.Errorf("remove ssh key: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.SSHKeyMutationResult) engine.HostResult {
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

// UserChangePassword creates a step that changes a user's password.
func (o *Orchestrator) UserChangePassword(
	target string,
	username string,
	password string,
) *Step {
	name := o.nextOpName("change-password")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.User.ChangePassword(ctx, target, username, password)
			if err != nil {
				return nil, fmt.Errorf("change password: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.UserMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Group
// ---------------------------------------------------------------------------

// GroupList creates a step that lists groups.
func (o *Orchestrator) GroupList(
	target string,
) *Step {
	name := o.nextOpName("list-group")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Group.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list groups: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.GroupInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// GroupGet creates a step that retrieves a specific group.
func (o *Orchestrator) GroupGet(
	target string,
	groupName string,
) *Step {
	name := o.nextOpName("get-group")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Group.Get(ctx, target, groupName)
			if err != nil {
				return nil, fmt.Errorf("get group: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.GroupInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// GroupCreate creates a step that creates a group.
func (o *Orchestrator) GroupCreate(
	target string,
	opts osapi.GroupCreateOpts,
) *Step {
	name := o.nextOpName("create-group")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Group.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create group: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.GroupMutationResult) engine.HostResult {
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

// GroupUpdate creates a step that updates a group.
func (o *Orchestrator) GroupUpdate(
	target string,
	groupName string,
	opts osapi.GroupUpdateOpts,
) *Step {
	name := o.nextOpName("update-group")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Group.Update(ctx, target, groupName, opts)
			if err != nil {
				return nil, fmt.Errorf("update group: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.GroupMutationResult) engine.HostResult {
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

// GroupDelete creates a step that deletes a group.
func (o *Orchestrator) GroupDelete(
	target string,
	groupName string,
) *Step {
	name := o.nextOpName("delete-group")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Group.Delete(ctx, target, groupName)
			if err != nil {
				return nil, fmt.Errorf("delete group: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.GroupMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Certificate
// ---------------------------------------------------------------------------

// CertificateList creates a step that lists CA certificates.
func (o *Orchestrator) CertificateList(
	target string,
) *Step {
	name := o.nextOpName("list-certificate")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Certificate.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list certificates: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.CertificateCAResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// CertificateCreate creates a step that creates a CA certificate.
func (o *Orchestrator) CertificateCreate(
	target string,
	opts osapi.CertificateCreateOpts,
) *Step {
	name := o.nextOpName("create-certificate")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Certificate.Create(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("create certificate: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.CertificateCAMutationResult) engine.HostResult {
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

// CertificateUpdate creates a step that updates a CA certificate.
func (o *Orchestrator) CertificateUpdate(
	target string,
	certName string,
	opts osapi.CertificateUpdateOpts,
) *Step {
	name := o.nextOpName("update-certificate")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Certificate.Update(ctx, target, certName, opts)
			if err != nil {
				return nil, fmt.Errorf("update certificate: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.CertificateCAMutationResult) engine.HostResult {
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

// CertificateDelete creates a step that deletes a CA certificate.
func (o *Orchestrator) CertificateDelete(
	target string,
	certName string,
) *Step {
	name := o.nextOpName("delete-certificate")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Certificate.Delete(ctx, target, certName)
			if err != nil {
				return nil, fmt.Errorf("delete certificate: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.CertificateCAMutationResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Process
// ---------------------------------------------------------------------------

// ProcessList creates a step that lists running processes.
func (o *Orchestrator) ProcessList(
	target string,
) *Step {
	name := o.nextOpName("list-process")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Process.List(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list processes: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ProcessInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// ProcessGet creates a step that retrieves a specific process by PID.
func (o *Orchestrator) ProcessGet(
	target string,
	pid int,
) *Step {
	name := o.nextOpName("get-process")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Process.Get(ctx, target, pid)
			if err != nil {
				return nil, fmt.Errorf("get process: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ProcessInfoResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// ProcessSignal creates a step that sends a signal to a process.
func (o *Orchestrator) ProcessSignal(
	target string,
	pid int,
	opts osapi.ProcessSignalOpts,
) *Step {
	name := o.nextOpName("signal-process")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Process.Signal(ctx, target, pid, opts)
			if err != nil {
				return nil, fmt.Errorf("signal process: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.ProcessSignalResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Power
// ---------------------------------------------------------------------------

// PowerReboot creates a step that initiates a reboot.
func (o *Orchestrator) PowerReboot(
	target string,
	opts osapi.PowerOpts,
) *Step {
	name := o.nextOpName("reboot")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Power.Reboot(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("reboot: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PowerResult) engine.HostResult {
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

// PowerShutdown creates a step that initiates a shutdown.
func (o *Orchestrator) PowerShutdown(
	target string,
	opts osapi.PowerOpts,
) *Step {
	name := o.nextOpName("shutdown")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Power.Shutdown(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("shutdown: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.PowerResult) engine.HostResult {
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

// ---------------------------------------------------------------------------
// Log
// ---------------------------------------------------------------------------

// LogQuery creates a step that queries journal log entries.
func (o *Orchestrator) LogQuery(
	target string,
	opts osapi.LogQueryOpts,
) *Step {
	name := o.nextOpName("query-log")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Log.Query(ctx, target, opts)
			if err != nil {
				return nil, fmt.Errorf("query log: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.LogEntryResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// LogSources creates a step that lists available log sources.
func (o *Orchestrator) LogSources(
	target string,
) *Step {
	name := o.nextOpName("list-log-sources")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Log.Sources(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("list log sources: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.LogSourceResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}

// LogQueryUnit creates a step that queries log entries for a systemd unit.
func (o *Orchestrator) LogQueryUnit(
	target string,
	unit string,
	opts osapi.LogQueryOpts,
) *Step {
	name := o.nextOpName("query-log-unit")

	task := o.plan.TaskFunc(
		name,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*engine.Result, error) {
			resp, err := c.Log.QueryUnit(ctx, target, unit, opts)
			if err != nil {
				return nil, fmt.Errorf("query log unit: %w", err)
			}

			return engine.CollectionResult(
				resp.Data, resp.RawJSON(),
				func(r osapi.LogEntryResult) engine.HostResult {
					return engine.HostResult{
						Hostname: r.Hostname,
						Status:   r.Status,
						Error:    r.Error,
					}
				},
			)
		},
	)

	return &Step{task: task}
}
