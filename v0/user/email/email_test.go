// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"testing"
)

var addressTests = []struct {
	in  string
	out error
}{
	{"name@example.com", nil},
	{"name@x.y.z.example.co.uk", nil},
	{"name@example.aero", nil},
	{"x@example.com", nil},
	{"first.last@example.com", nil},
	{"first.middle.last@example.com", nil},
	{"x+y@example.com", nil},
	{"Name@example.com", nil},
	{"NAme@x.y.z.example.co.uk", nil},
	{"NAME@example.aero", nil},
	{"X@example.com", nil},
	{"First.Last@example.com", nil},
	{"FIRST.middle.LAST@example.com", nil},
	{"X+Y@example.com", nil},
	{"name@.com", ErrInvalidAddress},
	{"name@example.x", ErrInvalidAddress},
	{"name@example.", ErrInvalidAddress},
	{"@example.com", ErrInvalidAddress},
	{".@example.com", ErrInvalidAddress},
}

func TestValidate(t *testing.T) {
	for _, tt := range addressTests {
		if err := Validate(tt.in); err != tt.out {
			t.Errorf(`Validate("%v") => %v, want "%v"`, tt.in, err, tt.out)
		}
	}
}
