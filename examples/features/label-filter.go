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

// Package main demonstrates filtering agents by labels and facts.
// Uses HasLabel to find agents with a specific label key-value pair
// and FactEquals to match agents by arbitrary fact values.
//
// DAG (per matching host):
//
//	health-check
//	    └── get-hostname (target=<host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run label-filter.go
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

	// Discover agents labeled as production web servers.
	agents, err := o.Discover(
		context.Background(),
		orchestrator.HasLabel("role", "web"),
		orchestrator.FactEquals("env", "prod"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d prod web agents\n", len(agents))

	for _, a := range agents {
		fmt.Printf("  %s (labels=%v)\n", a.Hostname, a.Labels)
	}

	health := o.HealthCheck("_any")

	for _, a := range agents {
		o.NodeHostnameGet(a.Hostname).After(health)
	}

	if _, err := o.Run(); err != nil {
		log.Fatal(err)
	}
}
