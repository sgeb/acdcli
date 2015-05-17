// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestDownloadCommand_implements(t *testing.T) {
	var _ cli.Command = &DownloadCommand{}
}
