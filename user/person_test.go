// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"fmt"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/password"
	"github.com/scotch/hal/person"
	"testing"
)

func TestCreateFromPerson(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	var err error
	var u *User

	// Save it.

	p := &person.Person{
		Name: &person.PersonName{
			GivenName:  "Kyle",
			FamilyName: "Finley",
		},
		Email: "test@examle.org",
		Password: &person.PersonPassword{
			New: "secret1",
		},
	}

	u, err = CreateFromPerson(c, p)

	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Check User

	if u.Email != p.Email {
		t.Errorf(`u.Email: %v, want %v`, u.Email, p.Email)
	}

	if u.Person.ID != fmt.Sprintf("%v", u.Key.StringID()) {
		t.Errorf(`u.Person.ID: %v, want %v`, u.Person.ID, u.Key.StringID())
	}

	// Check password

	if err = password.CompareHashAndPassword(u.Password, []byte("secret1")); err != nil {
		t.Errorf(`Password hash does not match`)
	}
	if u.Person.Password.IsSet != true {
		t.Errorf(`u.Person.Person.IsSet: %v, want %v`, u.Person.Password.IsSet, true)
	}
	if u.Person.Password.New != "" {
		t.Errorf(`u.Person.Person.New: %v, want ""`, u.Person.Password.New)
	}
	if u.Person.Password.Current != "" {
		t.Errorf(`u.Person.Person.Current: %v, want ""`, u.Person.Password.Current)
	}
}
