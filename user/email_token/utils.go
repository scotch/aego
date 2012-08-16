// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email_token

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	TOKEN_LENGTH = 30
)

func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

func genToken() string {
	t := make([]byte, TOKEN_LENGTH)
	rand.Read(t)
	return string(encode(t))
}
