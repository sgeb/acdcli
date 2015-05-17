// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"fmt"
	"strings"

	"github.com/sgeb/go-acd"
)

type DownloadCommand struct {
	Meta
}

func (c *DownloadCommand) Help() string {
	return ""
}

func (c *DownloadCommand) Run(args []string) int {
	client, err := c.NewAcdClient()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if len(args) == 0 {
		c.Ui.Error("Source path must be provided")
		return 4
	}
	if strings.HasSuffix(args[0], "/") {
		c.Ui.Error(fmt.Sprintf("Source path must not end with '/'"))
		return 4
	}

	// prepare the names to walk in node hierarchy
	names := make([]string, 0)
	for _, s := range strings.Split(args[0], "/") {
		if s != "" {
			names = append(names, s)
		}
	}

	// Start with root folder
	root, _, err := client.Nodes.GetRoot()
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

	// walk the hiearchy
	targetNode, _, err := root.WalkNodes(names...)
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}
	if !targetNode.IsFile() {
		c.Ui.Error(fmt.Sprintf("Node '%s' is not a file", *targetNode.Name))
		return 4
	}

	// determine output file path
	outputPath := *targetNode.Name
	if len(args) > 1 {
		outputPath = args[1]
	}

	// download
	_, err = targetNode.Typed().(*acd.File).Download(outputPath)
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

	return 0
}

func (c *DownloadCommand) Synopsis() string {
	return "Download files"
}
