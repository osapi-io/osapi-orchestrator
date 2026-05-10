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

// Package main demonstrates composing multiple fact predicates.
// Discovers agents that are Ubuntu, amd64, with at least 4 CPUs
// and 8GB memory, then queries their load averages.
//
// DAG (per matching host):
//
//	health-check
//	    └── get-load (target=<host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run fact-predicates.go
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

	// Compose predicates: Ubuntu + amd64 + 4 CPUs + 8GB.
	agents, err := o.Discover(
		context.Background(),
		orchestrator.OS("Ubuntu"),
		orchestrator.Arch("amd64"),
		orchestrator.MinCPU(4),
		orchestrator.MinMemory(8000),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d agents matching all predicates\n", len(agents))

	for _, a := range agents {
		fmt.Printf(
			"  %s (%s %s, %d CPUs)\n",
			a.Hostname,
			a.OSInfo.Distribution,
			a.Architecture,
			a.CPUCount,
		)
	}

	health := o.HealthCheck()

	for _, a := range agents {
		o.NodeLoadGet(a.Hostname).After(health)
	}

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
