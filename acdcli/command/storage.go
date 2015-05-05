// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"fmt"
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

	fmt.Printf("Last Calculated: %v\n\n", *accountUsage.LastCalculated)
	fmt.Printf("Other:\n Total: %v bytes (%v files)\n Billable: %v bytes (%v files)\n\n",
		*accountUsage.Other.Total.Bytes, *accountUsage.Other.Total.Count,
		*accountUsage.Other.Billable.Bytes, *accountUsage.Other.Billable.Count)
	fmt.Printf("Doc:\n Total: %v bytes (%v files)\n Billable: %v bytes (%v files)\n\n",
		*accountUsage.Doc.Total.Bytes, *accountUsage.Doc.Total.Count,
		*accountUsage.Doc.Billable.Bytes, *accountUsage.Doc.Billable.Count)
	fmt.Printf("Photo:\n Total: %v bytes (%v files)\n Billable: %v bytes (%v files)\n\n",
		*accountUsage.Photo.Total.Bytes, *accountUsage.Photo.Total.Count,
		*accountUsage.Photo.Billable.Bytes, *accountUsage.Photo.Billable.Count)
	fmt.Printf("Video:\n Total: %v bytes (%v files)\n Billable: %v bytes (%v files)\n\n",
		*accountUsage.Video.Total.Bytes, *accountUsage.Video.Total.Count,
		*accountUsage.Video.Billable.Bytes, *accountUsage.Video.Billable.Count)

	return 0
}

func (c *StorageCommand) Synopsis() string {
	return "Prints information on storage usage and quota"
}
