// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

// This code was heavily inspired by
// https://github.com/google/google-api-go-client/blob/master/examples/main.go
// which originally contained the following copyright notice:
//
// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// The original LICENSE file is appended in full to
// the LICENSE file in this repository.

// Main program entry point

package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	gitCommit string

	// The API keys for Amazon Cloud Drive. These will be filled in by the compiler.
	acdApiClientId string
	acdApiSecret   string
)

const (
	AppName = "acdcli"

	// The main version number that is being run at the moment.
	Version = "0.1"

	// A pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "dev"
)

func main() {
	os.Exit(RunCustom(os.Args[1:], Commands()))
}

func RunCustom(args []string, commands map[string]cli.CommandFactory) int {
	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: commands,
		Name:     AppName,
		Version:  Version,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return -1
	}

	return exitCode
}
