// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
)

const (
	PASSWORD_LENGTH_MIN = 4
	PASSWORD_LENGTH_MAX = 31
	BCRYPT_COST         = 10
)

var (
	ErrInvalidPassword  = errors.New("hal/user: invalid password")
	ErrPasswordMismatch = errors.New("hal/user: passwords do not match")
	ErrPasswordLength   = errors.New("hal/user: passwords must be between 4 and 31 charaters")
)

// validatePasswordLength returns true if the supplied string is
// between 4 and 31 character.
func Validate(pass string) error {
	if len(pass) < PASSWORD_LENGTH_MIN {
		return ErrPasswordLength
	}
	if len(pass) > PASSWORD_LENGTH_MAX {
		return ErrPasswordLength
	}
	return nil
}

func GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, BCRYPT_COST)
}

func CompareHashAndPassword(hash, password []byte) error {
	if bcrypt.CompareHashAndPassword(hash, password) != nil {
		return ErrInvalidPassword
	}
	return nil
}
