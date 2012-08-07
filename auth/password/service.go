// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	//"github.com/scotch/hal/api/code"
	"appengine"
	"github.com/scotch/hal/auth"
	"github.com/scotch/hal/auth/profile"
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

	c := appengine.NewContext(r)
	c.Errorf(`args: %v`, args)
	userID, _ := user.CurrentUserID(r)
	pf := profile.New("Password", r.URL.Host)
	err = authenticate(w, r, pf, args.Password, args.Person, userID)
	if _, err = auth.CreateAndLogin(w, r, pf); err != nil {
		return err
	}
	reply.Person = pf.Person
	c.Errorf(`err: %v`, err)
	c.Errorf(`pf: %v`, pf)
	return err
}
