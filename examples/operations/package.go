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

// Package main demonstrates package management:
// install a package, inspect it, list available updates, then remove it.
//
// Uses ContinueOnError since package operations require privileges
// and may not be available on all platforms.
//
// DAG:
//
//	install-package
//	    └── get-package
//	            └── list-package-updates
//	                    └── remove-package
//
// Run with: OSAPI_TOKEN="<jwt>" go run package.go
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

	const pkg = "sl"

	// Cleanup any leftover package from a previous run.
	oc := orchestrator.New(url, token)
	oc.PackageRemove("_any", pkg).ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// List → install → get → list updates → update → remove.
	o := orchestrator.New(url, token)

	list := o.PackageList("_any").ContinueOnError()

	install := o.PackageInstall("_any", pkg).After(list).ContinueOnError()

	get := o.PackageGet("_any", pkg).After(install).ContinueOnError()

	updates := o.PackageListUpdates("_any").After(get).ContinueOnError()

	update := o.PackageUpdate("_any").After(updates).ContinueOnError()

	o.PackageRemove("_any", pkg).After(update).ContinueOnError()

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var info osapi.PackageInfoResult
	if err := report.Decode("get-package", &info); err == nil && len(info.Packages) > 0 {
		p := info.Packages[0]
		fmt.Printf("Package %s: version=%s, status=%s\n", p.Name, p.Version, p.Status)
	} else {
		fmt.Println("Package operations require apt/dpkg (Debian-family only)")
	}
}
