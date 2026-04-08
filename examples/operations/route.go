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

// Package main demonstrates network route management:
// list routes on an interface and get a specific route.
//
// Uses ContinueOnError since route management requires Netplan
// and may not be available on all platforms.
//
// DAG:
//
//	list-route
//	    └── get-route
//
// Run with: OSAPI_TOKEN="<jwt>" OSAPI_INTERFACE="eth0" go run route.go
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

	list := o.RouteList("_any").ContinueOnError()

	get := o.RouteGet("_any", iface).After(list).ContinueOnError()

	// Create → update → delete a test route.
	create := o.RouteCreate("_any", iface, osapi.RouteConfigOpts{}).
		After(get).ContinueOnError()

	update := o.RouteUpdate("_any", iface, osapi.RouteConfigOpts{}).
		After(create).ContinueOnError()

	o.RouteDelete("_any", iface).After(update).ContinueOnError()

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var routes osapi.RouteListResult
	if err := report.Decode("list-route", &routes); err == nil {
		fmt.Printf("Routes on %s (%d):\n", iface, len(routes.Routes))
		for _, r := range routes.Routes {
			fmt.Printf("  %s via %s (metric: %d)\n",
				r.Destination, r.Gateway, r.Metric)
		}
	} else {
		fmt.Println("Route management requires Netplan (Debian-family only)")
	}
}
