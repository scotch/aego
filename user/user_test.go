// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine/datastore"
	"github.com/scotch/hal/context"
	//"github.com/scotch/hal/ds"
	"testing"
)

func setUp() {}

func tearDown() {
	context.Close()
}

func TestGet(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	// Save it.

	u := New()
	u.Email = "test@example.com"
	u.Key = datastore.NewKey(c, "User", "", 0, nil)
	err := u.Put(c)
	if err != nil {
		t.Errorf(`err: %q, want nil`, err)
	}

	// Get it.

	u2, err := Get(c, u.Key.IntID())
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if u2.Email != "test@example.com" {
		t.Errorf(`u2.Email: %v, want "test@example.com"`, u2.Email)
	}
}
