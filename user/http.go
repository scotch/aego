// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"github.com/scotch/hal/session"
	"net/http"
)

var config = map[string]string{
	"login_url": "/-/auth/login",
}

// LoginRequired is a wrapper for http.HandleFunc. If the requesting
// User is not logged in, they will be redirect to the login page.
func LoginRequired(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userid, _ := CurrentUserID(r)
		if userid == 0 {
			http.Redirect(w, r, config["login_url"], http.StatusFound)
		}
		fn(w, r)
	}
}

// AdminRequired is a wrapper for http.HandleFuc. It the requesting
// User is not an admin, they will redirect to the login page.
func AdminRequired(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		// Get the session
		sess, err := session.Store.Get(r, "auth")
		if err != nil {
			c.Errorf("There was an error retrieving the session: %v", err)
		}
		userid, _ := sess.Values["userid"].(int64)
		isAdmin, _ := sess.Values["isadmin"].(int64)
		if userid == 0 || isAdmin == 0 {
			http.Redirect(w, r, config["login_url"], http.StatusFound)
		}
		fn(w, r)
	}
}
