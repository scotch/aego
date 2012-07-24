// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"github.com/scotch/hal/acl"
)

// Can check if the user has permission to perform the action.
func (u *User) Can(c appengine.Context, perm string, key *datastore.Key) bool {
	// Users can do anything to their own user object.
	if key.Kind() == "User" && u.Key.StringID() == key.StringID() {
		return true
	}
	// Admins can do anything.
	if u.HasRole("admin") {
		return true
	}
	// Other permissions must be set.
	if ok, _ := acl.Can(c, u.Key.String(), perm, key); ok {
		return true
	}
	return false
}
