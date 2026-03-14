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

// Package main demonstrates fleet discovery with AgentList and AgentGet.
// Lists all registered agents, then retrieves comprehensive facts for the
// first agent — OS info, load averages, memory stats, network interfaces,
// labels, and lifecycle timestamps.
//
// DAG:
//
//	health-check
//	    └── list-agents
//	            └── print-agent-facts (TaskFunc: decodes enriched agent data)
//
// Run with: OSAPI_TOKEN="<jwt>" go run agent-facts.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"

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

	health := o.HealthCheck()
	agents := o.AgentList().After(health)

	// Use a TaskFunc to decode the agent list and display rich facts.
	o.TaskFunc(
		"print-agent-facts",
		func(_ context.Context, r orchestrator.Results) (*sdk.Result, error) {
			var list osapi.AgentList
			if err := r.Decode("list-agents", &list); err != nil {
				return &sdk.Result{Changed: false}, nil
			}

			fmt.Printf("  Agents: %d total\n\n", list.Total)

			for _, a := range list.Agents {
				fmt.Printf("  Agent: %s\n", a.Hostname)
				fmt.Printf("    Status:       %s\n", a.Status)
				fmt.Printf("    Architecture: %s\n", a.Architecture)
				fmt.Printf("    Kernel:       %s\n", a.KernelVersion)
				fmt.Printf("    CPUs:         %d\n", a.CPUCount)
				fmt.Printf("    FQDN:         %s\n", a.Fqdn)
				fmt.Printf("    Package Mgr:  %s\n", a.PackageMgr)
				fmt.Printf("    Service Mgr:  %s\n", a.ServiceMgr)
				fmt.Printf("    Uptime:       %s\n", a.Uptime)
				fmt.Printf("    Labels:       %v\n", a.Labels)

				if !a.StartedAt.IsZero() {
					fmt.Printf("    Started:      %s\n",
						a.StartedAt.Format("2006-01-02 15:04:05"))
				}

				if !a.RegisteredAt.IsZero() {
					fmt.Printf("    Registered:   %s\n",
						a.RegisteredAt.Format("2006-01-02 15:04:05"))
				}

				if a.OSInfo != nil {
					fmt.Printf("    OS:           %s %s\n",
						a.OSInfo.Distribution, a.OSInfo.Version)
				}

				if a.LoadAverage != nil {
					fmt.Printf("    Load:         %.2f %.2f %.2f\n",
						a.LoadAverage.OneMin,
						a.LoadAverage.FiveMin,
						a.LoadAverage.FifteenMin)
				}

				if a.Memory != nil {
					fmt.Printf("    Memory:       total=%d  used=%d  free=%d\n",
						a.Memory.Total, a.Memory.Used, a.Memory.Free)
				}

				if len(a.Interfaces) > 0 {
					fmt.Printf("    Interfaces:\n")
					for _, iface := range a.Interfaces {
						fmt.Printf("      %-12s ipv4=%-15s mac=%s\n",
							iface.Name, iface.IPv4, iface.MAC)
					}
				}

				fmt.Println()
			}

			return &sdk.Result{Changed: false}, nil
		},
	).After(agents)

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
