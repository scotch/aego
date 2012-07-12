// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine/datastore"
	"github.com/scotch/hal/context"
	"testing"
)

func TestCan(t *testing.T) {

	setUp()
	defer tearDown()

	c := context.NewContext(nil)

	u := New()
	// User key
	key := datastore.NewKey(c, "User", "1", 0, nil)
	u.Key = key
	if ok := u.Can(c, "write", key); ok {
		t.Error(`User should be able to "write" their own User object`)
	}
}
