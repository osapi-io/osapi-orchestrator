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

// Package main demonstrates Docker pull idempotency and the complete
// container lifecycle.
//
// Run 1: pre-cleanup removes image → pull downloads it (changed=true)
//
//	→ create → start → list → exec + inspect → stop → container remove
//
// Run 2: image is cached → pull is a no-op (changed=false)
//
//	→ create → start → list → exec + inspect → stop → container remove + image remove
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

	// Run 1: remove image first so pull downloads it (changed=true).
	fmt.Println("=== Run 1: fresh pull ===\n")
	runLifecycle(url, token, true, false)

	// Run 2: image is cached so pull is a no-op (changed=false).
	// Clean up the image at the end.
	fmt.Println("\n=== Run 2: cached pull ===\n")
	runLifecycle(url, token, false, true)
}

func runLifecycle(
	url string,
	token string,
	removeImageFirst bool,
	removeImageLast bool,
) {
	o := orchestrator.New(url, token)

	// Pre-cleanup: remove leftover container (always) and optionally
	// remove the image to force a fresh pull. Errors are swallowed.
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

			if removeImageFirst {
				_, _ = c.Docker.ImageRemove(
					ctx,
					targetHost,
					imageName,
					&osapi.DockerImageRemoveParams{Force: true},
				)
			}

			return &sdk.Result{Changed: false}, nil
		},
	)

	// Pull image — changed=true on first run, changed=false on second.
	pull := o.DockerPull(targetHost, osapi.DockerPullOpts{
		Image: imageName,
	}).After(preCleanup)

	// Create container (without auto-start).
	create := o.DockerCreate(targetHost, osapi.DockerCreateOpts{
		Image: imageName,
		Name:  containerName,
	}).After(pull)

	// Start the container explicitly.
	start := o.DockerStart(targetHost, containerName).After(create)

	// List running containers to show the new one.
	o.DockerList(targetHost, nil).After(start)

	// Exec: show nginx version.
	exec := o.DockerExec(
		targetHost,
		containerName,
		osapi.DockerExecOpts{
			Command: []string{"nginx", "-v"},
		},
	).After(start)

	// Inspect the running container.
	inspect := o.DockerInspect(targetHost, containerName).After(start)

	// Stop the container after exec and inspect finish.
	stop := o.DockerStop(
		targetHost,
		containerName,
		osapi.DockerStopOpts{},
	).After(exec, inspect)

	// Remove container after stopping.
	containerRemove := o.DockerRemove(
		targetHost,
		containerName,
		&osapi.DockerRemoveParams{Force: true},
	).After(stop)

	// Optionally remove the image at the end (cleanup on last run).
	if removeImageLast {
		o.DockerImageRemove(
			targetHost,
			imageName,
			&osapi.DockerImageRemoveParams{Force: true},
		).After(containerRemove)
	}

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
