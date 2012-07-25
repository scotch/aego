// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package api provides rpc service for Users.
*/

package api

import (
	"github.com/scotch/hal/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() {
}

func teardown() {
	context.Close()
}

func TestCurrent(t *testing.T) {
	setup()
	defer teardown()

	r, _ := http.NewRequest("GET", "/-/api/v1", nil)
	w := httptest.NewRecorder()

}
