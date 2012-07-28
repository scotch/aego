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

// Authenticate process the request and returns a populated Profile.
// If the Authenticate method can not authenticate the User based on the
// request, an error or a redirect URL wll be return.
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request,
	pf *profile.Profile) (url string, err error) {

	c := context.NewContext(r)
	var userID string
	var e *email.Email
	var passHash []byte
	var hasProfile bool

	p.Name = "Password"
	p.URL = r.URL.Host
	pf.Provider = p.Name
	pf.ProviderURL = p.URL

	// Validate

	// Validate pasword
	pass := &Password{
		New:     r.FormValue("Password.New"),
		Current: r.FormValue("Password.Current"),
	}
	if err := pass.Validate(); err != nil {
		return "", err
	}
	// Validate email
	if err := email.Validate(r.FormValue("Email")); err != nil {
		return "", err
	}

	// Search for existing account

	// Search Session
	userID, _ = user.CurrentUserID(r)
	// Search Email
	if e, err = email.Get(c, r.FormValue("Email")); err == nil {
		if e.UserID != "" {
			if e.UserID != userID {
				// TODO handle user accont merge.
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
			return "", ErrProfileNotFound
		}
		// Check
		if err := CompareHashAndPassword(pf.Auth, []byte(pass.Current)); err != nil {
			return "", err
		}
	}
	// Update or Create
	if pass.New != "" {
		// I there is an Existing Profile. Log chedk the password
		if hasProfile && pass.Current == "" {
			if err := CompareHashAndPassword(pf.Auth, []byte(pass.New)); err != nil {
				return "", err
			}
		} else {
			passHash, err = GenerateFromPassword([]byte(pass.New))
			pf.Auth = passHash
		}
		pf.ID = e.Address
		pf.UserID = userID
		pf.Person = decodePerson(r)
	}
	return "", nil
}

// func LoginByEmailAndPassword(w http.ResponseWriter, r *http.Request, emailAddress, password string) (u *User, err error) {
// 
// 	c := appengine.NewContext(r)
// 	// Get UserID
// 	e, err := email.Get(c, emailAddress)
// 	if err != nil {
// 		return
// 	}
// 	u, err = Get(c, e.UserID)
// 	if err != nil {
// 		return
// 	}
// 	// Compare pasword
// 	if err = CompareHashAndPassword(u.Password, []byte(password)); err != nil {
// 		return
// 	}
// 	// We made it. Log in the User
// 	err = CurrentUserSetID(w, r, u.Key.StringID())
// 	return
// }
