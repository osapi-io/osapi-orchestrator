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

// Package main demonstrates TaskFunc for custom logic and inter-task
// data passing. The summarize step reads typed results from prior
// steps and returns aggregated data available in Report.Decode().
//
// DAG:
//
//	health-check
//	    └── get-hostname
//	            └── summarize (TaskFunc: reads hostname data)
//
// Run with: OSAPI_TOKEN="<jwt>" go run task-func.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	health := o.HealthCheck("_any")
	hostname := o.NodeHostnameGet("_any").After(health)

	// TaskFunc receives Results from prior steps.
	o.TaskFunc(
		"summarize",
		func(_ context.Context, r orchestrator.Results) (*sdk.Result, error) {
			var h orchestrator.HostnameResult
			if err := r.Decode("get-hostname", &h); err != nil {
				return &sdk.Result{Changed: false}, nil
			}

			fmt.Printf("  Hostname: %s\n", h.Hostname)

			return &sdk.Result{
				Changed: true,
				Data: map[string]any{
					"host": h.Hostname,
				},
			}, nil
		},
	).After(hostname)

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Post-execution: decode typed results from the report.
	var h orchestrator.HostnameResult
	if err := report.Decode("get-hostname", &h); err == nil {
		fmt.Printf("Report hostname: %s\n", h.Hostname)
	}
}
