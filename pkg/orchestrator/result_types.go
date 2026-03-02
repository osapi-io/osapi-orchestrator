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
}

// PingResult holds typed ping output.
type PingResult struct {
	PacketsSent     int     `json:"packets_sent"`
	PacketsReceived int     `json:"packets_received"`
	PacketLoss      float64 `json:"packet_loss"`
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
}
