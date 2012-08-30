// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package code

// Error codes
const (
	NoLoggedInUser        = 100
	InvalidEmailAddress   = 101
	InvalidPassword       = 102
	InvalidPasswordLength = 103
)

var errorText = map[int]string{
	NoLoggedInUser:        "No Logged In User",
	InvalidEmailAddress:   "Invalid Email Address",
	InvalidPasswordLength: "Invalid Password Length",
	InvalidPassword:       "Invalid Password",
}

// ErrorText returns a text for the error code. It returns the empty
// string if the code is unknown.
func ErrorText(code int) string {
	return errorText[code]
}
