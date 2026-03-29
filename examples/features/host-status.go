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

// Package main demonstrates host-level status guards for broadcast
// operations.
//
// A hostname update broadcasts to all hosts. Hosts that don't support
// the operation report "skipped" status. The orchestrator distinguishes
// skipped from failed:
//
//   - OnlyIfAnyHostSkipped: runs when a host was skipped (unsupported).
//   - OnlyIfAnyHostFailed: runs when a host reported an error.
//
// DAG:
//
//	update-hostname (_all, continue on error)
//	    ├── notify-skipped (only-if-any-host-skipped)
//	    └── recover-failed (only-if-any-host-failed)
//
// Run with: OSAPI_TOKEN="<jwt>" go run host-status.go
package main

import (
	"context"
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

	// Hostname update is unsupported on macOS/containers, so those
	// hosts report "skipped" while supported hosts succeed or fail.
	update := o.NodeHostnameUpdate("_all", "test-hostname").
		ContinueOnError()

	// Runs only if at least one host was skipped (unsupported platform).
	o.CommandExec("_any", "echo", "some-hosts-were-skipped").
		Named("notify-skipped").
		After(update).
		OnlyIfAnyHostSkipped()

	// Runs only if at least one host reported an error.
	o.CommandExec("_any", "echo", "running-recovery").
		Named("recover-failed").
		After(update).
		OnlyIfAnyHostFailed()

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
