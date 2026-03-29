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

// Package main demonstrates the cron management lifecycle.
//
// Phase 1: upload a script, then create a cron entry that runs it daily.
// Phase 2: list all cron entries to confirm creation.
// Phase 3: cleanup — delete the cron entry.
//
// DAG (phase 1):
//
//	upload-script
//	    └── create-cron
//
// DAG (phase 2):
//
//	list-cron
//
// DAG (phase 3):
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

	// Phase 2: list cron entries to confirm creation.
	fmt.Println("\n=== Phase 2: List cron entries ===")

	o2 := orchestrator.New(url, token)
	o2.CronList("_any")

	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var cron osapi.CronEntryResult
	if err := report2.Decode("list-cron", &cron); err == nil {
		fmt.Printf("Cron entries on %s: found\n", cron.Hostname)
	}

	// Phase 3: cleanup — delete the cron entry.
	fmt.Println("\n=== Phase 3: Delete cron entry ===")

	o3 := orchestrator.New(url, token)
	o3.CronDelete("_any", "daily-cleanup")

	if _, err := o3.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cron entry deleted.")
}
