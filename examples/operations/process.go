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

// Package main demonstrates process management:
// list running processes and inspect PID 1.
//
// DAG:
//
//	list-process
//	    └── get-process
//
// Run with: OSAPI_TOKEN="<jwt>" go run process.go
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

	o := orchestrator.New(url, token)

	list := o.ProcessList("_any")

	get := o.ProcessGet("_any", 1).After(list)

	// Send SIGCONT to PID 1 (harmless no-op for init).
	o.ProcessSignal("_any", 1, osapi.ProcessSignalOpts{
		Signal: "SIGCONT",
	}).After(get).ContinueOnError()

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var procs osapi.ProcessInfoResult
	if err := report.Decode("list-process", &procs); err == nil {
		fmt.Printf("Running processes: %d\n", len(procs.Processes))
	}

	var init osapi.ProcessInfoResult
	if err := report.Decode("get-process", &init); err == nil && len(init.Processes) > 0 {
		p := init.Processes[0]
		fmt.Printf("PID 1: %s (user: %s, state: %s)\n", p.Name, p.User, p.State)
	}
}
