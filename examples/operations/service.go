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

// Package main demonstrates systemd service lifecycle:
// create → start → enable → get → list → restart → disable → stop → delete.
//
// Run with: OSAPI_TOKEN="<jwt>" go run service.go
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

	const svcName = "osapi-example"

	unitFile := []byte(`[Unit]
Description=OSAPI Example Service

[Service]
Type=oneshot
ExecStart=/bin/true
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
`)

	// Cleanup any leftover service from a previous run.
	oc := orchestrator.New(url, token)
	oc.ServiceDisable("_any", svcName).ContinueOnError()
	oc.ServiceStop("_any", svcName).ContinueOnError()
	oc.ServiceDelete("_any", svcName).ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// Upload → create → start → enable → get → disable → stop → delete.
	o := orchestrator.New(url, token)

	upload := o.FileUpload("osapi-example.service", "raw", unitFile,
		orchestrator.WithForce())

	create := o.ServiceCreate("_any", osapi.ServiceCreateOpts{
		Name:   svcName,
		Object: "osapi-example.service",
	}).After(upload)

	start := o.ServiceStart("_any", svcName).After(create)

	enable := o.ServiceEnable("_any", svcName).After(start)

	get := o.ServiceGet("_any", svcName).After(enable)

	update := o.ServiceUpdate("_any", svcName, osapi.ServiceUpdateOpts{
		Object: "osapi-example.service",
	}).After(get)

	list := o.ServiceList("_any").After(update)

	restart := o.ServiceRestart("_any", svcName).After(list)

	disable := o.ServiceDisable("_any", svcName).After(restart)

	stop := o.ServiceStop("_any", svcName).After(disable)

	o.ServiceDelete("_any", svcName).After(stop)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var svc osapi.ServiceGetResult
	if err := report.Decode("get-service", &svc); err == nil && svc.Service != nil {
		fmt.Printf("Service %q: status=%s, enabled=%v\n",
			svc.Service.Name, svc.Service.Status, svc.Service.Enabled)
	}
}
