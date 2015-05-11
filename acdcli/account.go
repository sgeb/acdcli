// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

// Command for the Account API
// See: https://developer.amazon.com/public/apis/experience/cloud-drive/content/account

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sgeb/go-acd"
)

var (
	api *acd.Client
)

func accountMain(client *http.Client, argv []string) {
	if len(argv) != 0 {
		fmt.Fprintln(os.Stderr, "Usage: account")
		return
	}

	api = acd.NewClient(client)
	printAccountInfo()
	printAccountQuota()
	printAccountUsage()
}

func printAccountInfo() {
	accountInfo, _, err := api.Account.GetInfo()
	if err != nil {
		fmt.Printf("\nerror: %v\n\n", err)
	}

	fmt.Printf("\nTerms of use: %v\nStatus: %v\n\n", *accountInfo.TermsOfUse, *accountInfo.Status)
}

func printAccountQuota() {
	accountQuota, _, err := api.Account.GetQuota()
	if err != nil {
		fmt.Printf("\nerror: %v\n\n", err)
	}

	fmt.Printf("\nQuota: %v\nLast Calculated: %v\nAvailable: %v\n\n",
		*accountQuota.Quota, *accountQuota.LastCalculated, *accountQuota.Available)
}

func printAccountUsage() {
	accountUsage, _, err := api.Account.GetUsage()
	if err != nil {
		fmt.Printf("\nerror: %v\n\n", err)
	}

	fmt.Printf("\nLast Calculated: %v\n", *accountUsage.LastCalculated)
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
}
