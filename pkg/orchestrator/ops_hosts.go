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

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
)

// HostnameGet creates a step that runs Hostname.Get on the target: get the hostname of the host.
func (o *Orchestrator) HostnameGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-hostname", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.HostnameEntry]], error) {
			return c.Hostname.Get(ctx, target)
		},
	)
}

// HostnameUpdate creates a step that runs Hostname.Update on the target: set the hostname of the host.
func (o *Orchestrator) HostnameUpdate(
	target string,
	input osapi.HostnameUpdateInput,
) *Step {
	return collectionStep(
		o, "update-hostname", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.HostnameResult]], error) {
			return c.Hostname.Update(ctx, target, input)
		},
	)
}

// StatusGet creates a step that runs Status.Get on the target: get the status summary of the host.
func (o *Orchestrator) StatusGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-status", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.StatusEntry]], error) {
			return c.Status.Get(ctx, target)
		},
	)
}

// UptimeGet creates a step that runs Uptime.Get on the target: get the uptime of the host.
func (o *Orchestrator) UptimeGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-uptime", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UptimeEntry]], error) {
			return c.Uptime.Get(ctx, target)
		},
	)
}

// DiskGet creates a step that runs Disk.Get on the target: get the disk usage of the host.
func (o *Orchestrator) DiskGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-disk", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DiskEntry]], error) {
			return c.Disk.Get(ctx, target)
		},
	)
}

// MemoryGet creates a step that runs Memory.Get on the target: get the memory usage of the host.
func (o *Orchestrator) MemoryGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-memory", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.MemoryEntry]], error) {
			return c.Memory.Get(ctx, target)
		},
	)
}

// LoadGet creates a step that runs Load.Get on the target: get the load averages of the host.
func (o *Orchestrator) LoadGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-load", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.LoadEntry]], error) {
			return c.Load.Get(ctx, target)
		},
	)
}

// OSGet creates a step that runs OS.Get on the target: get the operating system of the host.
func (o *Orchestrator) OSGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-os", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.OsEntry]], error) {
			return c.OS.Get(ctx, target)
		},
	)
}

// DNSGet creates a step that runs DNS.Get on the target: get the DNS resolvers of an interface.
func (o *Orchestrator) DNSGet(
	target string,
	interfaceName string,
) *Step {
	return collectionStep(
		o, "get-dns", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DnsEntry]], error) {
			return c.DNS.Get(ctx, target, interfaceName)
		},
	)
}

// DNSUpdate creates a step that runs DNS.Update on the target: set the DNS resolvers of an interface.
func (o *Orchestrator) DNSUpdate(
	target string,
	input osapi.DNSUpdateInput,
) *Step {
	return collectionStep(
		o, "update-dns", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DnsResult]], error) {
			return c.DNS.Update(ctx, target, input)
		},
	)
}

// DNSDelete creates a step that runs DNS.Delete on the target: drop the DNS resolvers osapi set on.
func (o *Orchestrator) DNSDelete(
	target string,
	interfaceName string,
) *Step {
	return collectionStep(
		o, "delete-dns", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DnsResult]], error) {
			return c.DNS.Delete(ctx, target, interfaceName)
		},
	)
}

// ConnectivityPing creates a step that runs Connectivity.Ping on the target: ping an address from the.
func (o *Orchestrator) ConnectivityPing(
	target string,
	input osapi.ConnectivityPingInput,
) *Step {
	return collectionStep(
		o, "ping-connectivity", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ConnectivityPingResult]], error) {
			return c.Connectivity.Ping(ctx, target, input)
		},
	)
}

// InterfaceList creates a step that runs Interface.List on the target: list the network interfaces of the.
func (o *Orchestrator) InterfaceList(
	target string,
) *Step {
	return collectionStep(
		o, "list-interface", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceEntry]], error) {
			return c.Interface.List(ctx, target)
		},
	)
}

// InterfaceGet creates a step that runs Interface.Get on the target: get a network interface of the.
func (o *Orchestrator) InterfaceGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-interface", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceEntry]], error) {
			return c.Interface.Get(ctx, target, name)
		},
	)
}

// InterfaceCreate creates a step that runs Interface.Create on the target: configure a network interface.
func (o *Orchestrator) InterfaceCreate(
	target string,
	input osapi.InterfaceCreateInput,
) *Step {
	return collectionStep(
		o, "create-interface", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceResult]], error) {
			return c.Interface.Create(ctx, target, input)
		},
	)
}

// InterfaceUpdate creates a step that runs Interface.Update on the target: change the config osapi has of a.
func (o *Orchestrator) InterfaceUpdate(
	target string,
	input osapi.InterfaceUpdateInput,
) *Step {
	return collectionStep(
		o, "update-interface", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceResult]], error) {
			return c.Interface.Update(ctx, target, input)
		},
	)
}

// InterfaceDelete creates a step that runs Interface.Delete on the target: drop the config osapi has of a.
func (o *Orchestrator) InterfaceDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-interface", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceResult]], error) {
			return c.Interface.Delete(ctx, target, name)
		},
	)
}

// InterfaceRouteList creates a step that runs InterfaceRoute.List on the target: list the routes of an.
func (o *Orchestrator) InterfaceRouteList(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "list-interface-route", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceRouteEntry]], error) {
			return c.InterfaceRoute.List(ctx, target, name)
		},
	)
}

// InterfaceRouteCreate creates a step that runs InterfaceRoute.Create on the target: add a static route to an.
func (o *Orchestrator) InterfaceRouteCreate(
	target string,
	input osapi.InterfaceRouteCreateInput,
) *Step {
	return collectionStep(
		o, "create-interface-route", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceRouteResult]], error) {
			return c.InterfaceRoute.Create(ctx, target, input)
		},
	)
}

// InterfaceRouteUpdate creates a step that runs InterfaceRoute.Update on the target: change the gateway or the.
func (o *Orchestrator) InterfaceRouteUpdate(
	target string,
	input osapi.InterfaceRouteUpdateInput,
) *Step {
	return collectionStep(
		o, "update-interface-route", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceRouteResult]], error) {
			return c.InterfaceRoute.Update(ctx, target, input)
		},
	)
}

// InterfaceRouteDelete creates a step that runs InterfaceRoute.Delete on the target: drop a route that osapi.
func (o *Orchestrator) InterfaceRouteDelete(
	target string,
	input osapi.InterfaceRouteDeleteInput,
) *Step {
	return collectionStep(
		o, "delete-interface-route", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.InterfaceRouteResult]], error) {
			return c.InterfaceRoute.Delete(ctx, target, input)
		},
	)
}

// NodeFileGet creates a step that runs NodeFile.Get on the target: show whether a deployed file holds what.
func (o *Orchestrator) NodeFileGet(
	target string,
	path string,
) *Step {
	return collectionStep(
		o, "get-node-file", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.FileStatus]], error) {
			return c.NodeFile.Get(ctx, target, path)
		},
	)
}

// NodeFileUpdate creates a step that runs NodeFile.Update on the target: deploy an object of the object store to.
func (o *Orchestrator) NodeFileUpdate(
	target string,
	input osapi.NodeFileUpdateInput,
) *Step {
	return collectionStep(
		o, "update-node-file", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.FileUpdateResult]], error) {
			return c.NodeFile.Update(ctx, target, input)
		},
	)
}

// NodeFileDelete creates a step that runs NodeFile.Delete on the target: remove a file that osapi deployed.
func (o *Orchestrator) NodeFileDelete(
	target string,
	input osapi.NodeFileDeleteInput,
) *Step {
	return collectionStep(
		o, "delete-node-file", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.FileDeleteResult]], error) {
			return c.NodeFile.Delete(ctx, target, input)
		},
	)
}

// DockerList creates a step that runs Docker.List on the target: list the containers of the host.
func (o *Orchestrator) DockerList(
	target string,
	params osapi.DockerListParams,
) *Step {
	return collectionStep(
		o, "list-docker", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DockerEntry]], error) {
			return c.Docker.List(ctx, target, params)
		},
	)
}

// DockerGet creates a step that runs Docker.Get on the target: get one container.
func (o *Orchestrator) DockerGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-docker", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DockerEntry]], error) {
			return c.Docker.Get(ctx, target, name)
		},
	)
}

// DockerCreate creates a step that runs Docker.Create on the target: create a container.
func (o *Orchestrator) DockerCreate(
	target string,
	input osapi.DockerCreateInput,
) *Step {
	return collectionStep(
		o, "create-docker", true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DockerResult]], error) {
			return c.Docker.Create(ctx, target, input)
		},
	)
}

// DockerStart creates a step that runs Docker.Start on the target: start a container.
func (o *Orchestrator) DockerStart(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "start-docker", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DockerResult]], error) {
			return c.Docker.Start(ctx, target, name)
		},
	)
}

// DockerStop creates a step that runs Docker.Stop on the target: stop a container.
func (o *Orchestrator) DockerStop(
	target string,
	input osapi.DockerStopInput,
) *Step {
	return collectionStep(
		o, "stop-docker", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DockerResult]], error) {
			return c.Docker.Stop(ctx, target, input)
		},
	)
}

// DockerExec creates a step that runs Docker.Exec on the target: run a command in a container.
func (o *Orchestrator) DockerExec(
	target string,
	input osapi.DockerExecInput,
) *Step {
	return collectionStep(
		o, "exec-docker", true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DockerExecResult]], error) {
			return c.Docker.Exec(ctx, target, input)
		},
	)
}

// DockerDelete creates a step that runs Docker.Delete on the target: remove a container.
func (o *Orchestrator) DockerDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-docker", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.DockerResult]], error) {
			return c.Docker.Delete(ctx, target, name)
		},
	)
}

// ImagePull creates a step that runs Image.Pull on the target: pull an image.
func (o *Orchestrator) ImagePull(
	target string,
	input osapi.ImagePullInput,
) *Step {
	return collectionStep(
		o, "pull-image", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ImagePullResult]], error) {
			return c.Image.Pull(ctx, target, input)
		},
	)
}

// ImageDelete creates a step that runs Image.Delete on the target: remove an image.
func (o *Orchestrator) ImageDelete(
	target string,
	input osapi.ImageDeleteInput,
) *Step {
	return collectionStep(
		o, "delete-image", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ImageResult]], error) {
			return c.Image.Delete(ctx, target, input)
		},
	)
}

// ScheduleList creates a step that runs Schedule.List on the target: list the cron entries that osapi.
func (o *Orchestrator) ScheduleList(
	target string,
) *Step {
	return collectionStep(
		o, "list-schedule", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ScheduleEntry]], error) {
			return c.Schedule.List(ctx, target)
		},
	)
}

// ScheduleGet creates a step that runs Schedule.Get on the target: get one cron entry that osapi.
func (o *Orchestrator) ScheduleGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-schedule", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ScheduleEntry]], error) {
			return c.Schedule.Get(ctx, target, name)
		},
	)
}

// ScheduleCreate creates a step that runs Schedule.Create on the target: deploy a cron entry.
func (o *Orchestrator) ScheduleCreate(
	target string,
	input osapi.ScheduleCreateInput,
) *Step {
	return collectionStep(
		o, "create-schedule", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ScheduleResult]], error) {
			return c.Schedule.Create(ctx, target, input)
		},
	)
}

// ScheduleUpdate creates a step that runs Schedule.Update on the target: deploy a cron entry that osapi.
func (o *Orchestrator) ScheduleUpdate(
	target string,
	input osapi.ScheduleUpdateInput,
) *Step {
	return collectionStep(
		o, "update-schedule", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ScheduleResult]], error) {
			return c.Schedule.Update(ctx, target, input)
		},
	)
}

// ScheduleDelete creates a step that runs Schedule.Delete on the target: remove a cron entry that osapi.
func (o *Orchestrator) ScheduleDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-schedule", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ScheduleResult]], error) {
			return c.Schedule.Delete(ctx, target, name)
		},
	)
}

// SysctlList creates a step that runs Sysctl.List on the target: list the kernel's parameters and which.
func (o *Orchestrator) SysctlList(
	target string,
	params osapi.SysctlListParams,
) *Step {
	return collectionStep(
		o, "list-sysctl", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.SysctlEntry]], error) {
			return c.Sysctl.List(ctx, target, params)
		},
	)
}

// SysctlGet creates a step that runs Sysctl.Get on the target: get a kernel parameter's current.
func (o *Orchestrator) SysctlGet(
	target string,
	key string,
) *Step {
	return collectionStep(
		o, "get-sysctl", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.SysctlEntry]], error) {
			return c.Sysctl.Get(ctx, target, key)
		},
	)
}

// SysctlCreate creates a step that runs Sysctl.Create on the target: manage a kernel parameter and apply.
func (o *Orchestrator) SysctlCreate(
	target string,
	input osapi.SysctlCreateInput,
) *Step {
	return collectionStep(
		o, "create-sysctl", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.SysctlMutationResult]], error) {
			return c.Sysctl.Create(ctx, target, input)
		},
	)
}

// SysctlUpdate creates a step that runs Sysctl.Update on the target: set a kernel parameter, manage it,.
func (o *Orchestrator) SysctlUpdate(
	target string,
	input osapi.SysctlUpdateInput,
) *Step {
	return collectionStep(
		o, "update-sysctl", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.SysctlMutationResult]], error) {
			return c.Sysctl.Update(ctx, target, input)
		},
	)
}

// SysctlDelete creates a step that runs Sysctl.Delete on the target: stop the management of a kernel.
func (o *Orchestrator) SysctlDelete(
	target string,
	key string,
) *Step {
	return collectionStep(
		o, "delete-sysctl", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.SysctlMutationResult]], error) {
			return c.Sysctl.Delete(ctx, target, key)
		},
	)
}

// NTPGet creates a step that runs NTP.Get on the target: get the sync status of chrony and its.
func (o *Orchestrator) NTPGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-ntp", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.NtpStatus]], error) {
			return c.NTP.Get(ctx, target)
		},
	)
}

// NTPCreate creates a step that runs NTP.Create on the target: set the NTP servers that chrony uses.
func (o *Orchestrator) NTPCreate(
	target string,
	input osapi.NTPCreateInput,
) *Step {
	return collectionStep(
		o, "create-ntp", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.NtpResult]], error) {
			return c.NTP.Create(ctx, target, input)
		},
	)
}

// NTPUpdate creates a step that runs NTP.Update on the target: replace the NTP servers that osapi set.
func (o *Orchestrator) NTPUpdate(
	target string,
	input osapi.NTPUpdateInput,
) *Step {
	return collectionStep(
		o, "update-ntp", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.NtpResult]], error) {
			return c.NTP.Update(ctx, target, input)
		},
	)
}

// NTPDelete creates a step that runs NTP.Delete on the target: remove the NTP servers that osapi.
func (o *Orchestrator) NTPDelete(
	target string,
) *Step {
	return collectionStep(
		o, "delete-ntp", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.NtpResult]], error) {
			return c.NTP.Delete(ctx, target)
		},
	)
}

// TimezoneGet creates a step that runs Timezone.Get on the target: get the timezone of the host.
func (o *Orchestrator) TimezoneGet(
	target string,
) *Step {
	return collectionStep(
		o, "get-timezone", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.TimezoneEntry]], error) {
			return c.Timezone.Get(ctx, target)
		},
	)
}

// TimezoneUpdate creates a step that runs Timezone.Update on the target: set the timezone of the host.
func (o *Orchestrator) TimezoneUpdate(
	target string,
	input osapi.TimezoneUpdateInput,
) *Step {
	return collectionStep(
		o, "update-timezone", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.TimezoneResult]], error) {
			return c.Timezone.Update(ctx, target, input)
		},
	)
}

// ServiceList creates a step that runs Service.List on the target: list the services of the host.
func (o *Orchestrator) ServiceList(
	target string,
) *Step {
	return collectionStep(
		o, "list-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceEntry]], error) {
			return c.Service.List(ctx, target)
		},
	)
}

// ServiceGet creates a step that runs Service.Get on the target: get one service.
func (o *Orchestrator) ServiceGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceEntry]], error) {
			return c.Service.Get(ctx, target, name)
		},
	)
}

// ServiceCreate creates a step that runs Service.Create on the target: deploy a service.
func (o *Orchestrator) ServiceCreate(
	target string,
	input osapi.ServiceCreateInput,
) *Step {
	return collectionStep(
		o, "create-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Create(ctx, target, input)
		},
	)
}

// ServiceUpdate creates a step that runs Service.Update on the target: deploy a new unit file for a.
func (o *Orchestrator) ServiceUpdate(
	target string,
	input osapi.ServiceUpdateInput,
) *Step {
	return collectionStep(
		o, "update-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Update(ctx, target, input)
		},
	)
}

// ServiceDelete creates a step that runs Service.Delete on the target: remove a service that osapi.
func (o *Orchestrator) ServiceDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Delete(ctx, target, name)
		},
	)
}

// ServiceStart creates a step that runs Service.Start on the target: start a service that osapi.
func (o *Orchestrator) ServiceStart(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "start-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Start(ctx, target, name)
		},
	)
}

// ServiceStop creates a step that runs Service.Stop on the target: stop a service that osapi.
func (o *Orchestrator) ServiceStop(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "stop-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Stop(ctx, target, name)
		},
	)
}

// ServiceRestart creates a step that runs Service.Restart on the target: restart a service that.
func (o *Orchestrator) ServiceRestart(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "restart-service", true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Restart(ctx, target, name)
		},
	)
}

// ServiceEnable creates a step that runs Service.Enable on the target: make a service that osapi.
func (o *Orchestrator) ServiceEnable(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "enable-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Enable(ctx, target, name)
		},
	)
}

// ServiceDisable creates a step that runs Service.Disable on the target: make a service that osapi.
func (o *Orchestrator) ServiceDisable(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "disable-service", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
			return c.Service.Disable(ctx, target, name)
		},
	)
}

// PackageList creates a step that runs Package.List on the target: list the packages installed on the.
func (o *Orchestrator) PackageList(
	target string,
) *Step {
	return collectionStep(
		o, "list-package", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PackageInstalled]], error) {
			return c.Package.List(ctx, target)
		},
	)
}

// PackageGet creates a step that runs Package.Get on the target: get one installed package.
func (o *Orchestrator) PackageGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-package", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PackageInstalled]], error) {
			return c.Package.Get(ctx, target, name)
		},
	)
}

// PackageCreate creates a step that runs Package.Create on the target: install a package.
func (o *Orchestrator) PackageCreate(
	target string,
	input osapi.PackageCreateInput,
) *Step {
	return collectionStep(
		o, "create-package", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PackageResult]], error) {
			return c.Package.Create(ctx, target, input)
		},
	)
}

// PackageDelete creates a step that runs Package.Delete on the target: remove a package.
func (o *Orchestrator) PackageDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-package", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PackageResult]], error) {
			return c.Package.Delete(ctx, target, name)
		},
	)
}

// PackageUpdatesList creates a step that runs PackageUpdates.List on the target: list the packages the index has.
func (o *Orchestrator) PackageUpdatesList(
	target string,
) *Step {
	return collectionStep(
		o, "list-package-updates", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PackageUpdatesAvailable]], error) {
			return c.PackageUpdates.List(ctx, target)
		},
	)
}

// PackageUpdatesUpdate creates a step that runs PackageUpdates.Update on the target: refresh the package index.
func (o *Orchestrator) PackageUpdatesUpdate(
	target string,
) *Step {
	return collectionStep(
		o, "update-package-updates", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PackageUpdatesResult]], error) {
			return c.PackageUpdates.Update(ctx, target)
		},
	)
}

// UserList creates a step that runs User.List on the target: list the users of the host.
func (o *Orchestrator) UserList(
	target string,
) *Step {
	return collectionStep(
		o, "list-user", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserAccount]], error) {
			return c.User.List(ctx, target)
		},
	)
}

// UserGet creates a step that runs User.Get on the target: get one user.
func (o *Orchestrator) UserGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-user", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserAccount]], error) {
			return c.User.Get(ctx, target, name)
		},
	)
}

// UserCreate creates a step that runs User.Create on the target: add a user.
func (o *Orchestrator) UserCreate(
	target string,
	input osapi.UserCreateInput,
) *Step {
	return collectionStep(
		o, "create-user", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserResult]], error) {
			return c.User.Create(ctx, target, input)
		},
	)
}

// UserUpdate creates a step that runs User.Update on the target: change a user.
func (o *Orchestrator) UserUpdate(
	target string,
	input osapi.UserUpdateInput,
) *Step {
	return collectionStep(
		o, "update-user", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserResult]], error) {
			return c.User.Update(ctx, target, input)
		},
	)
}

// UserDelete creates a step that runs User.Delete on the target: remove a user.
func (o *Orchestrator) UserDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-user", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserResult]], error) {
			return c.User.Delete(ctx, target, name)
		},
	)
}

// UserPasswordUpdate creates a step that runs UserPassword.Update on the target: set the password of a user.
func (o *Orchestrator) UserPasswordUpdate(
	target string,
	input osapi.UserPasswordUpdateInput,
) *Step {
	return collectionStep(
		o, "update-user-password", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserPasswordResult]], error) {
			return c.UserPassword.Update(ctx, target, input)
		},
	)
}

// UserSSHKeyList creates a step that runs UserSSHKey.List on the target: list the authorized SSH keys of a.
func (o *Orchestrator) UserSSHKeyList(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "list-user-ssh-key", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserSshkeyKey]], error) {
			return c.UserSSHKey.List(ctx, target, name)
		},
	)
}

// UserSSHKeyCreate creates a step that runs UserSSHKey.Create on the target: add an authorized SSH key to a.
func (o *Orchestrator) UserSSHKeyCreate(
	target string,
	input osapi.UserSSHKeyCreateInput,
) *Step {
	return collectionStep(
		o, "create-user-ssh-key", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserSshkeyResult]], error) {
			return c.UserSSHKey.Create(ctx, target, input)
		},
	)
}

// UserSSHKeyDelete creates a step that runs UserSSHKey.Delete on the target: remove an authorized SSH key.
func (o *Orchestrator) UserSSHKeyDelete(
	target string,
	input osapi.UserSSHKeyDeleteInput,
) *Step {
	return collectionStep(
		o, "delete-user-ssh-key", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.UserSshkeyResult]], error) {
			return c.UserSSHKey.Delete(ctx, target, input)
		},
	)
}

// GroupList creates a step that runs Group.List on the target: list the groups of the host.
func (o *Orchestrator) GroupList(
	target string,
) *Step {
	return collectionStep(
		o, "list-group", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.GroupAccount]], error) {
			return c.Group.List(ctx, target)
		},
	)
}

// GroupGet creates a step that runs Group.Get on the target: get one group.
func (o *Orchestrator) GroupGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-group", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.GroupAccount]], error) {
			return c.Group.Get(ctx, target, name)
		},
	)
}

// GroupCreate creates a step that runs Group.Create on the target: add a group.
func (o *Orchestrator) GroupCreate(
	target string,
	input osapi.GroupCreateInput,
) *Step {
	return collectionStep(
		o, "create-group", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.GroupResult]], error) {
			return c.Group.Create(ctx, target, input)
		},
	)
}

// GroupUpdate creates a step that runs Group.Update on the target: set the members of a group.
func (o *Orchestrator) GroupUpdate(
	target string,
	input osapi.GroupUpdateInput,
) *Step {
	return collectionStep(
		o, "update-group", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.GroupResult]], error) {
			return c.Group.Update(ctx, target, input)
		},
	)
}

// GroupDelete creates a step that runs Group.Delete on the target: remove a group.
func (o *Orchestrator) GroupDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-group", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.GroupResult]], error) {
			return c.Group.Delete(ctx, target, name)
		},
	)
}

// CertificateList creates a step that runs Certificate.List on the target: list the certificate authorities.
func (o *Orchestrator) CertificateList(
	target string,
) *Step {
	return collectionStep(
		o, "list-certificate", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.CertificateEntry]], error) {
			return c.Certificate.List(ctx, target)
		},
	)
}

// CertificateGet creates a step that runs Certificate.Get on the target: get one certificate authority.
func (o *Orchestrator) CertificateGet(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "get-certificate", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.CertificateEntry]], error) {
			return c.Certificate.Get(ctx, target, name)
		},
	)
}

// CertificateCreate creates a step that runs Certificate.Create on the target: add a certificate authority to.
func (o *Orchestrator) CertificateCreate(
	target string,
	input osapi.CertificateCreateInput,
) *Step {
	return collectionStep(
		o, "create-certificate", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.CertificateResult]], error) {
			return c.Certificate.Create(ctx, target, input)
		},
	)
}

// CertificateUpdate creates a step that runs Certificate.Update on the target: write a certificate authority.
func (o *Orchestrator) CertificateUpdate(
	target string,
	input osapi.CertificateUpdateInput,
) *Step {
	return collectionStep(
		o, "update-certificate", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.CertificateResult]], error) {
			return c.Certificate.Update(ctx, target, input)
		},
	)
}

// CertificateDelete creates a step that runs Certificate.Delete on the target: remove a certificate.
func (o *Orchestrator) CertificateDelete(
	target string,
	name string,
) *Step {
	return collectionStep(
		o, "delete-certificate", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.CertificateResult]], error) {
			return c.Certificate.Delete(ctx, target, name)
		},
	)
}

// ProcessList creates a step that runs Process.List on the target: list the processes of the host.
func (o *Orchestrator) ProcessList(
	target string,
) *Step {
	return collectionStep(
		o, "list-process", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ProcessEntry]], error) {
			return c.Process.List(ctx, target)
		},
	)
}

// ProcessGet creates a step that runs Process.Get on the target: get one process.
func (o *Orchestrator) ProcessGet(
	target string,
	pid int,
) *Step {
	return collectionStep(
		o, "get-process", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ProcessEntry]], error) {
			return c.Process.Get(ctx, target, pid)
		},
	)
}

// ProcessSignal creates a step that runs Process.Signal on the target: send a signal to a process.
func (o *Orchestrator) ProcessSignal(
	target string,
	input osapi.ProcessSignalInput,
) *Step {
	return collectionStep(
		o, "signal-process", true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.ProcessResult]], error) {
			return c.Process.Signal(ctx, target, input)
		},
	)
}

// PowerReboot creates a step that runs Power.Reboot on the target: reboot the host.
func (o *Orchestrator) PowerReboot(
	target string,
	input osapi.PowerRebootInput,
) *Step {
	return collectionStep(
		o, "reboot-power", true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PowerResult]], error) {
			return c.Power.Reboot(ctx, target, input)
		},
	)
}

// PowerShutdown creates a step that runs Power.Shutdown on the target: shut the host down.
func (o *Orchestrator) PowerShutdown(
	target string,
	input osapi.PowerShutdownInput,
) *Step {
	return collectionStep(
		o, "shutdown-power", true,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.PowerResult]], error) {
			return c.Power.Shutdown(ctx, target, input)
		},
	)
}

// LogList creates a step that runs Log.List on the target: list the entries of the journal.
func (o *Orchestrator) LogList(
	target string,
	params osapi.LogListParams,
) *Step {
	return collectionStep(
		o, "list-log", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.LogEntry]], error) {
			return c.Log.List(ctx, target, params)
		},
	)
}

// LogSourceList creates a step that runs LogSource.List on the target: list the sources of the journal.
func (o *Orchestrator) LogSourceList(
	target string,
) *Step {
	return collectionStep(
		o, "list-log-source", false,
		func(
			ctx context.Context,
			c *osapi.Client,
		) (*osapi.Response[osapi.Collection[osapi.LogSourceSource]], error) {
			return c.LogSource.List(ctx, target)
		},
	)
}
