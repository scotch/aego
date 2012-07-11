// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/session"
	"net/http"
)

var (
	ErrNoLoggedInUser = errors.New("user: no logged in user")
)

// CurrentUserID returns the userID of the requesting user.
func CurrentUserID(r *http.Request) (int64, error) {
	s, err := session.Store.Get(r, "user")
	if err != nil {
		return 0, err
	}
	userID, _ := s.Values["userid"].(int64)
	return userID, err
}

// CurrentUserSetID adds the provided userID to the current users session/cookie
func CurrentUserSetID(w http.ResponseWriter, r *http.Request, userID int64) error {
	s, err := session.Store.Get(r, "user")
	if err != nil {
		c := appengine.NewContext(r)
		c.Errorf("user: There was an error retrieving the session Error: %v", err)
		return err
	}
	s.Values["userid"] = userID

	return s.Save(r, w)
}

// Current checks the requesting User's session to see if they have an
// account. If they do, the provided User struct is populated with the
// information that is saved in the datastore. If they don't an error is
// returned.
func Current(r *http.Request) (*User, error) {
	userID, _ := CurrentUserID(r)

	if userID != 0 {
		c := context.NewContext(r)
		u := new(User)
		key := datastore.NewKey(c, "User", "", userID, nil)
		err := ds.Get(c, key, u)
		u.Key = key
		return u, err
	}
	return nil, ErrNoLoggedInUser
}

// Logout sets the session userid to 0, effectivly logging the user out.
func Logout(w http.ResponseWriter, r *http.Request) error {
	return CurrentUserSetID(w, r, 0)
}
