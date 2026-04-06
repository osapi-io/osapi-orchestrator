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

// Package main demonstrates NTP server management:
// create configuration, read status, update servers, then delete.
//
// Cleanup at the start ensures repeatability.
//
// DAG:
//
//	create-ntp
//	    └── get-ntp
//	            └── update-ntp
//	                    └── delete-ntp
//
// Run with: OSAPI_TOKEN="<jwt>" go run ntp.go
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

	// Cleanup any leftover NTP config from a previous run.
	oc := orchestrator.New(url, token)
	oc.NTPDelete("_any").ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// Create → get → update → delete.
	o := orchestrator.New(url, token)

	create := o.NTPCreate("_any", osapi.NtpCreateOpts{
		Servers: []string{"0.pool.ntp.org", "1.pool.ntp.org"},
	})

	get := o.NTPGet("_any").After(create)

	update := o.NTPUpdate("_any", osapi.NtpUpdateOpts{
		Servers: []string{"time.google.com", "time.cloudflare.com"},
	}).After(get)

	o.NTPDelete("_any").After(update)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var ntp osapi.NtpStatusResult
	if err := report.Decode("get-ntp", &ntp); err == nil {
		fmt.Printf("NTP servers: %v (synced: %v)\n", ntp.Servers, ntp.Synchronized)
	}
}
