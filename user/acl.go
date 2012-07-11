// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"github.com/scotch/hal/acl"
)

// Can check if the user has permission to perform the action.
func (u *User) Can(c appengine.Context, perm string, key *datastore.Key) bool {
	// Admins can do anything.
	if ok := u.HasRole("admin"); ok {
		return true
	}
	// Users can do anything to their own user object.
	if u.Key == key {
		return true
	}
	// Other permissions must be set.
	id := fmt.Sprintf("%s", u.Key.StringID())
	if ok, _ := acl.Can(c, id, perm, key); ok {
		return true
	}
	return false
}
