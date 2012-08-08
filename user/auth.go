// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"errors"
)

var (
	ErrAuthIDAlreadyAdded = errors.New("user: AuthID already added")
)

// AddAuthID appends the role to the User's Roles. Returns an error if the authID
// was already added.
func (u *User) AddAuthID(authID string) error {
	if u.HasAuthID(authID) {
		return ErrAuthIDAlreadyAdded
	}
	u.AuthIDs = append(u.AuthIDs, authID)
	return nil
}

// HasAuthID returns true if the user has the authID.
func (u *User) HasAuthID(authID string) bool {
	for _, a := range u.AuthIDs {
		if a == authID {
			return true
		}
	}
	return false
}
