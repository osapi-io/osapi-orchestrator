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

// Package main demonstrates the full cron management lifecycle.
//
// Phase 1: upload a script, then create a cron entry that runs it daily.
// Phase 2: get the specific cron entry to inspect it.
// Phase 3: update the cron entry to run hourly instead of daily.
// Phase 4: list all cron entries to confirm the update.
// Phase 5: cleanup — delete the cron entry.
//
// DAG (phase 1):
//
//	upload-script
//	    └── create-cron
//
// DAG (phase 2):
//
//	get-cron
//
// DAG (phase 3):
//
//	update-cron
//
// DAG (phase 4):
//
//	list-cron
//
// DAG (phase 5):
//
//	delete-cron
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

	scriptData := []byte("#!/bin/bash\necho \"cleanup ran at $(date)\" >> /tmp/cleanup.log\n")

	// Phase 1: upload the script and create a daily cron entry.
	fmt.Println("=== Phase 1: Upload script + create cron entry ===")

	o1 := orchestrator.New(url, token)

	upload := o1.FileUpload("cleanup.sh", "raw", scriptData, orchestrator.WithForce())

	o1.CronCreate("_any", osapi.CronCreateOpts{
		Name:     "daily-cleanup",
		Object:   "cleanup.sh",
		Interval: "daily",
		User:     "root",
	}).After(upload)

	if _, err := o1.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Phase 2: get the specific cron entry to inspect it.
	fmt.Println("\n=== Phase 2: Get cron entry ===")

	o2 := orchestrator.New(url, token)
	o2.CronGet("_any", "daily-cleanup")

	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var entry osapi.CronEntryResult
	if err := report2.Decode("get-cron", &entry); err == nil {
		schedule := entry.Interval
		if schedule == "" {
			schedule = entry.Schedule
		}
		fmt.Printf("Cron entry %q: schedule=%s, user=%s\n",
			entry.Name, schedule, entry.User)
	}

	// Phase 3: update the cron entry to use a cron expression instead.
	fmt.Println("\n=== Phase 3: Update cron entry (daily -> custom schedule) ===")

	o3 := orchestrator.New(url, token)
	o3.CronUpdate("_any", "daily-cleanup", osapi.CronUpdateOpts{
		Schedule: "0 */2 * * *",
	})

	if _, err := o3.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cron entry updated to every 2 hours.")

	// Phase 4: list cron entries to confirm the update.
	fmt.Println("\n=== Phase 4: List cron entries ===")

	o4 := orchestrator.New(url, token)
	o4.CronList("_any")

	report4, err := o4.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var cron osapi.CronEntryResult
	if err := report4.Decode("list-cron", &cron); err == nil {
		fmt.Printf("Cron entries on %s: found\n", cron.Hostname)
	}

	// Phase 5: cleanup — delete the cron entry.
	fmt.Println("\n=== Phase 5: Delete cron entry ===")

	o5 := orchestrator.New(url, token)
	o5.CronDelete("_any", "daily-cleanup")

	if _, err := o5.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cron entry deleted.")
}
