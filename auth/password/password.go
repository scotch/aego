// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hal/auth/password provides a password strategy using bcrypt.

*/
package password

import (
	"code.google.com/p/go.crypto/bcrypt"
	"code.google.com/p/gorilla/schema"
	"errors"
	"github.com/scotch/hal/auth/profile"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/person"
	"net/http"
	"strings"
)

const (
	MIN_PASS_LENGTH = 4
	MAX_PASS_LENGTH = 31
	BCRYPT_COST     = 10
)

//ErrMismatchedPassword is the the custom error for incorrect passwords.
var (
	ErrInvalidEmail     = errors.New("auth/password: invalid email address")
	ErrInvalidPassword  = errors.New("auth/password: invalid password")
	ErrPasswordMismatch = errors.New("auth/password: passwords do not match")
	ErrPasswordLength   = errors.New("auth/password: passwords must be between 4 and 31 charaters")
)

type Provider struct {
	Name, URL string
}

// New creates a New provider.
func New() *Provider {
	return &Provider{"Password", ""}
}

// validateEmail returns true if the supplied string contains
// an `@` and a `.`
func validateEmail(email string) error {
	// TODO maybe use a regex here instead?
	if ok := strings.Contains(email, "@"); ok == false {
		return ErrInvalidEmail
	}
	if strings.Contains(email, ".") == false {
		return ErrInvalidEmail
	}
	return nil
}

// validatePass returns true if the supplied string is
// between 4 and 31 character.
func validatePass(pass string) error {
	if len(pass) < MIN_PASS_LENGTH {
		return ErrPasswordLength
	}
	if len(pass) > MAX_PASS_LENGTH {
		return ErrPasswordLength
	}
	return nil
}

// Authenticate process the request and returns a populated Profile.
// If the Authenticate method can not authenticate the User based on the
// request, an error or a redirect URL wll be return.
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request,
	up *profile.Profile) (url string, err error) {

	c := context.NewContext(r)
	p.URL = r.URL.Host

	email := r.FormValue("Email")
	pass := r.FormValue("Password")

	up.Provider = p.Name
	up.ProviderURL = p.URL
	// Validate email
	if err := validateEmail(email); err != nil {
		return "", err
	}
	// Validate pasword
	if err := validatePass(pass); err != nil {
		return "", err
	}

	authID := profile.GenAuthID("Password", email)
	err = profile.Get(c, authID, up)

	passByte := []byte(pass)
	if err != nil {
		passHash, err := bcrypt.GenerateFromPassword(passByte, BCRYPT_COST)
		if err != nil {
			return "", err
		}
		up.Auth = passHash
		up.ID = email
		// Decode the form data and add the resulting Person type to the Profile.
		per := &person.Person{}
		decoder := schema.NewDecoder()
		decoder.Decode(per, r.Form)
		up.Person = per

	} else {
		// Verify the password.
		if bcrypt.CompareHashAndPassword(up.Auth, passByte) != nil {
			return "", ErrInvalidPassword
		}
	}
	return "", nil
}
