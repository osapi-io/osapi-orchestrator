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
// DAG:
//
//	health-check
//	    ├── get-hostname
//	    ├── get-disk
//	    ├── get-memory
//	    ├── get-load
//	    └── get-uptime
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

	// Level 0: health gate.
	health := o.HealthCheck("_any")

	// Level 1: five queries run in parallel — all share the same
	// dependency so the orchestrator schedules them concurrently.
	o.NodeHostnameGet("_any").After(health)
	o.NodeDiskGet("_any").After(health)
	o.NodeMemoryGet("_any").After(health)
	o.NodeLoadGet("_any").After(health)
	o.NodeUptimeGet("_any").After(health)

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s in %s\n", report.Summary(), report.Duration)
}
