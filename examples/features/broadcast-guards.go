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

// Package main demonstrates broadcast-aware guards. A deploy step
// broadcasts to all hosts with Continue error strategy. Downstream
// steps use host-level guards to react to per-host outcomes.
//
// DAG:
//
//	deploy (_all, continue on error)
//	    ├── cleanup (only-if-any-host-failed)
//	    └── verify (only-if-all-hosts-changed)
//
// Run with: OSAPI_TOKEN="<jwt>" go run broadcast-guards.go
package main

import (
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

	// Broadcast deploy to all hosts. Continue allows the plan to
	// proceed even if some hosts fail.
	deploy := o.CommandExec("_all", "echo", "deploying").
		Named("deploy").
		OnError(orchestrator.Continue)

	// Cleanup runs only if at least one host reported an error.
	o.CommandExec("_any", "echo", "running-cleanup").
		Named("cleanup").
		After(deploy).
		OnlyIfAnyHostFailed()

	// Verification runs only if every host reported a change.
	o.CommandExec("_all", "echo", "verifying").
		Named("verify").
		After(deploy).
		OnlyIfAllHostsChanged()

	if _, err := o.Run(); err != nil {
		log.Fatal(err)
	}
}
