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

// Package main demonstrates power management with a delayed reboot.
//
// WARNING: This example will actually reboot the target host!
// It uses a delay to give time to cancel. Uses ContinueOnError
// since power operations require privileges.
//
// DAG:
//
//	reboot
//
// Run with: OSAPI_TOKEN="<jwt>" go run power.go
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

	o := orchestrator.New(url, token)

	// Schedule a reboot with a 60-second delay.
	// WARNING: this will reboot the target host.
	o.PowerReboot("_any", osapi.PowerOpts{
		Delay:   60,
		Message: "OSAPI orchestrator example reboot",
	}).ContinueOnError()

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var result osapi.PowerResult
	if err := report.Decode("reboot", &result); err == nil {
		fmt.Printf("Power %s: delay=%ds, changed=%v\n",
			result.Action, result.Delay, result.Changed)
	} else {
		fmt.Println("Power operations require privileges")
	}
}
