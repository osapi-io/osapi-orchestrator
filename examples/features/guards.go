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

// Package main demonstrates When guards — conditional execution
// based on prior step results.
//
// Plan 1: guard passes — whoami runs because hostname is non-empty.
// Plan 2: guard blocks — whoami is skipped because hostname doesn't
// match "impossible".
//
// Run with: OSAPI_TOKEN="<jwt>" go run guards.go
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

	// Plan 1: Guard should PASS
	fmt.Println("=== Plan 1: Guard should PASS (hostname != \"\") ===")
	o1 := orchestrator.New(url, token)
	health1 := o1.HealthCheck()
	hostname1 := o1.NodeHostnameGet("_any").After(health1)
	o1.CommandExec("_any", "whoami").
		After(hostname1).
		When(func(r orchestrator.Results) bool {
			var h osapi.HostnameResult
			if err := r.Decode("get-hostname", &h); err != nil {
				return false
			}
			return h.Hostname != ""
		})
	report1, err := o1.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %s\n", report1.Summary())

	// Plan 2: Guard should BLOCK
	fmt.Println("\n=== Plan 2: Guard should BLOCK (hostname == \"impossible\") ===")
	o2 := orchestrator.New(url, token)
	health2 := o2.HealthCheck()
	hostname2 := o2.NodeHostnameGet("_any").After(health2)
	o2.CommandExec("_any", "whoami").
		After(hostname2).
		When(func(r orchestrator.Results) bool {
			var h osapi.HostnameResult
			if err := r.Decode("get-hostname", &h); err != nil {
				return false
			}
			return h.Hostname == "impossible"
		})
	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %s\n", report2.Summary())
}
