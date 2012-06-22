// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine/datastore"
	"github.com/scotch/hal/context"
	//"github.com/scotch/hal/ds"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setUp() {}

func tearDown() {
	context.Close()
}

func TestCurrentUserID(t *testing.T) {

	// Get CurrentUserID.

	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	userID, err := CurrentUserID(r)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if userID != 0 {
		t.Errorf(`userID: %v, want 0`, userID)
	}
}

func TestSetCurrentUserID(t *testing.T) {

	// Login User.

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	var userID int64 = 1
	err := CurrentUserSetID(w, r, userID)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Get CurrentUserID.

	id, err := CurrentUserID(r)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if id != 1 {
		t.Errorf(`userID: %v, want 1`, userID)
	}
}

func TestCurrent(t *testing.T) {
	setUp()
	defer tearDown()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	c := context.NewContext(r)

	// Create a User.

	u := New()
	u.Email = "test@example.com"
	key := datastore.NewKey(c, "User", "", 1, nil)
	_, err := u.Put(c, key)
	if err != nil {
		t.Fatalf(`err: %v, want nil`, err)
	}

	// Current.

	key, err = Current(r, u)
	if err != ErrNoLoggedInUser {
		t.Errorf(`err: %q, want %q`, err, ErrNoLoggedInUser)
	}

	// Login User.

	var userID int64 = 1
	err = CurrentUserSetID(w, r, userID)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Get Current.

	u2 := &User{}
	id, err := CurrentUserID(r)
	if id != 1 {
		t.Errorf(`userID: %v, want 1`, userID)
	}
	key2, err := Current(r, u2)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if key2.IntID() != 1 {
		t.Errorf(`key2.IntID(): %v, want 1`, key2.IntID())
	}
	if u2.Email != "test@example.com" {
		t.Errorf(`u2.Email: %v, want test@example.com`, u2.Email)
	}
	// Check Person
	if u2.Person.ID != "1" {
		t.Errorf(`u2.Person.ID: %v, want "1"`, u2.Person.ID)
	}
	if u2.Person.Created != u2.Created.Unix() {
		t.Errorf(`u2.Created: %v, want %v`, u2.Person.Created,
			u2.Created.Unix())
	}
	if u2.Person.Updated != u2.Updated.Unix() {
		t.Errorf(`u2.Updated: %v, want %v`, u2.Person.Updated,
			u2.Updated.Unix())
	}
	//if u2.Person.Email != u2.Email {
	//t.Errorf(`u2.Email: %v, want %v`, u2.Person.Email, u2.Email)
	//}
	if u2.Person.Password.IsSet != false {
		t.Errorf(`u2.Person.Password.IsSet: %v, want %v`,
			u2.Person.Password.IsSet, false)
	}
	// Logout User

	err = Logout(w, r)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Confirm logged out.

	key, err = Current(r, u)
	if err != ErrNoLoggedInUser {
		t.Errorf(`err: %q, want %q`, err, ErrNoLoggedInUser)
	}
}

func TestGet(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	// Save it.

	u := New()
	u.Email = "test@example.com"
	key := datastore.NewKey(c, "User", "", 0, nil)
	key, err := u.Put(c, key)
	if err != nil {
		t.Errorf(`err: %q, want nil`, err)
	}

	// Get it.

	u2, _, err := Get(c, key.IntID())
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if u2.Email != "test@example.com" {
		t.Errorf(`u2.Email: %v, want "test@example.com"`, u2.Email)
	}
}
