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
// Features shown:
//   - WithVerbose option for detailed output
//   - Typed operations (HealthCheck, NodeHostnameGet, CommandExec, etc.)
//   - Step chaining (After, Retry, OnlyIfChanged, When, OnError)
//   - TaskFunc for custom steps with inter-task data passing
//   - Status inspection in guards (Results.Status)
//   - OnlyIfAllChanged / OnlyIfFailed conditional execution
//   - Typed result decoding (Results.Decode, Report.Decode)
//   - Broadcast results via HostResults
//   - Error strategies (StopAll, Continue)
//
// DAG:
//
//	health-check
//	    ├── get-hostname ──┐
//	    ├── get-disk       │
//	    ├── get-memory     ├── whoami (only-if-changed, when hostname != "")
//	    ├── get-load [retry:2]
//	    └── run-uptime ────┘
//	    └── summarize (TaskFunc: reads hostname + uptime data) ──┐
//	                                                             └── notify (only-if-all-changed)
//	deploy-all (_all broadcast, continue on error)
//	    └── cleanup (only-if-failed: recovery step)
//
// Run with: OSAPI_TOKEN="<jwt>" go run main.go
// Verbose:  OSAPI_TOKEN="<jwt>" OSAPI_VERBOSE=1 go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"

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

	// --- Configuration ---
	// WithVerbose shows stdout, stderr, and response data for every task.
	var opts []orchestrator.Option
	if os.Getenv("OSAPI_VERBOSE") != "" {
		opts = append(opts, orchestrator.WithVerbose())
	}

	o := orchestrator.New(url, token, opts...)

	// --- Level 0: Health check ---
	health := o.HealthCheck("_any")

	// --- Level 1: Fleet discovery (parallel after health) ---
	hostname := o.NodeHostnameGet("_any").After(health)
	disk := o.NodeDiskGet("_any").After(health)
	memory := o.NodeMemoryGet("_any").After(health)

	load := o.NodeLoadGet("_any").
		After(health).
		Retry(2)

	uptime := o.CommandExec("_any", "uptime").
		After(health)

	// --- Level 2: Conditional command ---
	// Only runs if a dependency changed AND hostname was retrieved.
	// Uses Results.Decode for typed access in the When guard.
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

	// --- Level 2: Custom step with inter-task data passing ---
	// TaskFunc receives Results from prior steps. Reads hostname via
	// typed Decode and checks uptime status — returns aggregated data
	// that flows to downstream tasks and into Report.Tasks[].Data.
	summarize := o.TaskFunc(
		"summarize",
		func(_ context.Context, r orchestrator.Results) (*sdk.Result, error) {
			// Typed result decoding — no map[string]any digging.
			var h orchestrator.HostnameResult
			if err := r.Decode("get-hostname", &h); err != nil {
				return &sdk.Result{Changed: false}, nil
			}

			// Status inspection — check if uptime succeeded.
			uptimeOK := r.Status("run-uptime") == orchestrator.TaskStatusChanged

			return &sdk.Result{
				Changed: true,
				Data: map[string]any{
					"host":      h.Hostname,
					"uptime_ok": uptimeOK,
				},
			}, nil
		},
	).After(hostname, uptime)

	// --- Level 3: Only if all dependencies changed ---
	// This runs only if summarize reported changes (all deps changed).
	o.CommandExec("_any", "echo", "all-deps-changed").
		Named("notify").
		After(summarize).
		OnlyIfAllChanged()

	// --- Independent: Broadcast operation with error recovery ---
	// Targets all hosts via _all. After execution, per-host results
	// are displayed automatically by the renderer.
	deploy := o.CommandExec("_all", "echo", "deploying").
		Named("deploy-all").
		OnError(orchestrator.Continue)

	// Recovery step — only runs if deploy failed (OnlyIfFailed).
	// Demonstrates failure-triggered cleanup pattern.
	o.CommandExec("_any", "echo", "deploy-failed-cleanup").
		Named("cleanup").
		After(deploy).
		OnlyIfFailed()

	// --- Run ---
	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"\n%s in %s\n",
		report.Summary(),
		report.Duration,
	)

	// --- Post-execution: typed result decoding from Report ---
	// Report.Decode extracts typed results after the plan completes.
	var cmd orchestrator.CommandResult
	if err := report.Decode("run-uptime", &cmd); err == nil {
		fmt.Printf("uptime stdout: %s\n", cmd.Stdout)
		if cmd.Error != "" {
			fmt.Printf("uptime error:  %s\n", cmd.Error)
		}
	}
}
