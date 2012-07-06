// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/scotch/hal/email"
	"github.com/scotch/hal/password"
	"github.com/scotch/hal/user"
)

type Error struct {
	Code    int
	Message string
}

func ErrorCode(err error) int {
	switch err {
	// User
	case user.ErrNoLoggedInUser:
		return 100
	// Email
	case email.ErrInvalidAddress:
		return 101
	case user.ErrEmailInUse:
		return 102
	// Password
	case password.ErrPasswordLength:
		return 103
	case password.ErrPasswordMismatch:
		return 104
	}
	return 0
}

func ConvertError(err error) *Error {
	return &Error{Code: ErrorCode(err), Message: err.Error()}
}
