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

// Package main demonstrates a fleet inventory workflow that queries
// all node information in parallel.
//
// After a health check gate, seven node-info queries run concurrently:
// status, disk, memory, load, uptime, OS, and hostname.
//
// DAG:
//
//	health-check
//	    ├── get-status
//	    ├── get-disk
//	    ├── get-memory
//	    ├── get-load
//	    ├── get-uptime
//	    ├── get-os
//	    └── get-hostname
//
// Run with: OSAPI_TOKEN="<jwt>" go run node-info.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

func main() {
	token := os.Getenv("OSAPI_TOKEN")
	if token == "" {
		log.Fatal("OSAPI_TOKEN is required")
	}

	url := os.Getenv("OSAPI_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	o := orchestrator.New(url, token)

	health := o.HealthCheck()

	// All node queries run in parallel after the health check.
	o.StatusGet("_any").After(health)
	o.DiskGet("_any").After(health)
	o.MemoryGet("_any").After(health)
	o.LoadGet("_any").After(health)
	o.UptimeGet("_any").After(health)
	o.OSGet("_any").After(health)
	o.HostnameGet("_any").After(health)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Decode and display each result.
	var hostname osapi.HostnameEntry
	if err := report.Decode("get-hostname", &hostname); err == nil {
		fmt.Printf("Hostname: %s\n", hostname.Hostname)
	}

	var load osapi.LoadEntry
	if err := report.Decode("get-load", &load); err == nil && load.Hostname != "" {
		fmt.Printf("Load:     %.2f / %.2f / %.2f\n",
			load.LoadAverage.Min1, load.LoadAverage.Min5, load.LoadAverage.Min15)
	}

	var memory osapi.MemoryEntry
	if err := report.Decode("get-memory", &memory); err == nil && memory.Memory.Total > 0 {
		fmt.Printf("Memory:   %d total / %d free\n", memory.Memory.Total, memory.Memory.Free)
	}
}
