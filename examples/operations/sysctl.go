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

// Package main demonstrates sysctl parameter management:
// create a parameter, read it, update it, then delete it.
//
// Cleanup at the start ensures repeatability.
//
// DAG:
//
//	create-sysctl
//	    └── get-sysctl
//	            └── update-sysctl
//	                    └── delete-sysctl
//
// Run with: OSAPI_TOKEN="<jwt>" go run sysctl.go
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

	const key = "net.ipv4.ip_forward"

	// Cleanup any leftover parameter from a previous run.
	oc := orchestrator.New(url, token)
	oc.SysctlDelete("_any", key).ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// List → create → get → update → delete.
	o := orchestrator.New(url, token)

	syslist := o.SysctlList("_any")

	create := o.SysctlCreate("_any", osapi.SysctlCreateOpts{
		Key:   key,
		Value: "0",
	}).After(syslist)

	get := o.SysctlGet("_any", key).After(create)

	update := o.SysctlUpdate("_any", key, osapi.SysctlUpdateOpts{
		Value: "1",
	}).After(get)

	o.SysctlDelete("_any", key).After(update)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var entry osapi.SysctlEntryResult
	if err := report.Decode("get-sysctl", &entry); err == nil {
		fmt.Printf("Sysctl %s = %s\n", entry.Key, entry.Value)
	}
}
