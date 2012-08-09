// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hal/api provides rpc service for Users.
*/

package user

import (
	"appengine"
	"appengine/datastore"
	"errors"
	//"github.com/scotch/hal/api"
	"github.com/scotch/hal/person"
	"net/http"
)

var (
	ErrUserUnauthorized = errors.New("user: unauthorized")
)

type Empty struct{}

type Person struct {
	Person *person.Person
}

type Service struct{}

func (s *Service) Current(w http.ResponseWriter, r *http.Request,
	args *Empty, reply *Person) (err error) {

	var u *User
	if u, err = Current(r); err != nil {
		return err
	}
	reply.Person = u.Person
	return nil
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request,
	args *Person, reply *Person) (err error) {

	if err = Logout(w, r); err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request,
	args *Person, reply *Person) (err error) {

	c := appengine.NewContext(r)
	u, _ := Current(r)
	k := datastore.NewKey(c, "User", args.Person.ID, 0, nil)
	if can := u.Can(c, "write", k); can == false {
		return ErrUserUnauthorized
	}
	if u, err = UpdateFromPerson(c, args.Person); err != nil {
		return err
	}
	reply.Person = u.Person
	return nil
}
