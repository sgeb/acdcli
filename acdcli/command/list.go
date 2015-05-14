// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"fmt"

	"github.com/sgeb/go-acd"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Help() string {
	return ""
}

func (c *ListCommand) Run(_ []string) int {
	apiClient, err := c.NewAcdClient()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	root, _, err := apiClient.Nodes.GetRoot()
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

	opts := &acd.NodeListOptions{Sort: `["kind DESC","name ASC"]`}
	nodes, _, err := root.GetAllChildren(opts)
	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

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
