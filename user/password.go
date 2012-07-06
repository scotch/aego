// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"github.com/scotch/hal/email"
	pass "github.com/scotch/hal/password"
	"net/http"
)

func (u *User) setPassword(password string) (err error) {
	hash, err := pass.GenerateFromPassword([]byte(password))
	if err != nil {
		return
	}
	u.Password = hash
	return
}

func LoginByEmailAndPassword(w http.ResponseWriter, r *http.Request, emailAddress, password string) (u *User, err error) {

	c := appengine.NewContext(r)
	// Get UserId
	e, err := email.Get(c, emailAddress)
	if err != nil {
		return
	}
	u, err = Get(c, e.UserId)
	if err != nil {
		return
	}
	// Compare pasword
	if err = pass.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
		return
	}
	// We made it. Log in the User
	err = CurrentUserSetID(w, r, u.Key.IntID())
	return
}

func ChangePassword(c appengine.Context, emailAddress, currentPassword,
	newPassword string) (err error) {

	// Confirm the address is a valid email.
	err = email.Validate(emailAddress)
	if err != nil {
		return
	}
	// Get the UserId.
	// TODO(kylefinley) add status check confirm that the email has been
	// confirmed.
	e, err := email.Get(c, emailAddress)
	if err != nil {
		return
	}
	u, err := Get(c, e.UserId)
	if err != nil {
		return
	}
	// Compare pasword
	if err = pass.CompareHashAndPassword(u.Password,
		[]byte(currentPassword)); err != nil {
		return
	}
	// Set password hash to new value
	err = u.setPassword(newPassword)
	if err != nil {
		return
	}
	err = u.Put(c)
	if err != nil {
		return
	}
	return
}
