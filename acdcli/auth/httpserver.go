// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

// This code was heavily inspired by
// https://github.com/golang/go/blob/master/src/net/http/httptest/server.go
// which originally contained the following copyright notice:
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// The original LICENSE file is appended in full to
// the LICENSE file in this repository.

// Implementation of a gracefully closing HTTP server

package auth

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

// A Server is an HTTP server listening on a system-chosen port on the
// local loopback interface.
type Server struct {
	URL      string // base URL of form http://ipaddr:port with no trailing slash
	Listener net.Listener

	// Config may be changed after calling NewUnstartedServer and
	// before Start
	Config *http.Server

	// wg counts the number of outstanding HTTP requests on this server.
	// Close blocks until all requests are finished.
	wg sync.WaitGroup
}

// historyListener keeps track of all connections that it's ever
// accepted.
type historyListener struct {
	net.Listener
	sync.Mutex // protects history
	history    []net.Conn
}

func (hs *historyListener) Accept() (c net.Conn, err error) {
	c, err = hs.Listener.Accept()
	if err == nil {
		hs.Lock()
		hs.history = append(hs.history, c)
		hs.Unlock()
	}
	return
}

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	return l
}

// NewUnstartedServer returns a new Server but doesn't start it.
//
// After changing its configuration, the caller should call Start.
//
// The caller should call Close when finished, to shut it down.
func NewUnstartedServer(handler http.Handler) *Server {
	return &Server{
		Listener: newLocalListener(),
		Config:   &http.Server{Handler: handler},
	}
}

// Start starts a server from NewUnstartedServer.
func (s *Server) Start() {
	if s.URL != "" {
		panic("Server already started")
	}
	s.Listener = &historyListener{Listener: s.Listener}
	s.URL = "http://" + s.Listener.Addr().String()
	s.wrapHandler()
	go s.Config.Serve(s.Listener)
}

func (s *Server) wrapHandler() {
	h := s.Config.Handler
	if h == nil {
		h = http.DefaultServeMux
	}
	s.Config.Handler = &waitGroupHandler{
		s: s,
		h: h,
	}
}

// Close shuts down the server and blocks until all outstanding
// requests on this server have completed.
func (s *Server) Close() {
	s.Listener.Close()
	s.wg.Wait()
	s.CloseClientConnections()
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		t.CloseIdleConnections()
	}
}

// CloseClientConnections closes any currently open HTTP connections
// to the test Server.
func (s *Server) CloseClientConnections() {
	hl, ok := s.Listener.(*historyListener)
	if !ok {
		return
	}
	hl.Lock()
	for _, conn := range hl.history {
		conn.Close()
	}
	hl.Unlock()
}

// waitGroupHandler wraps a handler, incrementing and decrementing a
// sync.WaitGroup on each request, to enable Server.Close to block
// until outstanding requests are finished.
type waitGroupHandler struct {
	s *Server
	h http.Handler // non-nil
}

func (h *waitGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.s.wg.Add(1)
	defer h.s.wg.Done() // a defer, in case ServeHTTP below panics
	h.h.ServeHTTP(w, r)
}
