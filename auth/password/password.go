// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	//"appengine"
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	//"github.com/scotch/hal/email"
)

const (
	PASSWORD_LENGTH_MIN = 4
	PASSWORD_LENGTH_MAX = 31
	BCRYPT_COST         = 10
)

var (
	ErrPasswordMismatch = errors.New("auth/password: passwords do not match")
	ErrPasswordLength   = errors.New("auth/password: passwords must be between 4 and 31 charaters")
)

type Password struct {
	// New: the new password.
	New string `json:"new,omitempty"`

	// Current: The current password.
	Current string `json:"current,omitempty"`

	// IsSet: Indictor that the User has created a password.
	IsSet bool `json:"isSet"`
}

// validatePasswordLength returns true if the supplied string is
// between 4 and 31 character.
func Validate(p string) error {
	if len(p) < PASSWORD_LENGTH_MIN {
		return ErrPasswordLength
	}
	if len(p) > PASSWORD_LENGTH_MAX {
		return ErrPasswordLength
	}
	return nil
}

func (p *Password) Validate() (err error) {
	if p.New != "" {
		if err = Validate(p.New); err != nil {
			return
		}
	}
	if p.Current != "" {
		if err = Validate(p.Current); err != nil {
			return
		}
	}
	return
}

func GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, BCRYPT_COST)
}

func CompareHashAndPassword(hash, password []byte) error {
	if bcrypt.CompareHashAndPassword(hash, password) != nil {
		return ErrPasswordMismatch
	}
	return nil
}

// func ChangePassword(c appengine.Context,
// 	emailAddress, currentPassword, newPassword string) (err error) {
// 
// 	// Confirm the address is a valid email.
// 	err = email.Validate(emailAddress)
// 	if err != nil {
// 		return
// 	}
// 	// Get the UserID.
// 	// TODO(kylefinley) add status check confirm that the email has been
// 	// confirmed.
// 	e, err := email.Get(c, emailAddress)
// 	if err != nil {
// 		return
// 	}
// 	u, err := Get(c, e.UserID)
// 	if err != nil {
// 		return
// 	}
// 	// Compare pasword
// 	if err = pass.CompareHashAndPassword(u.Password,
// 		[]byte(currentPassword)); err != nil {
// 		return
// 	}
// 	// Set password hash to new value
// 	err = u.setPassword(newPassword)
// 	if err != nil {
// 		return
// 	}
// 	err = u.Put(c)
// 	if err != nil {
// 		return
// 	}
// 	return
// }
