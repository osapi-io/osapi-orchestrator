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

//go:build ignore

// Package main demonstrates conditional file upload using FileChanged.
// Checks whether local content differs from the Object Store version,
// then only uploads and deploys if changes are detected.
//
// DAG:
//
//	health-check
//	    └── check-file
//	            └── upload-file (OnlyIfChanged)
//	                    └── deploy-file (OnlyIfChanged)
//
// Run with: OSAPI_TOKEN="<jwt>" go run file-changed.go
package main

import (
	"fmt"
	"log"
	"os"

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

	configData := []byte("server:\n  port: 8080\n  debug: false\n")

	o := orchestrator.New(url, token)

	// Level 0: verify the API is reachable.
	health := o.HealthCheck("_any")

	// Level 1: check if the file content has changed.
	check := o.FileChanged("app-config.yaml", configData).After(health)

	// Level 2: upload only if the content differs from stored version.
	upload := o.FileUpload("app-config.yaml", "raw", configData).
		After(check).
		OnlyIfChanged()

	// Level 3: deploy only if a new version was uploaded.
	o.FileDeploy("_any", orchestrator.FileDeployOpts{
		ObjectName:  "app-config.yaml",
		Path:        "/tmp/app-config.yaml",
		ContentType: "raw",
		Mode:        "0644",
	}).After(upload).
		OnlyIfChanged()

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Decode the change check result.
	var changed orchestrator.FileChangedResult
	if err := report.Decode("check-file", &changed); err == nil {
		fmt.Printf("File %s changed: %v (sha256: %s)\n",
			changed.Name, changed.Changed, changed.SHA256[:12])
	}
}
