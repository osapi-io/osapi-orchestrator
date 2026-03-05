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

// Package main demonstrates agent discovery with fact predicates.
// Discovers agents running Ubuntu on amd64, then retrieves the
// hostname from each matching host.
//
// DAG (per discovered host):
//
//	health-check
//	    └── get-hostname (target=<discovered host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run main.go
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

	// Discover Ubuntu agents at plan-build time.
	agents, err := o.Discover(
		context.Background(),
		orchestrator.OS("Ubuntu"),
		orchestrator.Arch("amd64"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Discovered %d matching agents\n", len(agents))

	health := o.HealthCheck("_any")

	// Create a hostname step for each discovered agent.
	for _, a := range agents {
		o.NodeHostnameGet(a.Hostname).After(health)
	}

	if _, err := o.Run(); err != nil {
		log.Fatal(err)
	}
}
