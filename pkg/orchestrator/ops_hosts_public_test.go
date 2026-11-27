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

package orchestrator_test

import (
	"net/http"
	"strings"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

// TestHostOps runs each op that targets hosts against a fake controller:
// the request has the op's method and the host's path, and the report
// has the host, the job, the step's name and whether the host changed.
func (s *OpsPublicTestSuite) TestHostOps() {
	tests := []struct {
		name    string
		build   func(*orchestrator.Orchestrator) *orchestrator.Step
		method  string
		step    string
		changed bool
	}{
		{
			name:    "HostnameGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.HostnameGet("web-01") },
			method:  http.MethodGet,
			step:    "get-hostname",
			changed: false,
		},
		{
			name: "HostnameUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.HostnameUpdate("web-01", osapi.HostnameUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-hostname",
			changed: true,
		},
		{
			name:    "StatusGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.StatusGet("web-01") },
			method:  http.MethodGet,
			step:    "get-status",
			changed: false,
		},
		{
			name:    "UptimeGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.UptimeGet("web-01") },
			method:  http.MethodGet,
			step:    "get-uptime",
			changed: false,
		},
		{
			name:    "DiskGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.DiskGet("web-01") },
			method:  http.MethodGet,
			step:    "get-disk",
			changed: false,
		},
		{
			name:    "MemoryGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.MemoryGet("web-01") },
			method:  http.MethodGet,
			step:    "get-memory",
			changed: false,
		},
		{
			name:    "LoadGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.LoadGet("web-01") },
			method:  http.MethodGet,
			step:    "get-load",
			changed: false,
		},
		{
			name:    "OSGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.OSGet("web-01") },
			method:  http.MethodGet,
			step:    "get-os",
			changed: false,
		},
		{
			name:    "DNSGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.DNSGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-dns",
			changed: false,
		},
		{
			name: "DNSUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.DNSUpdate("web-01", osapi.DNSUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-dns",
			changed: true,
		},
		{
			name:    "DNSDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.DNSDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-dns",
			changed: true,
		},
		{
			name: "ConnectivityPing",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ConnectivityPing("web-01", osapi.ConnectivityPingInput{})
			},
			method:  http.MethodPost,
			step:    "ping-connectivity",
			changed: true,
		},
		{
			name:    "InterfaceList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.InterfaceList("web-01") },
			method:  http.MethodGet,
			step:    "list-interface",
			changed: false,
		},
		{
			name:    "InterfaceGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.InterfaceGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-interface",
			changed: false,
		},
		{
			name: "InterfaceCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.InterfaceCreate("web-01", osapi.InterfaceCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-interface",
			changed: true,
		},
		{
			name: "InterfaceUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.InterfaceUpdate("web-01", osapi.InterfaceUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-interface",
			changed: true,
		},
		{
			name:    "InterfaceDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.InterfaceDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-interface",
			changed: true,
		},
		{
			name:    "InterfaceRouteList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.InterfaceRouteList("web-01", "") },
			method:  http.MethodGet,
			step:    "list-interface-route",
			changed: false,
		},
		{
			name: "InterfaceRouteCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.InterfaceRouteCreate("web-01", osapi.InterfaceRouteCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-interface-route",
			changed: true,
		},
		{
			name: "InterfaceRouteUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.InterfaceRouteUpdate("web-01", osapi.InterfaceRouteUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-interface-route",
			changed: true,
		},
		{
			name: "InterfaceRouteDelete",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.InterfaceRouteDelete("web-01", osapi.InterfaceRouteDeleteInput{})
			},
			method:  http.MethodDelete,
			step:    "delete-interface-route",
			changed: true,
		},
		{
			name:    "NodeFileGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.NodeFileGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-node-file",
			changed: false,
		},
		{
			name: "NodeFileUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.NodeFileUpdate("web-01", osapi.NodeFileUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-node-file",
			changed: true,
		},
		{
			name: "NodeFileDelete",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.NodeFileDelete("web-01", osapi.NodeFileDeleteInput{})
			},
			method:  http.MethodDelete,
			step:    "delete-node-file",
			changed: true,
		},
		{
			name: "DockerList",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.DockerList("web-01", osapi.DockerListParams{})
			},
			method:  http.MethodGet,
			step:    "list-docker",
			changed: false,
		},
		{
			name:    "DockerGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.DockerGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-docker",
			changed: false,
		},
		{
			name: "DockerCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.DockerCreate("web-01", osapi.DockerCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-docker",
			changed: true,
		},
		{
			name:    "DockerStart",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.DockerStart("web-01", "") },
			method:  http.MethodPost,
			step:    "start-docker",
			changed: true,
		},
		{
			name: "DockerStop",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.DockerStop("web-01", osapi.DockerStopInput{})
			},
			method:  http.MethodPost,
			step:    "stop-docker",
			changed: true,
		},
		{
			name: "DockerExec",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.DockerExec("web-01", osapi.DockerExecInput{})
			},
			method:  http.MethodPost,
			step:    "exec-docker",
			changed: true,
		},
		{
			name:    "DockerDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.DockerDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-docker",
			changed: true,
		},
		{
			name: "ImagePull",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ImagePull("web-01", osapi.ImagePullInput{})
			},
			method:  http.MethodPost,
			step:    "pull-image",
			changed: true,
		},
		{
			name: "ImageDelete",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ImageDelete("web-01", osapi.ImageDeleteInput{})
			},
			method:  http.MethodDelete,
			step:    "delete-image",
			changed: true,
		},
		{
			name:    "ScheduleList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ScheduleList("web-01") },
			method:  http.MethodGet,
			step:    "list-schedule",
			changed: false,
		},
		{
			name:    "ScheduleGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ScheduleGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-schedule",
			changed: false,
		},
		{
			name: "ScheduleCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ScheduleCreate("web-01", osapi.ScheduleCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-schedule",
			changed: true,
		},
		{
			name: "ScheduleUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ScheduleUpdate("web-01", osapi.ScheduleUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-schedule",
			changed: true,
		},
		{
			name:    "ScheduleDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ScheduleDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-schedule",
			changed: true,
		},
		{
			name: "SysctlList",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.SysctlList("web-01", osapi.SysctlListParams{})
			},
			method:  http.MethodGet,
			step:    "list-sysctl",
			changed: false,
		},
		{
			name:    "SysctlGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.SysctlGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-sysctl",
			changed: false,
		},
		{
			name: "SysctlCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.SysctlCreate("web-01", osapi.SysctlCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-sysctl",
			changed: true,
		},
		{
			name: "SysctlUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.SysctlUpdate("web-01", osapi.SysctlUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-sysctl",
			changed: true,
		},
		{
			name:    "SysctlDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.SysctlDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-sysctl",
			changed: true,
		},
		{
			name:    "NTPGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.NTPGet("web-01") },
			method:  http.MethodGet,
			step:    "get-ntp",
			changed: false,
		},
		{
			name: "NTPCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.NTPCreate("web-01", osapi.NTPCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-ntp",
			changed: true,
		},
		{
			name: "NTPUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.NTPUpdate("web-01", osapi.NTPUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-ntp",
			changed: true,
		},
		{
			name:    "NTPDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.NTPDelete("web-01") },
			method:  http.MethodDelete,
			step:    "delete-ntp",
			changed: true,
		},
		{
			name:    "TimezoneGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.TimezoneGet("web-01") },
			method:  http.MethodGet,
			step:    "get-timezone",
			changed: false,
		},
		{
			name: "TimezoneUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.TimezoneUpdate("web-01", osapi.TimezoneUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-timezone",
			changed: true,
		},
		{
			name:    "ServiceList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceList("web-01") },
			method:  http.MethodGet,
			step:    "list-service",
			changed: false,
		},
		{
			name:    "ServiceGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-service",
			changed: false,
		},
		{
			name: "ServiceCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ServiceCreate("web-01", osapi.ServiceCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-service",
			changed: true,
		},
		{
			name: "ServiceUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ServiceUpdate("web-01", osapi.ServiceUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-service",
			changed: true,
		},
		{
			name:    "ServiceDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-service",
			changed: true,
		},
		{
			name:    "ServiceStart",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceStart("web-01", "") },
			method:  http.MethodPost,
			step:    "start-service",
			changed: true,
		},
		{
			name:    "ServiceStop",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceStop("web-01", "") },
			method:  http.MethodPost,
			step:    "stop-service",
			changed: true,
		},
		{
			name:    "ServiceRestart",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceRestart("web-01", "") },
			method:  http.MethodPost,
			step:    "restart-service",
			changed: true,
		},
		{
			name:    "ServiceEnable",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceEnable("web-01", "") },
			method:  http.MethodPost,
			step:    "enable-service",
			changed: true,
		},
		{
			name:    "ServiceDisable",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ServiceDisable("web-01", "") },
			method:  http.MethodPost,
			step:    "disable-service",
			changed: true,
		},
		{
			name:    "PackageList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.PackageList("web-01") },
			method:  http.MethodGet,
			step:    "list-package",
			changed: false,
		},
		{
			name:    "PackageGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.PackageGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-package",
			changed: false,
		},
		{
			name: "PackageCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.PackageCreate("web-01", osapi.PackageCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-package",
			changed: true,
		},
		{
			name:    "PackageDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.PackageDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-package",
			changed: true,
		},
		{
			name:    "PackageUpdatesList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.PackageUpdatesList("web-01") },
			method:  http.MethodGet,
			step:    "list-package-updates",
			changed: false,
		},
		{
			name:    "PackageUpdatesUpdate",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.PackageUpdatesUpdate("web-01") },
			method:  http.MethodPut,
			step:    "update-package-updates",
			changed: true,
		},
		{
			name:    "UserList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.UserList("web-01") },
			method:  http.MethodGet,
			step:    "list-user",
			changed: false,
		},
		{
			name:    "UserGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.UserGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-user",
			changed: false,
		},
		{
			name: "UserCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.UserCreate("web-01", osapi.UserCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-user",
			changed: true,
		},
		{
			name: "UserUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.UserUpdate("web-01", osapi.UserUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-user",
			changed: true,
		},
		{
			name:    "UserDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.UserDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-user",
			changed: true,
		},
		{
			name: "UserPasswordUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.UserPasswordUpdate("web-01", osapi.UserPasswordUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-user-password",
			changed: true,
		},
		{
			name:    "UserSSHKeyList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.UserSSHKeyList("web-01", "") },
			method:  http.MethodGet,
			step:    "list-user-ssh-key",
			changed: false,
		},
		{
			name: "UserSSHKeyCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.UserSSHKeyCreate("web-01", osapi.UserSSHKeyCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-user-ssh-key",
			changed: true,
		},
		{
			name: "UserSSHKeyDelete",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.UserSSHKeyDelete("web-01", osapi.UserSSHKeyDeleteInput{})
			},
			method:  http.MethodDelete,
			step:    "delete-user-ssh-key",
			changed: true,
		},
		{
			name:    "GroupList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.GroupList("web-01") },
			method:  http.MethodGet,
			step:    "list-group",
			changed: false,
		},
		{
			name:    "GroupGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.GroupGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-group",
			changed: false,
		},
		{
			name: "GroupCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.GroupCreate("web-01", osapi.GroupCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-group",
			changed: true,
		},
		{
			name: "GroupUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.GroupUpdate("web-01", osapi.GroupUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-group",
			changed: true,
		},
		{
			name:    "GroupDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.GroupDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-group",
			changed: true,
		},
		{
			name:    "CertificateList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.CertificateList("web-01") },
			method:  http.MethodGet,
			step:    "list-certificate",
			changed: false,
		},
		{
			name:    "CertificateGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.CertificateGet("web-01", "") },
			method:  http.MethodGet,
			step:    "get-certificate",
			changed: false,
		},
		{
			name: "CertificateCreate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.CertificateCreate("web-01", osapi.CertificateCreateInput{})
			},
			method:  http.MethodPost,
			step:    "create-certificate",
			changed: true,
		},
		{
			name: "CertificateUpdate",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.CertificateUpdate("web-01", osapi.CertificateUpdateInput{})
			},
			method:  http.MethodPut,
			step:    "update-certificate",
			changed: true,
		},
		{
			name:    "CertificateDelete",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.CertificateDelete("web-01", "") },
			method:  http.MethodDelete,
			step:    "delete-certificate",
			changed: true,
		},
		{
			name:    "ProcessList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ProcessList("web-01") },
			method:  http.MethodGet,
			step:    "list-process",
			changed: false,
		},
		{
			name:    "ProcessGet",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.ProcessGet("web-01", 0) },
			method:  http.MethodGet,
			step:    "get-process",
			changed: false,
		},
		{
			name: "ProcessSignal",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ProcessSignal("web-01", osapi.ProcessSignalInput{})
			},
			method:  http.MethodPost,
			step:    "signal-process",
			changed: true,
		},
		{
			name: "PowerReboot",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.PowerReboot("web-01", osapi.PowerRebootInput{})
			},
			method:  http.MethodPost,
			step:    "reboot-power",
			changed: true,
		},
		{
			name: "PowerShutdown",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.PowerShutdown("web-01", osapi.PowerShutdownInput{})
			},
			method:  http.MethodPost,
			step:    "shutdown-power",
			changed: true,
		},
		{
			name: "LogList",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.LogList("web-01", osapi.LogListParams{})
			},
			method:  http.MethodGet,
			step:    "list-log",
			changed: false,
		},
		{
			name:    "LogSourceList",
			build:   func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.LogSourceList("web-01") },
			method:  http.MethodGet,
			step:    "list-log-source",
			changed: false,
		},
		{
			name: "CommandExec",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.CommandExec("web-01", osapi.CommandExecInput{Command: "/usr/bin/uptime"})
			},
			method:  http.MethodPost,
			step:    "run-uptime",
			changed: true,
		},
		{
			name: "CommandShell",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.CommandShell(
					"web-01",
					osapi.CommandShellInput{Command: "/bin/echo hi | wc"},
				)
			},
			method:  http.MethodPost,
			step:    "shell-echo",
			changed: true,
		},
		{
			name: "CommandShell with no words",
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.CommandShell("web-01", osapi.CommandShellInput{})
			},
			method:  http.MethodPost,
			step:    "shell-shell",
			changed: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server, got := s.serve(http.StatusOK, hostRow)

			report, err := s.run(server, tc.build)
			s.Require().NoError(err)
			s.Require().Len(*got, 1)
			s.Equal(tc.method, (*got)[0].method)
			s.True(
				strings.HasPrefix((*got)[0].path, "/api/node/web-01/"),
				"path %s", (*got)[0].path,
			)

			s.Require().Len(report.Tasks, 1)
			task := report.Tasks[0]
			s.Equal(tc.step, task.Name)
			s.Equal("job-1", task.JobID)
			s.Equal(tc.changed, task.Changed)
			s.Require().Len(task.HostResults, 1)
			s.Equal("web-01", task.HostResults[0].Hostname)
			s.Equal("ok", task.HostResults[0].Status)
			s.Equal(tc.changed, task.HostResults[0].Changed)
			s.Equal("web-01", task.HostResults[0].Data["hostname"])
		})
	}
}

// TestHostOpRefusals covers what a host op does with an answer that is
// not a success.
func (s *OpsPublicTestSuite) TestHostOpRefusals() {
	tests := []struct {
		name      string
		status    int
		body      string
		build     func(*orchestrator.Orchestrator) *orchestrator.Step
		errSubstr string
		hostError string
	}{
		{
			name:      "The controller refuses the call",
			status:    http.StatusUnauthorized,
			body:      `{"error":"unauthorized"}`,
			build:     func(o *orchestrator.Orchestrator) *orchestrator.Step { return o.HostnameGet("web-01") },
			errSubstr: "get-hostname",
		},
		{
			name:   "A host fails",
			status: http.StatusOK,
			body:   `{"results":[{"hostname":"web-01","status":"failed","error":"no such unit"}],"job_id":"j"}`,
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.ServiceStart("web-01", "nginx")
			},
			errSubstr: "no such unit",
			hostError: "no such unit",
		},
		{
			name:   "A command exits non-zero",
			status: http.StatusOK,
			body:   `{"results":[{"hostname":"web-01","status":"ok","exit_code":3}],"job_id":"j"}`,
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.CommandExec("web-01", osapi.CommandExecInput{Command: "false"})
			},
			errSubstr: "exit code 3",
			hostError: "exit code 3",
		},
		{
			name:   "A command's error takes precedence over its exit code",
			status: http.StatusOK,
			body:   `{"results":[{"hostname":"web-01","status":"failed","exit_code":3,"error":"timed out"}],"job_id":"j"}`,
			build: func(o *orchestrator.Orchestrator) *orchestrator.Step {
				return o.CommandShell("web-01", osapi.CommandShellInput{Command: "sleep 9"})
			},
			errSubstr: "timed out",
			hostError: "timed out",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server, _ := s.serve(tc.status, tc.body)

			report, err := s.run(server, tc.build)
			s.Require().Error(err)
			s.Contains(err.Error(), tc.errSubstr)

			if tc.hostError != "" {
				s.Require().Len(report.Tasks, 1)
				s.Require().Len(report.Tasks[0].HostResults, 1)
				s.Equal(tc.hostError, report.Tasks[0].HostResults[0].Error)
			}
		})
	}
}
