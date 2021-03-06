// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package api provides rpc service for hal.
*/
package api

import (
	"code.google.com/p/gorilla/rpc"
	"code.google.com/p/gorilla/rpc/json"
	"github.com/scotch/aego/v1/auth/password"
	"github.com/scotch/aego/v1/auth/profile"
	"github.com/scotch/aego/v1/user"
	"net/http"
)

const (
	API_URL string = "/-/api/v1"
)

func init() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(user.Service), "User")
	s.RegisterService(new(password.Service), "Password")
	s.RegisterService(new(profile.Service), "AuthProfile")
	http.Handle(API_URL, s)
}
