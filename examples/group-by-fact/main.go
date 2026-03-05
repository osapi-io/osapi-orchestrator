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

// Package main demonstrates grouping agents by a fact value.
// Groups the fleet by OS distribution and runs a distro-specific
// package update command on each group.
//
// DAG (per group, per host):
//
//	health-check
//	    └── shell-<update-cmd> (target=<host>)
//
// Run with: OSAPI_TOKEN="<jwt>" go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

func installCmd(
	distro string,
) string {
	switch distro {
	case "Ubuntu", "Debian":
		return "apt-get update -qq"
	case "CentOS", "Rocky", "AlmaLinux":
		return "yum check-update -q"
	default:
		return "echo unsupported"
	}
}

func main() {
	token := os.Getenv("OSAPI_TOKEN")
	if token == "" {
		log.Fatal("OSAPI_TOKEN is required")
	}

	url := os.Getenv("OSAPI_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	o := orchestrator.New(url, token)

	groups, err := o.GroupByFact(
		context.Background(),
		"os.distribution",
	)
	if err != nil {
		log.Fatal(err)
	}

	health := o.HealthCheck("_any")

	for distro, agents := range groups {
		cmd := installCmd(distro)
		fmt.Printf("Group %s (%d hosts): %s\n", distro, len(agents), cmd)

		for _, a := range agents {
			o.CommandShell(a.Hostname, cmd).After(health)
		}
	}

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s in %s\n", report.Summary(), report.Duration)
}
