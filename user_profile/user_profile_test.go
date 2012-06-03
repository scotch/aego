// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user_profile

import (
	"appengine/datastore"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/types"
	"testing"
)

func tearDown() {
	context.Close()
}

func TestNewKey(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	k1 := datastore.NewKey(c, "UserProfile", "google|12345", 0, nil)
	k2 := NewKey(c, "Google", "12345")
	if k1.String() != k2.String() {
		t.Errorf("k2: %q, want %q.", k2, k1)
		t.Errorf("k1:", k1)
		t.Errorf("k2:", k2)
	}
}

func TestGet(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	// Save it.

	u := New()
	u.ID = "12345"
	u.Provider = "Google"
	key := NewKey(c, "google", "12345")
	u.Key = key
	err := u.Put(c)
	if err != nil {
		t.Errorf(`err: %q, want nil`, err)
	}

	// Get it.

	u2 := &UserProfile{}
	id := "google|12345"
	key = datastore.NewKey(c, "UserProfile", id, 0, nil)
	err = ds.Get(c, key, u2)
	if err != nil {
		t.Errorf(`err: %q, want nil`, err)
	}
	err = Get(c, id, u2)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if u2.ID != "12345" {
		t.Errorf(`u2.ID: %v, want "1"`, u2.ID)
	}
	if u2.Key.StringID() != "google|12345" {
		t.Errorf(`uKey.StringID(): %v, want "google|12345"`, u2.Key.StringID())
	}
}

func TestSetPerson(t *testing.T) {

	// Encode it.

	p := &types.Person{
		Name: &types.PersonName{
			GivenName:  "Barack",
			FamilyName: "Obama",
		},
	}
	u := New()
	err := u.SetPerson(p)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Confirm.

	p, err = u.Person()

	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := p.Name.GivenName; x != "Barack" {
		t.Errorf(`p.Name.GivenName: %q, want %v`, x, "Barack")
	}
	if x := p.Name.FamilyName; x != "Obama" {
		t.Errorf(`p.Name.FamilyName: %q, want %v`, x, "Obama")
	}
}
