// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"fmt"
	"strings"
)

type InfoCommand struct {
	Meta
}

func (c *InfoCommand) Help() string {
	return ""
}

func (c *InfoCommand) Run(args []string) int {
	client, err := c.NewAcdClient()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	// Start with root folder
	root, _, err := client.Nodes.GetRoot()
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

	// prepare the names to walk in node hierarchy
	names := make([]string, 0)
	targetArg := ""
	if len(args) > 0 {
		targetArg = args[0]
	}
	for _, s := range strings.Split(targetArg, "/") {
		if s != "" {
			names = append(names, s)
		}
	}

	// walk the hiearchy
	targetNode, _, err := root.WalkNodes(names...)
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

	if !targetNode.IsFolder() && strings.HasSuffix(targetArg, "/") {
		c.Ui.Error(fmt.Sprintf("Node '%s' is not a folder", *targetNode.Name))
		return 4
	}

	// retrieve and print the metadata
	md, err := targetNode.GetMetadata()
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}
	c.Ui.Output(md)

	return 0
}

func (c *InfoCommand) Synopsis() string {
	return "Display a node's metadata"
}
