// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
This file contains User methods connected to UserProfiles.

*/

package user

// AddAuthID Adds an AuthID to the User's AuthIDs list. It returns true
// if this is a new authID
func (u *User) AddAuthID(authID string) bool {
	for _, id := range u.AuthIDs {
		if id == authID {
			return false
		}
	}
	u.AuthIDs = append(u.AuthIDs, authID)
	return true
}
