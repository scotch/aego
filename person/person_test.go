// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package person

import (
	"github.com/scotch/hal/context"
	"testing"
)

func setUp() {}

func tearDown() {
	context.Close()
}

func TestPerson(t *testing.T) {

	setUp()
	defer tearDown()

	c := context.NewContext(nil)

	p := New()

	if !p.Can(c, "write", key) {
		t.Error(`User should be able to "write" their own User object`)
	}
}
