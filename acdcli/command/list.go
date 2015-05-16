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

type ListCommand struct {
	Meta
}

func (c *ListCommand) Help() string {
	return ""
}

func (c *ListCommand) Run(args []string) int {
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

	// prepare list of nodes to be printed
	nodes := []*acd.Node{targetNode}
	switch t := targetNode.Typed().(type) {
	case *acd.Folder:
		opts := &acd.NodeListOptions{Sort: `["kind DESC","name ASC"]`}
		ns, _, err := t.GetAllChildren(opts)
		if err != nil {
			c.Ui.Error(err.Error())
			return 3
		}
		nodes = ns
	}

	// print list of nodes
	for _, node := range nodes {
		if node.Name != nil {
			name := *node.Name
			if _, ok := node.Typed().(*acd.Folder); ok {
				name += "/"
			}
			c.Ui.Output(name)
		} else if node.Id != nil {
			c.Ui.Output(fmt.Sprintf("-> %v", *node.Id))
		} else {
			c.Ui.Output(fmt.Sprintf("?? %v", node))
		}
	}

	return 0
}

func (c *ListCommand) Synopsis() string {
	return "List files and folder in the root folder of the drive"
}
