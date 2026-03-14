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

// Package main demonstrates Docker container lifecycle management
// with idempotency. The example runs the full lifecycle twice to show
// that pull returns changed=false on the second run (image already
// present).
//
// Phase 1: pre-cleanup (swallow errors)
// Phase 2: pull → create → exec + inspect → cleanup
// Phase 3: run again to demonstrate idempotent pull
//
// Run with: OSAPI_TOKEN="<jwt>" go run docker.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

const (
	containerName = "osapi-example-nginx"
	imageName     = "nginx:alpine"
	targetHost    = "_any"
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

	// Run the lifecycle twice to demonstrate idempotency.
	// First run: pull downloads the image (changed=true).
	// Second run: pull finds the image already present (changed=false).
	for i := range 2 {
		fmt.Printf("=== Run %d ===\n\n", i+1)
		runLifecycle(url, token)
		fmt.Println()
	}
}

func runLifecycle(
	url string,
	token string,
) {
	o := orchestrator.New(url, token)

	// Pre-cleanup: remove any leftover container. Swallow errors
	// so a missing container does not block the run.
	preCleanup := o.TaskFunc(
		"pre-cleanup",
		func(
			ctx context.Context,
			c *osapi.Client,
			_ orchestrator.Results,
		) (*sdk.Result, error) {
			_, _ = c.Docker.Remove(
				ctx,
				targetHost,
				containerName,
				&osapi.DockerRemoveParams{Force: true},
			)

			return &sdk.Result{Changed: false}, nil
		},
	)

	// Pull image.
	pull := o.DockerPull(targetHost, osapi.DockerPullOpts{
		Image: imageName,
	}).After(preCleanup)

	// Create container with auto-start.
	autoStart := true
	create := o.DockerCreate(targetHost, osapi.DockerCreateOpts{
		Image:     imageName,
		Name:      containerName,
		AutoStart: &autoStart,
	}).After(pull)

	// Exec: show nginx version.
	exec := o.DockerExec(
		targetHost,
		containerName,
		osapi.DockerExecOpts{
			Command: []string{"nginx", "-v"},
		},
	).After(create)

	// Inspect the running container.
	inspect := o.DockerInspect(targetHost, containerName).After(create)

	// Cleanup: force remove after exec and inspect finish.
	o.DockerRemove(
		targetHost,
		containerName,
		&osapi.DockerRemoveParams{Force: true},
	).After(exec, inspect)

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
