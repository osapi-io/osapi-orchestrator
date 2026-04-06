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

// Package main demonstrates CA certificate management:
// upload a PEM file, create the certificate, update it, then delete.
//
// Cleanup at the start ensures repeatability.
//
// DAG:
//
//	upload-file
//	    └── create-certificate
//	            └── list-certificate
//	                    └── delete-certificate
//
// Run with: OSAPI_TOKEN="<jwt>" go run certificate.go
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

	const certName = "osapi-example-ca"

	// Placeholder PEM content (not a real certificate).
	pem := []byte("-----BEGIN CERTIFICATE-----\nMIIBExample\n-----END CERTIFICATE-----\n")

	// Cleanup any leftover certificate from a previous run.
	oc := orchestrator.New(url, token)
	oc.CertificateDelete("_any", certName).ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// Upload PEM → create cert → list certs → delete cert.
	o := orchestrator.New(url, token)

	upload := o.FileUpload("osapi-example-ca.pem", "raw", pem,
		orchestrator.WithForce())

	create := o.CertificateCreate("_any", osapi.CertificateCreateOpts{
		Name:   certName,
		Object: "osapi-example-ca.pem",
	}).After(upload)

	list := o.CertificateList("_any").After(create)

	o.CertificateDelete("_any", certName).After(list)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var certs osapi.CertificateCAResult
	if err := report.Decode("list-certificate", &certs); err == nil {
		fmt.Printf("CA certificates (%d):\n", len(certs.Certificates))
		for _, c := range certs.Certificates {
			fmt.Printf("  %s (source: %s)\n", c.Name, c.Source)
		}
	}
}
