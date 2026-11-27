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

// Package main demonstrates the Docker container lifecycle:
// pull → create → start → exec + inspect → stop → remove → image remove.
//
// Run with: OSAPI_TOKEN="<jwt>" go run docker.go
package main

import (
	"context"
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

	const (
		image  = "nginx:alpine"
		name   = "osapi-example-nginx"
		target = "_any"
	)

	// Cleanup: remove leftover container and image.
	cleanup := orchestrator.New(url, token)
	cleanup.DockerDelete(target, name).
		ContinueOnError()
	cleanup.ImageDelete(target, osapi.ImageDeleteInput{Image: image, Force: true}).
		ContinueOnError()
	_, _ = cleanup.Run(context.Background())

	// Main workflow.
	o := orchestrator.New(url, token)

	pull := o.ImagePull(target, osapi.ImagePullInput{Image: image})

	create := o.DockerCreate(target, osapi.DockerCreateInput{
		Image: image,
		Name:  name,
	}).After(pull)

	start := o.DockerStart(target, name).After(create)

	exec := o.DockerExec(target, osapi.DockerExecInput{
		Name:    name,
		Command: []string{"nginx", "-v"},
	}).After(start)

	inspect := o.DockerGet(target, name).After(start)
	o.DockerList(target, osapi.DockerListParams{}).After(start)

	stop := o.DockerStop(target, osapi.DockerStopInput{Name: name}).
		After(exec, inspect)

	remove := o.DockerDelete(target, name).
		After(stop)

	o.ImageDelete(target, osapi.ImageDeleteInput{Image: image, Force: true}).
		After(remove)

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
