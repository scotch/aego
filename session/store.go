// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"code.google.com/p/gorilla/sessions"
	//"config"
)

// Store is a gorilla/session store.
var Store = sessions.NewCookieStore([]byte("123456789"))

//var Store = sessions.NewCookieStore([]byte(config.SecretKey))
