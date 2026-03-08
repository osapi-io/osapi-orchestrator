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

//go:build ignore

// Package main demonstrates condition-based agent filtering.
// Discovers agents without DiskPressure, then retrieves the
// hostname from each healthy host.
//
// DAG (per discovered host):
//
//	health-check
//	    └── get-hostname (target=<healthy host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run condition-filter.go
package main

import (
	"context"
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

	// Discover agents without DiskPressure at plan-build time.
	agents, err := o.Discover(
		context.Background(),
		orchestrator.NoCondition("DiskPressure"),
		orchestrator.Healthy(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Discovered %d healthy agents\n", len(agents))

	health := o.HealthCheck("_any")

	// Create a hostname step for each healthy agent.
	for _, a := range agents {
		o.NodeHostnameGet(a.Hostname).After(health)
	}

	if _, err := o.Run(); err != nil {
		log.Fatal(err)
	}
}
