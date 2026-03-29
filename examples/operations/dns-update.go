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

// Package main demonstrates the read-then-write DNS pattern.
//
// Reads current DNS config, then updates with new servers. All steps
// use OnError(Continue) so the example runs on any platform — on
// containers or macOS the DNS steps may be skipped or fail gracefully.
//
// DAG:
//
//	health-check
//	    └── get-dns (continue on error)
//	            └── update-dns (continue on error)
//
// Run with: OSAPI_TOKEN="<jwt>" go run dns-update.go
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

	iface := os.Getenv("OSAPI_INTERFACE")
	if iface == "" {
		iface = "eth0"
	}

	o := orchestrator.New(url, token)

	health := o.HealthCheck()

	// Read current DNS configuration.
	getDNS := o.NetworkDNSGet("_any", iface).
		After(health).
		OnError(orchestrator.Continue)

	// Write new DNS servers after reading the current config.
	o.NetworkDNSUpdate(
		"_any",
		iface,
		[]string{"8.8.8.8", "8.8.4.4"},
		[]string{"example.com"},
	).After(getDNS).OnError(orchestrator.Continue)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var dns osapi.DNSConfig
	if err := report.Decode("get-dns", &dns); err == nil {
		fmt.Printf("DNS servers:     %v\n", dns.Servers)
		fmt.Printf("Search domains:  %v\n", dns.SearchDomains)
	} else {
		fmt.Println("DNS operations require a valid network interface (set OSAPI_INTERFACE)")
	}
}
