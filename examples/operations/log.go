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

// Package main demonstrates journal log querying:
// query recent logs, list available sources, and query by unit.
//
// All steps use ContinueOnError since log access may require
// systemd-journald on the target.
//
// DAG:
//
//	query-log
//	    └── list-log-sources
//	            └── query-log-unit
//
// Run with: OSAPI_TOKEN="<jwt>" go run log.go
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

	lines := 10

	o := orchestrator.New(url, token)

	query := o.LogQuery("_any", osapi.LogQueryOpts{
		Lines: &lines,
	}).ContinueOnError()

	sources := o.LogSources("_any").After(query).ContinueOnError()

	o.LogQueryUnit("_any", "sshd", osapi.LogQueryOpts{
		Lines: &lines,
	}).After(sources).ContinueOnError()

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var logs osapi.LogEntryResult
	if err := report.Decode("query-log", &logs); err == nil {
		fmt.Printf("Recent log entries: %d\n", len(logs.Entries))
		for _, e := range logs.Entries {
			fmt.Printf("  [%s] %s: %s\n", e.Priority, e.Unit, e.Message)
		}
	} else {
		fmt.Println("Log operations require systemd-journald")
	}

	var src osapi.LogSourceResult
	if err := report.Decode("list-log-sources", &src); err == nil {
		fmt.Printf("Log sources: %d units\n", len(src.Sources))
	}
}
