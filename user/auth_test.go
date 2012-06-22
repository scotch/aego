// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine/datastore"
	"github.com/scotch/hal/context"
	"testing"
)

func TestAddAuthID(t *testing.T) {
	u := New()
	a1 := "a|1"
	a2 := "a|2"
	a3 := "a|1"

	u.AddAuthID(a1)
	if len(u.AuthIDs) != 1 {
		t.Errorf(`len(u.AuthIDs): %v, want %v`, len(u.AuthIDs), 1)
		t.Errorf(`u.AuthIDs: %q`, u.AuthIDs)
	}
	u.AddAuthID(a2)
	if len(u.AuthIDs) != 2 {
		t.Errorf(`len(u.AuthIDs): %v, want %v`, len(u.AuthIDs), 2)
		t.Errorf(`u.AuthIDs: %q`, u.AuthIDs)
	}
	u.AddAuthID(a3)
	if len(u.AuthIDs) != 2 {
		t.Errorf(`len(u.AuthIDs): %v, want %v`, len(u.AuthIDs), 2)
		t.Errorf(`u.AuthIDs: %q`, u.AuthIDs)
	}
}

func TestGetOrInsertByAuthID(t *testing.T) {

	setUp()
	defer tearDown()

	c := context.NewContext(nil)

	authID := "example|1"

	// Create.

	u, err := GetOrInsertByAuthID(c, authID)
	if err != nil {
		t.Errorf(`err: %q, want: %q`, err, nil)
	}
	if len(u.AuthIDs) != 1 {
		t.Errorf(`len(u.AuthIDs): %v, want %v`, len(u.AuthIDs), 1)
		t.Errorf(`u.AuthIDs: %q`, u.AuthIDs)
	}

	// Confirm.

	u2, err := Get(c, u.Key.IntID())
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if len(u2.AuthIDs) != 1 {
		t.Errorf(`len(u2.AuthIDs): %v, want %v`, len(u2.AuthIDs), 1)
		t.Errorf(`u2.AuthIDs: %q`, u2.AuthIDs)
		t.Errorf(`u2: %q`, u2)
	}
	if u2.AuthIDs[0] != authID {
		t.Errorf(`u2.AuthIDs[0]: %v, want %v`, u2.AuthIDs[0], authID)
	}
	q2 := datastore.NewQuery("User")
	if cnt, _ := q2.Count(c); cnt != 1 {
		t.Errorf(`User cnt: %v, want 1`, cnt)
	}

	// Again.

	u = New()
	u.Email = "test@example.com"
	u, err = GetOrInsertByAuthID(c, authID)
	if err != nil {
		t.Errorf(`err: %q, want: %q`, err, nil)
	}

	// Confirm.

	u2, err = Get(c, u.Key.IntID())
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if len(u2.AuthIDs) != 1 {
		t.Errorf(`len(u2.AuthIDs): %v, want %v`, len(u2.AuthIDs), 1)
		t.Errorf(`u2.AuthIDs: %q`, u2.AuthIDs)
		t.Errorf(`u2: %q`, u2)
	}
	if u2.AuthIDs[0] != authID {
		t.Errorf(`u2.AuthIDs[0]: %v, want %v`, u2.AuthIDs[0], authID)
	}
	q2 = datastore.NewQuery("User")
	if cnt, _ := q2.Count(c); cnt != 1 {
		t.Errorf(`User cnt: %v, want 1`, cnt)
	}
}
