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

// Package main demonstrates error recovery at two levels.
//
// Plan 1: Infrastructure failure — TaskFunc returns an error.
//
//	OnlyIfFailed cleanup runs because the step failed at the
//	SDK level.
//
// Plan 2: Command failure — a command exits non-zero.
//
//	OnlyIfAnyHostFailed cleanup runs because non-zero exit codes
//	are treated as host errors by the orchestrator.
//
// Run with: OSAPI_TOKEN="<jwt>" go run error-recovery.go
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

	// Plan 1: Infrastructure failure (TaskFunc returns error)
	fmt.Println("=== Plan 1: Infrastructure failure (OnlyIfFailed) ===")
	o1 := orchestrator.New(url, token)
	deploy1 := o1.TaskFunc(
		"deploy",
		func(_ context.Context, _ *osapi.Client, _ orchestrator.Results) (*orchestrator.Result, error) {
			return nil, fmt.Errorf("simulated deployment failure")
		},
	).OnError(orchestrator.Continue)
	o1.CommandExec("_any", "echo", "running-infra-cleanup").
		Named("cleanup").
		After(deploy1).
		OnlyIfFailed()
	if _, err := o1.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Plan 2: Command failure (non-zero exit code)
	// The orchestrator treats non-zero exit codes as host errors,
	// so OnlyIfAnyHostFailed works naturally for command steps.
	fmt.Println("\n=== Plan 2: Command failure (OnlyIfAnyHostFailed) ===")
	o2 := orchestrator.New(url, token)
	deploy2 := o2.CommandShell("_any", "cat /nonexistent-file").
		Named("deploy").
		OnError(orchestrator.Continue)
	o2.CommandExec("_any", "echo", "running-command-cleanup").
		Named("cleanup").
		After(deploy2).
		OnlyIfAnyHostFailed()
	if _, err := o2.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
