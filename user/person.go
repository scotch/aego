// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"errors"
	dserrors "github.com/scotch/hal/ds/errors"
	"github.com/scotch/hal/user/email"
	"github.com/scotch/hal/person"
)

var (
	ErrEmailInUse = errors.New("user: email in use")
)

func validatePerson(p *person.Person) (err error) {
	// Ensure that the email is an actually email.
	if err = email.Validate(p.Email); err != nil {
		return
	}
	return
}

// CreateFromPerson Creates a User from a Person object
func CreateFromPerson(c appengine.Context, p *person.Person) (u *User, err error) {

	if err = validatePerson(p); err != nil {
		return
	}

	// Transaction Action
	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		// Get the email
		e, err := email.Get(c, p.Email)
		// An error that is not an ErrNoSuchEntity indicates an an internal error
		// and it should be returned.
		if err != nil && err != dserrors.ErrNoSuchEntity {
			return err
		}
		// Lack of an error indicates that the email existing in the ds.
		if err == nil {
			return ErrEmailInUse
		}
		// Create a new User
		u = New()
		u.Person = p
		u.Email = p.Email
		if err = u.Put(c); err != nil {
			return err
		}
		// Update the Email with UserID.
		e.UserID = u.Key.StringID()
		return e.Put(c)
		// XG transation
	}, &datastore.TransactionOptions{XG: true})

	return u, err
}

func UpdateFromPerson(c appengine.Context, p *person.Person) (u *User, err error) {

	if err = validatePerson(p); err != nil {
		return
	}

	// Transaction Action
	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		// Get the user
		u, err = Get(c, p.ID)
		if err != nil {
			return err
		}
		u.Person = p
		// TODO more care needs to be taken when changing emails.
		u.Email = p.Email
		return u.Put(c)
		// XG transation
	}, &datastore.TransactionOptions{XG: false})

	return u, err
}
