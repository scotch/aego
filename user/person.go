// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"errors"
	dserrors "github.com/scotch/hal/ds/errors"
	"github.com/scotch/hal/email"
	"github.com/scotch/hal/password"
	"github.com/scotch/hal/types"
)

var (
	ErrEmailInUse = errors.New("user: email in use")
)

func CreateFromPerson(c appengine.Context, p *types.Person) (*User, error) {

	var err error
	var u *User

	// Ensure that the email is an actually email.
	if err = email.Validate(p.Email); err != nil {
		return nil, err
	}
	// Ensure that the password is the approprate length
	if err = password.Validate(p.Password.New); err != nil {
		return nil, err
	}

	// Transaction Action
	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		var e *email.Email
		// Get the email
		e, err = email.Get(c, p.Email)
		// An error that is not an ErrNoSuchEntity indicates an an internal error
		// and it should be returned.
		if err != nil && err != dserrors.ErrNoSuchEntity {
			return err
		}
		// Lack of an error indicates that the email existing in the ds.
		if err == nil {
			return ErrEmailInUse
		}
		hash, err := password.GenerateFromPassword([]byte(p.Password.New))
		if err != nil {
			return err
		}
		// Create a new User
		u = New()
		u.Person = p
		u.Password = hash
		u.Email = p.Email
		if err = u.Put(c); err != nil {
			return err
		}
		// Update the Email with UserID.
		e.UserID = u.Key.IntID()
		return e.Put(c)
		// XG transation
	}, &datastore.TransactionOptions{XG: true})

	return u, err
}
