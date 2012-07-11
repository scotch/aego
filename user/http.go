// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"net/http"
)

var config = map[string]string{
	"login_url": "/-/auth/login",
}

// LoginRequired is a wrapper for http.HandleFunc. If the requesting
// User is not logged in, they will be redirect to the login page.
func LoginRequired(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if id, _ := CurrentUserID(r); id == "" {
			http.Redirect(w, r, config["login_url"], http.StatusForbidden)
		}
		fn(w, r)
	}
}

// AdminRequired is a wrapper for http.HandleFuc. If the requesting
// User is *not* an admin, they will redirect to the login page.
func AdminRequired(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if hasRole, err := CurrentUserHasRole(w, r, "admin"); hasRole == false || err != nil {
			http.Redirect(w, r, config["login_url"], http.StatusForbidden)
		}
		fn(w, r)
	}
}
