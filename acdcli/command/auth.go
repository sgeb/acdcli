// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"fmt"
	"github.com/sgeb/acdcli/acdcli/auth"
	"github.com/sgeb/acdcli/acdcli/client"
)

type AuthCommand struct {
	Meta
}

func (c *AuthCommand) Help() string {
	return ""
}

func (c *AuthCommand) Run(_ []string) int {
	ctx := client.Context()
	config := client.Config(c.AcdApiClientId, c.AcdApiSecret)

	token, err := auth.TokenFromWeb(ctx, config)
	if err != nil {
		c.Ui.Output(fmt.Sprintf("Cannot authorize: %s", err.Error()))
		return 1
	}

	cache := auth.NewCache(config)
	err = cache.SaveToken(token)
	if err != nil {
		c.Ui.Output(fmt.Sprintf("Cannot authorize: %s", err.Error()))
		return 2
	}

	c.Ui.Output("Successfully authorized")
	return 0
}

func (c *AuthCommand) Synopsis() string {
	return fmt.Sprintf("Authorizes access to your Amazon Cloud Drive account")
}
