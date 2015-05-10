// Copyright (c) 2015 Serge Gebhardt. All rights reserved.
//
// Use of this source code is governed by the ISC
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"strings"
)

// BearerAuthTransport wraps a RoundTripper. It capitalized bearer token
// authorization headers.
type BearerAuthTransport struct {
	rt http.RoundTripper
}

// RoundTrip satisfies the RoundTripper interface. It replaces authorization
// headers of scheme `bearer` by capitalized `Bearer` (as per OAuth 2.0 spec).
func (t *BearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	auth := req.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		auth = "Bearer " + auth[7:]
	}

	req2 := cloneRequest(req) // per RoundTripper contract
	req2.Header.Set("Authorization", auth)

	return t.rt.RoundTrip(req2)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
