// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

package context

import (
	"appengine"
	"config"
	"net/http"
	"sync"
)

var startOnce sync.Once

type Context appengine.Context

// NewContext returns a new context for an in-flight HTTP request.
func NewContext(req *http.Request) appengine.Context {

	startOnce.Do(func() {
		config.Start(req)
	})

	return appengine.NewContext(req)
}
