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

// Package main demonstrates agent drain/undrain for maintenance.
//
// Discovers agents, then drains the first one (prevents new jobs),
// runs a maintenance command, and undrains it. Uses OnError(Continue)
// so repeated runs don't fail if the agent is already drained/undrained.
//
// DAG:
//
//	health-check
//	    └── undrain-agent (cleanup, continue on error)
//	            └── drain-agent
//	                    └── run-maintenance
//	                            └── undrain-agent
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

	// Discover the first available agent.
	o1 := orchestrator.New(url, token)
	o1.AgentList()

	report1, err := o1.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var agents osapi.AgentList
	if err := report1.Decode("list-agents", &agents); err != nil || len(agents.Agents) == 0 {
		log.Fatal("no agents registered")
	}

	host := os.Getenv("OSAPI_HOST")
	if host == "" {
		host = agents.Agents[0].Hostname
	}

	fmt.Printf("Target: %s\n\n", host)

	// Cleanup: ensure the agent is undrained from any previous run.
	oc := orchestrator.New(url, token)
	oc.AgentUndrain(host).ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// Drain → maintenance → undrain.
	o2 := orchestrator.New(url, token)

	drain := o2.AgentDrain(host)

	maint := o2.CommandExec(host, "echo", "maintenance complete").
		After(drain)

	o2.AgentUndrain(host).After(maint)

	if _, err := o2.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
