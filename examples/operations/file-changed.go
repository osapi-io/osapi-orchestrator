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

// Package main demonstrates drift detection with FileChanged.
//
// Phase 1: upload known content to establish baseline.
// Phase 2: FileChanged with same content — no drift (changed=false),
//
//	OnlyIfChanged steps are skipped.
//
// Phase 3: FileChanged with different content — drift detected
//
//	(changed=true), OnlyIfChanged steps run.
//
// Run with: OSAPI_TOKEN="<jwt>" go run file-changed.go
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

	contentA := []byte("server:\n  port: 8080\n  debug: false\n")
	contentB := []byte("server:\n  port: 9090\n  debug: true\n")

	deployOpts := osapi.FileDeployOpts{
		ObjectName:  "app-config.yaml",
		Path:        "/tmp/app-config.yaml",
		ContentType: "raw",
		Mode:        "0644",
	}

	// Phase 1: upload known content to establish baseline.
	fmt.Println("=== Phase 1: Upload baseline ===")

	o1 := orchestrator.New(url, token)

	o1.FileUpload("app-config.yaml", "raw", contentA, orchestrator.WithForce())

	if _, err := o1.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Phase 2: FileChanged with same content — no drift expected.
	fmt.Println("\n=== Phase 2: Check with same content (expect no drift) ===")

	o2 := orchestrator.New(url, token)

	check2 := o2.FileChanged("app-config.yaml", contentA)

	upload2 := o2.FileUpload("app-config.yaml", "raw", contentA).
		After(check2).
		OnlyIfChanged()

	o2.FileDeploy("_any", deployOpts).
		After(upload2).
		OnlyIfChanged()

	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var changed2 osapi.FileChanged
	if err := report2.Decode("check-file", &changed2); err == nil {
		fmt.Printf("File %s changed: %v (sha256: %s)\n",
			changed2.Name, changed2.Changed, changed2.SHA256[:12])
	}

	// Phase 3: FileChanged with different content — drift expected.
	fmt.Println("\n=== Phase 3: Check with different content (expect drift) ===")

	o3 := orchestrator.New(url, token)

	check3 := o3.FileChanged("app-config.yaml", contentB)

	upload3 := o3.FileUpload("app-config.yaml", "raw", contentB, orchestrator.WithForce()).
		After(check3).
		OnlyIfChanged()

	o3.FileDeploy("_any", deployOpts).
		After(upload3).
		OnlyIfChanged()

	report3, err := o3.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var changed3 osapi.FileChanged
	if err := report3.Decode("check-file", &changed3); err == nil {
		fmt.Printf("File %s changed: %v (sha256: %s)\n",
			changed3.Name, changed3.Changed, changed3.SHA256[:12])
	}
}
