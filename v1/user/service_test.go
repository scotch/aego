// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package api provides rpc service for Users.
*/

package user

import (
	"github.com/scotch/aego/v1/context"
	//"net/http"
	//"net/http/httptest"
	"testing"
)

func setup() {
}

func teardown() {
	context.Close()
}

func TestServiceCurrent(t *testing.T) {
	setup()
	defer teardown()

	//r, _ := http.NewRequest("GET", "/-/api/v1", nil)
	//w := httptest.NewRecorder()

}
