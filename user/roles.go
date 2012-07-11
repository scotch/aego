// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"github.com/scotch/hal/session"
	"net/http"
)

// CurrentUserHasRole checks for the presents of a role listed under the current user.
func CurrentUserHasRole(w http.ResponseWriter, r *http.Request, role string) (bool, error) {

	// Confirm we have a user.
	if id, err := CurrentUserID(r); id != 0 || err != nil {
		return false, err
	}
	// 1st Check the session.
	s, err := session.Store.Get(r, "user|roles")
	if err != nil {
		c := appengine.NewContext(r)
		c.Errorf("user: There was an error retrieving the session Error: %v", err)
		return false, err
	}
	if s.Values[role] == true {
		return true, nil
	}
	// 2nd Check the ds.
	u, err := Current(r)
	if err != nil {
		return false, err
	}
	if u.HasRole(role) == true {
		// Set the role to true in the session to avoid this look up in the future.
		err = CurrentUserSetRole(w, r, role, true)
		return true, err
	}
	return false, err
}

func CurrentUserSetRole(w http.ResponseWriter, r *http.Request, role string,
	value bool) (err error) {

	s, err := session.Store.Get(r, "user")
	if err != nil {
		c := appengine.NewContext(r)
		c.Errorf("hal/user: There was an error retrieving the session Error: %v", err)
		return
	}
	s.Values[role] = value
	return s.Save(r, w)
}
