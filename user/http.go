// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/session"
	"net/http"
)

var config = map[string]string{
	"login_url": "/-/auth/login",
}

// CurrentUserID returns the userID of the requesting user.
func CurrentUserID(r *http.Request) (int64, error) {
	s, err := session.Store.Get(r, "hal/user")
	if err != nil {
		return 0, err
	}
	userID, _ := s.Values["userid"].(int64)
	return userID, err
}

// CurrentUserSetID adds the provided userID to the current users session/cookie
func CurrentUserSetID(w http.ResponseWriter, r *http.Request, userID int64) error {
	s, err := session.Store.Get(r, "hal/user")
	if err != nil {
		c := appengine.NewContext(r)
		c.Errorf("hal/user: There was an error retrieving the session\n Error: %v", err)
		return err
	}
	s.Values["userid"] = userID
	s.Save(r, w)
	return nil
}

// CurrentUserHasRole checks for the presents of a role listed under the current user.
func CurrentUserHasRole(w http.ResponseWriter, r *http.Request, role string) (bool, error) {

	// Confirm we have a user.

	if id, err := CurrentUserID(r); id != 0 || err != nil {
		return false, err
	}

	// 1st Check the session.

	s, err := session.Store.Get(r, "hal/user|roles")
	if err != nil {
		c := appengine.NewContext(r)
		c.Errorf("hal/user: There was an error retrieving the session\n Error: %v", err)
		return false, err
	}

	if s.Values[role] == true {
		return true, nil
	}

	// 2nd Check the ds.

	u := new(User)
	_, _ = Current(r, u)

	if u.HasRole(role) == true {
		// Set the role to true in the session to avoid this look up in the future.
		err = CurrentUserSetRole(w, r, role, true)
		return true, err
	}
	return false, err
}

func CurrentUserSetRole(w http.ResponseWriter, r *http.Request, role string,
	value bool) (err error) {

	s, err := session.Store.Get(r, "hal/user")
	if err != nil {
		c := appengine.NewContext(r)
		c.Errorf("hal/user: There was an error retrieving the session\n Error: %v", err)
		return
	}
	s.Values[role] = value
	s.Save(r, w)
	return
}

// Current checks the requesting User's session to see if they have an
// account. If they do, the provided User struct is populated with the
// information that is saved in the datastore. If they don't an error is
// returned.
func Current(r *http.Request, u *User) (*datastore.Key, error) {
	userID, _ := CurrentUserID(r)
	if userID != 0 {
		c := context.NewContext(r)
		key := datastore.NewKey(c, "User", "", userID, nil)
		err := ds.Get(c, key, u)
		return key, err
	}
	return nil, ErrNoLoggedInUser
}

// Logout sets the session userid to 0, effectivly logging the user out.
func Logout(w http.ResponseWriter, r *http.Request) error {
	return CurrentUserSetID(w, r, 0)
}

// LoginRequired is a wrapper for http.HandleFunc. If the requesting
// User is not logged in, they will be redirect to the login page.
func LoginRequired(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if userid, _ := CurrentUserID(r); userid == 0 {
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
