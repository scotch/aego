// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"appengine"
	"code.google.com/p/gorilla/rpc"
	"code.google.com/p/gorilla/rpc/json"
	"errors"
	"github.com/scotch/hal/email"
	"github.com/scotch/hal/types"
	"github.com/scotch/hal/user"
	"net/http"
)

var (
	ErrEmailInUse = errors.New("hal/user: email in use")
)

func init() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(UserService), "User")
	http.Handle("/-/api/v1/users", s)
}

type APIError struct {
	Code    int
	Message string
}

type ChangePasswordArgs struct {
	Current, New string
}

type ChangePasswordReply struct {
	Message string
}

type Empty struct{}

type Person struct {
	Person *types.Person
	Error  *Error
}

type UserService struct{}

func (us *UserService) Current(r *http.Request,
	args *Empty, reply *Person) error {

	u, err := user.Current(r)

	if err != nil {
		reply.Error = err.Error()
		return nil
	}
	reply.Person = u.Person
	return nil
}

func (us *UserService) Create(r *http.Request,
	args *Person, reply *Person) error {

	c := appengine.NewContext(r)
	u, err := user.CreateFromPerson(c, args.Person)
	if err != nil {
		var e *APIError
		switch err {
		case email.ErrInvalidAddress:
			e = &APIError{"ErrInvalidAddress", Message: err.Error()}
		}
		reply.Error = err.Error()
		return nil
	}
	reply.Person = u.Person
	return nil
}

func (us *UserService) ChangePassword(r *http.Request,
	args *ChangePasswordArgs, reply *ChangePasswordReply) error {

	reply.Message = "Current Pasword: " + args.Current + " New Password: " + args.New
	return nil
}
