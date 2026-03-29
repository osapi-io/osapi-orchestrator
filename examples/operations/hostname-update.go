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

// Package main demonstrates hostname read and update operations.
//
// Phase 1: reads the current hostname from a target host.
// Phase 2: updates the hostname on all hosts (broadcast).
// Both phases gate on a health check.
//
// DAG (phase 1):
//
//	health-check
//	    └── get-hostname
//
// DAG (phase 2):
//
//	health-check
//	    └── update-hostname (_all, continue on error)
//
// Run with: OSAPI_TOKEN="<jwt>" go run hostname-update.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	osapi "github.com/retr0h/osapi/pkg/sdk/client"

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

	fmt.Println("=== Phase 1: read current hostname ===")

	o1 := orchestrator.New(url, token)

	health1 := o1.HealthCheck()
	o1.NodeHostnameGet("_any").After(health1)

	report1, err := o1.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var hostname osapi.HostnameResult
	if err := report1.Decode("get-hostname", &hostname); err == nil {
		fmt.Printf("Current hostname: %s\n", hostname.Hostname)
	}

	fmt.Println("\n=== Phase 2: update hostname on all hosts ===")

	o2 := orchestrator.New(url, token)

	health2 := o2.HealthCheck()

	// Broadcast hostname update to all hosts. Continue allows the
	// plan to proceed even if some hosts are unsupported (skipped).
	o2.NodeHostnameUpdate("_all", "new-hostname").
		After(health2).
		OnError(orchestrator.Continue)

	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s in %s\n", report2.Summary(), report2.Duration)
}
