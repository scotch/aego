// Copyright 2012 The AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"errors"
	"github.com/scotch/aego/v1/user/email"
)

var (
	ErrKeyNotSet = errors.New("user: key not set")
)

// AddEmail appends the email to the User's emails. Returns true if role
// was added.
func (u *User) AddEmail(c appengine.Context, address string, status int64) (e *email.Email, err error) {
	if u.Key.StringID() == "" {
		return nil, ErrKeyNotSet
	}
	if u.HasEmail(address) {
		return nil, email.ErrAddressAlreadyAdded
	}
	e, err = email.AddForUser(c, address, u.Key.StringID(), status)
	if err != nil {
		return e, err
	}
	u.Emails = append(u.Emails, address)
	return e, nil
}

// HasEmail returns true if the user has the email.
func (u *User) HasEmail(address string) bool {
	for _, a := range u.Emails {
		if a == address {
			return true
		}
	}
	return false
}
