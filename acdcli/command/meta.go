// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"errors"
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/sgeb/acdcli/acdcli/client"
	"github.com/sgeb/go-acd"
)

// Meta contains the meta-options and functionality that nearly every
// command inherits.
type Meta struct {
	AppName string
	Ui      cli.Ui

	AcdApiClientId string
	AcdApiSecret   string
	CallbackPort   string
}

// Creates a new client for Amazon Cloud Drive.
func (m *Meta) NewAcdClient() (*acd.Client, error) {
	var err error = nil

	c, err := client.NewClient(m.AcdApiClientId, m.AcdApiSecret)
	if c == nil || err != nil {
		err = errors.New(fmt.Sprintf("Could not get a ACD client. Please authorize the app with `%s auth`.", m.AppName))
	}

	return c, err
}
