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

import osapi "github.com/retr0h/osapi/pkg/sdk/client"

// Type aliases map orchestrator types to their SDK equivalents.
// This avoids duplicating types while preserving the public API.

// HostnameResult holds typed hostname output.
type HostnameResult = osapi.HostnameResult

// DiskResult holds typed disk usage output.
type DiskResult = osapi.DiskResult

// DiskUsage represents a single disk's usage.
type DiskUsage = osapi.Disk

// MemoryResult holds typed memory stats.
type MemoryResult = osapi.MemoryResult

// LoadResult holds typed load averages.
type LoadResult = osapi.LoadResult

// CommandResult holds typed command execution output.
type CommandResult = osapi.CommandResult

// PingResult holds typed ping output.
type PingResult = osapi.PingResult

// DNSConfigResult holds typed DNS configuration output.
type DNSConfigResult = osapi.DNSConfig

// DNSUpdateResult holds typed DNS update output.
type DNSUpdateResult = osapi.DNSUpdateResult

// AgentListResult holds typed agent list output.
type AgentListResult = osapi.AgentList

// AgentResult holds typed agent details.
type AgentResult = osapi.Agent

// ConditionResult holds a node condition from the agent.
type ConditionResult = osapi.Condition

// AgentLoadAverage holds system load averages from agent heartbeat.
type AgentLoadAverage = osapi.LoadAverage

// AgentMemory holds memory usage from agent heartbeat.
type AgentMemory = osapi.Memory

// AgentOSInfo holds operating system info from agent heartbeat.
type AgentOSInfo = osapi.OSInfo

// InterfaceResult holds typed network interface info.
type InterfaceResult = osapi.NetworkInterface

// FileDeployOpts holds parameters for the FileDeploy operation.
type FileDeployOpts = osapi.FileDeployOpts

// FileDeployResult holds typed file deploy output.
type FileDeployResult = osapi.FileDeployResult

// FileStatusResult holds typed file status output.
type FileStatusResult = osapi.FileStatusResult

// FileUploadResult holds typed file upload output.
type FileUploadResult = osapi.FileUpload

// FileChangedResult holds typed file change detection output.
type FileChangedResult = osapi.FileChanged
