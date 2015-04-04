// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

// This code was heavily inspired by
// https://github.com/google/google-api-go-client/blob/master/examples/main.go
// which originally contained the following copyright notice:
//
// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// The original LICENSE file is appended in full to
// the LICENSE file in this repository.

// Main program entry point

package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Flags
var (
	clientID     = flag.String("clientid", "", "OAuth 2.0 Client ID. If non-empty, overrides --clientid_file")
	clientIDFile = flag.String("clientid-file", "clientid.dat",
		"Name of a file containing just the project's OAuth 2.0 Client ID from https://developer.amazon.com/lwa/sp/overview.html")
	secret     = flag.String("secret", "", "OAuth 2.0 Client Secret. If non-empty, overrides --secret_file")
	secretFile = flag.String("secret-file", "clientsecret.dat",
		"Name of a file containing just the project's OAuth 2.0 Client Secret from https://developer.amazon.com/lwa/sp/overview.html")
	cacheToken = flag.Bool("cachetoken", true, "cache the OAuth 2.0 token")
	debug      = flag.Bool("debug", false, "show HTTP traffic")
	port       = flag.String("port", "", "if non-empty, sets the redirect port")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <command> [command args]\n\nPossible commands:\n\n", os.Args[0])
	for n := range commandFunc {
		fmt.Fprintf(os.Stderr, "  * %s\n", n)
	}
	os.Exit(2)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	name := flag.Arg(0)
	command, ok := commandFunc[name]
	if !ok {
		usage()
	}

	config := &oauth2.Config{
		ClientID:     valueOrFileContents(*clientID, *clientIDFile),
		ClientSecret: valueOrFileContents(*secret, *secretFile),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.amazon.com/ap/oa",
			TokenURL: "https://api.amazon.com/auth/o2/token",
		},
		Scopes: []string{"clouddrive:read"},
	}

	ctx := context.Background()
	if *debug {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &logTransport{http.DefaultTransport},
		})
	}

	c := newOAuthClient(ctx, config)
	command(c, flag.Args()[1:])
}

var (
	commandFunc = make(map[string]func(*http.Client, []string))
)

func registerCommand(name string, main func(c *http.Client, argv []string)) {
	if commandFunc[name] != nil {
		panic(name + " already registered")
	}
	commandFunc[name] = main
}

func osUserCacheDir() string {
	//	switch runtime.GOOS {
	//	case "darwin":
	//		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	//	case "linux", "freebsd":
	//		return filepath.Join(os.Getenv("HOME"), ".cache")
	//	}
	//	return filepath.Join(os.Getenv("HOME"), "Desktop")
	//	log.Printf("TODO: osUserCacheDir on GOOS %q", runtime.GOOS)
	return "."
}

func tokenCacheFile(config *oauth2.Config) string {
	hash := fnv.New32a()
	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))
	fn := fmt.Sprintf("acdcli-tok%v", hash.Sum32())
	return filepath.Join(osUserCacheDir(), url.QueryEscape(fn))
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	if !*cacheToken {
		return nil, errors.New("--cachetoken is false")
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Warning: failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}

func newOAuthClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile := tokenCacheFile(config)
	token, err := tokenFromFile(cacheFile)
	if err != nil {
		token = tokenFromWeb(ctx, config)
		saveToken(cacheFile, token)
	} else {
		log.Printf("Using cached token %#v from %q", token, cacheFile)
	}

	return config.Client(ctx, token)
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	if *debug {
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

	if *port != "" {
		if ts.Listener != nil {
			ts.Listener.Close()
		}
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", *port))
		if err != nil {
			panic(fmt.Sprintf("Failed to listen on a port: %v", err))
		}
		ts.Listener = l
	}

	ts.Start()
	defer ts.Close()

	config.RedirectURL = ts.URL
	if *debug {
		log.Printf("Redirect URL: %s", config.RedirectURL)
	}
	authURL := config.AuthCodeURL(randState)
	log.Printf("Trying to authorize this app. If your browser does not open, please navigate directly to: %s", authURL)
	go openURL(authURL)
	code := <-ch
	if *debug {
		log.Printf("Got authorization code: %s", code)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}

	// Amazon Cloud Drive might return lowercase "bearer", but
	// request only accepted with capitalized "Bearer"
	if token.TokenType == "bearer" {
		token.TokenType = "Bearer"
	}

	return token
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

func valueOrFileContents(value string, filename string) string {
	if value != "" {
		return value
	}
	slurp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading %q: %v", filename, err)
	}
	return strings.TrimSpace(string(slurp))
}
