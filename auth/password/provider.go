// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hal/auth/password provides a password strategy using bcrypt.

auth/password stategy takes a POST with the following keys:

Email (required)
Password.New (required/optional)
Password.Current (required/optional)
Name.GivenName
Name.FamilyName
* (Any other Person attributes)

Based on the supplied attributes auth/password will do one of three things:

1. Create new User and log them in. POST:
  - "Password.New" (present)
  - "Password.Current" (NOT present)
  - + Person attributes, E.g. "Name.GivenName", "Name.FamilyName"

2. Login User or return error if password does not match. POST:
  - "Password.Current" (present)
  - "Password.New" (NOT present)

3. Update Password / Person details. POST:
  - "Password.New" (present)
  - "Password.Current" (present)
  - + Person attributes, E.g. "Name.GivenName", "Name.FamilyName"

*/
package password

import (
	"code.google.com/p/gorilla/schema"
	"errors"
	"github.com/scotch/hal/auth/profile"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/email"
	"github.com/scotch/hal/person"
	"github.com/scotch/hal/user"
	"net/http"
)

var (
	ErrProfileNotFound = errors.New("auth/password: profile not found for email address")
)

// Provider represents the auth.Provider
type Provider struct {
	Name, URL string
}

// New creates a New provider.
func New() *Provider {
	return &Provider{"Password", ""}
}

func decodePerson(r *http.Request) *person.Person {
	// Decode the form data and add the resulting Person type to the Profile.
	p := &person.Person{}
	decoder := schema.NewDecoder()
	decoder.Decode(p, r.Form)
	return p
}

// func validate(e *email.Email, p *Password) error {
// 	// Validate pasword
// 	if err := p.Validate(); err != nil {
// 		return err
// 	}
// 	// Validate email
// 	if err := email.Validate(pers.Email); err != nil {
// 		return err
// 	}
// 	return
// }
// 
// func getUserID(w http.ResponseWriter, r *http.Request, emailAddress *email.Email) (id string) {
// 	// TODO: User merge if the session UserID is different then the email UserID
// 	// search session
// 	sessID, _ := user.CurrentUserID(r)
// 	if sessID != "" {
// 		return sessID
// 	}
// 	// search by email
// 	c := context.NewContext(r)
// 	e, err := email.Get(c, emailAddress)
// 	if err != nil {
// 		return ""
// 	}
// }
// 
// func create(w http.ResponseWriter, r *http.Request,
// 	pf *profile.Profile, pass *Password, pers *person.Person, userID string) (err error) {
// 
// 	c := context.NewContext(r)
// 	if err = validate(pers.Email, pass); err != nil {
// 		return err
// 	}
// 	userID := getUserID(w, r, pers.Email)
// 	// if we have a user ID check for a profile
// 	if userID != "" {
// 		pid := profile.GenAuthID("Password", userID)
// 		if err = profile.Get(c, pid, pf); err == nil {
// 			hasProfile = true
// 		}
// 	}
// 	userID, _ = user.AllocateID(c)
// 	passHash, _ := GenerateFromPassword([]byte(pass.New))
// 
// 	pf.ID = userID
// 	pf.UserID = userID
// 	pf.Auth = passHash
// 	pf.Person = pers
// 
// 	return
// }
// 
// func login(w http.ResponseWriter, r *http.Request,
// 	pf *profile.Profile, pass *Password, pers *person.Person, userID string) (err error) {
// 
// 	c := context.NewContext(r)
// 	if err = validate(pers.Email, pass); err != nil {
// 		return err
// 	}
// 	userID := getUserID(w, r, pers.Email)
// 	// if we have a user ID check for a profile
// 	if userID == "" {
// 		return ErrProfileNotFound
// 	}
// 	pid := profile.GenAuthID("Password", userID)
// 	if err = profile.Get(c, pid, pf); err != nil {
// 		return ErrProfileNotFound
// 	}
// 	if err := CompareHashAndPassword(pf.Auth, []byte(pass.Current)); err != nil {
// 		return err
// 	}
// 	return nil
// 	userID, _ = user.AllocateID(c)
// 	passHash, _ := GenerateFromPassword([]byte(pass.New))
// 
// 	pf.ID = userID
// 	pf.UserID = userID
// 	pf.Auth = passHash
// 	pf.Person = pers
// 
// 	if !hasProfile {
// 		return ErrProfileNotFound
// 	}
// 	// Check
// 	if err := CompareHashAndPassword(pf.Auth, []byte(pass.Current)); err != nil {
// 		return err
// 	}
// }

func authenticate(w http.ResponseWriter, r *http.Request,
	pf *profile.Profile, pass *Password, pers *person.Person, userID string) (err error) {

	c := context.NewContext(r)
	//var userID string
	var e *email.Email
	var passHash []byte
	var hasProfile bool

	// Validate

	// Validate pasword
	if err := pass.Validate(); err != nil {
		return err
	}
	// Validate email
	if err := email.Validate(pers.Email); err != nil {
		return err
	}

	// Search for existing account

	// Search for Session
	// Search Email
	if e, err = email.Get(c, pers.Email); err == nil {
		if e.UserID != "" {
			if e.UserID != userID {
				// TODO handle user account merge.
			}
			// Set the useID to the UserID stored in the saved Email
			userID = e.UserID
		}
	}
	// Search for Profile
	if userID != "" {
		pid := profile.GenAuthID("Password", userID)
		if err = profile.Get(c, pid, pf); err == nil {
			hasProfile = true
		}
	}
	// Update or Login
	if pass.Current != "" {
		if !hasProfile {
			return ErrProfileNotFound
		}
		// Check
		if err := CompareHashAndPassword(pf.Auth, []byte(pass.Current)); err != nil {
			return err
		}
	}
	// Update or Create
	if pass.New != "" {
		// If there is an existing profile, check the password
		if hasProfile && pass.Current == "" {
			if err := CompareHashAndPassword(pf.Auth, []byte(pass.New)); err != nil {
				return err
			}
		} else {
			passHash, err = GenerateFromPassword([]byte(pass.New))
			pf.Auth = passHash
		}
		pf.ID = e.Address
		if userID != "" {
			pf.UserID = userID
		}
		pf.Person = pers
	}
	pf.ID = pf.UserID
	return
}

// Authenticate process the request and returns a populated Profile.
// If the Authenticate method can not authenticate the User based on the
// request, an error or a redirect URL wll be return.
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request) (
	pf *profile.Profile, url string, err error) {

	p.URL = r.URL.Host
	pf = profile.New(p.Name, p.URL)

	pass := &Password{
		New:     r.FormValue("Password.New"),
		Current: r.FormValue("Password.Current"),
	}
	pers := decodePerson(r)
	userID, _ := user.CurrentUserID(r)
	err = authenticate(w, r, pf, pass, pers, userID)
	return pf, "", err
}
