// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"github.com/scotch/hal/session"
	"net/http"
)

// AddRole appends the role to the User's Roles. Returns true if role
// was added.
// TODO does this make sense? Should it instead return true if the role
// was present. Or should this method be removed?
func (u *User) AddRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return false
		}
	}
	u.Roles = append(u.Roles, role)
	return true
}

// HasRole returns true if the user has the role.
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// CurrentUserHasRole checks for the presents of a role listed under the current user.
func CurrentUserHasRole(w http.ResponseWriter, r *http.Request, role string) bool {

	// Confirm we have a user.
	if id, err := CurrentUserID(r); id != "" || err != nil {
		return false
	}
	// 1st Check the session.
	s, err := session.Store.Get(r, "user|roles")
	if err != nil {
		c := appengine.NewContext(r)
		c.Criticalf("user: There was an error retrieving the session Error: %v", err)
		return false
	}
	if s.Values[role] == true {
		return true
	}
	// 2nd Check the ds.
	u, err := Current(r)
	if err != nil {
		return false
	}
	if u.HasRole(role) {
		// Set the role to true in the session to avoid this look up in the future.
		if err = CurrentUserSetRole(w, r, role, true); err != nil {
			return false
		}
		return true
	}
	return false
}

func CurrentUserSetRole(w http.ResponseWriter, r *http.Request, role string,
	value bool) (err error) {

	s, err := session.Store.Get(r, "user")
	if err != nil {
		c := appengine.NewContext(r)
		c.Criticalf("user: There was an error retrieving the session Error: %v", err)
		return
	}
	// If the user is already an admin then there's no need to
	// re-add the that role.
	// if !user.CurrentUserHasRole(w, r, "admin") {
	//    u.AddRole("admin")
	// }
	s.Values[role] = value
	return s.Save(r, w)
}
