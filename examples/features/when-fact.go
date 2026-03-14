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

// Package main demonstrates the WhenFact execution-time guard.
// Lists agents, then conditionally runs a command only if the target
// agent is running Ubuntu.
//
// DAG:
//
//	health-check
//	    └── list-agents
//	            └── shell-apt-update (when-fact: os == Ubuntu)
//
// Run with: OSAPI_TOKEN="<jwt>" OSAPI_TARGET="<hostname>" go run when-fact.go
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

	target := os.Getenv("OSAPI_TARGET")
	if target == "" {
		target = "_any"
	}

	o := orchestrator.New(url, token)

	health := o.HealthCheck()
	agents := o.AgentList().After(health)

	// Guard: only run if the target agent is Ubuntu.
	o.CommandShell(target, "apt-get update -qq").
		After(agents).
		WhenFact("list-agents", func(a orchestrator.AgentResult) bool {
			return a.OSInfo != nil &&
				a.OSInfo.Distribution == "Ubuntu"
		})

	if _, err := o.Run(); err != nil {
		log.Fatal(err)
	}
}
