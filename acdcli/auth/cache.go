// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/oauth2"
)

type Cache struct {
	tokenFile string
}

func NewCache(config *oauth2.Config) Cache {
	return Cache{
		tokenFile: tokenCacheFile(config),
	}
}

func (c Cache) Token() (*oauth2.Token, error) {
	f, err := os.Open(c.tokenFile)
	if err != nil {
		return nil, err
	}
	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func (c Cache) SaveToken(token *oauth2.Token) error {
	f, err := os.Create(c.tokenFile)
	if err != nil {
		return err
	}
	defer f.Close()

	gob.NewEncoder(f).Encode(token)
	return nil
}

func tokenCacheFile(config *oauth2.Config) string {
	hash := fnv.New32a()
	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))
	fn := fmt.Sprintf("acdcli-tok%v", hash.Sum32())
	return filepath.Join(osUserCacheDir(), url.QueryEscape(fn))
}

func osUserCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd":
		return filepath.Join(os.Getenv("HOME"), ".cache")
	}
	log.Printf("TODO: osUserCacheDir on GOOS %q. Using current directory as fallback.", runtime.GOOS)
	return "."
}
