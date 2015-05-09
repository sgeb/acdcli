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
	fmt.Printf("%v\n", newStorageRow("Photos", accountUsage.Photo))
	fmt.Printf("%v\n", newStorageRow("Video", accountUsage.Video))
	fmt.Printf("%v\n", newStorageRow("Doc", accountUsage.Doc))
	fmt.Printf("%v\n", newStorageRow("Other", accountUsage.Other))

	return 0
}

func (c *StorageCommand) Synopsis() string {
	return "Prints information on storage usage and quota"
}

type storageRow struct {
	title         string
	size          uint64
	count         uint64
	billableSize  uint64
	billableCount uint64
}

func newStorageRow(title string, c *acd.CategoryUsage) storageRow {
	return storageRow{
		title:         title,
		size:          *c.Total.Bytes,
		count:         *c.Total.Count,
		billableSize:  *c.Billable.Bytes,
		billableCount: *c.Billable.Count,
	}
}

func (r storageRow) String() string {
	return fmt.Sprintf(" %8v %8v %9v  %8v %9v",
		r.title,
		humanize.IBytes(r.size), humanize.Comma(int64(r.count)),
		humanize.IBytes(r.billableSize), humanize.Comma(int64(r.billableCount)))
}
