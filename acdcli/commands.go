// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/mitchellh/cli"

	"github.com/sgeb/acdcli/acdcli/command"
)

// Commands returns the mapping of CLI commands. The meta
// parameter lets you set meta options for all commands.
func Commands() map[string]cli.CommandFactory {
	meta := &command.Meta{
		AppName: AppName,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},

		AcdApiClientId: acdApiClientId,
		AcdApiSecret:   acdApiSecret,
		CallbackPort:   "56789",
	}

	return map[string]cli.CommandFactory{
		"auth": func() (cli.Command, error) {
			return &command.AuthCommand{
				Meta: *meta,
			}, nil
		},
		"storage": func() (cli.Command, error) {
			return &command.StorageCommand{
				Meta: *meta,
			}, nil
		},
		"ls": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *meta,
			}, nil
		},
		"info": func() (cli.Command, error) {
			return &command.InfoCommand{
				Meta: *meta,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				AppName:           AppName,
				Revision:          gitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
				Ui:                meta.Ui,
			}, nil
		},
	}
}
