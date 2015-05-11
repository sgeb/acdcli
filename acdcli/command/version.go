// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/sgeb/go-acd"
)

type VersionCommand struct {
	AppName           string
	Revision          string
	Version           string
	VersionPrerelease string
	Ui                cli.Ui
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "%s v%s", c.AppName, c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	fmt.Fprintln(&versionString, "")
	fmt.Fprintf(&versionString, "using go-acd v%s", acd.LibraryVersion)

	c.Ui.Output(versionString.String())
	return 0

}

func (c *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Prints the %s version", c.AppName)
}
