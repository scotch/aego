// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	"github.com/scotch/hal/auth"
	"github.com/scotch/hal/auth/profile"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/person"
	"github.com/scotch/hal/user"
	"net/http"
)

type Service struct{}

type Args struct {
	Password *Password
	Person   *person.Person
}

type Reply struct {
	Person *person.Person
}

func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request,
	args *Args, reply *Reply) (err error) {

	c := context.NewContext(r)
	userID, _ := user.CurrentUserIDByEmail(r, args.Password.Email)
	pf, err := authenticate(c, args.Password, args.Person, userID)
	if err != nil {
		return err
	}
	if _, err = auth.CreateAndLogin(w, r, pf); err != nil {
		return err
	}
	reply.Person = pf.Person
	return nil
}

func (s *Service) IsSet(w http.ResponseWriter, r *http.Request,
	args *Args, reply *Args) (err error) {

	c := context.NewContext(r)
	var isSet bool
	userID, _ := user.CurrentUserID(r)
	_, err = profile.Get(c, profile.GenAuthID("Password", userID))
	if err == nil {
		isSet = true
	}
	reply.Password = &Password{IsSet: isSet}
	return nil
}
