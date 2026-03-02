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

// Package main demonstrates the osapi-orchestrator user experience.
//
// Users declare what they want, where, in what order, and under what
// conditions. The orchestrator handles display, logging, DAG execution,
// polling, and error handling.
//
// DAG:
//
//	health-check
//	    ├── get-hostname ──┐
//	    ├── get-disk       │
//	    ├── get-memory     ├── whoami (only-if-changed, when hostname != "")
//	    ├── get-load [retry:2]
//	    └── run-uptime ────┘
//	configure-dns (independent)
//
// Run with: OSAPI_TOKEN="<jwt>" go run main.go
package main

import (
	"fmt"
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

	o := orchestrator.New(url, token)

	// Health check — must pass before anything else runs.
	health := o.HealthCheck("_any")

	// Fleet discovery — all run in parallel after health.
	hostname := o.NodeHostnameGet("_any").After(health)
	disk := o.NodeDiskGet("_any").After(health)
	memory := o.NodeMemoryGet("_any").After(health)

	load := o.NodeLoadGet("_any").
		After(health).
		Retry(2)

	uptime := o.CommandExec("_any", "uptime", "-s").
		After(health)

	// Conditional command — only if something changed, and only if
	// hostname was successfully retrieved.
	o.CommandExec("_any", "whoami").
		After(hostname, disk, memory, load, uptime).
		OnlyIfChanged().
		When(func(r orchestrator.Results) bool {
			var h orchestrator.HostnameResult
			if err := r.Decode("get-hostname", &h); err != nil {
				return false
			}

			return h.Hostname != ""
		})

	// Independent task — runs in parallel with everything above.
	o.NetworkDNSUpdate("group:web.prod", "eth0",
		[]string{"8.8.8.8", "8.8.4.4"},
		[]string{"example.com"},
	).OnError(orchestrator.Continue)

	// Run the plan.
	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"%s in %s\n",
		report.Summary(),
		report.Duration,
	)
}
