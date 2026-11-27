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
	"strconv"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/osapi-io/osapi/pkg/sdk/types"

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
	// Signal a process only when OSAPI_PID names one: there is no signal that
	// is safe to send to an arbitrary process, such as PID 1.
	if pid, err := strconv.Atoi(os.Getenv("OSAPI_PID")); err == nil {
		o.ProcessSignal("_any", osapi.ProcessSignalInput{PID: pid, Signal: types.SignalUsr1}).
			After(get).ContinueOnError()
	}

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if procs := report.HostResults("list-process"); len(procs) > 0 {
		fmt.Printf("Running processes: %d\n", len(procs))
	}

	var init osapi.ProcessEntry
	if err := report.Decode("get-process", &init); err == nil && init.Name != "" {
		fmt.Printf("PID 1: %s (user: %s, state: %s)\n", init.Name, init.User, init.State)
	}
}
