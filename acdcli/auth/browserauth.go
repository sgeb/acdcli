// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/sgeb/acdcli/acdcli/debug"
)

func TokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	if debug.Debug {
		log.Printf("Starting redirect server")
	}
	ts := NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	}))

	if config.RedirectURL != "" && strings.HasPrefix(config.RedirectURL, "http://") {
		if ts.Listener != nil {
			err := ts.Listener.Close()
			if err != nil {
				return nil, err
			}
		}
		laddress := config.RedirectURL[7:]
		l, err := net.Listen("tcp", laddress)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to listen on address %s: %v", laddress, err)
			return nil, errors.New(errMsg)
		}
		ts.Listener = l
	}

	ts.Start()
	defer ts.Close()

	if debug.Debug {
		log.Printf("Redirect URL: %s", config.RedirectURL)
	}
	authURL := config.AuthCodeURL(randState)
	fmt.Printf("Trying to authorize this app. If your browser does not open, please navigate directly to:\n\n%s\n\n", authURL)
	go openURL(authURL)
	code := <-ch
	if debug.Debug {
		log.Printf("Got authorization code: %s", code)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Token exchange error: %v", err))
	}

	// Amazon Cloud Drive might return lowercase "bearer", but
	// request only accepted with capitalized "Bearer"
	if token.TokenType == "bearer" {
		token.TokenType = "Bearer"
	}

	return token, nil
}

func openURL(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
	log.Printf("Error opening URL in browser.")
}
