// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/sgeb/go-acd/acd"
)

type StorageCommand struct {
	Meta
}

func (c *StorageCommand) Help() string {
	return ""
}

func (c *StorageCommand) Run(_ []string) int {
	apiClient, err := c.NewAcdClient()
	if err != nil {
		c.Ui.Output(err.Error())
		return 1
	}

	accountUsage, _, err := apiClient.Account.GetUsage()
	if err != nil {
		fmt.Printf("\nerror: %v\n\n", err)
		return 3
	}

	fmt.Printf("Last Calculated: %v (%v)\n",
		humanize.Time(*accountUsage.LastCalculated), *accountUsage.LastCalculated)
	fmt.Printf("%v\n", categoryUsageString("Photos", accountUsage.Photo))
	fmt.Printf("%v\n", categoryUsageString("Video", accountUsage.Video))
	fmt.Printf("%v\n", categoryUsageString("Doc", accountUsage.Doc))
	fmt.Printf("%v\n", categoryUsageString("Other", accountUsage.Other))

	return 0
}

func (c *StorageCommand) Synopsis() string {
	return "Prints information on storage usage and quota"
}

func categoryUsageString(category string, c *acd.CategoryUsage) string {
	return fmt.Sprintf(" %8v %8v %9v  %8v %9v",
		category,
		humanize.IBytes(*c.Total.Bytes), humanize.Comma(int64(*c.Total.Count)),
		humanize.IBytes(*c.Billable.Bytes), humanize.Comma(int64(*c.Billable.Count)))
}
