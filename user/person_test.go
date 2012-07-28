// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

// import (
// 	"github.com/scotch/hal/context"
// 	//"github.com/scotch/hal/email"
// 	//"github.com/scotch/hal/password"
// 	"github.com/scotch/hal/person"
// 	"testing"
// )
// 
// var p1 *person.Person = &person.Person{
// 	Name: &person.PersonName{
// 		GivenName:  "Kyle",
// 		FamilyName: "Finley",
// 	},
// 	Email: "1@examle.org",
// 	Password: &person.PersonPassword{
// 		New: "secret1",
// 	},
// }
// 
// var p2 *person.Person = &person.Person{
// 	Name: &person.PersonName{
// 		GivenName:  "Kyle",
// 		FamilyName: "Finley",
// 	},
// 	Email: "2@examle.org",
// 	Password: &person.PersonPassword{
// 		Current: "secret1",
// 		New:     "secret2",
// 	},
// }
// 
// var p3 *person.Person = &person.Person{
// 	Name: &person.PersonName{
// 		GivenName:  "Kyle",
// 		FamilyName: "Finley",
// 	},
// 	Email: "3@examle.org",
// 	Password: &person.PersonPassword{
// 		New: "secret1",
// 	},
// }
// 
// func TestPutByPerson(t *testing.T) {
// 	c := context.NewContext(nil)
// 	defer tearDown()
// 
// 	var err error
// 	var u *User
// 
// 	// Round #1 New User with email & password
// 	// Save it.
// 	if u, err = PutByPerson(c, p1); err != nil {
// 		t.Errorf(`err: %v, want nil`, err)
// 	}
// 
// 	// Check User
// 	// Get from ds to confirm save
// 	if u, err = Get(c, u.Key.StringID()); err != nil {
// 		t.Errorf(`err: %v, want nil`, err)
// 	}
// 	if u.Email != p1.Email {
// 		t.Errorf(`u.Email: %v, want %v`, u.Email, p1.Email)
// 	}
// 	if u.Person.ID != u.Key.StringID() {
// 		t.Errorf(`u.Person.ID: %v, want %v`, u.Person.ID, u.Key.StringID())
// 	}
// 	// Check password
// 	if err = password.CompareHashAndPassword(u.Password, []byte("secret1")); err != nil {
// 		t.Errorf(`Password hash does not match`)
// 	}
// 	if u.Person.Password.IsSet != true {
// 		t.Errorf(`u.Person.Password.IsSet: %v, want true`, u.Person.Password.IsSet)
// 	}
// 	if u.Person.Password.New != "" {
// 		t.Errorf(`u.Person.Password.New: %v, want ""`, u.Person.Password.New)
// 	}
// 	if u.Person.Password.Current != "" {
// 		t.Errorf(`u.Person.Password.Current: %v, want ""`, u.Person.Password.Current)
// 	}
// 	// Check Email - An email entity should have been created for each email
// 	//if _, err = email.Get(c, p1.Email); err != nil {
// 	//t.Errorf(`err: %v, want nil`, err)
// 	//}
// }
