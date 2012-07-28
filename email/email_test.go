// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"testing"
)

func TestValidate(t *testing.T) {
	if x := Validate("fakemail"); x != ErrInvalidAddress {
		t.Errorf(`validateEmail("fakeemail") = %v, want false`, x)
	}
	if x := Validate("fake@email"); x != ErrInvalidAddress {
		t.Errorf(`validateEmail("fake@email") = %v, want false`, x)
	}
	if x := Validate("fake@email.com"); x != nil {
		t.Errorf(`validateEmail("fake@email.com") = %v, want true`, x)
	}
}
