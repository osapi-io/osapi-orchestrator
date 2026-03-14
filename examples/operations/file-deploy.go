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

// Package main demonstrates file deployment: upload a file to the
// Object Store, deploy it to a target agent, then verify its status.
//
// Set OSAPI_FORCE=1 to bypass the SHA-256 pre-check and always upload
// even if content hasn't changed.
//
// DAG:
//
//	health-check
//	    └── upload-file
//	            └── deploy-file
//	                    └── file-status
//
// Run with: OSAPI_TOKEN="<jwt>" go run file-deploy.go
package main

import (
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

	o := orchestrator.New(url, token)

	// Level 0: verify the API is reachable.
	health := o.HealthCheck()

	// Level 1: upload the file to the Object Store.
	// Use WithForce() to bypass SHA-256 pre-check if requested.
	var uploadOpts []orchestrator.UploadOption
	if os.Getenv("OSAPI_FORCE") == "1" {
		uploadOpts = append(uploadOpts, orchestrator.WithForce())
	}

	upload := o.FileUpload(
		"app-config.yaml",
		"raw",
		[]byte("server:\n  port: 8080\n  debug: false\n"),
		uploadOpts...,
	).After(health)

	// Level 2: deploy the uploaded file to the target agent.
	deploy := o.FileDeploy("_any", osapi.FileDeployOpts{
		ObjectName:  "app-config.yaml",
		Path:        "/tmp/app-config.yaml",
		ContentType: "raw",
		Mode:        "0644",
		Owner:       "root",
		Group:       "root",
	}).After(upload)

	// Level 3: verify the deployed file is in sync.
	o.FileStatusGet("_any", "/tmp/app-config.yaml").After(deploy)

	report, err := o.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Decode typed results from the report.
	var status osapi.FileStatusResult
	if err := report.Decode("file-status", &status); err == nil {
		fmt.Printf("File %s status: %s\n", status.Path, status.Status)
	}
}
