// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package client

import (
	"log"
	"net/http"

	"github.com/sgeb/acdcli/acdcli/auth"
	"github.com/sgeb/go-acd/acd"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const debug = true

func NewClient(acdApiClientId, acdApiSecret string) (*acd.Client, error) {
	ctx := Context()
	config := Config(acdApiClientId, acdApiSecret)

	httpClient, err := newOAuthClient(ctx, config)
	if err != nil {
		return nil, err
	}

	return acd.NewClient(httpClient), nil
}

func Context() context.Context {
	ctx := context.Background()
	if debug {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &logTransport{http.DefaultTransport},
		})
	}
	return ctx
}

func Config(acdApiClientId, acdApiSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     acdApiClientId,
		ClientSecret: acdApiSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.amazon.com/ap/oa",
			TokenURL: "https://api.amazon.com/auth/o2/token",
		},
		Scopes:      []string{"clouddrive:read"},
		RedirectURL: "http://127.0.0.1:56789",
	}
}

func newOAuthClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	token, err := auth.NewCache(config).Token()
	if err != nil {
		return nil, err
	} else {
		if debug {
			log.Printf("Using cached token %#v", token)
		}
	}

	return config.Client(ctx, token), nil
}
