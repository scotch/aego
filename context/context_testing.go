// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

package context

import (
	"net/http"

	"appengine"
	"appenginetesting"
)

type Context appengine.Context

var cntx appengine.Context

// New returns a new testing context.
func NewContext(r *http.Request) appengine.Context {
	if cntx == nil {
		var err error
		cntx, err = appenginetesting.NewContext(nil)
		if err != nil {
			panic(err)
		}
	}
	return cntx
}

// Close closes a testing context registered when New() is called.
func Close() {
	if cntx != nil {
		cntx.(*appenginetesting.Context).Close()
		cntx = nil
	}
}
