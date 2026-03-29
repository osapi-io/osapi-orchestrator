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

// Package main demonstrates the cron management lifecycle:
// create an entry, inspect it, then delete it.
//
// Cleanup at the start ensures repeatability. The create step
// depends on a file upload (the script the cron entry will run).
//
// DAG:
//
//	upload-file
//	    └── create-cron
//	            └── get-cron
//	                    └── delete-cron
//
// Run with: OSAPI_TOKEN="<jwt>" go run cron.go
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

	script := []byte("#!/bin/bash\necho \"cleanup ran at $(date)\" >> /tmp/cleanup.log\n")

	// Cleanup any leftover entry from a previous run.
	oc := orchestrator.New(url, token)
	oc.CronDelete("_any", "daily-cleanup").ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// Create → inspect → delete in one plan.
	o := orchestrator.New(url, token)

	upload := o.FileUpload("cleanup.sh", "raw", script, orchestrator.WithForce())

	create := o.CronCreate("_any", osapi.CronCreateOpts{
		Name:     "daily-cleanup",
		Object:   "cleanup.sh",
		Interval: "daily",
		User:     "root",
	}).After(upload)

	get := o.CronGet("_any", "daily-cleanup").After(create)

	o.CronDelete("_any", "daily-cleanup").After(get)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var entry osapi.CronEntryResult
	if err := report.Decode("get-cron", &entry); err == nil {
		schedule := entry.Interval
		if schedule == "" {
			schedule = entry.Schedule
		}
		fmt.Printf("Entry %q: schedule=%s, user=%s\n", entry.Name, schedule, entry.User)
	}
}
