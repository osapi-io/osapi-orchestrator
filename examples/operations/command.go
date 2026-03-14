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

// Package main demonstrates command execution and error handling.
//
// Phase 1: exec and shell commands run in parallel — both always
// report changed=true (commands are non-idempotent by design).
//
// Phase 2: a command that fails shows how to inspect
// CommandResult.Stderr and ExitCode.
//
// DAG (phase 1):
//
//	health-check
//	    ├── run-uptime (exec)
//	    └── shell-uname -a (shell)
//
// DAG (phase 2):
//
//	run-ls /nonexistent (exec, error expected)
//
// Run with: OSAPI_TOKEN="<jwt>" go run command.go
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

	fmt.Println("=== Phase 1: parallel exec and shell commands ===")

	o := orchestrator.New(url, token)

	health := o.HealthCheck()

	// Direct execution: runs the binary with args.
	o.CommandExec("_any", "uptime").After(health)

	// Shell execution: interpreted by /bin/sh, supports pipes.
	o.CommandShell("_any", "uname -a").After(health)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var cmd osapi.CommandResult
	if err := report.Decode("run-uptime", &cmd); err == nil {
		fmt.Printf("uptime stdout: %s\n", cmd.Stdout)
	}

	if err := report.Decode("shell-uname -a", &cmd); err == nil {
		fmt.Printf("uname stdout:  %s\n", cmd.Stdout)
	}

	fmt.Println("=== Phase 2: failing command — inspect error and exit code ===")

	o2 := orchestrator.New(url, token)

	o2.CommandExec("_any", "ls", "/nonexistent").
		OnError(orchestrator.Continue)

	report2, err := o2.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var failCmd osapi.CommandResult
	if err := report2.Decode("run-ls", &failCmd); err == nil {
		fmt.Printf("stderr:    %s", failCmd.Stderr)
		fmt.Printf("exit code: %d\n", failCmd.ExitCode)
	}
}
