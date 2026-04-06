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

// Package main demonstrates local group management:
// list → create → get → update → delete.
//
// Run with: OSAPI_TOKEN="<jwt>" go run group.go
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

	const groupName = "osapi-example-group"

	// Cleanup any leftover group from a previous run.
	oc := orchestrator.New(url, token)
	oc.GroupDelete("_any", groupName).ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// List → create → get → update → delete.
	o := orchestrator.New(url, token)

	list := o.GroupList("_any")

	create := o.GroupCreate("_any", osapi.GroupCreateOpts{
		Name: groupName,
	}).After(list)

	get := o.GroupGet("_any", groupName).After(create)

	update := o.GroupUpdate("_any", groupName, osapi.GroupUpdateOpts{}).
		After(get).ContinueOnError()

	o.GroupDelete("_any", groupName).After(update)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var info osapi.GroupInfoResult
	if err := report.Decode("get-group", &info); err == nil && len(info.Groups) > 0 {
		g := info.Groups[0]
		fmt.Printf("Group %q: gid=%d\n", g.Name, g.GID)
	}
}
