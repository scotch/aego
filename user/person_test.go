// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/user/email"
	"github.com/scotch/hal/person"
	"testing"
)

var p1 *person.Person = &person.Person{
	Name: &person.PersonName{
		GivenName:  "Kyle",
		FamilyName: "Finley",
	},
	Email: "1@examle.org",
	Password: &person.PersonPassword{
		New: "secret1",
	},
}

var p2 *person.Person = &person.Person{
	Name: &person.PersonName{
		GivenName:  "Kyle",
		FamilyName: "Finley",
	},
	Email: "2@examle.org",
	Password: &person.PersonPassword{
		Current: "secret1",
		New:     "secret2",
	},
}

var p3 *person.Person = &person.Person{
	Name: &person.PersonName{
		GivenName:  "Kyle",
		FamilyName: "Finley",
	},
	Email: "3@examle.org",
	Password: &person.PersonPassword{
		New: "secret1",
	},
}

func TestCreateFromPerson(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	var err error
	var u *User
	var u2 *User
	var e *email.Email

	// Round #1 New User with email & password
	// Save it.
	if u, err = CreateFromPerson(c, p1); err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Check User
	// Get from ds to confirm save
	if u, err = Get(c, u.Key.StringID()); err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if u.Email != p1.Email {
		t.Errorf(`u.Email: %v, want %v`, u.Email, p1.Email)
	}
	if u.Person.ID != u.Key.StringID() {
		t.Errorf(`u.Person.ID: %v, want %v`, u.Person.ID, u.Key.StringID())
	}
	// Check Email
	if e, err = email.Get(c, u.Email); err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if e.UserID != u.Key.StringID() {
		t.Errorf(`u.UserID: %v, want %v`, e.UserID, u.Key.StringID())
	}

	// Round #2 Existing User with email & password
	// Get it
	if u2, err = UpdateFromPerson(c, u.Person); err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	if u2.Key.StringID() != u.Key.StringID() {
		t.Errorf(`u2.Key.StringID: %v, want %v`, u2.Key.StringID(), u.Key.StringID())
	}
}
