// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	"github.com/scotch/hal/auth"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/user_profile"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setUp() {}

func tearDown() {
	context.Close()
}

func TestValidatePass(t *testing.T) {
	//ds.Register("UserProfile", true, false, false)
}

func TestValidateEmail(t *testing.T) {
	if q := validateEmail("fakemail"); q != ErrInvalidEmail {
		t.Errorf(`validateEmail("fakeemail") = %q, want false`, q)
	}
	if q := validateEmail("fake@email"); q != ErrInvalidEmail {
		t.Errorf(`validateEmail("fake@email") = %q, want false`, q)
	}
	if q := validateEmail("fake@email.com"); q != nil {
		t.Errorf(`validateEmail("fake@email.com") = %q, want true`, q)
	}
	if q := validatePass("pas"); q != ErrPasswordLength {
		t.Errorf(`validatePass("pas") = %q, want false`, q)
	}
	if q := validatePass("passw"); q != nil {
		t.Errorf(`validatePass("passw") = %q, want true`, q)
	}
}

func TestAuthenticate(t *testing.T) {
	setUp()
	defer tearDown()
	c := context.NewContext(nil)
	w := httptest.NewRecorder()

	// Register.

	pro := New()
	auth.Register("password", pro)

	// Post.
	v := url.Values{}
	v.Set("Email", "test@example.com")
	v.Set("Password", "secret1")
	v.Set("Gender", "male")
	v.Set("Name.GivenName", "Barack")
	v.Set("Name.FamilyName", "Obama")
	v.Set("AboutMe", "This is a bio about me.")
	body := strings.NewReader(v.Encode())

	req, _ := http.NewRequest("POST",
		"http://localhost:8080/-/auth/password", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	// Process.

	up := user_profile.New()
	u, err := pro.Authenticate(w, req, up)

	// Check.

	if u != "" {
		t.Errorf(`url: %v, want: ""`, u)
	}
	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}

	per := up.Person

	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := per.Name.GivenName; x != "Barack" {
		t.Errorf(`per.Name.GivenName: %q, want %v`, x, "Barack")
	}
	if x := per.Name.FamilyName; x != "Obama" {
		t.Errorf(`per.Name.FamilyName: %q, want %v`, x, "Obama")
	}

	// Round 2: Existing Profile. Correct password.

	_ = up.Put(c)

	v = url.Values{}
	v.Set("Email", "test@example.com")
	v.Set("Password", "secret1")
	v.Set("Name.GivenName", "Berry")
	body = strings.NewReader(v.Encode())

	req, _ = http.NewRequest("POST",
		"http://localhost:8080/-/auth/password", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	// Process.

	up = user_profile.New()
	u, err = pro.Authenticate(w, req, up)

	// Check.

	if u != "" {
		t.Errorf(`url: %v, want: ""`, u)
	}
	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}

	per = up.Person

	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := per.Name.GivenName; x != "Barack" {
		t.Errorf(`per.Name.GivenName: %q, want %v. When confirming an
    existing account the UserProfile should not be modified.`, x, "Barack")
	}

	// Round 3: Existing Profile. In-Correct password.

	v = url.Values{}
	v.Set("Email", "test@example.com")
	v.Set("Password", "fakepass")
	body = strings.NewReader(v.Encode())

	req, _ = http.NewRequest("POST",
		"http://localhost:8080/-/auth/password", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	// Process.

	up = user_profile.New()
	u, err = pro.Authenticate(w, req, up)

	// Check.

	if err != ErrInvalidPassword {
		t.Errorf(`err: %v, want: %v`, err, ErrInvalidPassword)
	}
}
