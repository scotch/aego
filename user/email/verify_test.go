// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"testing"
)

func setup() {}

func teardown() {
	//context.Close()
}

func Test_genVerifyURL(t *testing.T) {
	url := "http://localhost:8080"
	token := "wildthing12345"
	have := genVerifyURL(url, token)
	want := "http://localhost:8080/-/email/verify?token=" + token
	if have != want {
		t.Errorf(`url: %s, want: %s"`, have, want)
	}
}
