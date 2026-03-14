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

// Package main demonstrates Docker container lifecycle management.
//
// Phase 1: pre-cleanup — force remove any leftover container (errors
// are swallowed so a missing container does not block the run).
//
// Phase 2: lifecycle — pull image, create and auto-start the container,
// exec a command inside it, then inspect the running container.
//
// Phase 3: cleanup — force remove the container.
//
// DAG (phase 2):
//
//	docker-pull
//	    └── docker-create
//	            ├── docker-exec
//	            └── docker-inspect
//
// Run with: OSAPI_TOKEN="<jwt>" go run docker.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	osapi "github.com/retr0h/osapi/pkg/sdk/client"

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

	// Phase 1: pre-cleanup — remove any leftover container, swallow errors.
	fmt.Println("=== Phase 1: Pre-cleanup ===")

	o1 := orchestrator.New(url, token)
	o1.DockerRemove(
		targetHost,
		containerName,
		&osapi.DockerRemoveParams{Force: true},
	).OnError(orchestrator.Continue)

	//nolint:errcheck
	o1.Run(context.Background())

	// Phase 2: lifecycle — pull, create, exec, inspect.
	fmt.Println("\n=== Phase 2: Container lifecycle ===")

	autoStart := true
	o2 := orchestrator.New(url, token)

	pull := o2.DockerPull(targetHost, osapi.DockerPullOpts{
		Image: imageName,
	})

	create := o2.DockerCreate(targetHost, osapi.DockerCreateOpts{
		Image:     imageName,
		Name:      containerName,
		AutoStart: &autoStart,
	}).After(pull)

	exec := o2.DockerExec(
		targetHost,
		containerName,
		osapi.DockerExecOpts{
			Command: []string{"nginx", "-v"},
		},
	).After(create)

	o2.DockerInspect(targetHost, containerName).After(create)

	report, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var pullResult osapi.DockerPullResult
	if err := report.Decode("docker-pull", &pullResult); err == nil {
		fmt.Printf("Pulled image ID: %s\n", pullResult.ImageID)
	}

	var createResult osapi.DockerResult
	if err := report.Decode("docker-create", &createResult); err == nil {
		fmt.Printf(
			"Container %q created (id=%s state=%s)\n",
			createResult.Name,
			createResult.ID,
			createResult.State,
		)
	}

	_ = exec

	var execResult osapi.DockerExecResult
	if err := report.Decode("docker-exec", &execResult); err == nil {
		fmt.Printf("nginx version: %s\n", execResult.Stdout)
	}

	var inspectResult osapi.DockerDetailResult
	if err := report.Decode("docker-inspect", &inspectResult); err == nil {
		fmt.Printf("Container state: %s\n", inspectResult.State)
	}

	// Phase 3: cleanup — force remove the container.
	fmt.Println("\n=== Phase 3: Cleanup ===")

	o3 := orchestrator.New(url, token)
	o3.DockerRemove(
		targetHost,
		containerName,
		&osapi.DockerRemoveParams{Force: true},
	)

	if _, err := o3.Run(context.Background()); err != nil {
		log.Printf("cleanup failed: %v", err)
	}

	fmt.Println("Done.")
}
