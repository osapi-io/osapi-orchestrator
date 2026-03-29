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

// Package main demonstrates an agent drain/undrain maintenance workflow.
//
// Phase 1: lists all registered agents to discover the fleet.
//
// Phase 2: gets detailed info about the target agent, then drains it
// (cordon), runs a maintenance command, and undrains it (uncordon).
//
// DAG (phase 1):
//
//	health-check
//	    └── list-agents
//
// DAG (phase 2):
//
//	health-check
//	    └── get-agent
//	        └── drain-agent
//	            └── run-apt-get (maintenance command)
//	                └── undrain-agent
//
// Run with: OSAPI_TOKEN="<jwt>" go run agent-drain.go
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

	host := os.Getenv("OSAPI_HOST")
	if host == "" {
		host = "web-01"
	}

	fmt.Println("=== Phase 1: discover registered agents ===")

	o1 := orchestrator.New(url, token)

	health1 := o1.HealthCheck()
	o1.AgentList().After(health1)

	report1, err := o1.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var agents osapi.AgentList
	if err := report1.Decode("list-agents", &agents); err == nil {
		fmt.Printf("Registered agents: %d\n", agents.Total)
		for _, a := range agents.Agents {
			fmt.Printf("  - %s (%s)\n", a.Hostname, a.Status)
		}
	}

	fmt.Printf("\n=== Phase 2: inspect %s, drain, run maintenance, undrain ===\n", host)

	o2 := orchestrator.New(url, token)

	health2 := o2.HealthCheck()

	// Get detailed agent info before draining.
	getAgent := o2.AgentGet(host).After(health2)

	// Drain the agent — prevents new jobs from being scheduled.
	drain := o2.AgentDrain(host).After(getAgent)

	// Run a maintenance command while the agent is drained.
	maint := o2.CommandExec(host, "apt-get", "update", "-y").After(drain)

	// Undrain the agent — allows new jobs again.
	o2.AgentUndrain(host).After(maint)

	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var cmd osapi.CommandResult
	if err := report2.Decode("run-apt-get", &cmd); err == nil {
		fmt.Printf("maintenance stdout: %s\n", cmd.Stdout)
	}

	fmt.Printf("\n%s in %s\n", report2.Summary(), report2.Duration)
}
