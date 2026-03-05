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

import "time"

// HostnameResult holds typed hostname output.
type HostnameResult struct {
	Hostname string            `json:"hostname"`
	Labels   map[string]string `json:"labels,omitempty"`
}

// DiskResult holds typed disk usage output.
type DiskResult struct {
	Disks []DiskUsage `json:"disks"`
}

// DiskUsage represents a single disk's usage.
type DiskUsage struct {
	Name  string `json:"name"`
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
}

// MemoryResult holds typed memory stats.
type MemoryResult struct {
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Cached uint64 `json:"cached"`
}

// LoadResult holds typed load averages.
type LoadResult struct {
	Load1  float32 `json:"load1"`
	Load5  float32 `json:"load5"`
	Load15 float32 `json:"load15"`
}

// CommandResult holds typed command execution output.
type CommandResult struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
}

// PingResult holds typed ping output.
type PingResult struct {
	PacketsSent     int     `json:"packets_sent"`
	PacketsReceived int     `json:"packets_received"`
	PacketLoss      float64 `json:"packet_loss"`
	Error           string  `json:"error,omitempty"`
}

// DNSConfigResult holds typed DNS configuration output.
type DNSConfigResult struct {
	DNSServers    []string `json:"dns_servers"`
	SearchDomains []string `json:"search_domains"`
}

// DNSUpdateResult holds typed DNS update output.
type DNSUpdateResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// AgentListResult holds typed agent list output.
type AgentListResult struct {
	Agents []AgentResult `json:"agents"`
	Total  int           `json:"total"`
}

// AgentResult holds typed agent details.
type AgentResult struct {
	Hostname      string            `json:"hostname"`
	Status        string            `json:"status"`
	Labels        map[string]string `json:"labels,omitempty"`
	Architecture  string            `json:"architecture,omitempty"`
	KernelVersion string            `json:"kernel_version,omitempty"`
	CPUCount      int               `json:"cpu_count,omitempty"`
	FQDN          string            `json:"fqdn,omitempty"`
	ServiceMgr    string            `json:"service_mgr,omitempty"`
	PackageMgr    string            `json:"package_mgr,omitempty"`
	LoadAverage   *AgentLoadAverage `json:"load_average,omitempty"`
	Memory        *AgentMemory      `json:"memory,omitempty"`
	OSInfo        *AgentOSInfo      `json:"os_info,omitempty"`
	Interfaces    []InterfaceResult `json:"interfaces,omitempty"`
	Uptime        string            `json:"uptime,omitempty"`
	StartedAt     time.Time         `json:"started_at,omitempty"`
	RegisteredAt  time.Time         `json:"registered_at,omitempty"`
	Facts         map[string]any    `json:"facts,omitempty"`
}

// AgentLoadAverage holds system load averages from agent heartbeat.
type AgentLoadAverage struct {
	OneMin     float32 `json:"one_min"`
	FiveMin    float32 `json:"five_min"`
	FifteenMin float32 `json:"fifteen_min"`
}

// AgentMemory holds memory usage from agent heartbeat.
type AgentMemory struct {
	Total int `json:"total"`
	Used  int `json:"used"`
	Free  int `json:"free"`
}

// AgentOSInfo holds operating system info from agent heartbeat.
type AgentOSInfo struct {
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
}

// InterfaceResult holds typed network interface info.
type InterfaceResult struct {
	Name   string `json:"name"`
	IPv4   string `json:"ipv4,omitempty"`
	IPv6   string `json:"ipv6,omitempty"`
	MAC    string `json:"mac,omitempty"`
	Family string `json:"family,omitempty"`
}
