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
// pull → create → start → list + exec + inspect → stop → remove → image remove.
//
// Pre-cleanup removes any leftover container and image so the example
// is repeatable. Requires Docker on the target host.
//
// DAG:
//
//	pre-cleanup
//	    └── docker-pull
//	            └── docker-create
//	                    └── docker-start
//	                            ├── docker-list
//	                            ├── docker-exec
//	                            └── docker-inspect
//	                                    └── docker-stop
//	                                            └── docker-remove
//	                                                    └── docker-image-remove
//
// Run with: OSAPI_TOKEN="<jwt>" go run docker.go
package main

import (
	"context"
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

	const (
		image  = "nginx:alpine"
		name   = "osapi-example-nginx"
		target = "_any"
	)

	o := orchestrator.New(url, token)

	// Pre-cleanup: remove leftover container and image.
	preCleanup := o.TaskFunc(
		"pre-cleanup",
		func(
			ctx context.Context,
			c *osapi.Client,
			_ orchestrator.Results,
		) (*orchestrator.Result, error) {
			_, _ = c.Docker.Remove(ctx, target, name, &osapi.DockerRemoveParams{Force: true})
			_, _ = c.Docker.ImageRemove(
				ctx,
				target,
				image,
				&osapi.DockerImageRemoveParams{Force: true},
			)

			return &orchestrator.Result{Changed: false}, nil
		},
	)

	pull := o.DockerPull(target, osapi.DockerPullOpts{Image: image}).
		After(preCleanup)

	create := o.DockerCreate(target, osapi.DockerCreateOpts{
		Image: image,
		Name:  name,
	}).After(pull)

	start := o.DockerStart(target, name).After(create)

	// List, exec, and inspect run in parallel after start.
	o.DockerList(target, nil).After(start)

	exec := o.DockerExec(target, name, osapi.DockerExecOpts{
		Command: []string{"nginx", "-v"},
	}).After(start)

	inspect := o.DockerInspect(target, name).After(start)

	stop := o.DockerStop(target, name, osapi.DockerStopOpts{}).
		After(exec, inspect)

	remove := o.DockerRemove(target, name, &osapi.DockerRemoveParams{Force: true}).
		After(stop)

	o.DockerImageRemove(target, image, &osapi.DockerImageRemoveParams{Force: true}).
		After(remove)

	if _, err := o.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
