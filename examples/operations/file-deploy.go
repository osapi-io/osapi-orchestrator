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

// Package main demonstrates file deployment with idempotency proof.
//
// Phase 1: cleanup — remove any previously deployed file.
// Phase 2: first deploy — upload + deploy + verify (expect changed).
// Phase 3: idempotency — same upload + deploy + verify (expect unchanged).
// Phase 4: cleanup — remove the deployed file.
//
// Run with: OSAPI_TOKEN="<jwt>" go run file-deploy.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	osapi "github.com/retr0h/osapi/pkg/sdk/client"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

func newOrchestrator(
	url string,
	token string,
) *orchestrator.Orchestrator {
	return orchestrator.New(url, token)
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

	configData := []byte("server:\n  port: 8080\n  debug: false\n")

	deployOpts := osapi.FileDeployOpts{
		ObjectName:  "app-config.yaml",
		Path:        "/tmp/app-config.yaml",
		ContentType: "raw",
		Mode:        "0644",
		Owner:       "root",
		Group:       "root",
	}

	// Phase 1: cleanup — remove any previously deployed file.
	fmt.Println("=== Phase 1: Cleanup ===")

	o1 := newOrchestrator(url, token)
	o1.CommandShell("_any", "rm -f /tmp/app-config.yaml").OnError(orchestrator.Continue)
	//nolint:errcheck
	o1.Run(context.Background()) //nolint:errcheck

	// Phase 2: first deploy — upload + deploy + verify (expect changed).
	fmt.Println("\n=== Phase 2: First deploy (expect changed) ===")

	o2 := newOrchestrator(url, token)
	upload2 := o2.FileUpload("app-config.yaml", "raw", configData, orchestrator.WithForce())
	deploy2 := o2.FileDeploy("_any", deployOpts).After(upload2)
	o2.FileStatusGet("_any", "/tmp/app-config.yaml").After(deploy2)

	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var status2 osapi.FileStatusResult
	if err := report2.Decode("file-status", &status2); err == nil {
		fmt.Printf("File %s status: %s\n", status2.Path, status2.Status)
	}

	// Phase 3: idempotency — same upload + deploy + verify (expect unchanged).
	fmt.Println("\n=== Phase 3: Idempotency check (expect unchanged) ===")

	o3 := newOrchestrator(url, token)
	upload3 := o3.FileUpload("app-config.yaml", "raw", configData)
	deploy3 := o3.FileDeploy("_any", deployOpts).After(upload3)
	o3.FileStatusGet("_any", "/tmp/app-config.yaml").After(deploy3)

	report3, err := o3.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var status3 osapi.FileStatusResult
	if err := report3.Decode("file-status", &status3); err == nil {
		fmt.Printf("File %s status: %s\n", status3.Path, status3.Status)
	}

	// Phase 4: cleanup — remove the deployed file.
	fmt.Println("\n=== Phase 4: Cleanup ===")

	o4 := newOrchestrator(url, token)
	o4.CommandShell("_any", "rm -f /tmp/app-config.yaml").OnError(orchestrator.Continue)
	//nolint:errcheck
	o4.Run(context.Background()) //nolint:errcheck
}
