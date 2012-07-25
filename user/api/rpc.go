// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package api provides rpc service for Users.
*/

package api

import (
	"appengine"
	"appengine/datastore"
	"code.google.com/p/gorilla/rpc"
	"code.google.com/p/gorilla/rpc/json"
	"errors"
	"github.com/scotch/hal/api"
	"github.com/scotch/hal/types"
	"github.com/scotch/hal/user"
	"net/http"
)

var (
	ErrEmailInUse = errors.New("user: email in use")
)

func init() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(UserService), "User")
	http.Handle("/-/api/v1/users", s)
}

type Empty struct{}

type ErrorReply struct {
	Error *api.Error
}

type Person struct {
	Person *types.Person
	Error  *api.Error
}

type UserService struct{}

func (us *UserService) Current(w http.ResponseWriter, r *http.Request,
	args *Empty, reply *Person) error {

	c := appengine.NewContext(r)
	u, err := user.Current(r)
	if err != nil {
		reply.Error = api.ConvertError(err)
		return nil
	}
	reply.Person = u.Person
	return nil
}

func (us *UserService) Login(w http.ResponseWriter, r *http.Request,
	args *Person, reply *Person) error {

	u, err := user.LoginByEmailAndPassword(w, r,
		args.Person.Email, args.Person.Password.New)
	if err != nil {
		reply.Error = api.ConvertError(err)
		return nil
	}
	reply.Person = u.Person
	return nil
}

func (us *UserService) Logout(w http.ResponseWriter, r *http.Request,
	args *Person, reply *Person) error {

	if err := user.Logout(w, r); err != nil {
		reply.Error = api.ConvertError(err)
		return nil
	}
	return nil
}

func (us *UserService) Create(w http.ResponseWriter, r *http.Request,
	args *Person, reply *Person) error {

	c := appengine.NewContext(r)
	u, err := user.CreateFromPerson(c, args.Person)
	if err != nil {
		reply.Error = api.ConvertError(err)
		return nil
	}
	reply.Person = u.Person
	return nil
}

func (us *UserService) Update(w http.ResponseWriter, r *http.Request,
	args *Person, reply *Person) error {

	c := appengine.NewContext(r)
	u, err := user.Current(r)
	k := datastore.NewKey(c, "User", args.Person.ID, 0, nil)
	if can := u.Can(c, "write", k); can == false {
		reply.Error = &api.Error{Code: 401, Message: "user: unauthorized"}
		return nil
	}
	if u, err = user.UpdateFromPerson(c, args.Person); err != nil {
		reply.Error = api.ConvertError(err)
		return nil
	}
	reply.Person = u.Person
	return nil
}

type ChangePasswordArgs struct {
	Email, Current, New string
}

func (us *UserService) ChangePassword(w http.ResponseWriter, r *http.Request,
	args *ChangePasswordArgs, reply *ErrorReply) error {

	c := appengine.NewContext(r)
	if err := user.ChangePassword(c, args.Email, args.Current, args.New); err != nil {
		reply.Error = api.ConvertError(err)
		return nil
	}
	return nil
}
