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

// Package main demonstrates OnlyIfChanged with file operations.
//
// Phase 1: cleanup — remove deployed file.
// Phase 2: deploy file (changed=true) → OnlyIfChanged command runs.
// Phase 3: deploy same file (changed=false) → OnlyIfChanged command
//
//	is skipped.
//
// Phase 4: cleanup.
//
// Run with: OSAPI_TOKEN="<jwt>" go run only-if-changed.go
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

	content := []byte("only-if-changed test content\n")

	// Phase 1: Cleanup
	fmt.Println("=== Phase 1: Cleanup ===")
	o1 := orchestrator.New(url, token)
	o1.CommandShell("_any", "rm -f /tmp/only-if-changed.txt").
		OnError(orchestrator.Continue)
	o1.Run(context.Background()) //nolint:errcheck

	// Phase 2: Deploy (expect changed, post-deploy runs)
	fmt.Println("\n=== Phase 2: Deploy file (expect changed, post-deploy runs) ===")
	o2 := orchestrator.New(url, token)
	upload2 := o2.FileUpload("only-if-changed.txt", "raw", content, orchestrator.WithForce())
	deploy2 := o2.FileDeploy("_any", osapi.FileDeployOpts{
		ObjectName:  "only-if-changed.txt",
		Path:        "/tmp/only-if-changed.txt",
		ContentType: "raw",
		Mode:        "0644",
	}).After(upload2)
	o2.CommandExec("_any", "echo", "post-deploy-hook").
		Named("post-deploy").
		After(deploy2).
		OnlyIfChanged()
	if _, err := o2.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Phase 3: Same deploy (expect unchanged, post-deploy skipped)
	fmt.Println("\n=== Phase 3: Same deploy (expect unchanged, post-deploy skipped) ===")
	o3 := orchestrator.New(url, token)
	upload3 := o3.FileUpload("only-if-changed.txt", "raw", content)
	deploy3 := o3.FileDeploy("_any", osapi.FileDeployOpts{
		ObjectName:  "only-if-changed.txt",
		Path:        "/tmp/only-if-changed.txt",
		ContentType: "raw",
		Mode:        "0644",
	}).After(upload3)
	o3.CommandExec("_any", "echo", "post-deploy-hook").
		Named("post-deploy").
		After(deploy3).
		OnlyIfChanged()
	if _, err := o3.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Phase 4: Cleanup
	fmt.Println("\n=== Phase 4: Cleanup ===")
	o4 := orchestrator.New(url, token)
	o4.CommandShell("_any", "rm -f /tmp/only-if-changed.txt").
		OnError(orchestrator.Continue)
	o4.Run(context.Background()) //nolint:errcheck
}
