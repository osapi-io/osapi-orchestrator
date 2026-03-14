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

// Package main demonstrates parallel task execution. Multiple operations
// at the same DAG level run concurrently.
//
// The final level runs cat /proc/version which succeeds on Linux hosts
// but fails on macOS, demonstrating per-host error rendering.
//
// DAG:
//
//	health-check
//	    ├── get-hostname
//	    ├── get-disk
//	    ├── get-memory
//	    ├── get-load
//	    └── get-uptime
//	          └── shell-cat /proc/version (continue on error)
//
// Run with: OSAPI_TOKEN="<jwt>" go run parallel.go
package main

import (
	"context"
	"log"
	"os"

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

	// o := orchestrator.New(url, token)
	o := orchestrator.New(url, token, orchestrator.WithVerbose())

	// Level 0: health gate.
	health := o.HealthCheck()

	// Level 1: five queries run in parallel — all share the same
	// dependency so the orchestrator schedules them concurrently.
	o.NodeHostnameGet("_all").After(health)
	o.NodeDiskGet("_all").After(health)
	o.NodeMemoryGet("_all").After(health)
	o.NodeLoadGet("_all").After(health)
	uptime := o.NodeUptimeGet("_all").After(health)

	// Level 2: read /proc/version — exists on Linux, missing on
	// macOS. Demonstrates per-host error rendering with Continue.
	o.CommandShell("_all", "cat /proc/version").
		After(uptime).
		OnError(orchestrator.Continue)

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
