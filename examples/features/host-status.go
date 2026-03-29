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

// Package main demonstrates host-level status guards for broadcast operations.
//
// A broadcast command targets all hosts. Some agents may report
// "skipped" (unsupported platform) or "failed" (error during execution).
// The orchestrator distinguishes these via host-level status:
//
//   - OnlyIfAnyHostSkipped: runs only when at least one host was skipped.
//   - OnlyIfAnyHostFailed: runs only when at least one host reported an error.
//
// DAG:
//
//	deploy (_all, continue on error)
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

	// Broadcast a command to all hosts. Continue allows the plan
	// to proceed even when some hosts fail or skip.
	deploy := o.CommandShell("_all", "cat /nonexistent-file").
		Named("deploy").
		OnError(orchestrator.Continue)

	// Notification step runs only if at least one host was skipped
	// (e.g., the operation is unsupported on that platform).
	o.CommandExec("_any", "echo", "some-hosts-were-skipped").
		Named("notify-skipped").
		After(deploy).
		OnlyIfAnyHostSkipped()

	// Recovery step runs only if at least one host reported an error
	// (e.g., the command exited non-zero).
	o.CommandExec("_any", "echo", "running-recovery").
		Named("recover-failed").
		After(deploy).
		OnlyIfAnyHostFailed()

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
