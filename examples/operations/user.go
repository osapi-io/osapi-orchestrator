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

// Package main demonstrates user account management:
// list → create → get → update → keys → change password → delete.
//
// Run with: OSAPI_TOKEN="<jwt>" go run user.go
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

	const username = "osapi-example-user"

	// Example SSH public key (not a real key).
	const sshKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample osapi-example"

	// Cleanup any leftover user from a previous run.
	oc := orchestrator.New(url, token)
	oc.UserDelete("_any", username).ContinueOnError()
	oc.Run(context.Background()) //nolint:errcheck

	// List → create → get → update → add key → list keys → remove key →
	// change password → delete.
	o := orchestrator.New(url, token)

	userList := o.UserList("_any")

	create := o.UserCreate("_any", osapi.UserCreateOpts{
		Name:  username,
		Shell: "/bin/bash",
	}).After(userList)

	get := o.UserGet("_any", username).After(create)

	update := o.UserUpdate("_any", username, osapi.UserUpdateOpts{
		Shell: "/bin/sh",
	}).After(get).ContinueOnError()

	addKey := o.UserAddKey("_any", username, osapi.SSHKeyAddOpts{
		Key: sshKey,
	}).After(update)

	listKeys := o.UserListKeys("_any", username).After(addKey)

	// Use a placeholder fingerprint — in practice you would decode
	// the key's fingerprint from list-ssh-key results.
	removeKey := o.UserRemoveKey("_any", username, "SHA256:example").
		After(listKeys).ContinueOnError()

	chpass := o.UserChangePassword("_any", username, "temp-pass-123").
		After(removeKey).ContinueOnError()

	o.UserDelete("_any", username).After(chpass)

	report, err := o.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var keys osapi.SSHKeyInfoResult
	if err := report.Decode("list-ssh-key", &keys); err == nil {
		fmt.Printf("User %q has %d SSH key(s)\n", username, len(keys.Keys))
		for _, k := range keys.Keys {
			fmt.Printf("  %s %s\n", k.Type, k.Fingerprint)
		}
	}
}
