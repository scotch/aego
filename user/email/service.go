// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

// import (
// 	"errors"
// 	"github.com/scotch/hal/user"
// 	"net/http"
// )
// 
// var (
// 	ErrUserUnauthorized = errors.New("user: unauthorized")
// )
// 
// type Empty struct{}
// 
// type Emails struct {
// 	Emails *[]Email
// }
// 
// type Service struct{}
// 
// func (s *Service) All(w http.ResponseWriter, r *http.Request,
// 	args *Empty, reply *Emails) (err error) {
// 
// 	var u *user.User
// 	var ee []*Emails
// 	if u, err = Current(r); err != nil {
// 		return err
// 	}
// 	if ee, err := GetMulti(c, u.Emails); err != nil {
// 		return err
// 	}
// 	reply.Emails = ee
// 	return nil
// }
