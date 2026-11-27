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

// Package main demonstrates network interface management:
// list interfaces, then get details on a specific one.
//
// Uses ContinueOnError since Netplan-based interface management
// is only available on Debian-family systems.
//
// DAG:
//
//	list-interface
//	    └── get-interface
//
// Run with: OSAPI_TOKEN="<jwt>" OSAPI_INTERFACE="eth0" go run interface.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"

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

	list := o.InterfaceList("_any").ContinueOnError()

	get := o.InterfaceGet("_any", iface).After(list).ContinueOnError()

	// Create → update → delete a test interface.
	create := o.InterfaceCreate("_any", osapi.InterfaceCreateInput{Name: "osapi-test0"}).
		After(get).ContinueOnError()

	update := o.InterfaceUpdate("_any", osapi.InterfaceUpdateInput{Name: "osapi-test0"}).
		After(create).ContinueOnError()

	o.InterfaceDelete("_any", "osapi-test0").After(update).ContinueOnError()

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if ifaces := report.HostResults("list-interface"); len(ifaces) > 0 {
		fmt.Printf("Network interfaces (%d):\n", len(ifaces))
		for _, hr := range ifaces {
			var i osapi.InterfaceEntry
			if err := hr.Decode(&i); err == nil {
				fmt.Printf("  %s: ipv4=%s, state=%s\n", i.Name, i.IPv4, i.State)
			}
		}
	} else {
		fmt.Println("Interface management requires Netplan (Debian-family only)")
	}

	var detail osapi.InterfaceEntry
	if err := report.Decode("get-interface", &detail); err == nil && detail.Name != "" {
		fmt.Printf("\nInterface %s: dhcp4=%v, mtu=%d, mac=%s\n",
			detail.Name, detail.DHCP4, detail.MTU, detail.MAC)
	}
}
