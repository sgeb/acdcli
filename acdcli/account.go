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

	"github.com/sgeb/go-acd/acd"
)

func init() {
	registerCommand("account", accountMain)
}

func accountMain(client *http.Client, argv []string) {
	if len(argv) != 0 {
		fmt.Fprintln(os.Stderr, "Usage: account")
		return
	}

	api := acd.NewClient(client)
	accountInfo, _, err := api.Account.GetInfo()
	if err != nil {
		fmt.Printf("\nerror: %v\n\n", err)
	}

	fmt.Printf("\nTerms of use: %s\nStatus: %s\n\n", *accountInfo.TermsOfUse, *accountInfo.Status)
}
